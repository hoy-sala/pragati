package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/rs/zerolog/log"
)

type MentorHandler struct {
	db *pgxpool.Pool
}

func NewMentorHandler(db *pgxpool.Pool) *MentorHandler {
	return &MentorHandler{db: db}
}

func (h *MentorHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	yearID := r.URL.Query().Get("academic_year_id")
	if yearID == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "academic_year_id is required"))
		return
	}
	query := `SELECT ma.id, ma.mentor_id, ma.student_id, ma.academic_year_id, u.name as mentor_name,
		s.first_name || ' ' || COALESCE(s.last_name,'') as student_name, s.sats_number, c.name as class_name
		FROM mentor_assignments ma
		JOIN users u ON u.id = ma.mentor_id
		JOIN students s ON s.id = ma.student_id
		JOIN classes c ON c.id = s.class_id
		WHERE ma.academic_year_id = $1 AND s.deleted_at IS NULL`
	args := []interface{}{yearID}
	n := 2
	if mentorID := r.URL.Query().Get("mentor_id"); mentorID != "" {
		query += fmt.Sprintf(" AND ma.mentor_id = $%d", n)
		args = append(args, mentorID)
		n++
	}
	if classID := r.URL.Query().Get("class_id"); classID != "" {
		query += fmt.Sprintf(" AND s.class_id = $%d", n)
		args = append(args, classID)
		n++
	}
	query += " ORDER BY u.name, s.first_name"
	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch assignments"))
		return
	}
	defer rows.Close()
	type view struct {
		ID string `json:"id"`; MentorID string `json:"mentor_id"`; MentorName string `json:"mentor_name"`
		StudentID string `json:"student_id"`; StudentName string `json:"student_name"`; SATSNumber string `json:"sats_number"`
		ClassName string `json:"class_name"`; AcademicYearID string `json:"academic_year_id"`
	}
	var results []view
	for rows.Next() {
		var a view
		if err := rows.Scan(&a.ID, &a.MentorID, &a.StudentID, &a.AcademicYearID, &a.MentorName, &a.StudentName, &a.SATSNumber, &a.ClassName); err != nil { continue }
		results = append(results, a)
	}
	renderJSON(w, http.StatusOK, apiOK(results))
}

