package models

import "time"

type MentorAssignment struct {
	ID             string    `json:"id" db:"id"`
	MentorID       string    `json:"mentor_id" db:"mentor_id"`
	StudentID      string    `json:"student_id" db:"student_id"`
	AcademicYearID string    `json:"academic_year_id" db:"academic_year_id"`
	AssignedAt     time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy     string    `json:"assigned_by,omitempty" db:"assigned_by"`
}

type MentorAttendance struct {
	ID               string     `json:"id" db:"id"`
	MentorID         string     `json:"mentor_id" db:"mentor_id"`
	StudentID        string     `json:"student_id" db:"student_id"`
	Date             string     `json:"date" db:"date"`
	Status           string     `json:"status" db:"status"`
	ParentContacted  bool       `json:"parent_contacted" db:"parent_contacted"`
	ParentContactTime *time.Time `json:"parent_contact_time,omitempty" db:"parent_contact_time"`
	Remarks          string     `json:"remarks,omitempty" db:"remarks"`
}

type MentorLog struct {
	ID                 string    `json:"id" db:"id"`
	MentorID           string    `json:"mentor_id" db:"mentor_id"`
	StudentID          string    `json:"student_id" db:"student_id"`
	LogDate            string    `json:"log_date" db:"log_date"`
	Category           string    `json:"category" db:"category"`
	Severity           string    `json:"severity" db:"severity"`
	Description        string    `json:"description" db:"description"`
	ActionTaken        string    `json:"action_taken,omitempty" db:"action_taken"`
	ParentInformed     bool      `json:"parent_informed" db:"parent_informed"`
	ReviewedByPrincipal bool     `json:"reviewed_by_principal" db:"reviewed_by_principal"`
	PrincipalNotes     string    `json:"principal_notes,omitempty" db:"principal_notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type MentorRosterStudent struct {
	Student
	MentorID       string `json:"mentor_id"`
	FatherName     string `json:"father_name"`
	MotherName     string `json:"mother_name"`
	ParentPhone    string `json:"parent_phone"`
	Address        string `json:"address"`
}

type AlertLog struct {
	MentorLog
	StudentName string `json:"student_name"`
	MentorName  string `json:"mentor_name"`
}

type MonthlySummary struct {
	Month              string  `json:"month"`
	TotalStudents      int     `json:"total_students"`
	AvgAttendancePct   float64 `json:"avg_attendance_pct"`
	TotalLogs          int     `json:"total_logs"`
	UrgentLogs         int     `json:"urgent_logs"`
	ParentContactsMade int     `json:"parent_contacts_made"`
}
