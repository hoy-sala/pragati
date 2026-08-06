package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/rs/zerolog/log"
)

type ReportsHandler struct {
	db *pgxpool.Pool
}

func NewReportsHandler(db *pgxpool.Pool) *ReportsHandler {
	return &ReportsHandler{db: db}
}

type gradeBoundary struct {
	grade string
	max   float64
	min   float64
	gp    float64
}

var defaultGrades = []gradeBoundary{
	{"A1", 100, 91, 10},
	{"A2", 90, 81, 9},
	{"B1", 80, 71, 8},
	{"B2", 70, 61, 7},
	{"C1", 60, 51, 6},
	{"C2", 50, 41, 5},
	{"D", 40, 33, 4},
	{"E", 32, 0, 0},
}

func gradeForPct(pct float64) (string, float64) {
	for _, g := range defaultGrades {
		if pct >= g.min && pct <= g.max {
			return g.grade, g.gp
		}
	}
	if pct > 100 {
		return "A1", 10
	}
	return "E", 0
}

// MarkSheetAssessment describes one assessment column.
type MarkSheetAssessment struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	SubjectID    string  `json:"subject_id"`
	SubjectCode  string  `json:"subject_code"`
	SubjectName  string  `json:"subject_name"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	MaxMarks     float64 `json:"max_marks"`
	Date         string  `json:"date,omitempty"`
}

// MarkSheetStudent is one student row.
type MarkSheetStudent struct {
	StudentID  string             `json:"student_id"`
	SATSNumber string             `json:"sats_number"`
	Name       string             `json:"name"`
	RollNo     int                `json:"roll_no"`
	Marks      []MarkCell         `json:"marks"`
	Total      float64            `json:"total"`
	MaxTotal   float64            `json:"max_total"`
	Pct        float64            `json:"percentage"`
	Grade      string             `json:"grade"`
	Rank       int                `json:"rank"`
	Subjects   []SubjectAggregate `json:"subjects"`
}

// MarkCell holds a single assessment mark for a student.
type MarkCell struct {
	AssessmentID  string  `json:"assessment_id"`
	Value         float64 `json:"value"`
	IsAbsent      bool    `json:"is_absent"`
	HasMark       bool    `json:"has_mark"`
}

// SubjectAggregate rolls up a student's marks per subject.
type SubjectAggregate struct {
	SubjectID   string  `json:"subject_id"`
	SubjectCode string  `json:"subject_code"`
	SubjectName string  `json:"subject_name"`
	Total       float64 `json:"total"`
	MaxTotal    float64 `json:"max_total"`
	Pct         float64 `json:"percentage"`
	Grade       string  `json:"grade"`
}

// MarkSheetResponse is the full class mark sheet.
type MarkSheetResponse struct {
	ClassID      string              `json:"class_id"`
	ClassName    string              `json:"class_name"`
	AcademicYear string              `json:"academic_year_id"`
	Subjects     []SubjectGroup      `json:"subjects"`
	Assessments  []MarkSheetAssessment `json:"assessments"`
	Students     []MarkSheetStudent  `json:"students"`
}

// SubjectGroup collects assessments by subject for column grouping.
type SubjectGroup struct {
	SubjectID   string              `json:"subject_id"`
	SubjectCode string              `json:"subject_code"`
	SubjectName string              `json:"subject_name"`
	Assessments []MarkSheetAssessment `json:"assessments"`
}

// GET /api/v1/reports/mark-sheet?class_id=&academic_year_id=
func (h *ReportsHandler) MarkSheet(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	if classID == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "class_id is required"))
		return
	}
	academicYearID := r.URL.Query().Get("academic_year_id")

	var className string
	err := h.db.QueryRow(r.Context(),
		`SELECT name FROM classes WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		classID, claims.SchoolID).Scan(&className)
	if err != nil {
		renderJSON(w, http.StatusNotFound, apiErr("NOT_FOUND", "class not found"))
		return
	}

	assyQuery := `SELECT a.id, COALESCE(a.name,''), a.subject_id, s.code, s.name,
		a.category_id, c.name, a.max_marks::double precision, COALESCE(a.date::text,'')
		FROM assessments a
		JOIN subjects s ON s.id = a.subject_id AND s.deleted_at IS NULL
		JOIN assessment_categories c ON c.id = a.category_id
		WHERE a.class_id = $1 AND a.deleted_at IS NULL AND a.is_published = true`
	assyArgs := []interface{}{classID}
	n := 2
	if academicYearID != "" {
		assyQuery += fmt.Sprintf(" AND a.academic_year_id = $%d", n)
		assyArgs = append(assyArgs, academicYearID)
		n++
	}
	assyQuery += " ORDER BY s.name, c.sort_order NULLS LAST, a.created_at"

	rows, err := h.db.Query(r.Context(), assyQuery, assyArgs...)
	if err != nil {
		log.Error().Err(err).Msg("mark sheet: query assessments failed")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch assessments"))
		return
	}
	defer rows.Close()

	assessments := []MarkSheetAssessment{}
	for rows.Next() {
		var a MarkSheetAssessment
		if err := rows.Scan(&a.ID, &a.Name, &a.SubjectID, &a.SubjectCode, &a.SubjectName,
			&a.CategoryID, &a.CategoryName, &a.MaxMarks, &a.Date); err != nil {
			continue
		}
		assessments = append(assessments, a)
	}

	studQuery := `SELECT id, sats_number, first_name, COALESCE(last_name,''), roll_no
		FROM students WHERE class_id = $1 AND deleted_at IS NULL AND is_active = true
		ORDER BY roll_no NULLS LAST, first_name`
	studRows, err := h.db.Query(r.Context(), studQuery, classID)
	if err != nil {
		log.Error().Err(err).Msg("mark sheet: query students failed")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch students"))
		return
	}
	defer studRows.Close()

	type rawStudent struct {
		id, sats, first, last string
		roll                  int
	}
	rawStudents := []rawStudent{}
	for studRows.Next() {
		var s rawStudent
		if err := studRows.Scan(&s.id, &s.sats, &s.first, &s.last, &s.roll); err != nil {
			continue
		}
		rawStudents = append(rawStudents, s)
	}

	markQuery := `SELECT m.student_id, m.assessment_id, m.marks_obtained::double precision, m.is_absent
		FROM marks m
		JOIN assessments a ON a.id = m.assessment_id
		WHERE a.class_id = $1 AND a.deleted_at IS NULL AND a.is_published = true`
	markRows, err := h.db.Query(r.Context(), markQuery, classID)
	if err != nil {
		log.Error().Err(err).Msg("mark sheet: query marks failed")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch marks"))
		return
	}
	defer markRows.Close()

	type markKey struct{ sid, aid string }
	marks := map[markKey]MarkCell{}
	for markRows.Next() {
		var sid, aid string
		var val float64
		var absent bool
		if err := markRows.Scan(&sid, &aid, &val, &absent); err != nil {
			continue
		}
		marks[markKey{sid, aid}] = MarkCell{AssessmentID: aid, Value: val, IsAbsent: absent, HasMark: true}
	}

	students := make([]MarkSheetStudent, 0, len(rawStudents))
	for _, rs := range rawStudents {
		cells := make([]MarkCell, len(assessments))
		var total, maxTotal float64
		subAgg := map[string]*SubjectAggregate{}
		for i, a := range assessments {
			cell := MarkCell{AssessmentID: a.ID}
			if m, ok := marks[markKey{rs.id, a.ID}]; ok {
				cell = m
				if !m.IsAbsent {
					total += m.Value
				}
			}
			maxTotal += a.MaxMarks
			cells[i] = cell

			if subAgg[a.SubjectID] == nil {
				subAgg[a.SubjectID] = &SubjectAggregate{
					SubjectID: a.SubjectID, SubjectCode: a.SubjectCode, SubjectName: a.SubjectName,
				}
			}
			sa := subAgg[a.SubjectID]
			sa.MaxTotal += a.MaxMarks
			if !cell.IsAbsent && cell.HasMark {
				sa.Total += cell.Value
			}
		}
		subList := make([]SubjectAggregate, 0, len(subAgg))
		for _, sa := range subAgg {
			if sa.MaxTotal > 0 {
				sa.Pct = (sa.Total / sa.MaxTotal) * 100
			}
			sa.Grade, _ = gradeForPct(sa.Pct)
			subList = append(subList, *sa)
		}
		sort.Slice(subList, func(i, j int) bool { return subList[i].SubjectName < subList[j].SubjectName })

		var pct float64
		if maxTotal > 0 {
			pct = (total / maxTotal) * 100
		}
		grade, _ := gradeForPct(pct)
		students = append(students, MarkSheetStudent{
			StudentID: rs.id, SATSNumber: rs.sats,
			Name: rs.first, RollNo: rs.roll, Marks: cells,
			Total: total, MaxTotal: maxTotal, Pct: pct, Grade: grade, Subjects: subList,
		})
		if rs.last != "" {
			students[len(students)-1].Name += " " + rs.last
		}
	}

	sort.Slice(students, func(i, j int) bool {
		if students[i].Total != students[j].Total {
			return students[i].Total > students[j].Total
		}
		return students[i].Name < students[j].Name
	})
	for i := range students {
		students[i].Rank = i + 1
	}

	subjGroups := map[string]*SubjectGroup{}
	subjOrder := []string{}
	for _, a := range assessments {
		if subjGroups[a.SubjectID] == nil {
			subjGroups[a.SubjectID] = &SubjectGroup{
				SubjectID: a.SubjectID, SubjectCode: a.SubjectCode, SubjectName: a.SubjectName,
			}
			subjOrder = append(subjOrder, a.SubjectID)
		}
		subjGroups[a.SubjectID].Assessments = append(subjGroups[a.SubjectID].Assessments, a)
	}
	groups := make([]SubjectGroup, 0, len(subjOrder))
	for _, sid := range subjOrder {
		groups = append(groups, *subjGroups[sid])
	}

	renderJSON(w, http.StatusOK, apiOK(MarkSheetResponse{
		ClassID: classID, ClassName: className, AcademicYear: academicYearID,
		Subjects: groups, Assessments: assessments, Students: students,
	}))
}

