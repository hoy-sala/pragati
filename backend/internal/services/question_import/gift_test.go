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
