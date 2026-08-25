package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/cce"
	"github.com/pragati/backend/internal/middleware"
	"github.com/rs/zerolog/log"
)

type ReportsHandler struct {
	db *pgxpool.Pool
}

func NewReportsHandler(db *pgxpool.Pool) *ReportsHandler {
	return &ReportsHandler{db: db}
}

// subjectType classifies a subject by its code.
// Curricular: Kannada, English, Hindi, Mathematics, Science, Social Science.
// Co-curricular: Physical Education, Computer Science, Music, Drawing.
type subjectType string

const (
	curricular    subjectType = "curricular"
	coCurricular  subjectType = "co_curricular"
)

func classifySubject(code string) subjectType {
	switch strings.ToUpper(code) {
	case "KAN", "ENG", "HIN", "MAT", "SCI", "SOC":
		return curricular
	}
	return coCurricular
}

// classNum extracts the numeric standard from a class name (e.g. "Class 9" → 9).
func classNum(className string) int {
	fields := strings.Fields(className)
	for _, f := range fields {
		if n, err := strconv.Atoi(f); err == nil {
			return n
		}
	}
	return 0
}

// gradeInfo holds a computed grade.
type gradeInfo struct {
	grade string
	label string
}

// resolveYear returns the given academic year id, or the school's current year when empty.
func (h *ReportsHandler) resolveYear(ctx context.Context, schoolID, ayID string) string {
	if ayID != "" {
		return ayID
	}
	var id string
	_ = h.db.QueryRow(ctx,
		`SELECT id FROM academic_years WHERE school_id = $1 AND is_current = true AND deleted_at IS NULL`,
		schoolID).Scan(&id)
	return id
}

// computeGrade returns the grade for a percentage given subject type and class number.
func computeGrade(pct float64, st subjectType, class int) gradeInfo {
	if pct > 100 {
		pct = 100
	}
	if st == coCurricular {
		switch {
		case pct >= 80:
			return gradeInfo{"A", "Excellent"}
		case pct >= 50:
			return gradeInfo{"B", "Good"}
		case pct >= 35:
			return gradeInfo{"C", "Progressive"}
		default:
			return gradeInfo{"F", "Needs Improvement"}
		}
	}
	// curricular
	if class >= 9 {
		switch {
		case pct >= 90:
			return gradeInfo{"A+", ""}
		case pct >= 80:
			return gradeInfo{"A", ""}
		case pct >= 70:
			return gradeInfo{"B+", ""}
		case pct >= 60:
			return gradeInfo{"B", ""}
		case pct >= 50:
			return gradeInfo{"C+", ""}
		case pct >= 34:
			return gradeInfo{"C", ""}
		default:
			return gradeInfo{"D", ""}
		}
	}
	// classes 6-8
	switch {
	case pct >= 90:
		return gradeInfo{"A+", ""}
	case pct >= 70:
		return gradeInfo{"A", ""}
	case pct >= 50:
		return gradeInfo{"B+", ""}
	case pct >= 30:
		return gradeInfo{"B", ""}
	default:
		return gradeInfo{"C", ""}
	}
}

// splitFASA separates assessment marks into FA and SA slices based on category code.
func splitFASA(assessments []ReportAssessment) (fa []float64, sa []float64) {
	for _, a := range assessments {
		if !a.HasMark || a.Absent {
			continue
		}
		code := strings.ToUpper(strings.TrimSpace(a.Category))
		if strings.HasPrefix(code, "FA") {
			fa = append(fa, a.Value)
		} else if strings.HasPrefix(code, "SA") {
			sa = append(sa, a.Value)
		}
	}
	return
}

// assessmentTerm maps an assessment name (FA1, SA1, etc.) to its term.
func assessmentTerm(assessmentName string) string {
	switch strings.ToUpper(strings.TrimSpace(assessmentName)) {
	case "FA1", "FA2", "SA1":
		return "Term 1"
	case "FA3", "FA4", "SA2":
		return "Term 2"
	}
	return ""
}

// subjectOrder defines the canonical display order for subjects.
// Curricular first (by language sequence), then co-curricular.
var subjectOrder = map[string]int{
	"KAN": 0, "ENG": 1, "HIN": 2, "MAT": 3, "SCI": 4, "SOC": 5,
	"PE": 10, "CS": 11, "MUS": 12, "DRW": 13,
}

