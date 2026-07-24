package models

import "time"

type HPCConfiguration struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	Stage            string    `json:"stage" db:"stage"`
	ClassID          string    `json:"class_id,omitempty" db:"class_id"`
	AcademicYearID   string    `json:"academic_year_id,omitempty" db:"academic_year_id"`
	GradingScheme    []GradeDefinition `json:"grading_scheme" db:"grading_scheme"`
	ProficiencyScale []ProficiencyLevel `json:"proficiency_scale" db:"proficiency_scale"`
	CoScholasticAreas map[string]interface{} `json:"co_scholastic_areas" db:"co_scholastic_areas"`
	HealthParams     map[string]interface{} `json:"health_params" db:"health_params"`
	Terms            []string `json:"terms" db:"terms"`
	IsActive         bool     `json:"is_active" db:"is_active"`
	AuditInfo
}

type GradeDefinition struct {
	Grade      string `json:"grade"`
	MinPct     float64 `json:"min_pct"`
	MaxPct     float64 `json:"max_pct"`
	Descriptor string `json:"descriptor"`
}

type ProficiencyLevel struct {
	Level      int    `json:"level"`
	Label      string `json:"label"`
	Descriptor string `json:"descriptor"`
}

type LearningOutcome struct {
	ID            string    `json:"id" db:"id"`
	SchoolID      string    `json:"school_id" db:"school_id"`
	SubjectID     string    `json:"subject_id" db:"subject_id"`
	ClassID       string    `json:"class_id" db:"class_id"`
	Code          string    `json:"code" db:"code"`
	Description   string    `json:"description" db:"description"`
	Domain        string    `json:"domain" db:"domain"`
	ExpectedLevel int       `json:"expected_level" db:"expected_level"`
	SortOrder     int       `json:"sort_order" db:"sort_order"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type HPCEntry struct {
	ID               string                 `json:"id" db:"id"`
	StudentID        string                 `json:"student_id" db:"student_id"`
	AcademicYearID   string                 `json:"academic_year_id" db:"academic_year_id"`
	Term             string                 `json:"term" db:"term"`
	Status           string                 `json:"status" db:"status"`
	Scholastic       []ScholasticEntry      `json:"scholastic" db:"scholastic"`
	CoScholastic     map[string]interface{} `json:"co_scholastic" db:"co_scholastic"`
	HealthPE         map[string]interface{} `json:"health_pe" db:"health_pe"`
	WorkEducation    map[string]interface{} `json:"work_education" db:"work_education"`
	SelfAssessment   map[string]interface{} `json:"self_assessment" db:"self_assessment"`
	PeerAssessment   map[string]interface{} `json:"peer_assessment" db:"peer_assessment"`
	ParentFeedback   map[string]interface{} `json:"parent_feedback" db:"parent_feedback"`
	TeacherRemarks   string                 `json:"teacher_remarks" db:"teacher_remarks"`
	AttendanceSummary map[string]interface{} `json:"attendance_summary" db:"attendance_summary"`
	GeneratedPdfURL  string                 `json:"generated_pdf_url,omitempty" db:"generated_pdf_url"`
	Version          int                    `json:"version" db:"version"`
	LockedAt         *time.Time             `json:"locked_at,omitempty" db:"locked_at"`
	LockedBy         string                 `json:"locked_by,omitempty" db:"locked_by"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

type ScholasticEntry struct {
	SubjectID    string              `json:"subject_id"`
	SubjectName  string              `json:"subject_name,omitempty"`
	MaxMarks     float64             `json:"max_marks"`
	MarksScored  float64             `json:"marks_scored"`
	Grade        string              `json:"grade"`
	GradePoint   float64             `json:"grade_point,omitempty"`
	LOAttainment []LOAttainmentEntry `json:"lo_attainment,omitempty"`
}

type LOAttainmentEntry struct {
	LOID            string `json:"lo_id"`
	Code            string `json:"code"`
	Description     string `json:"description"`
	ProficiencyLevel int   `json:"proficiency_level"`
}

type LearningOutcomeAssessment struct {
	ID               string    `json:"id" db:"id"`
	StudentID        string    `json:"student_id" db:"student_id"`
	LearningOutcomeID string   `json:"learning_outcome_id" db:"learning_outcome_id"`
	SubjectID        string    `json:"subject_id" db:"subject_id"`
	Term             string    `json:"term" db:"term"`
	ProficiencyLevel int       `json:"proficiency_level" db:"proficiency_level"`
	AssessedBy       string    `json:"assessed_by" db:"assessed_by"`
	AssessmentDate   string    `json:"assessment_date" db:"assessment_date"`
	Remarks          string    `json:"remarks" db:"remarks"`
	AuditInfo
}

type HPCConfigInput struct {
	Stage            string                   `json:"stage"`
	ClassID          string                   `json:"class_id"`
	AcademicYearID   string                   `json:"academic_year_id"`
	GradingScheme    []GradeDefinition        `json:"grading_scheme"`
	ProficiencyScale []ProficiencyLevel       `json:"proficiency_scale"`
	CoScholasticAreas map[string]interface{}  `json:"co_scholastic_areas"`
	HealthParams     map[string]interface{}   `json:"health_params"`
	Terms            []string                 `json:"terms"`
}

type HPCEntryInput struct {
	StudentID        string                 `json:"student_id"`
	AcademicYearID   string                 `json:"academic_year_id"`
	Term             string                 `json:"term"`
	Scholastic       []ScholasticEntry      `json:"scholastic"`
	CoScholastic     map[string]interface{} `json:"co_scholastic"`
	HealthPE         map[string]interface{} `json:"health_pe"`
	WorkEducation    map[string]interface{} `json:"work_education"`
	SelfAssessment   map[string]interface{} `json:"self_assessment"`
	PeerAssessment   map[string]interface{} `json:"peer_assessment"`
	ParentFeedback   map[string]interface{} `json:"parent_feedback"`
	TeacherRemarks   string                 `json:"teacher_remarks"`
	AttendanceSummary map[string]interface{} `json:"attendance_summary"`
}

type HPCGridRow struct {
	StudentID        string `json:"student_id"`
	SATSNumber       string `json:"sats_number"`
	Name             string `json:"name"`
	RollNo           int    `json:"roll_no"`
	EntryID          string `json:"entry_id,omitempty"`
	Status           string `json:"status"`
	HasPDF           bool   `json:"has_pdf"`
}

type LOBatchAssessmentInput struct {
	StudentID    string `json:"student_id"`
	SubjectID    string `json:"subject_id"`
	Term         string `json:"term"`
	Assessments  []struct {
		LearningOutcomeID string `json:"learning_outcome_id"`
		ProficiencyLevel  int    `json:"proficiency_level"`
	} `json:"assessments"`
}

type HPCReportSummary struct {
	TotalStudents   int `json:"total_students"`
	PublishedCount  int `json:"published_count"`
	DraftCount      int `json:"draft_count"`
	ByGrade         map[string]int `json:"by_grade,omitempty"`
	ByProficiency   map[int]int  `json:"by_proficiency,omitempty"`
}
