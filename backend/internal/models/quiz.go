package models

import "time"

type QuizAssignment struct {
	ID               string     `json:"id" db:"id"`
	SchoolID         string     `json:"school_id" db:"school_id"`
	Title            string     `json:"title" db:"title"`
	Description      string     `json:"description" db:"description"`
	TargetType       string     `json:"target_type" db:"target_type"`
	TargetID         string     `json:"target_id,omitempty" db:"target_id"`
	PassPct          int        `json:"pass_pct" db:"pass_pct"`
	MaxAttempts      int        `json:"max_attempts" db:"max_attempts"`
	DurationMin      *int       `json:"duration_min,omitempty" db:"duration_min"`
	ShuffleQuestions bool       `json:"shuffle_questions" db:"shuffle_questions"`
	ShuffleOptions   bool       `json:"shuffle_options" db:"shuffle_options"`
	ShowResult       bool       `json:"show_result" db:"show_result"`
	StartAt          *time.Time `json:"start_at,omitempty" db:"start_at"`
	EndAt            *time.Time `json:"end_at,omitempty" db:"end_at"`
	IsPublished      bool       `json:"is_published" db:"is_published"`
	CreatedBy        string     `json:"created_by" db:"created_by"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type QuizQuestion struct {
	QuizID     string `json:"quiz_id" db:"quiz_id"`
	QuestionID string `json:"question_id" db:"question_id"`
	SortOrder  int    `json:"sort_order" db:"sort_order"`
	Marks      int    `json:"marks" db:"marks"`
}

type QuizAttempt struct {
	ID          string     `json:"id" db:"id"`
	QuizID      string     `json:"quiz_id" db:"quiz_id"`
	UserID      string     `json:"user_id" db:"user_id"`
	AttemptNo   int        `json:"attempt_no" db:"attempt_no"`
	Status      string     `json:"status" db:"status"`
	Score       *float64   `json:"score,omitempty" db:"score"`
	Percentage  *float64   `json:"percentage,omitempty" db:"percentage"`
	Passed      *bool      `json:"passed,omitempty" db:"passed"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty" db:"submitted_at"`
	GradedAt    *time.Time `json:"graded_at,omitempty" db:"graded_at"`
	GradedBy    string     `json:"graded_by,omitempty" db:"graded_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type QuizResponse struct {
	ID              string   `json:"id" db:"id"`
	AttemptID       string   `json:"attempt_id" db:"attempt_id"`
	QuestionID      string   `json:"question_id" db:"question_id"`
	SelectedOptions []string `json:"selected_options,omitempty" db:"selected_options"`
	TextAnswer      string   `json:"text_answer" db:"text_answer"`
	IsCorrect       *bool    `json:"is_correct,omitempty" db:"is_correct"`
	MarksAwarded    float64  `json:"marks_awarded" db:"marks_awarded"`
	MarksTotal      float64  `json:"marks_total" db:"marks_total"`
	GradedAt        *time.Time `json:"graded_at,omitempty" db:"graded_at"`
	GradedBy        string   `json:"graded_by,omitempty" db:"graded_by"`
}

type QuizCreateInput struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	TargetType       string `json:"target_type"`
	TargetID         string `json:"target_id"`
	PassPct          int    `json:"pass_pct"`
	MaxAttempts      int    `json:"max_attempts"`
	DurationMin      *int   `json:"duration_min"`
	ShuffleQuestions bool   `json:"shuffle_questions"`
	ShuffleOptions   bool   `json:"shuffle_options"`
	ShowResult       bool   `json:"show_result"`
	StartAt          string `json:"start_at"`
	EndAt            string `json:"end_at"`
}

type QuizAnswerInput struct {
	QuestionID      string   `json:"question_id"`
	SelectedOptions []string `json:"selected_options"`
	TextAnswer      string   `json:"text_answer"`
}

type QuizSubmitInput struct {
	Responses []QuizAnswerInput `json:"responses"`
}

type QuizListItem struct {
	QuizAssignment
	QuestionCount int    `json:"question_count"`
	AttemptCount  int    `json:"attempt_count"`
	CreatedByName string `json:"created_by_name"`
}

type AvailableQuizItem struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DurationMin     *int       `json:"duration_min"`
	PassPct         int        `json:"pass_pct"`
	MaxAttempts     int        `json:"max_attempts"`
	QuestionCount   int        `json:"question_count"`
	StartAt         *time.Time `json:"start_at"`
	EndAt           *time.Time `json:"end_at"`
	AttemptsUsed    int        `json:"attempts_used"`
	LastStatus      string     `json:"last_status,omitempty"`
	LastScore       *float64   `json:"last_score,omitempty"`
	LastPassed      *bool      `json:"last_passed,omitempty"`
}

type QuizResultData struct {
	Attempt      QuizAttempt          `json:"attempt"`
	Quiz         QuizAssignment       `json:"quiz"`
	Responses    []QuizResponseDetail `json:"responses"`
	TotalMarks   float64              `json:"total_marks"`
	TotalAwarded float64              `json:"total_awarded"`
}

type QuizResponseDetail struct {
	QuizResponse
	QuestionText string          `json:"question_text"`
	QuestionType string          `json:"question_type"`
	Options      []Option        `json:"options"`
	CorrectAnswer string         `json:"correct_answer,omitempty"`
}
