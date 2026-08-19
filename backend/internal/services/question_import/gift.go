package question_import

import (
	"fmt"
	"regexp"
	"strconv"
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
	Imported int              `json:"imported"`
	Errors   []ImportRowError `json:"errors,omitempty"`
}

type ImportRowError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// ParseGIFT parses GIFT-format content into questions. LaTeX math delimited by
// \(...\) (inline) or \[...\] (display) is treated as opaque text: braces,
// tildes and equals signs inside a math zone carry no GIFT meaning and are
// preserved verbatim for the frontend renderer. GIFT special characters
// outside math zones can be escaped with a backslash (\{ \} \~ \= \\).
func ParseGIFT(content string) (*ImportResult, []ImportedQuestion) {
	questions := []ImportedQuestion{}
	result := &ImportResult{}
	pos := 0
	for pos < len(content) {
		block, next := nextBlock(content, pos)
		block = strings.TrimSpace(block)
		line := 1 + strings.Count(content[:pos], "\n")
		if block != "" {
			q, err := parseGIFTBlock(block, line)
			if err != nil {
				result.Errors = append(result.Errors, ImportRowError{Line: line, Message: err.Error()})
			} else if q != nil {
				questions = append(questions, *q)
				result.Imported++
			}
		}
		if next <= pos {
			break
		}
		pos = next
	}
	return result, questions
}

// nextBlock advances past blank lines and // comments, then collects the next
// question block. A block ends at the first unescaped '}' that closes an
// answer block opened by an unescaped '{' outside a math zone.
func nextBlock(content string, pos int) (string, int) {
	n := len(content)
	i := pos
	for i < n {
		if content[i] == '/' && i+1 < n && content[i+1] == '/' {
			for i < n && content[i] != '\n' {
				i++
			}
			continue
		}
		if c := content[i]; c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			i++
			continue
		}
		break
	}

	var sb strings.Builder
	var inline, display bool
	braceDepth := 0
	for i < n {
		c := content[i]
		if c == '\\' && i+1 < n {
			nc := content[i+1]
			if display {
				if nc == ']' {
					display = false
				}
			} else if inline {
				if nc == ')' {
					inline = false
				}
			} else {
				switch nc {
				case '(':
					inline = true
				case '[':
					display = true
				}
			}
			sb.WriteByte(c)
			sb.WriteByte(nc)
			i += 2
			continue
		}
		if !inline && !display {
			switch c {
			case '{':
				braceDepth++
			case '}':
				braceDepth--
				if braceDepth <= 0 {
					sb.WriteByte(c)
					return sb.String(), i + 1
				}
			}
		}
		sb.WriteByte(c)
		i++
	}
	return sb.String(), n
}