// StudentReportResponse is an individual student's report card data.
type StudentReportResponse struct {
	Student      ReportStudent        `json:"student"`
	AcademicYear string               `json:"academic_year_id"`
	Term         string               `json:"term,omitempty"`
	Subjects     []ReportSubject      `json:"subjects"`
	GrandTotal   float64              `json:"grand_total"`
	GrandMax     float64              `json:"grand_max"`
	Pct          float64              `json:"percentage"`
	Grade        string               `json:"grade"`
	Attendance   *AttendanceSummary   `json:"attendance,omitempty"`
	Remarks      string               `json:"remarks,omitempty"`
}

type ReportStudent struct {
	Name       string `json:"name"`
	Class      string `json:"class"`
	Section    string `json:"section,omitempty"`
	RollNo     int    `json:"roll_no"`
	SATSNumber string `json:"sats_number"`
	Gender     string `json:"gender,omitempty"`
	DOB        string `json:"date_of_birth,omitempty"`
}

type ReportSubject struct {
	SubjectID   string           `json:"subject_id"`
	SubjectCode string           `json:"subject_code"`
	SubjectName string           `json:"subject_name"`
	Assessments []ReportAssessment `json:"assessments"`
	Total       float64          `json:"total"`
	MaxTotal    float64          `json:"max_max"`
	Pct         float64          `json:"percentage"`
	Grade       string           `json:"grade"`
}