// termOrder maps term name to a sort key.
var termOrder = map[string]int{
	"Term 1": 0, "Term 2": 1, "Term 3": 2, "Term 4": 3,
}

// sortAssessments orders assessments by term → canonical subject → category sort order.
func sortAssessments(a []MarkSheetAssessment) {
	sort.SliceStable(a, func(i, j int) bool {
		si, sj := subjectOrder[a[i].SubjectCode], subjectOrder[a[j].SubjectCode]
		if si != sj {
			return si < sj
		}
		ti, tj := termOrder[a[i].Term], termOrder[a[j].Term]
		if ti != tj {
			return ti < tj
		}
		return a[i].Name < a[j].Name
	})
}
type MarkSheetAssessment struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	SubjectID    string  `json:"subject_id"`
	SubjectCode  string  `json:"subject_code"`
	SubjectName  string  `json:"subject_name"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryCode string  `json:"category_code"`
	MaxMarks     float64 `json:"max_marks"`
	Date         string  `json:"date,omitempty"`
	Term         string  `json:"term"`
	SubjectType  string  `json:"subject_type"`
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
	Grade         string  `json:"grade,omitempty"`
}

// SubjectAggregate rolls up a student's marks per subject.
type SubjectAggregate struct {
	SubjectID    string  `json:"subject_id"`
	SubjectCode  string  `json:"subject_code"`
	SubjectName  string  `json:"subject_name"`
	SubjectType  string  `json:"subject_type"`
	Total        float64 `json:"total"`
	MaxTotal     float64 `json:"max_total"`
	Pct          float64 `json:"percentage"`
	Grade        string  `json:"grade"`
	GradeLabel   string  `json:"grade_label,omitempty"`
}

// MarkSheetResponse is the full class mark sheet.
type MarkSheetResponse struct {
	ClassID      string              `json:"class_id"`
	ClassName    string              `json:"class_name"`
	AcademicYear string              `json:"academic_year_id"`
	Term         string              `json:"term,omitempty"`
	Subjects     []SubjectGroup      `json:"subjects"`
	Terms        []TermGroup         `json:"terms"`
	Assessments  []MarkSheetAssessment `json:"assessments"`
	Students     []MarkSheetStudent  `json:"students"`
}

// TermGroup collects assessments by term then subject.
type TermGroup struct {
	Term      string         `json:"term"`
	Subjects  []SubjectGroup `json:"subjects"`
}

// SubjectGroup collects assessments by subject for column grouping.
type SubjectGroup struct {
	SubjectID    string                `json:"subject_id"`
	SubjectCode  string                `json:"subject_code"`
	SubjectName  string                `json:"subject_name"`
	SubjectType  string                `json:"subject_type"`
	Assessments  []MarkSheetAssessment `json:"assessments"`
}

// GET /api/v1/reports/mark-sheet?class_id=&academic_year_id=
func (h *ReportsHandler) MarkSheet(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	if classID == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "class_id is required"))
		return
	}
	academicYearID := h.resolveYear(r.Context(), claims.SchoolID, r.URL.Query().Get("academic_year_id"))
	term := r.URL.Query().Get("term")

	var className string
	err := h.db.QueryRow(r.Context(),
		`SELECT name FROM classes WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		classID, claims.SchoolID).Scan(&className)
	if err != nil {
		renderJSON(w, http.StatusNotFound, apiErr("NOT_FOUND", "class not found"))
		return
	}

	assyQuery := `SELECT a.id, COALESCE(a.name,''), a.subject_id, s.code, s.name,
		a.category_id, c.name, c.code, a.max_marks::double precision, COALESCE(a.date::text,'')
		FROM assessments a
		JOIN subjects s ON s.id = a.subject_id AND s.deleted_at IS NULL
		JOIN assessment_categories c ON c.id = a.category_id
		WHERE a.class_id = $1 AND a.deleted_at IS NULL`
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
		var catCode string
		if err := rows.Scan(&a.ID, &a.Name, &a.SubjectID, &a.SubjectCode, &a.SubjectName,
			&a.CategoryID, &a.CategoryName, &catCode, &a.MaxMarks, &a.Date); err != nil {
			continue
		}
		a.CategoryCode = catCode
		a.Term = assessmentTerm(a.Name)
		a.SubjectType = string(classifySubject(a.SubjectCode))
		if term != "" && a.Term != term {
			continue
		}
		assessments = append(assessments, a)
	}
	sortAssessments(assessments)

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
		WHERE a.class_id = $1 AND a.deleted_at IS NULL`
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

		cn := classNum(className)
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
				if cell.HasMark && !cell.IsAbsent && a.MaxMarks > 0 {
					g := computeGrade((cell.Value/a.MaxMarks)*100, subjectType(a.SubjectType), cn)
					cell.Grade = g.grade
				}
				cells[i] = cell

				if subAgg[a.SubjectID] == nil {
					subAgg[a.SubjectID] = &SubjectAggregate{
						SubjectID: a.SubjectID, SubjectCode: a.SubjectCode, SubjectName: a.SubjectName,
						SubjectType: a.SubjectType,
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
			g := computeGrade(sa.Pct, subjectType(sa.SubjectType), cn)
			sa.Grade = g.grade
			sa.GradeLabel = g.label
			subList = append(subList, *sa)
		}
		sort.Slice(subList, func(i, j int) bool {
			return subjectOrder[subList[i].SubjectCode] < subjectOrder[subList[j].SubjectCode]
		})

			var pct float64
			if maxTotal > 0 {
				pct = (total / maxTotal) * 100
			}
			grade := computeGrade(pct, curricular, cn)
			students = append(students, MarkSheetStudent{
			StudentID: rs.id, SATSNumber: rs.sats,
			Name: rs.first, RollNo: rs.roll, Marks: cells,
			Total: total, MaxTotal: maxTotal, Pct: pct, Grade: grade.grade, Subjects: subList,
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
				SubjectType: a.SubjectType,
			}
			subjOrder = append(subjOrder, a.SubjectID)
		}
		subjGroups[a.SubjectID].Assessments = append(subjGroups[a.SubjectID].Assessments, a)
	}
	groups := make([]SubjectGroup, 0, len(subjOrder))
	for _, sid := range subjOrder {
		groups = append(groups, *subjGroups[sid])
	}

	termGroups := buildTermGroups(assessments)

	renderJSON(w, http.StatusOK, apiOK(MarkSheetResponse{
		ClassID: classID, ClassName: className, AcademicYear: academicYearID, Term: term,
		Subjects: groups, Terms: termGroups, Assessments: assessments, Students: students,
	}))
}

// buildTermGroups nests assessments under Term > Subject.
func buildTermGroups(assessments []MarkSheetAssessment) []TermGroup {
	termMap := map[string]map[string]*SubjectGroup{}
	termOrder := []string{}
	termSubjectOrder := map[string][]string{}
	for _, a := range assessments {
		if a.Term == "" {
			continue
		}
		if termMap[a.Term] == nil {
			termMap[a.Term] = map[string]*SubjectGroup{}
			termSubjectOrder[a.Term] = []string{}
			termOrder = append(termOrder, a.Term)
		}
		sm := termMap[a.Term]
		if sm[a.SubjectID] == nil {
			sm[a.SubjectID] = &SubjectGroup{
				SubjectID: a.SubjectID, SubjectCode: a.SubjectCode, SubjectName: a.SubjectName,
			}
			termSubjectOrder[a.Term] = append(termSubjectOrder[a.Term], a.SubjectID)
		}
		sm[a.SubjectID].Assessments = append(sm[a.SubjectID].Assessments, a)
	}
	groups := make([]TermGroup, 0, len(termOrder))
	for _, t := range termOrder {
		subjects := make([]SubjectGroup, 0, len(termMap[t]))
		for _, sid := range termSubjectOrder[t] {
			subjects = append(subjects, *termMap[t][sid])
		}
		groups = append(groups, TermGroup{Term: t, Subjects: subjects})
	}
	return groups
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
	SubjectID    string           `json:"subject_id"`
	SubjectCode  string           `json:"subject_code"`
	SubjectName  string           `json:"subject_name"`
	SubjectType  string           `json:"subject_type"`
	Assessments  []ReportAssessment `json:"assessments"`
	Total        float64          `json:"total"`
	MaxTotal     float64          `json:"max_max"`
	Pct          float64          `json:"percentage"`
	Grade        string           `json:"grade"`
	GradeLabel   string           `json:"grade_label,omitempty"`
	CCE          *cce.ConversionResult `json:"cce,omitempty"`
}

type ReportAssessment struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Category string `json:"category"`
	Term    string  `json:"term"`
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
	academicYearID := h.resolveYear(r.Context(), claims.SchoolID, r.URL.Query().Get("academic_year_id"))
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
		c.name, c.code, a.max_marks::double precision
		FROM assessments a
		JOIN subjects s ON s.id = a.subject_id AND s.deleted_at IS NULL
		JOIN assessment_categories c ON c.id = a.category_id
		WHERE a.class_id = $1 AND a.deleted_at IS NULL`
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
		id, name, sid, scode, sname, cat, catCode string
		max                                       float64
	}
	allAssy := []rawAssy{}
	for rows.Next() {
		var a rawAssy
		if err := rows.Scan(&a.id, &a.name, &a.sid, &a.scode, &a.sname, &a.cat, &a.catCode, &a.max); err != nil {
			continue
		}
		allAssy = append(allAssy, a)
	}

	marksQuery := `SELECT m.assessment_id, m.marks_obtained::double precision, m.is_absent
		FROM marks m JOIN assessments a ON a.id = m.assessment_id
		WHERE m.student_id = $1 AND a.class_id = $2 AND a.deleted_at IS NULL`
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
		t := assessmentTerm(a.name)
		if term != "" && t != term {
			continue
		}
		stype := classifySubject(a.scode)
		if subjMap[a.sid] == nil {
			subjMap[a.sid] = &ReportSubject{
				SubjectID: a.sid, SubjectCode: a.scode, SubjectName: a.sname,
				SubjectType: string(stype),
			}
			subjOrder = append(subjOrder, a.sid)
		}
		rs := subjMap[a.sid]
		cell := ReportAssessment{ID: a.id, Name: a.name, Category: a.cat, Max: a.max, Term: t}
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

	classRange := cce.DetermineClassRange(st.Class)
	subjects := make([]ReportSubject, 0, len(subjOrder))

	for _, sid := range subjOrder {
		rs := subjMap[sid]
		if rs.MaxTotal > 0 {
			rs.Pct = (rs.Total / rs.MaxTotal) * 100
		}

		isFL := cce.IsFirstLanguage(rs.SubjectCode)
		subjType := cce.OtherSubject
		if isFL {
			subjType = cce.FirstLanguage
		}

		faMarks, saMarks := splitFASA(rs.Assessments)

		var cceResult cce.ConversionResult
		if classRange == cce.Class6to8 {
			cceResult = cce.ConvertClasses68(faMarks, saMarks)
		} else {
			cceResult = cce.ConvertClasses910(faMarks, saMarks, subjType)
		}
		rs.CCE = &cceResult
		rs.Grade = cceResult.FinalGrade
		rs.GradeLabel = ""
		subjects = append(subjects, *rs)
	}

	var finalMarks float64
	for _, s := range subjects {
		if s.CCE != nil {
			finalMarks += s.CCE.FinalMarks
		}
	}
	count := float64(len(subjects))
	if count == 0 {
		count = 1
	}
	overallPct := finalMarks / count
	overallGrade := cce.GradeFromMarks(overallPct)

	resp := StudentReportResponse{
		Student: st, AcademicYear: academicYearID, Term: term,
		Subjects: subjects, GrandTotal: grandTotal, GrandMax: grandMax, Pct: overallPct, Grade: overallGrade,
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

// GET /api/v1/reports/student-me — authenticated student's own report card
func (h *ReportsHandler) StudentSelf(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	studentID := claims.UserID
	academicYearID := h.resolveYear(r.Context(), claims.SchoolID, r.URL.Query().Get("academic_year_id"))

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
		c.name, c.code, a.max_marks::double precision
		FROM assessments a
		JOIN subjects s ON s.id = a.subject_id AND s.deleted_at IS NULL
		JOIN assessment_categories c ON c.id = a.category_id
		WHERE a.class_id = $1 AND a.deleted_at IS NULL`
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

	type rawAssy struct{ id, name, sid, scode, sname, cat, catCode string; max float64 }
	allAssy := []rawAssy{}
	for rows.Next() {
		var a rawAssy
		if err := rows.Scan(&a.id, &a.name, &a.sid, &a.scode, &a.sname, &a.cat, &a.catCode, &a.max); err != nil {
			continue
		}
		allAssy = append(allAssy, a)
	}

	marksQuery := `SELECT m.assessment_id, m.marks_obtained::double precision, m.is_absent
		FROM marks m JOIN assessments a ON a.id = m.assessment_id
		WHERE m.student_id = $1 AND a.class_id = $2 AND a.deleted_at IS NULL`
	mkRows, err := h.db.Query(r.Context(), marksQuery, studentID, classID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch marks"))
		return
	}
	defer mkRows.Close()

	markMap := map[string]MarkCell{}
	for mkRows.Next() {
		var aid string; var val float64; var absent bool
		if err := mkRows.Scan(&aid, &val, &absent); err != nil { continue }
		markMap[aid] = MarkCell{AssessmentID: aid, Value: val, IsAbsent: absent, HasMark: true}
	}

	subjMap := map[string]*ReportSubject{}
	subjOrder := []string{}
	var grandTotal, grandMax float64
	for _, a := range allAssy {
		atype := classifySubject(a.scode)
		if subjMap[a.sid] == nil {
			subjMap[a.sid] = &ReportSubject{SubjectID: a.sid, SubjectCode: a.scode, SubjectName: a.sname, SubjectType: string(atype)}
			subjOrder = append(subjOrder, a.sid)
		}
		rs := subjMap[a.sid]
		cell := ReportAssessment{ID: a.id, Name: a.name, Category: a.cat, Max: a.max}
		if m, ok := markMap[a.id]; ok {
			cell.Value = m.Value; cell.Absent = m.IsAbsent; cell.HasMark = true
			if !m.IsAbsent { rs.Total += m.Value; grandTotal += m.Value }
		}
		rs.MaxTotal += a.max; grandMax += a.max
		rs.Assessments = append(rs.Assessments, cell)
	}

	classRange := cce.DetermineClassRange(st.Class)
	subjects := make([]ReportSubject, 0, len(subjOrder))
	for _, sid := range subjOrder {
		rs := subjMap[sid]
		if rs.MaxTotal > 0 { rs.Pct = (rs.Total / rs.MaxTotal) * 100 }

		isFL := cce.IsFirstLanguage(rs.SubjectCode)
		subjType := cce.OtherSubject
		if isFL { subjType = cce.FirstLanguage }

		faMarks, saMarks := splitFASA(rs.Assessments)

		var cceResult cce.ConversionResult
		if classRange == cce.Class6to8 {
			cceResult = cce.ConvertClasses68(faMarks, saMarks)
		} else {
			cceResult = cce.ConvertClasses910(faMarks, saMarks, subjType)
		}
		rs.CCE = &cceResult
		rs.Grade = cceResult.FinalGrade
		rs.GradeLabel = ""
		subjects = append(subjects, *rs)
	}

	var finalMarks float64
	for _, s := range subjects {
		if s.CCE != nil { finalMarks += s.CCE.FinalMarks }
	}
	count := float64(len(subjects))
	if count == 0 { count = 1 }
	overallPct := finalMarks / count
	overallGrade := cce.GradeFromMarks(overallPct)

	renderJSON(w, http.StatusOK, apiOK(StudentReportResponse{
		Student: st, AcademicYear: academicYearID, Subjects: subjects,
		GrandTotal: grandTotal, GrandMax: grandMax, Pct: overallPct, Grade: overallGrade,
	}))
}

// MentorReportStudent is one student in a mentor's group.
type MentorReportStudent struct {
	StudentID    string  `json:"student_id"`
	SATSNumber   string  `json:"sats_number"`
	Name         string  `json:"name"`
	RollNo       int     `json:"roll_no"`
	Gender       string  `json:"gender"`
	ClassName    string  `json:"class_name"`
	CognitivePct float64 `json:"cognitive_pct"`
	Tier         int     `json:"tier"`
}

// MentorReportGroup is one mentor's list of students.
type MentorReportGroup struct {
	MentorID        string                `json:"mentor_id"`
	MentorName      string                `json:"mentor_name"`
	StudentCount    int                   `json:"student_count"`
	AvgCognitivePct float64               `json:"avg_cognitive_pct"`
	Students        []MentorReportStudent `json:"students"`
}

// GET /api/v1/reports/mentors?academic_year_id=
func (h *ReportsHandler) MentorReport(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	academicYearID := h.resolveYear(r.Context(), claims.SchoolID, r.URL.Query().Get("academic_year_id"))
	if academicYearID == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "no current academic year configured"))
		return
	}

	var yearName string
	if err := h.db.QueryRow(r.Context(),
		`SELECT name FROM academic_years WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		academicYearID, claims.SchoolID,
	).Scan(&yearName); err != nil {
		renderJSON(w, http.StatusNotFound, apiErr("NOT_FOUND", "academic year not found"))
		return
	}

	rows, err := h.db.Query(r.Context(), `
		WITH per_student AS (
			SELECT m.student_id, ROUND(AVG(m.marks_obtained / a.max_marks * 100)::numeric, 2) AS cognitive_pct
			FROM marks m
			JOIN assessments a ON a.id = m.assessment_id
			WHERE m.marks_obtained IS NOT NULL AND a.max_marks > 0
			GROUP BY m.student_id
		),
		ranked AS (
			SELECT s.id AS student_id, s.class_id, p.cognitive_pct,
				ROW_NUMBER() OVER (PARTITION BY s.class_id ORDER BY p.cognitive_pct DESC, s.first_name) AS rank_in_class,
				COUNT(*) OVER (PARTITION BY s.class_id) AS n_in_class
			FROM students s
			JOIN per_student p ON p.student_id = s.id
			WHERE s.deleted_at IS NULL AND s.is_active = true AND s.school_id = $2
		),
		tiered AS (
			SELECT student_id, class_id, cognitive_pct,
				LEAST(4, CEIL(rank_in_class * 4.0 / n_in_class)) AS tier
			FROM ranked
		)
		SELECT u.id AS mentor_id, u.name AS mentor_name,
			s.id AS student_id, s.sats_number,
			s.first_name || ' ' || COALESCE(s.last_name,'') AS student_name,
			COALESCE(s.roll_no, 0), COALESCE(s.gender, ''),
			c.name AS class_name, t.cognitive_pct, t.tier
		FROM mentor_assignments ma
		JOIN users u ON u.id = ma.mentor_id
		JOIN students s ON s.id = ma.student_id
		JOIN classes c ON c.id = s.class_id
		LEFT JOIN tiered t ON t.student_id = s.id
		WHERE ma.academic_year_id = $1
		  AND u.school_id = $2 AND u.deleted_at IS NULL AND u.is_active = true
		  AND s.deleted_at IS NULL AND s.is_active = true
		ORDER BY u.name ASC, t.cognitive_pct DESC NULLS LAST, s.first_name ASC`,
		academicYearID, claims.SchoolID)
	if err != nil {
		log.Error().Err(err).Msg("mentor report query failed")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch mentor report"))
		return
	}
	defer rows.Close()

	groups := map[string]*MentorReportGroup{}
	var order []string
	for rows.Next() {
		var mid, mname, sid, sats, sname, gender, cname string
		var roll int
		var cog float64
		var tier int
		if err := rows.Scan(&mid, &mname, &sid, &sats, &sname, &roll, &gender, &cname, &cog, &tier); err != nil {
			continue
		}
		g, ok := groups[mid]
		if !ok {
			g = &MentorReportGroup{MentorID: mid, MentorName: mname}
			groups[mid] = g
			order = append(order, mid)
		}
		g.Students = append(g.Students, MentorReportStudent{
			StudentID: sid, SATSNumber: sats, Name: sname, RollNo: roll,
			Gender: gender, ClassName: cname, CognitivePct: cog, Tier: tier,
		})
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("mentor report rows error")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch mentor report"))
		return
	}

	mentors := make([]MentorReportGroup, 0, len(order))
	for _, mid := range order {
		g := groups[mid]
		g.StudentCount = len(g.Students)
		var sum float64
		var cnt int
		for _, s := range g.Students {
			if s.CognitivePct > 0 {
				sum += s.CognitivePct
				cnt++
			}
		}
		if cnt > 0 {
			g.AvgCognitivePct = sum / float64(cnt)
		}
		mentors = append(mentors, *g)
	}

	renderJSON(w, http.StatusOK, apiOK(map[string]interface{}{
		"academic_year_id":   academicYearID,
		"academic_year_name": yearName,
		"mentors":            mentors,
	}))
}

func apiOK(data interface{}) map[string]interface{} {
	return map[string]interface{}{"data": data}
}

func apiErr(code, msg string) map[string]interface{} {
	return map[string]interface{}{"error": map[string]string{"code": code, "message": msg}}
}