func parseGIFTBlock(block string, line int) (*ImportedQuestion, error) {
	block = strings.TrimSpace(block)
	if block == "" {
		return nil, nil
	}

	block, marks, difficulty, chapters, tags := extractMetadata(block)
	if marks == 0 {
		marks = 1
	}
	if difficulty == "" {
		difficulty = "medium"
	}

	start := scanFor(block, 0, "{")
	if start < 0 {
		return nil, fmt.Errorf("unrecognized question format at line %d (no answer block)", line)
	}
	end := scanFor(block, start+1, "}")
	if end < 0 {
		return nil, fmt.Errorf("unrecognized question format at line %d (unclosed answer block)", line)
	}

	stem := strings.TrimSpace(normalizeText(block[:start]))
	inside := block[start+1 : end]

	// True/False
	switch strings.ToUpper(strings.TrimSpace(inside)) {
	case "TRUE", "T":
		return &ImportedQuestion{
			QuestionType: "true_false",
			QuestionText: stem,
			Answer:       "TRUE",
			Marks:        marks,
			Difficulty:   difficulty,
			Chapters:     chapters,
			Tags:         tags,
			Line:         line,
		}, nil
	case "FALSE", "F":
		return &ImportedQuestion{
			QuestionType: "true_false",
			QuestionText: stem,
			Answer:       "FALSE",
			Marks:        marks,
			Difficulty:   difficulty,
			Chapters:     chapters,
			Tags:         tags,
			Line:         line,
		}, nil
	}

	// MCQ if the answer block contains unescaped ~ or = option markers.
	if scanFor(inside, 0, "~=") >= 0 {
		opts := splitOptions(inside)
		if len(opts) == 0 {
			return nil, fmt.Errorf("invalid MCQ at line %d (no options)", line)
		}
		answerKey := ""
		for _, o := range opts {
			if o.Correct {
				answerKey = o.Key
				break
			}
		}
		if answerKey == "" {
			return nil, fmt.Errorf("no correct answer (marked with =) found in MCQ at line %d", line)
		}
		return &ImportedQuestion{
			QuestionType: "mcq",
			QuestionText: stem,
			Options:      opts,
			Answer:       answerKey,
			Marks:        marks,
			Difficulty:   difficulty,
			Chapters:     chapters,
			Tags:         tags,
			Line:         line,
		}, nil
	}

	// Fill in the blank: single answer block, alternatives separated by |.
	text := stem + "_____" + normalizeText(block[end+1:])
	parts := splitPipe(inside)
	opts := []Option{}
	correctIdx := -1
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		marked := false
		if strings.HasPrefix(p, "=") {
			marked = true
			p = strings.TrimPrefix(p, "=")
		} else if strings.HasPrefix(p, "~") {
			p = strings.TrimPrefix(p, "~")
		}
		if len(parts) == 1 {
			marked = true
		}
		opt := Option{Key: string(rune('A' + len(opts))), Value: normalizeText(strings.TrimSpace(p)), Correct: marked}
		opts = append(opts, opt)
		if marked && correctIdx < 0 {
			correctIdx = len(opts) - 1
		}
	}
	if len(opts) == 0 {
		return nil, fmt.Errorf("invalid fill-blank at line %d (no answer)", line)
	}
	if correctIdx < 0 {
		opts[0].Correct = true
		correctIdx = 0
	}
	return &ImportedQuestion{
		QuestionType: "fill_blank",
		QuestionText: text,
		Options:      opts,
		Answer:       opts[correctIdx].Key,
		Marks:        marks,
		Difficulty:   difficulty,
		Chapters:     chapters,
		Tags:         tags,
		Line:         line,
	}, nil
}

// scanFor returns the index of the first unescaped byte in targets that
// appears outside a math zone (starting at from), or -1 if none.
func scanFor(s string, from int, targets string) int {
	var inline, display bool
	for i := from; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			nc := s[i+1]
			if display {
				if nc == ']' {
					display = false
				}
			} else if inline {
				if nc == ')' {
					inline = false
				}
			} else {
				switch nc {
				case '(':
					inline = true
				case '[':
					display = true
				}
			}
			i++
			continue
		}
		if !inline && !display && strings.IndexByte(targets, c) >= 0 {
			return i
		}
	}
	return -1
}

// splitOptions splits an MCQ answer-block body on unescaped ~/= markers
// outside math zones, returning one Option per marker. The first option
// marked with = becomes the correct answer.
func splitOptions(s string) []Option {
	opts := []Option{}
	var inline, display bool
	keyIdx := 0
	cur := strings.Builder{}
	curMarker := byte(0)
	started := false

	flush := func() {
		if !started {
			return
		}
		v := strings.TrimSpace(cur.String())
		if v == "" {
			return
		}
		opts = append(opts, Option{Key: string(rune('A' + keyIdx)), Value: normalizeText(v), Correct: curMarker == '='})
		keyIdx++
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			nc := s[i+1]
			if display {
				if nc == ']' {
					display = false
				}
			} else if inline {
				if nc == ')' {
					inline = false
				}
			} else {
				switch nc {
				case '(':
					inline = true
				case '[':
					display = true
				case '~', '=':
					cur.WriteByte('\\')
					cur.WriteByte(nc)
					i++
					continue
				}
			}
			cur.WriteByte(c)
			cur.WriteByte(nc)
			i++
			continue
		}
		if !inline && !display && (c == '~' || c == '=') {
			flush()
			curMarker = c
			cur.Reset()
			started = true
			continue
		}
		cur.WriteByte(c)
	}
	flush()
	return opts
}

