package question_import

import (
	"testing"
)

func TestParseGIFT_MCQ(t *testing.T) {
	input := `Who wrote "Hamlet"? {
	=Shakespeare
	~Milton
	~Chaucer
	~Dickens
	}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	if len(questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(questions))
	}

	q := questions[0]
	if q.QuestionType != "mcq" {
		t.Errorf("expected type 'mcq', got '%s'", q.QuestionType)
	}
	if q.Answer != "A" {
		t.Errorf("expected answer 'A', got '%s'", q.Answer)
	}
	if len(q.Options) != 4 {
		t.Errorf("expected 4 options, got %d", len(q.Options))
	}
	if q.Options[0].Value != "Shakespeare" || !q.Options[0].Correct {
		t.Errorf("first option should be 'Shakespeare' and correct")
	}
}

func TestParseGIFT_TrueFalse(t *testing.T) {
	input := `The Earth is round.{TRUE}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionType != "true_false" {
		t.Errorf("expected type 'true_false', got '%s'", q.QuestionType)
	}
	if q.Answer != "TRUE" {
		t.Errorf("expected answer 'TRUE', got '%s'", q.Answer)
	}
}

func TestParseGIFT_TrueFalseShort(t *testing.T) {
	input := `The sky is blue.{T}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.Answer != "TRUE" {
		t.Errorf("expected answer 'TRUE', got '%s'", q.Answer)
	}
}

func TestParseGIFT_FillBlank(t *testing.T) {
	input := `The capital of France is {Paris}.`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionType != "fill_blank" {
		t.Errorf("expected type 'fill_blank', got '%s'", q.QuestionType)
	}
	if q.Answer != "A" {
		t.Errorf("expected answer 'A', got '%s'", q.Answer)
	}
	if len(q.Options) != 1 || q.Options[0].Value != "Paris" {
		t.Errorf("expected option 'Paris', got '%s'", q.Options[0].Value)
	}
}

func TestParseGIFT_MultipleQuestions(t *testing.T) {
	input := `Q1? {
	=Correct
	~Wrong
	}

	The answer is true.{TRUE}`

	result, questions := ParseGIFT(input)
	if result.Imported != 2 {
		t.Fatalf("expected 2 imported, got %d", result.Imported)
	}
	if len(questions) != 2 {
		t.Fatalf("expected 2 questions, got %d", len(questions))
	}
	if questions[0].QuestionType != "mcq" {
		t.Errorf("first question should be 'mcq', got '%s'", questions[0].QuestionType)
	}
	if questions[1].QuestionType != "true_false" {
		t.Errorf("second question should be 'true_false', got '%s'", questions[1].QuestionType)
	}
}

func TestParseGIFT_EmptyContent(t *testing.T) {
	result, questions := ParseGIFT("")
	if result.Imported != 0 {
		t.Errorf("expected 0 imported for empty content, got %d", result.Imported)
	}
	if len(questions) != 0 {
		t.Errorf("expected 0 questions for empty content, got %d", len(questions))
	}
}

func TestParseGIFT_CommentsIgnored(t *testing.T) {
	input := `// This is a comment
	Q1?{=A}{~B}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	if len(questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(questions))
	}
}

func TestParseMCQ_MultipleCorrectAnswers(t *testing.T) {
	input := `Which are programming languages? {
	=Python
	=Go
	~English
	~French
	}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.Answer != "A" {
		t.Errorf("expected answer 'A' (first correct), got '%s'", q.Answer)
	}
	if !q.Options[0].Correct || !q.Options[1].Correct {
		t.Errorf("both Python and Go should be marked correct")
	}
}

func TestParseGIFT_InvalidFormat(t *testing.T) {
	input := `This has no braces or recognizable format`
	result, _ := ParseGIFT(input)
	if result.Imported != 0 {
		t.Errorf("expected 0 imported for invalid format, got %d", result.Imported)
	}
}

func TestParseGIFT_InlineMath(t *testing.T) {
	input := `What is \(E=mc^2\)? {
	=\(m=E/c^2\)
	~\(m=c^2/E\)
	~\(m=E\cdot c^2\)
	}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionType != "mcq" {
		t.Fatalf("expected mcq, got %s", q.QuestionType)
	}
	if q.QuestionText != "What is \\(E=mc^2\\)?" {
		t.Errorf("question text math not preserved: %q", q.QuestionText)
	}
	if q.Answer != "A" {
		t.Errorf("expected answer A, got %s", q.Answer)
	}
	if q.Options[0].Value != `\(m=E/c^2\)` {
		t.Errorf("option 0 math not preserved: %q", q.Options[0].Value)
	}
	if len(q.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(q.Options))
	}
}

