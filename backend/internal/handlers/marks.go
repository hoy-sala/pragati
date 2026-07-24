package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type MarkHandler struct {
	db *pgxpool.Pool
}

func NewMarkHandler(db *pgxpool.Pool) *MarkHandler {
	return &MarkHandler{db: db}
}

func (h *MarkHandler) GetGrid(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	assessmentID := r.URL.Query().Get("assessment_id")
	if assessmentID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "assessment_id is required"}})
		return
	}

	var assessment models.Assessment
	err := h.db.QueryRow(r.Context(),
		`SELECT id, class_id, section_id, max_marks, is_locked, version
		 FROM assessments WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		assessmentID, claims.SchoolID,
	).Scan(&assessment.ID, &assessment.ClassID, &assessment.SectionID, &assessment.MaxMarks, &assessment.IsLocked, &assessment.Version)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "assessment not found"}})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT s.id, s.sats_number, s.first_name, COALESCE(s.last_name, ''), s.roll_no,
			m.id, m.assessment_id, COALESCE(m.marks_obtained, -1), COALESCE(m.is_absent, false), COALESCE(m.remarks, '')
		FROM students s
		LEFT JOIN marks m ON m.student_id = s.id AND m.assessment_id = $1
		WHERE s.class_id = $2 AND s.deleted_at IS NULL AND s.is_active = true
		ORDER BY s.roll_no ASC, s.first_name ASC`,
		assessmentID, assessment.ClassID,
	)
	if err != nil {
		log.Error().Err(err).Msg("get marks grid failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch marks"}})
		return
	}
	defer rows.Close()

	type StudentRow struct {
		StudentID     string  `json:"student_id"`
		SATSNumber    string  `json:"sats_number"`
		Name          string  `json:"name"`
		RollNo        int     `json:"roll_no"`
		MarkID        string  `json:"mark_id,omitempty"`
		MarksObtained float64 `json:"marks_obtained"`
		IsAbsent      bool    `json:"is_absent"`
		Remarks       string  `json:"remarks"`
	}

	grid := []StudentRow{}
	for rows.Next() {
		var r StudentRow
		var lastName, markID, remarks string
		var marksObtained float64
		var isAbsent bool
		if err := rows.Scan(&r.StudentID, &r.SATSNumber, &r.Name, &lastName, &r.RollNo,
			&markID, &r.MarkID, &marksObtained, &isAbsent, &remarks); err != nil {
			log.Error().Err(err).Msg("scan marks row failed")
			continue
		}
		if lastName != "" {
			r.Name += " " + lastName
		}
		r.MarkID = markID
		r.MarksObtained = marksObtained
		r.IsAbsent = isAbsent
		r.Remarks = remarks
		grid = append(grid, r)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"assessment": assessment,
		"students":   grid,
	}})
}

