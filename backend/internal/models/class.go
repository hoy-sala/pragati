package models

type Class struct {
	ID             string `json:"id" db:"id"`
	SchoolID       string `json:"school_id" db:"school_id"`
	AcademicYearID string `json:"academic_year_id" db:"academic_year_id"`
	Name           string `json:"name" db:"name"`
	Code           string `json:"code,omitempty" db:"code"`
	SortOrder      int    `json:"sort_order" db:"sort_order"`
	AuditInfo
}

type Section struct {
	ID             string `json:"id" db:"id"`
	ClassID        string `json:"class_id" db:"class_id"`
	Name           string `json:"name" db:"name"`
	ClassTeacherID string `json:"class_teacher_id,omitempty" db:"class_teacher_id"`
	AuditInfo
}

type Subject struct {
	ID        string `json:"id" db:"id"`
	SchoolID  string `json:"school_id" db:"school_id"`
	Name      string `json:"name" db:"name"`
	Code      string `json:"code,omitempty" db:"code"`
	IsLanguage bool  `json:"is_language" db:"is_language"`
	IsCore    bool   `json:"is_core" db:"is_core"`
	AuditInfo
}

type ClassSubject struct {
	ClassID   string `json:"class_id" db:"class_id"`
	SubjectID string `json:"subject_id" db:"subject_id"`
}

type AcademicYear struct {
	ID         string `json:"id" db:"id"`
	SchoolID   string `json:"school_id" db:"school_id"`
	Name       string `json:"name" db:"name"`
	StartDate  string `json:"start_date" db:"start_date"`
	EndDate    string `json:"end_date" db:"end_date"`
	IsCurrent  bool   `json:"is_current" db:"is_current"`
	AuditInfo
}

type House struct {
	ID       string `json:"id" db:"id"`
	SchoolID string `json:"school_id" db:"school_id"`
	Name     string `json:"name" db:"name"`
	Code     string `json:"code,omitempty" db:"code"`
	Color    string `json:"color,omitempty" db:"color"`
	LogoURL  string `json:"logo_url,omitempty" db:"logo_url"`
	AuditInfo
}

type TeacherSubject struct {
	TeacherID string `json:"teacher_id" db:"teacher_id"`
	SubjectID string `json:"subject_id" db:"subject_id"`
	ClassID   string `json:"class_id" db:"class_id"`
}