func (h *MentorHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MentorID string `json:"mentor_id"`; StudentID string `json:"student_id"`; AcademicYearID string `json:"academic_year_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, apiErr("INVALID_INPUT", "invalid request"))
		return
	}
	if req.MentorID == "" || req.StudentID == "" || req.AcademicYearID == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "mentor_id, student_id, academic_year_id required"))
		return
	}
	var count int
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM mentor_assignments WHERE mentor_id = $1 AND academic_year_id = $2`, req.MentorID, req.AcademicYearID).Scan(&count)
	if count >= 25 {
		renderJSON(w, http.StatusBadRequest, apiErr("RATIO_EXCEEDED", "mentor already has 25 students (max)"))
		return
	}
	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO mentor_assignments (id, mentor_id, student_id, academic_year_id, assigned_by)
		 VALUES ($1,$2,$3,$4,$5) ON CONFLICT (mentor_id, student_id, academic_year_id) DO NOTHING`,
		id, req.MentorID, req.StudentID, req.AcademicYearID, claims.UserID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to create assignment"))
		return
	}
	renderJSON(w, http.StatusCreated, apiOK(map[string]string{"id": id}))
}

func (h *MentorHandler) DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `DELETE FROM mentor_assignments WHERE id = $1`, id)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to delete assignment"))
		return
	}
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}

func (h *MentorHandler) Stats(w http.ResponseWriter, r *http.Request) {
	yearID := r.URL.Query().Get("academic_year_id")
	if yearID == "" { renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "academic_year_id is required")); return }
	classID := r.URL.Query().Get("class_id")
	var unassignedQuery string; var args []interface{}
	if classID != "" {
		unassignedQuery = `SELECT COUNT(*) FROM students s WHERE s.class_id = $1 AND s.deleted_at IS NULL AND s.is_active = true
			AND s.id NOT IN (SELECT student_id FROM mentor_assignments WHERE academic_year_id = $2)`
		args = []interface{}{classID, yearID}
	} else {
		unassignedQuery = `SELECT COUNT(*) FROM students s WHERE s.deleted_at IS NULL AND s.is_active = true
			AND s.id NOT IN (SELECT student_id FROM mentor_assignments WHERE academic_year_id = $1)`
		args = []interface{}{yearID}
	}
	var unassigned int
	h.db.QueryRow(r.Context(), unassignedQuery, args...).Scan(&unassigned)

	mentorRows, err := h.db.Query(r.Context(),
		`SELECT u.id, u.name, COUNT(ma.id) as student_count
		 FROM users u LEFT JOIN mentor_assignments ma ON ma.mentor_id = u.id AND ma.academic_year_id = $1
		 WHERE u.role IN ('teacher','special_educator') AND u.is_active = true AND u.deleted_at IS NULL
		 GROUP BY u.id, u.name ORDER BY u.name`, yearID)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch stats")); return }
	defer mentorRows.Close()
	type ms struct { ID string `json:"id"`; Name string `json:"name"`; StudentCount int `json:"student_count"` }
	var mentors []ms
	for mentorRows.Next() {
		var m ms
		if err := mentorRows.Scan(&m.ID, &m.Name, &m.StudentCount); err != nil { continue }
		mentors = append(mentors, m)
	}
	renderJSON(w, http.StatusOK, apiOK(map[string]interface{}{"unassigned": unassigned, "mentors": mentors}))
}

func (h *MentorHandler) Roster(w http.ResponseWriter, r *http.Request) {
	mentorID := r.URL.Query().Get("mentor_id")
	if mentorID == "" { mentorID = claims.UserID }
	yearID := r.URL.Query().Get("academic_year_id")
	if yearID == "" { renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "academic_year_id is required")); return }
	rows, err := h.db.Query(r.Context(),
		`SELECT s.id, s.sats_number, s.first_name || ' ' || COALESCE(s.last_name,'') as name, s.roll_no, s.gender,
			COALESCE(s.father_name,''), COALESCE(s.mother_name,''), COALESCE(s.parent_phone,''), COALESCE(s.address,''),
			s.class_id, c.name as class_name
		 FROM mentor_assignments ma
		 JOIN students s ON s.id = ma.student_id
		 JOIN classes c ON c.id = s.class_id
		 WHERE ma.mentor_id = $1 AND ma.academic_year_id = $2 AND s.deleted_at IS NULL
		 ORDER BY c.sort_order, s.first_name`, mentorID, yearID)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch roster")); return }
	defer rows.Close()
	type rs struct {
		ID string `json:"id"`; SATSNumber string `json:"sats_number"`; Name string `json:"name"`
		RollNo int `json:"roll_no"`; Gender string `json:"gender"`; FatherName string `json:"father_name"`
		MotherName string `json:"mother_name"`; ParentPhone string `json:"parent_phone"`; Address string `json:"address"`
		ClassID string `json:"class_id"`; ClassName string `json:"class_name"`
	}
	var students []rs
	for rows.Next() {
		var s rs
		if err := rows.Scan(&s.ID, &s.SATSNumber, &s.Name, &s.RollNo, &s.Gender,
			&s.FatherName, &s.MotherName, &s.ParentPhone, &s.Address, &s.ClassID, &s.ClassName); err != nil { continue }
		students = append(students, s)
	}
	renderJSON(w, http.StatusOK, apiOK(students))
}

func (h *MentorHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	mentorID := r.URL.Query().Get("mentor_id")
	if mentorID == "" { mentorID = claims.UserID }
	date := r.URL.Query().Get("date")
	if date == "" { date = time.Now().Format("2006-01-02") }
	yearID := r.URL.Query().Get("academic_year_id")
	rows, err := h.db.Query(r.Context(),
		`SELECT ma.student_id, s.first_name || ' ' || COALESCE(s.last_name,'') as name,
			COALESCE(a.status, 'present') as status, a.parent_contacted, a.remarks
		 FROM mentor_assignments ma
		 JOIN students s ON s.id = ma.student_id
		 LEFT JOIN mentor_attendance a ON a.student_id = ma.student_id AND a.mentor_id = ma.mentor_id AND a.date = $1
		 WHERE ma.mentor_id = $2 AND ma.academic_year_id = $3 AND s.deleted_at IS NULL
		 ORDER BY s.first_name`, date, mentorID, yearID)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch attendance")); return }
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var sid, name, status, remarks string; var parentContacted bool
		if err := rows.Scan(&sid, &name, &status, &parentContacted, &remarks); err != nil { continue }
		results = append(results, map[string]interface{}{
			"student_id": sid, "name": name, "status": status, "parent_contacted": parentContacted, "remarks": remarks})
	}
	renderJSON(w, http.StatusOK, apiOK(results))
}

func (h *MentorHandler) SaveAttendance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MentorID string `json:"mentor_id"`; Date string `json:"date"`; StudentID string `json:"student_id"`
		Status string `json:"status"`; Remarks string `json:"remarks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, apiErr("INVALID_INPUT", "invalid request")); return
	}
	if req.MentorID == "" || req.StudentID == "" || req.Status == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "mentor_id, student_id, status required")); return
	}
	if req.Date == "" { req.Date = time.Now().Format("2006-01-02") }
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO mentor_attendance (id, mentor_id, student_id, date, status, remarks)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		 ON CONFLICT (mentor_id, student_id, date) DO UPDATE SET status = EXCLUDED.status, remarks = EXCLUDED.remarks`,
		req.MentorID, req.StudentID, req.Date, req.Status, req.Remarks)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to save attendance")); return }
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}

func (h *MentorHandler) ContactParent(w http.ResponseWriter, r *http.Request) {
	var req struct { MentorID string `json:"mentor_id"`; StudentID string `json:"student_id"`; Date string `json:"date"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, apiErr("INVALID_INPUT", "invalid request")); return
	}
	if req.Date == "" { req.Date = time.Now().Format("2006-01-02") }
	var id string
	err := h.db.QueryRow(r.Context(),
		`SELECT id FROM mentor_attendance WHERE mentor_id = $1 AND student_id = $2 AND date = $3`,
		req.MentorID, req.StudentID, req.Date).Scan(&id)
	if err != nil {
		h.db.Exec(r.Context(),
			`INSERT INTO mentor_attendance (id, mentor_id, student_id, date, parent_contacted, parent_contact_time)
			 VALUES (gen_random_uuid(), $1, $2, $3, true, NOW())`, req.MentorID, req.StudentID, req.Date)
	} else {
		h.db.Exec(r.Context(), `UPDATE mentor_attendance SET parent_contacted = true, parent_contact_time = NOW() WHERE id = $1`, id)
	}
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}

