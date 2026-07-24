package question_import

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type CSVQuestionRow struct {
	QuestionType string
	QuestionText string
	Options      []string // A, B, C, D values
	CorrectAns   string   // A, B, C, D or TRUE/FALSE or text
	Marks        float64
	Difficulty   string
	Chapter      string
	Tags         string
}

func ParseCSV(reader io.Reader) (*ImportResult, []ImportedQuestion) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return &ImportResult{Errors: []ImportRowError{{Line: 0, Message: "failed to parse CSV: " + err.Error()}}}, nil
	}

	result := &ImportResult{}
	questions := []ImportedQuestion{}

	if len(records) < 2 {
		result.Errors = append(result.Errors, ImportRowError{Line: 0, Message: "CSV must have header and at least one data row"})
		return result, nil
	}

	headers := records[0]
	colMap := make(map[string]int)
	for i, h := range headers {
		colMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	for i := 1; i < len(records); i++ {
		row := records[i]
		get := func(col string) string {
			if idx, ok := colMap[col]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		qType := strings.ToLower(get("type"))
		text := get("question")

		if text == "" {
			result.Errors = append(result.Errors, ImportRowError{Line: i + 1, Message: "empty question text"})
			continue
		}

		q := ImportedQuestion{
			QuestionText: text,
			QuestionType: mapType(qType),
			Difficulty:   get("difficulty"),
			Line:         i + 1,
		}

		if q.QuestionType == "" {
			result.Errors = append(result.Errors, ImportRowError{Line: i + 1, Message: "invalid type: " + qType})
			continue
		}
		if q.Difficulty == "" {
			q.Difficulty = "medium"
		}
		if m, err := strconv.ParseFloat(get("marks"), 64); err == nil && m > 0 {
			q.Marks = m
		} else {
			q.Marks = 1
		}

		if chap := get("chapter"); chap != "" {
			q.Chapters = []string{chap}
		}
		if tags := get("tags"); tags != "" {
			q.Tags = strings.Split(tags, ",")
			for i, t := range q.Tags {
				q.Tags[i] = strings.TrimSpace(t)
			}
		}

		switch q.QuestionType {
		case "mcq":
			opts := []Option{}
			correctKey := ""
			for _, key := range []string{"a", "b", "c", "d"} {
				val := get(key)
				if val == "" {
					continue
				}
				upperKey := strings.ToUpper(key)
				isCorrect := strings.EqualFold(get("answer"), upperKey) || strings.EqualFold(get("correct"), upperKey)
				opts = append(opts, Option{Key: upperKey, Value: val, Correct: isCorrect})
				if isCorrect {
					correctKey = upperKey
				}
			}
			if correctKey == "" {
				result.Errors = append(result.Errors, ImportRowError{Line: i + 1, Message: "no correct answer specified for MCQ"})
				continue
			}
			q.Options = opts
			q.Answer = correctKey

		case "true_false":
			ans := strings.ToUpper(get("answer"))
			if ans != "TRUE" && ans != "FALSE" {
				result.Errors = append(result.Errors, ImportRowError{Line: i + 1, Message: "answer must be TRUE or FALSE"})
				continue
			}
			q.Answer = ans

		case "fill_blank", "short_answer":
			ans := get("answer")
			if ans == "" {
				result.Errors = append(result.Errors, ImportRowError{Line: i + 1, Message: "answer is required"})
				continue
			}
			q.Answer = ans
		}

		questions = append(questions, q)
		result.Imported++
	}

	return result, questions
}

func mapType(t string) string {
	switch t {
	case "mcq", "multiple choice", "multiple-choice":
		return "mcq"
	case "true/false", "true_false", "tf", "boolean":
		return "true_false"
	case "fill blank", "fill_blank", "fill-in", "fill":
		return "fill_blank"
	case "short answer", "short_answer", "short":
		return "short_answer"
	}
	return ""
}