func TestParseGIFT_FractionInStemAndOptions(t *testing.T) {
	input := `Simplify \[\frac{1}{2} + \frac{1}{3}\]. {
	=\(\frac{5}{6}\)
	~\(\frac{2}{5}\)
	}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionText != `Simplify \[\frac{1}{2} + \frac{1}{3}\].` {
		t.Errorf("display math not preserved: %q", q.QuestionText)
	}
	if q.Options[0].Value != `\(\frac{5}{6}\)` {
		t.Errorf("option math not preserved: %q", q.Options[0].Value)
	}
	if q.Answer != "A" {
		t.Errorf("expected answer A, got %s", q.Answer)
	}
}

func TestParseGIFT_TildeEqualsInsideMathNotMarkers(t *testing.T) {
	input := `Convert 1 g to milligrams. {
	=\(\mathrm{1000~mg}\)
	~\(\mathrm{10~mg}\)
	}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionType != "mcq" {
		t.Fatalf("expected mcq, got %s", q.QuestionType)
	}
	if len(q.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(q.Options))
	}
	if q.Options[0].Value != `\(\mathrm{1000~mg}\)` {
		t.Errorf("tilde inside math treated as marker: %q", q.Options[0].Value)
	}
	if q.Options[1].Value != `\(\mathrm{10~mg}\)` {
		t.Errorf("option 1 mangled: %q", q.Options[1].Value)
	}
	if q.Answer != "A" {
		t.Errorf("expected answer A, got %s", q.Answer)
	}
}

func TestParseGIFT_DisplayMathMultiline(t *testing.T) {
	input := `Evaluate \[
\sum_{i=1}^{n} i = \frac{n(n+1)}{2}
\]. {
	=TRUE
	~FALSE
	}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionType != "mcq" {
		t.Fatalf("expected mcq, got %s", q.QuestionType)
	}
	want := "Evaluate \\[\n\\sum_{i=1}^{n} i = \\frac{n(n+1)}{2}\n\\]."
	if q.QuestionText != want {
		t.Errorf("multiline display math not preserved:\n got: %q\nwant: %q", q.QuestionText, want)
	}
}

func TestParseGIFT_EscapedBraces(t *testing.T) {
	input := `What is the set \{1, 2, 3\}? {=A ~B}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.QuestionType != "mcq" {
		t.Fatalf("expected mcq, got %s", q.QuestionType)
	}
	if q.QuestionText != "What is the set {1, 2, 3}?" {
		t.Errorf("escaped braces not unescaped: %q", q.QuestionText)
	}
}

func TestParseGIFT_Metadata(t *testing.T) {
	input := `::Algebra Q1::
Simplify.
[difficulty:hard]
[chapter:Linear Equations]
[tags:algebra,formula]
{=\(x^2\) #2 ~\(x^3\)}`
	result, questions := ParseGIFT(input)
	if result.Imported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.Imported)
	}
	q := questions[0]
	if q.Marks != 2 {
		t.Errorf("expected marks 2, got %v", q.Marks)
	}
	if q.Difficulty != "hard" {
		t.Errorf("expected difficulty hard, got %s", q.Difficulty)
	}
	if len(q.Chapters) != 1 || q.Chapters[0] != "Linear Equations" {
		t.Errorf("expected chapter, got %v", q.Chapters)
	}
	if len(q.Tags) != 2 || q.Tags[0] != "algebra" || q.Tags[1] != "formula" {
		t.Errorf("expected tags, got %v", q.Tags)
	}
	if q.QuestionText != "Simplify." {
		t.Errorf("metadata leaked into question text: %q", q.QuestionText)
	}
	if q.Answer != "A" {
		t.Errorf("expected answer A, got %s", q.Answer)
	}
}