func (h *MentorHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	mentorID := r.URL.Query().Get("mentor_id")
	if mentorID == "" { mentorID = claims.UserID }
	rows, err := h.db.Query(r.Context(),
		`SELECT ml.id, ml.student_id, s.first_name || ' ' || COALESCE(s.last_name,'') as student_name,
			ml.log_date, ml.category, ml.severity, ml.description, ml.action_taken, ml.parent_informed,
			ml.reviewed_by_principal
		 FROM mentor_logs ml JOIN students s ON s.id = ml.student_id
		 WHERE ml.mentor_id = $1 ORDER BY ml.created_at DESC LIMIT 100`, mentorID)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch logs")); return }
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var id, sid, sname, ldate, cat, sev, desc, action string; var parentInformed, reviewed bool
		if err := rows.Scan(&id, &sid, &sname, &ldate, &cat, &sev, &desc, &action, &parentInformed, &reviewed); err != nil { continue }
		results = append(results, map[string]interface{}{
			"id": id, "student_id": sid, "student_name": sname, "log_date": ldate, "category": cat,
			"severity": sev, "description": desc, "action_taken": action, "parent_informed": parentInformed,
			"reviewed_by_principal": reviewed})
	}
	renderJSON(w, http.StatusOK, apiOK(results))
}

