package models

import "time"

type Student struct {
	ID             string     `json:"id" db:"id"`
	SchoolID       string     `json:"school_id" db:"school_id"`
	UserID         *string    `json:"user_id,omitempty" db:"user_id"`
	SATSNumber     string     `json:"sats_number" db:"sats_number"`
	AdmissionNo    string     `json:"admission_no,omitempty" db:"admission_no"`
	RollNo         int        `json:"roll_no,omitempty" db:"roll_no"`
	FirstName      string     `json:"first_name" db:"first_name"`
	LastName       string     `json:"last_name,omitempty" db:"last_name"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender         string     `json:"gender,omitempty" db:"gender"`
	PhotoURL       string     `json:"photo_url,omitempty" db:"photo_url"`
	BloodGroup     string     `json:"blood_group,omitempty" db:"blood_group"`
	Address        string     `json:"address,omitempty" db:"address"`
	Phone          string     `json:"phone,omitempty" db:"phone"`
	Email          string     `json:"email,omitempty" db:"email"`
	ClassID        string     `json:"class_id" db:"class_id"`
	SectionID      string     `json:"section_id,omitempty" db:"section_id"`
	HouseID        *string    `json:"house_id,omitempty" db:"house_id"`
	AcademicYearID string     `json:"academic_year_id" db:"academic_year_id"`
	ParentName     string     `json:"parent_name,omitempty" db:"parent_name"`
	ParentPhone    string     `json:"parent_phone,omitempty" db:"parent_phone"`
	ParentEmail    string     `json:"parent_email,omitempty" db:"parent_email"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	AuditInfo
}

type CreateStudentRequest struct {
	SATSNumber     string  `json:"sats_number" validate:"required,len=9"`
	AdmissionNo    string  `json:"admission_no"`
	RollNo         int     `json:"roll_no"`
	FirstName      string  `json:"first_name" validate:"required"`
	LastName       string  `json:"last_name"`
	DateOfBirth    string  `json:"date_of_birth"`
	Gender         string  `json:"gender"`
	BloodGroup     string  `json:"blood_group"`
	Address        string  `json:"address"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	ClassID        string  `json:"class_id" validate:"required"`
	SectionID      string  `json:"section_id"`
	HouseID        string  `json:"house_id"`
	AcademicYearID string  `json:"academic_year_id" validate:"required"`
	ParentName     string  `json:"parent_name"`
	ParentPhone    string  `json:"parent_phone"`
	ParentEmail    string  `json:"parent_email"`
}

type UpdateStudentRequest struct {
	AdmissionNo *string `json:"admission_no"`
	RollNo      *int    `json:"roll_no"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	PhotoURL    *string `json:"photo_url"`
	BloodGroup  *string `json:"blood_group"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	ClassID     *string `json:"class_id"`
	SectionID   *string `json:"section_id"`
	HouseID     *string `json:"house_id"`
	ParentName  *string `json:"parent_name"`
	ParentPhone *string `json:"parent_phone"`
	ParentEmail *string `json:"parent_email"`
	IsActive    *bool   `json:"is_active"`
}

type ImportStudentRow struct {
	SATSNumber     string `csv:"sats_number" excel:"sats_number"`
	AdmissionNo    string `csv:"admission_no" excel:"admission_no"`
	RollNo         int    `csv:"roll_no" excel:"roll_no"`
	FirstName      string `csv:"first_name" excel:"first_name"`
	LastName       string `csv:"last_name" excel:"last_name"`
	DateOfBirth    string `csv:"date_of_birth" excel:"date_of_birth"`
	Gender         string `csv:"gender" excel:"gender"`
	Phone          string `csv:"phone" excel:"phone"`
	Email          string `csv:"email" excel:"email"`
	ClassCode      string `csv:"class" excel:"class"`
	SectionName    string `csv:"section" excel:"section"`
	AcademicYear   string `csv:"academic_year" excel:"academic_year"`
	ParentName     string `csv:"parent_name" excel:"parent_name"`
	ParentPhone    string `csv:"parent_phone" excel:"parent_phone"`
	ParentEmail    string `csv:"parent_email" excel:"parent_email"`
}

type ImportResult struct {
	Imported int              `json:"imported"`
	Skipped  int              `json:"skipped"`
	Errors   []ImportRowError `json:"errors,omitempty"`
}

type ImportRowError struct {
	Row     int    `json:"row"`
	SATS    string `json:"sats_number"`
	Field   string `json:"field"`
	Message string `json:"message"`
}
