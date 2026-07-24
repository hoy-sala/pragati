package question_import

import (
	"fmt"
	"regexp"
	"strings"
)

type ImportedQuestion struct {
	QuestionType string   `json:"question_type"`
	QuestionText string   `json:"question_text"`
	Options      []Option `json:"options"`
	Answer       string   `json:"answer"`
	Marks        float64  `json:"marks"`
	Difficulty   string   `json:"difficulty"`
	Chapters     []string `json:"chapters"`
	Tags         []string `json:"tags"`
	Line         int      `json:"line"`
}

type Option struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Correct bool   `json:"correct"`
}

type ImportResult struct {
	Imported int                 `json:"imported"`
	Errors   []ImportRowError    `json:"errors,omitempty"`
}

type ImportRowError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

var mcqPattern = regexp.MustCompile(`(?s)^(.+?)\{\s*(.*?)\s*\}$`)
var tfPattern = regexp.MustCompile(`(?i)\{(TRUE|FALSE|T|F)\}`)
var fillPattern = regexp.MustCompile(`\{([^}]+)\}`)
var optionPattern = regexp.MustCompile(`([~=])([^~=]*)`)

func ParseGIFT(content string) (*ImportResult, []ImportedQuestion) {
	questions := []ImportedQuestion{}
	result := &ImportResult{}
	lines := strings.Split(content, "\n")

	currentLine := 1
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		currentLine++

		if line == "" || strings.HasPrefix(line, "//") {
			i++
			continue
		}

		block := ""
		for i < len(lines) {
			block += lines[i] + "\n"
			i++
			if strings.Contains(lines[i-1], "}") && !strings.Contains(lines[i-1], "{{") {
				break
			}
		}
		block = strings.TrimSpace(block)

		q, err := parseGIFTBlock(block, currentLine)
		if err != nil {
			result.Errors = append(result.Errors, ImportRowError{Line: currentLine, Message: err.Error()})
			continue
		}
		if q != nil {
			questions = append(questions, *q)
			result.Imported++
		}
	}

	return result, questions
}

func parseGIFTBlock(block string, line int) (*ImportedQuestion, error) {
	block = strings.TrimSpace(block)
	if block == "" {
		return nil, nil
	}

	// True/False
	if matches := tfPattern.FindStringSubmatch(block); matches != nil {
		text := tfPattern.ReplaceAllString(block, "")
		answer := strings.ToUpper(matches[1])
		if answer == "T" {
			answer = "TRUE"
		} else if answer == "F" {
			answer = "FALSE"
		}
		return &ImportedQuestion{
			QuestionType: "true_false",
			QuestionText: strings.TrimSpace(text),
			Answer:       answer,
			Difficulty:   "medium",
			Line:         line,
		}, nil
	}

	// Has options -> MCQ or Fill/Short
	if strings.Contains(block, "{") && strings.Contains(block, "}") {
		// Try MCQ first (has ~ or = markers for options)
		if strings.Contains(block, "~") || strings.Contains(block, "=") {
			return parseMCQ(block, line)
		}

		// Fill in blank (simple text inside braces)
		text, answers := parseFillBlank(block)
		if len(answers) > 0 {
			options := []Option{}
			correctAnswer := ""
			for _, a := range answers {
				opt := Option{Key: string(rune('A' + len(options))), Value: a.value, Correct: a.correct}
				options = append(options, opt)
				if a.correct {
					correctAnswer = opt.Key
				}
			}
			return &ImportedQuestion{
				QuestionType: "fill_blank",
				QuestionText: text,
				Options:      options,
				Answer:       correctAnswer,
				Difficulty:   "medium",
				Line:         line,
			}, nil
		}
	}

	return nil, fmt.Errorf("unrecognized question format at line %d", line)
}

type fillAnswer struct {
	value   string
	correct bool
}

func parseFillBlank(block string) (string, []fillAnswer) {
	text := block
	answers := []fillAnswer{}

	for {
		start := strings.Index(text, "{")
		if start < 0 {
			break
		}
		end := strings.Index(text[start:], "}")
		if end < 0 {
			break
		}
		inner := text[start+1 : start+end]
		text = text[:start] + "_____" + text[start+end+1:]

		parts := strings.Split(inner, "|")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if strings.HasPrefix(p, "=") {
				answers = append(answers, fillAnswer{value: strings.TrimPrefix(p, "="), correct: true})
			} else if strings.HasPrefix(p, "~") {
				answers = append(answers, fillAnswer{value: strings.TrimPrefix(p, "~"), correct: false})
			} else if len(parts) == 1 {
				answers = append(answers, fillAnswer{value: p, correct: true})
			}
		}
	}

	return strings.TrimSpace(text), answers
}

func parseMCQ(block string, line int) (*ImportedQuestion, error) {
	braceStart := strings.Index(block, "{")
	braceEnd := strings.LastIndex(block, "}")

	if braceStart < 0 || braceEnd < 0 || braceStart >= braceEnd {
		return nil, fmt.Errorf("invalid MCQ format")
	}

	questionText := strings.TrimSpace(block[:braceStart])
	optionsPart := block[braceStart+1 : braceEnd]

	options := []Option{}
	matches := optionPattern.FindAllStringSubmatch(optionsPart, -1)
	answerKey := ""

	keyIdx := 0
	for _, m := range matches {
		key := string(rune('A' + keyIdx))
		value := strings.TrimSpace(m[2])
		isCorrect := m[1] == "="

		options = append(options, Option{Key: key, Value: value, Correct: isCorrect})
		if isCorrect {
			answerKey = key
		}
		keyIdx++
	}

	if answerKey == "" {
		return nil, fmt.Errorf("no correct answer (marked with =) found in MCQ")
	}

	return &ImportedQuestion{
		QuestionType: "mcq",
		QuestionText: questionText,
		Options:      options,
		Answer:       answerKey,
		Difficulty:   "medium",
		Line:         line,
	}, nil
}