// splitPipe splits a fill-blank answer body on unescaped | separators outside
// math zones (e.g. {=a | ~b}).
func splitPipe(s string) []string {
	parts := []string{}
	var inline, display bool
	cur := strings.Builder{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			nc := s[i+1]
			if display {
				if nc == ']' {
					display = false
				}
			} else if inline {
				if nc == ')' {
					inline = false
				}
			} else {
				switch nc {
				case '(':
					inline = true
				case '[':
					display = true
				}
			}
			cur.WriteByte(c)
			cur.WriteByte(nc)
			i++
			continue
		}
		if !inline && !display && c == '|' {
			parts = append(parts, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteByte(c)
	}
	parts = append(parts, cur.String())
	return parts
}

// normalizeText produces the stored question/option text: math zones are kept
// verbatim (delimiters included), GIFT escapes outside math are unescaped.
func normalizeText(s string) string {
	var sb strings.Builder
	var inline, display bool
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			nc := s[i+1]
			if display {
				if nc == ']' {
					display = false
				}
			} else if inline {
				if nc == ')' {
					inline = false
				}
			} else {
				switch nc {
				case '(':
					inline = true
					sb.WriteByte('\\')
					sb.WriteByte('(')
					i++
					continue
				case '[':
					display = true
					sb.WriteByte('\\')
					sb.WriteByte('[')
					i++
					continue
				case '{', '}', '~', '=', '\\':
					sb.WriteByte(nc)
					i++
					continue
				}
			}
			sb.WriteByte(c)
			sb.WriteByte(nc)
			i++
			continue
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

var (
	titleRe     = regexp.MustCompile(`^\s*::.*?::\s*`)
	marksRe = regexp.MustCompile(`\{?#\s*([0-9]+(?:\.[0-9]+)?)\}?`)
	diffRe      = regexp.MustCompile(`\[difficulty\s*[:=]\s*(easy|medium|hard)\]`)
	diffBareRe  = regexp.MustCompile(`\[(easy|medium|hard)\]`)
	chapterRe   = regexp.MustCompile(`\[chapter\s*[:=]\s*([^\]]+)\]`)
	tagsRe      = regexp.MustCompile(`\[tags\s*[:=]\s*([^\]]+)\]`)
	tagsSpaceRe = regexp.MustCompile(`\[tags\s+([^\]]+)\]`)
)

// extractMetadata strips GIFT metadata (::title::, {#marks}, [difficulty],
// [chapter:...], [tags:...]) from the block and returns the cleaned block plus
// parsed marks, difficulty, chapters and tags. Defaults: marks 0 (meaning "use
// 1" downstream), difficulty "medium".
func extractMetadata(block string) (string, float64, string, []string, []string) {
	block = titleRe.ReplaceAllString(block, "")

	marks := 0.0
	if m := marksRe.FindStringSubmatch(block); m != nil {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			marks = v
		}
		block = marksRe.ReplaceAllString(block, "")
	}

	difficulty := ""
	if m := diffRe.FindStringSubmatch(block); m != nil {
		difficulty = m[1]
		block = diffRe.ReplaceAllString(block, "")
	} else if m := diffBareRe.FindStringSubmatch(block); m != nil {
		difficulty = m[1]
		block = diffBareRe.ReplaceAllString(block, "")
	}

	chapters := []string{}
	if m := chapterRe.FindStringSubmatch(block); m != nil {
		chapters = append(chapters, strings.TrimSpace(m[1]))
		block = chapterRe.ReplaceAllString(block, "")
	}

	tags := []string{}
	if m := tagsRe.FindStringSubmatch(block); m != nil {
		tags = splitTagList(m[1])
		block = tagsRe.ReplaceAllString(block, "")
	} else if m := tagsSpaceRe.FindStringSubmatch(block); m != nil {
		tags = splitTagList(m[1])
		block = tagsSpaceRe.ReplaceAllString(block, "")
	}

	return block, marks, difficulty, chapters, tags
}

func splitTagList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' || r == ';' || r == '\t' })
	tags := []string{}
	for _, f := range fields {
		if f = strings.TrimSpace(f); f != "" {
			tags = append(tags, f)
		}
	}
	return tags
}