func (h *MarkHandler) BatchSave(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	var req struct {
		AssessmentID string          `json:"assessment_id"`
		Version      int             `json:"version"`
		Marks        []models.MarkInput `json:"marks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to start transaction"}})
		return
	}
	defer tx.Rollback(r.Context())

	var currentVersion int
	err = tx.QueryRow(r.Context(),
		`SELECT version FROM assessments WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL FOR UPDATE`,
		req.AssessmentID, claims.SchoolID,
	).Scan(&currentVersion)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "assessment not found"}})
		return
	}

	if currentVersion != req.Version {
		renderJSON(w, http.StatusConflict, models.APIResponse{Error: &models.APIError{Code: "VERSION_CONFLICT", Message: "marks were updated by another user. Please refresh."}})
		return
	}

	updated := 0
	errors := []map[string]interface{}{}

	for _, m := range req.Marks {
		var maxMarks float64
		tx.QueryRow(r.Context(), "SELECT max_marks FROM assessments WHERE id = $1", req.AssessmentID).Scan(&maxMarks)

		if m.MarksObtained > maxMarks {
			errors = append(errors, map[string]interface{}{
				"student_id":    m.StudentID,
				"marks_obtained": m.MarksObtained,
				"max_marks":     maxMarks,
				"message":       fmt.Sprintf("marks (%.2f) exceed maximum (%.2f)", m.MarksObtained, maxMarks),
			})
			continue
		}

		_, err := tx.Exec(r.Context(),
			`INSERT INTO marks (id, assessment_id, student_id, marks_obtained, is_absent, remarks, entered_by, entered_at, updated_at)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW())
			 ON CONFLICT (assessment_id, student_id) DO UPDATE SET
				marks_obtained = EXCLUDED.marks_obtained,
				is_absent = EXCLUDED.is_absent,
				remarks = EXCLUDED.remarks,
				entered_by = EXCLUDED.entered_by,
				updated_at = NOW()`,
			req.AssessmentID, m.StudentID, m.MarksObtained, m.IsAbsent, m.Remarks, claims.UserID,
		)
		if err != nil {
			log.Error().Err(err).Str("student", m.StudentID).Msg("save mark failed")
			errors = append(errors, map[string]interface{}{
				"student_id": m.StudentID,
				"message":    "database error",
			})
			continue
		}
		updated++
	}

	if _, err := tx.Exec(r.Context(),
		`UPDATE assessments SET version = version + 1, updated_at = NOW() WHERE id = $1`,
		req.AssessmentID,
	); err != nil {
		log.Error().Err(err).Msg("increment version failed")
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Error().Err(err).Msg("commit marks batch failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save marks"}})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"updated": updated,
		"errors":  errors,
	}})
}

func (h *MarkHandler) ImportExcel(w http.ResponseWriter, r *http.Request) {
	assessmentID := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "file is required"}})
		return
	}
	defer file.Close()

	rows, err := parseCSV(file)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "failed to parse file"}})
		return
	}

	if len(rows) < 1 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "file is empty"}})
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "transaction failed"}})
		return
	}
	defer tx.Rollback(r.Context())

	headers := rows[0]
	dataRows := rows[1:]

	satsIdx := indexOf(headers, "sats_number")
	marksIdx := indexOf(headers, "marks")
	absentIdx := indexOf(headers, "absent")
	remarksIdx := indexOf(headers, "remarks")

	if satsIdx < 0 || marksIdx < 0 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "file must have 'sats_number' and 'marks' columns"}})
		return
	}

	updated := 0
	for _, row := range dataRows {
		if len(row) <= satsIdx || len(row) <= marksIdx {
			continue
		}
		satsNo := row[satsIdx]
		if satsNo == "" {
			continue
		}

		var studentID string
		err := tx.QueryRow(r.Context(),
			`SELECT id FROM students WHERE sats_number = $1 AND school_id = $2 AND deleted_at IS NULL`,
			satsNo, claims.SchoolID,
		).Scan(&studentID)
		if err != nil {
			continue
		}

		marks := parseFloat(row[marksIdx])
		isAbsent := false
		if absentIdx >= 0 && len(row) > absentIdx {
			isAbsent = row[absentIdx] == "true" || row[absentIdx] == "yes" || row[absentIdx] == "1"
		}
		remarks := ""
		if remarksIdx >= 0 && len(row) > remarksIdx {
			remarks = row[remarksIdx]
		}

		if isAbsent {
			marks = 0
		}

		if _, err := tx.Exec(r.Context(),
			`INSERT INTO marks (id, assessment_id, student_id, marks_obtained, is_absent, remarks, entered_by, entered_at, updated_at)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW())
			 ON CONFLICT (assessment_id, student_id) DO UPDATE SET
				marks_obtained = EXCLUDED.marks_obtained,
				is_absent = EXCLUDED.is_absent,
				remarks = EXCLUDED.remarks,
				entered_by = EXCLUDED.entered_by,
				updated_at = NOW()`,
			assessmentID, studentID, marks, isAbsent, remarks, claims.UserID,
		); err != nil {
			log.Error().Err(err).Msg("import mark row failed")
		}
		updated++
	}

	if _, err := tx.Exec(r.Context(),
		`UPDATE assessments SET version = version + 1 WHERE id = $1`, assessmentID,
	); err != nil {
		log.Error().Err(err).Msg("increment version failed")
	}
	if err := tx.Commit(r.Context()); err != nil {
		log.Error().Err(err).Msg("commit import failed")
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"imported": updated,
	}})
}

func parseCSV(file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func parseFloat(s string) float64 {
	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}

func indexOf(arr []string, target string) int {
	for i, s := range arr {
		if s == target {
			return i
		}
	}
	return -1
}