func (h *MentorHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID string `json:"student_id"`; Category string `json:"category"`
		Severity string `json:"severity"`; Description string `json:"description"`
		ActionTaken string `json:"action_taken"`; ParentInformed bool `json:"parent_informed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, apiErr("INVALID_INPUT", "invalid request")); return
	}
	if req.StudentID == "" || req.Category == "" || req.Description == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "student_id, category, description required")); return
	}
	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO mentor_logs (id, mentor_id, student_id, category, severity, description, action_taken, parent_informed)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		id, claims.UserID, req.StudentID, req.Category, req.Severity, req.Description, req.ActionTaken, req.ParentInformed)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to create log")); return }
	renderJSON(w, http.StatusCreated, apiOK(map[string]string{"id": id}))
}

func (h *MentorHandler) PrincipalAlerts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT ml.id, ml.student_id, s.first_name || ' ' || COALESCE(s.last_name,'') as student_name,
			ml.log_date, ml.category, ml.severity, ml.description, ml.action_taken, ml.parent_informed,
			u.name as mentor_name
		 FROM mentor_logs ml
		 JOIN students s ON s.id = ml.student_id
		 JOIN users u ON u.id = ml.mentor_id
		 WHERE ml.severity IN ('high','urgent') AND ml.reviewed_by_principal = false
		 ORDER BY ml.created_at DESC LIMIT 50`)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch alerts")); return }
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var id, sid, sname, ldate, cat, sev, desc, action, mentor string; var parentInformed bool
		if err := rows.Scan(&id, &sid, &sname, &ldate, &cat, &sev, &desc, &action, &parentInformed, &mentor); err != nil { continue }
		results = append(results, map[string]interface{}{
			"id": id, "student_id": sid, "student_name": sname, "log_date": ldate, "category": cat,
			"severity": sev, "description": desc, "action_taken": action, "parent_informed": parentInformed,
			"mentor_name": mentor})
	}
	renderJSON(w, http.StatusOK, apiOK(results))
}

func (h *MentorHandler) ReviewLog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct { PrincipalNotes string `json:"principal_notes"` }
	json.NewDecoder(r.Body).Decode(&req)
	_, err := h.db.Exec(r.Context(),
		`UPDATE mentor_logs SET reviewed_by_principal = true, principal_notes = $1 WHERE id = $2`, req.PrincipalNotes, id)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to review log")); return }
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}

func (h *MentorHandler) MonthlySummary(w http.ResponseWriter, r *http.Request) {
	yearID := r.URL.Query().Get("academic_year_id")
	if yearID == "" { renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "academic_year_id is required")); return }
	rows, err := h.db.Query(r.Context(),
		`SELECT u.id, u.name,
			COUNT(DISTINCT ma.student_id) as student_count,
			COUNT(DISTINCT CASE WHEN a.status = 'present' THEN a.student_id END) as present_days,
			COUNT(DISTINCT a.student_id) as total_attendance_records,
			COUNT(DISTINCT ml.id) as log_count,
			COUNT(DISTINCT CASE WHEN ml.severity IN ('high','urgent') THEN ml.id END) as urgent_count,
			COUNT(DISTINCT CASE WHEN a.parent_contacted = true THEN a.id END) as parent_contacts
		 FROM users u
		 LEFT JOIN mentor_assignments ma ON ma.mentor_id = u.id AND ma.academic_year_id = $1
		 LEFT JOIN mentor_attendance a ON a.mentor_id = u.id
		 LEFT JOIN mentor_logs ml ON ml.mentor_id = u.id
		 WHERE u.role IN ('teacher','special_educator') AND u.is_active = true AND u.deleted_at IS NULL
		 GROUP BY u.id, u.name ORDER BY u.name`, yearID)
	if err != nil { renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch summary")); return }
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var id, name string; var sc, pd, tar, lc, uc, pc int
		if err := rows.Scan(&id, &name, &sc, &pd, &tar, &lc, &uc, &pc); err != nil { continue }
		avgAtt := 0.0
		if tar > 0 { avgAtt = float64(pd) / float64(tar) * 100 }
		results = append(results, map[string]interface{}{
			"mentor_id": id, "mentor_name": name, "student_count": sc,
			"avg_attendance_pct": avgAtt, "log_count": lc, "urgent_count": uc, "parent_contacts": pc})
	}
	renderJSON(w, http.StatusOK, apiOK(results))
}