type ReportAssessment struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Category string `json:"category"`
	Max     float64 `json:"max"`
	Value   float64 `json:"value"`
	Absent  bool    `json:"absent"`
	HasMark bool    `json:"has_mark"`
}

type AttendanceSummary struct {
	Present int     `json:"present"`
	Total   int     `json:"total"`
	Pct     float64 `json:"percentage"`
}

// GET /api/v1/reports/student?student_id=&academic_year_id=&term=
func (h *ReportsHandler) StudentReport(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "student_id is required"))
		return
	}
	academicYearID := r.URL.Query().Get("academic_year_id")
	term := r.URL.Query().Get("term")

	var st ReportStudent
	var classID, dob, gender string
	err := h.db.QueryRow(r.Context(),
		`SELECT s.first_name || ' ' || COALESCE(s.last_name,''),
			c.name, COALESCE(s.section_id::text,''), s.roll_no, s.sats_number,
			COALESCE(s.gender,''), COALESCE(s.date_of_birth::text,''), s.class_id
		 FROM students s JOIN classes c ON c.id = s.class_id
		 WHERE s.id = $1 AND s.school_id = $2 AND s.deleted_at IS NULL`,
		studentID, claims.SchoolID,
	).Scan(&st.Name, &st.Class, &st.Section, &st.RollNo, &st.SATSNumber, &gender, &dob, &classID)
	if err != nil {
		renderJSON(w, http.StatusNotFound, apiErr("NOT_FOUND", "student not found"))
		return
	}
	st.Gender = gender
	st.DOB = dob

	assyQuery := `SELECT a.id, COALESCE(a.name,''), a.subject_id, s.code, s.name,
		c.name, a.max_marks::double precision
		FROM assessments a
		JOIN subjects s ON s.id = a.subject_id AND s.deleted_at IS NULL
		JOIN assessment_categories c ON c.id = a.category_id
		WHERE a.class_id = $1 AND a.deleted_at IS NULL AND a.is_published = true`
	assyArgs := []interface{}{classID}
	n := 2
	if academicYearID != "" {
		assyQuery += fmt.Sprintf(" AND a.academic_year_id = $%d", n)
		assyArgs = append(assyArgs, academicYearID)
		n++
	}
	assyQuery += " ORDER BY s.name, c.sort_order NULLS LAST, a.created_at"

	rows, err := h.db.Query(r.Context(), assyQuery, assyArgs...)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch assessments"))
		return
	}
	defer rows.Close()

	type rawAssy struct {
		id, name, sid, scode, sname, cat string
		max                              float64
	}
	allAssy := []rawAssy{}
	for rows.Next() {
		var a rawAssy
		if err := rows.Scan(&a.id, &a.name, &a.sid, &a.scode, &a.sname, &a.cat, &a.max); err != nil {
			continue
		}
		allAssy = append(allAssy, a)
	}

	marksQuery := `SELECT m.assessment_id, m.marks_obtained::double precision, m.is_absent
		FROM marks m JOIN assessments a ON a.id = m.assessment_id
		WHERE m.student_id = $1 AND a.class_id = $2 AND a.deleted_at IS NULL AND a.is_published = true`
	mkRows, err := h.db.Query(r.Context(), marksQuery, studentID, classID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch marks"))
		return
	}
	defer mkRows.Close()

	markMap := map[string]MarkCell{}
	for mkRows.Next() {
		var aid string
		var val float64
		var absent bool
		if err := mkRows.Scan(&aid, &val, &absent); err != nil {
			continue
		}
		markMap[aid] = MarkCell{AssessmentID: aid, Value: val, IsAbsent: absent, HasMark: true}
	}

	subjMap := map[string]*ReportSubject{}
	subjOrder := []string{}
	var grandTotal, grandMax float64
	for _, a := range allAssy {
		if subjMap[a.sid] == nil {
			subjMap[a.sid] = &ReportSubject{SubjectID: a.sid, SubjectCode: a.scode, SubjectName: a.sname}
			subjOrder = append(subjOrder, a.sid)
		}
		rs := subjMap[a.sid]
		cell := ReportAssessment{ID: a.id, Name: a.name, Category: a.cat, Max: a.max}
		if m, ok := markMap[a.id]; ok {
			cell.Value = m.Value
			cell.Absent = m.IsAbsent
			cell.HasMark = true
			if !m.IsAbsent {
				rs.Total += m.Value
				grandTotal += m.Value
			}
		}
		rs.MaxTotal += a.max
		grandMax += a.max
		rs.Assessments = append(rs.Assessments, cell)
	}

	subjects := make([]ReportSubject, 0, len(subjOrder))
	for _, sid := range subjOrder {
		rs := subjMap[sid]
		if rs.MaxTotal > 0 {
			rs.Pct = (rs.Total / rs.MaxTotal) * 100
		}
		rs.Grade, _ = gradeForPct(rs.Pct)
		subjects = append(subjects, *rs)
	}

	var pct float64
	if grandMax > 0 {
		pct = (grandTotal / grandMax) * 100
	}
	grade, _ := gradeForPct(pct)

	resp := StudentReportResponse{
		Student: st, AcademicYear: academicYearID, Term: term,
		Subjects: subjects, GrandTotal: grandTotal, GrandMax: grandMax, Pct: pct, Grade: grade,
	}

	if term != "" && academicYearID != "" {
		var attJSON []byte
		err := h.db.QueryRow(r.Context(),
			`SELECT attendance_summary FROM hpc_entries
			 WHERE student_id = $1 AND term = $2 AND academic_year_id = $3 AND deleted_at IS NULL`,
			studentID, term, academicYearID).Scan(&attJSON)
		if err == nil && len(attJSON) > 0 {
			var att AttendanceSummary
			if json.Unmarshal(attJSON, &att) == nil {
				if att.Total > 0 {
					att.Pct = (float64(att.Present) / float64(att.Total)) * 100
				}
				resp.Attendance = &att
			}
		}
	}

	renderJSON(w, http.StatusOK, apiOK(resp))
}

func apiOK(data interface{}) map[string]interface{} {
	return map[string]interface{}{"data": data}
}

func apiErr(code, msg string) map[string]interface{} {
	return map[string]interface{}{"error": map[string]string{"code": code, "message": msg}}
}
