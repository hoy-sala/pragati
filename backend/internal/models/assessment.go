package models

import "time"

type Assessment struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	CategoryID     string    `json:"category_id" db:"category_id"`
	SubjectID      string    `json:"subject_id" db:"subject_id"`
	TeacherID      string    `json:"teacher_id" db:"teacher_id"`
	ClassID        string    `json:"class_id" db:"class_id"`
	SectionID      string    `json:"section_id,omitempty" db:"section_id"`
	Name           string    `json:"name" db:"name"`
	MaxMarks       float64   `json:"max_marks" db:"max_marks"`
	Weightage      int       `json:"weightage" db:"weightage"`
	Date           string    `json:"date" db:"date"`
	AcademicYearID string    `json:"academic_year_id" db:"academic_year_id"`
	IsPublished    bool      `json:"is_published" db:"is_published"`
	IsLocked       bool      `json:"is_locked" db:"is_locked"`
	Version        int       `json:"version" db:"version"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type AssessmentCategory struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	Weightage int       `json:"weightage" db:"weightage"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MarkInput struct {
	StudentID     string  `json:"student_id"`
	MarksObtained float64 `json:"marks_obtained"`
	IsAbsent      bool    `json:"is_absent"`
	Remarks       string  `json:"remarks"`
}
