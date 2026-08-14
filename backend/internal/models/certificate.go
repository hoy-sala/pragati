package models

import "time"

type CertificateEvent struct {
	ID             string     `json:"id" db:"id"`
	SchoolID       string     `json:"school_id" db:"school_id"`
	AcademicYearID *string    `json:"academic_year_id,omitempty" db:"academic_year_id"`
	Name           string     `json:"name" db:"name"`
	Category       string     `json:"category" db:"category"`
	HeldDate       *time.Time `json:"held_date,omitempty" db:"held_date"`
	Venue          string     `json:"venue,omitempty" db:"venue"`
	Description    string     `json:"description,omitempty" db:"description"`
	AuditInfo
}

type Certificate struct {
	ID          string     `json:"id" db:"id"`
	SchoolID    string     `json:"school_id" db:"school_id"`
	EventID     string     `json:"event_id" db:"event_id"`
	StudentID   string     `json:"student_id" db:"student_id"`
	Position    string     `json:"position" db:"position"`
	PrizeTitle  string     `json:"prize_title,omitempty" db:"prize_title"`
	IssueDate   *time.Time `json:"issue_date,omitempty" db:"issue_date"`
	AuditInfo
}

type CertificateSignatory struct {
	ID           string `json:"id" db:"id"`
	EventID      string `json:"event_id" db:"event_id"`
	Name         string `json:"name" db:"name"`
	Role         string `json:"role" db:"role"`
	Title        string `json:"title,omitempty" db:"title"`
	SignatureURL string `json:"signature_url,omitempty" db:"signature_url"`
	SortOrder    int    `json:"sort_order" db:"sort_order"`
	AuditInfo
}

type CertificateCreateInput struct {
	AcademicYearID string `json:"academic_year_id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	HeldDate       string `json:"held_date"`
	Venue          string `json:"venue"`
	Description    string `json:"description"`
}

type CertificateParticipantInput struct {
	StudentID  string `json:"student_id"`
	Position   string `json:"position"`
	PrizeTitle string `json:"prize_title"`
	IssueDate  string `json:"issue_date"`
}

type CertificateSignatoryInput struct {
	Name         string `json:"name"`
	Role         string `json:"role"`
	Title        string `json:"title"`
	SignatureURL string `json:"signature_url"`
	SortOrder    int    `json:"sort_order"`
}

type CertificateParticipant struct {
	ID         string     `json:"id"`
	StudentID  string     `json:"student_id"`
	StudentName string    `json:"student_name"`
	SATSNumber string     `json:"sats_number"`
	ClassName  string     `json:"class_name"`
	Position   string     `json:"position"`
	PrizeTitle string     `json:"prize_title,omitempty"`
	IssueDate  *time.Time `json:"issue_date,omitempty"`
}

type CertificateDetail struct {
	Certificate
	StudentName string                   `json:"student_name"`
	SATSNumber  string                   `json:"sats_number"`
	ClassName   string                   `json:"class_name"`
	Event       *CertificateEvent        `json:"event"`
	Signatories []CertificateSignatory   `json:"signatories"`
}