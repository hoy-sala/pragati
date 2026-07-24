package models

import "time"

type Question struct {
	ID           string   `json:"id" db:"id"`
	SchoolID     string   `json:"school_id" db:"school_id"`
	SubjectID    string   `json:"subject_id" db:"subject_id"`
	TeacherID    string   `json:"teacher_id" db:"teacher_id"`
	QuestionType string   `json:"question_type" db:"question_type"`
	QuestionText string   `json:"question_text" db:"question_text"`
	QuestionImage string  `json:"question_image,omitempty" db:"question_image"`
	Options      []Option `json:"options" db:"options"`
	Answer       string   `json:"answer" db:"answer"`
	Marks        float64  `json:"marks" db:"marks"`
	Difficulty   string   `json:"difficulty" db:"difficulty"`
	Chapters     []string `json:"chapters" db:"chapters"`
	Tags         []string `json:"tags" db:"tags"`
	Explanation  string   `json:"explanation,omitempty" db:"explanation"`
	IsActive     bool     `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Option struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Correct bool   `json:"correct"`
}

type CreateQuestionRequest struct {
	SubjectID    string   `json:"subject_id" validate:"required"`
	QuestionType string   `json:"question_type" validate:"required,oneof=mcq true_false fill_blank short_answer"`
	QuestionText string   `json:"question_text" validate:"required"`
	QuestionImage string  `json:"question_image"`
	Options      []Option `json:"options"`
	Answer       string   `json:"answer" validate:"required"`
	Marks        float64  `json:"marks"`
	Difficulty   string   `json:"difficulty" validate:"oneof=easy medium hard"`
	Chapters     []string `json:"chapters"`
	Tags         []string `json:"tags"`
	Explanation  string   `json:"explanation"`
}
