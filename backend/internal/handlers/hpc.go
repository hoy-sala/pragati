package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/pragati/backend/internal/services/hpc"
	"github.com/rs/zerolog/log"
)

type HPCHandler struct {
	db *pgxpool.Pool
}

func NewHPCHandler(db *pgxpool.Pool) *HPCHandler {
	return &HPCHandler{db: db}
}

// GET /api/v1/hpc/config?class_id=&stage=
func (h *HPCHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	stage := r.URL.Query().Get("stage")

	query := `SELECT id, school_id, stage, class_id, academic_year_id,
		grading_scheme, proficiency_scale, co_scholastic_areas, health_params, terms,
		is_active, created_at, updated_at
		FROM hpc_configurations
		WHERE school_id = $1 AND deleted_at IS NULL`
	args := []interface{}{claims.SchoolID}
	n := 2

	if classID != "" {
		query += fmt.Sprintf(" AND (class_id = $%d OR class_id IS NULL)", n)
		args = append(args, classID)
		n++
	}
	if stage != "" {
		query += fmt.Sprintf(" AND stage = $%d", n)
		args = append(args, stage)
	}

	query += " ORDER BY class_id NULLS LAST LIMIT 1"

	var cfg models.HPCConfiguration
	var schemeJSON, profJSON, coSchJSON, healthJSON []byte
	err := h.db.QueryRow(r.Context(), query, args...).Scan(
		&cfg.ID, &cfg.SchoolID, &cfg.Stage, &cfg.ClassID, &cfg.AcademicYearID,
		&schemeJSON, &profJSON, &coSchJSON, &healthJSON, &cfg.Terms,
		&cfg.IsActive, &cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "no HPC configuration found"},
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("get hpc config failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch HPC config"},
		})
		return
	}

	json.Unmarshal(schemeJSON, &cfg.GradingScheme)
	json.Unmarshal(profJSON, &cfg.ProficiencyScale)
	json.Unmarshal(coSchJSON, &cfg.CoScholasticAreas)
	json.Unmarshal(healthJSON, &cfg.HealthParams)

	renderJSON(w, http.StatusOK, models.APIResponse{Data: cfg})
}

// PUT /api/v1/hpc/config
func (h *HPCHandler) SaveConfig(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.HPCConfigInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	schemeJSON, _ := json.Marshal(req.GradingScheme)
	profJSON, _ := json.Marshal(req.ProficiencyScale)
	coSchJSON, _ := json.Marshal(req.CoScholasticAreas)
	healthJSON, _ := json.Marshal(req.HealthParams)
	termsJSON, _ := json.Marshal(req.Terms)

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO hpc_configurations (id, school_id, stage, class_id, academic_year_id,
			grading_scheme, proficiency_scale, co_scholastic_areas, health_params, terms,
			is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,true,NOW(),NOW())
		 ON CONFLICT ON CONSTRAINT idx_hpc_config_unique DO UPDATE SET
			grading_scheme = EXCLUDED.grading_scheme,
			proficiency_scale = EXCLUDED.proficiency_scale,
			co_scholastic_areas = EXCLUDED.co_scholastic_areas,
			health_params = EXCLUDED.health_params,
			terms = EXCLUDED.terms,
			updated_at = NOW()`,
		id, claims.SchoolID, req.Stage, nullIfEmpty(req.ClassID), nullIfEmpty(req.AcademicYearID),
		string(schemeJSON), string(profJSON), string(coSchJSON), string(healthJSON), string(termsJSON),
	)
	if err != nil {
		log.Error().Err(err).Msg("save hpc config failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save HPC config"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": id}})
}

// GET /api/v1/hpc/learning-outcomes?subject_id=&class_id=
func (h *HPCHandler) ListLearningOutcomes(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	subjectID := r.URL.Query().Get("subject_id")
	classID := r.URL.Query().Get("class_id")

	query := `SELECT id, school_id, subject_id, class_id, code, description, domain,
		expected_level, sort_order, is_active, created_at, updated_at
		FROM learning_outcomes
		WHERE school_id = $1 AND deleted_at IS NULL`
	args := []interface{}{claims.SchoolID}
	n := 2

	if subjectID != "" {
		query += fmt.Sprintf(" AND subject_id = $%d", n)
		args = append(args, subjectID)
		n++
	}
	if classID != "" {
		query += fmt.Sprintf(" AND class_id = $%d", n)
		args = append(args, classID)
		n++
	}
	query += " ORDER BY subject_id, sort_order ASC, code ASC"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("list learning outcomes failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch learning outcomes"},
		})
		return
	}
	defer rows.Close()

	los := []models.LearningOutcome{}
	for rows.Next() {
		var lo models.LearningOutcome
		if err := rows.Scan(&lo.ID, &lo.SchoolID, &lo.SubjectID, &lo.ClassID, &lo.Code,
			&lo.Description, &lo.Domain, &lo.ExpectedLevel, &lo.SortOrder,
			&lo.IsActive, &lo.CreatedAt, &lo.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan learning outcome failed")
			continue
		}
		los = append(los, lo)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: los})
}

// POST /api/v1/hpc/learning-outcomes/import
func (h *HPCHandler) ImportLearningOutcomes(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		SubjectID string                   `json:"subject_id"`
		ClassID   string                   `json:"class_id"`
		Outcomes  []models.LearningOutcome `json:"outcomes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"},
		})
		return
	}

	imported := 0
	for _, lo := range req.Outcomes {
		id := uuid.New().String()
		_, err := h.db.Exec(r.Context(),
			`INSERT INTO learning_outcomes (id, school_id, subject_id, class_id, code, description,
				domain, expected_level, sort_order, is_active, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,true,NOW(),NOW())
			 ON CONFLICT ON CONSTRAINT idx_lo_code_unique DO UPDATE SET
				description = EXCLUDED.description,
				domain = EXCLUDED.domain,
				expected_level = EXCLUDED.expected_level,
				sort_order = EXCLUDED.sort_order,
				updated_at = NOW()`,
			id, claims.SchoolID, req.SubjectID, req.ClassID, lo.Code, lo.Description,
			lo.Domain, lo.ExpectedLevel, lo.SortOrder,
		)
		if err != nil {
			log.Error().Err(err).Str("code", lo.Code).Msg("import LO failed")
			continue
		}
		imported++
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]int{"imported": imported}})
}

// GET /api/v1/hpc/grid?class_id=&term=&academic_year_id=
func (h *HPCHandler) GetGrid(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	term := r.URL.Query().Get("term")
	academicYearID := r.URL.Query().Get("academic_year_id")

	if classID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "class_id is required"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT s.id, s.sats_number, s.first_name, COALESCE(s.last_name, ''), s.roll_no,
			e.id, COALESCE(e.status, ''), CASE WHEN e.generated_pdf_url != '' THEN true ELSE false END
		FROM students s
		LEFT JOIN hpc_entries e ON e.student_id = s.id AND e.term = $1
			AND e.academic_year_id = $2 AND e.deleted_at IS NULL
		WHERE s.class_id = $3 AND s.deleted_at IS NULL AND s.is_active = true
		ORDER BY s.roll_no ASC, s.first_name ASC`,
		term, academicYearID, classID,
	)
	if err != nil {
		log.Error().Err(err).Msg("get hpc grid failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch HPC grid"},
		})
		return
	}
	defer rows.Close()

	grid := []models.HPCGridRow{}
	for rows.Next() {
		var r models.HPCGridRow
		var lastName string
		if err := rows.Scan(&r.StudentID, &r.SATSNumber, &r.Name, &lastName, &r.RollNo,
			&r.EntryID, &r.Status, &r.HasPDF); err != nil {
			log.Error().Err(err).Msg("scan hpc grid row failed")
			continue
		}
		if lastName != "" {
			r.Name += " " + lastName
		}
		grid = append(grid, r)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: grid})
}

// GET /api/v1/hpc/entries?student_id=&term=&academic_year_id=
func (h *HPCHandler) GetEntry(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	studentID := r.URL.Query().Get("student_id")
	term := r.URL.Query().Get("term")
	academicYearID := r.URL.Query().Get("academic_year_id")

	// always return student info
	var studentName string
	err := h.db.QueryRow(r.Context(),
		`SELECT first_name || ' ' || COALESCE(last_name, '') FROM students
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		studentID, claims.SchoolID,
	).Scan(&studentName)
	if err == pgx.ErrNoRows {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "student not found"},
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("get student for hpc entry failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch student"},
		})
		return
	}

	var entry models.HPCEntry
	var scholasticJSON, coSchJSON, healthJSON, workJSON []byte
	var selfJSON, peerJSON, parentJSON, attJSON []byte

	err = h.db.QueryRow(r.Context(),
		`SELECT id, student_id, academic_year_id, term, status,
			scholastic, co_scholastic, health_pe, work_education,
			self_assessment, peer_assessment, parent_feedback, teacher_remarks,
			attendance_summary, generated_pdf_url, version,
			locked_at, locked_by, created_at, updated_at
		 FROM hpc_entries
		 WHERE student_id = $1 AND term = $2 AND academic_year_id = $3 AND deleted_at IS NULL`,
		studentID, term, academicYearID,
	).Scan(&entry.ID, &entry.StudentID, &entry.AcademicYearID, &entry.Term, &entry.Status,
		&scholasticJSON, &coSchJSON, &healthJSON, &workJSON,
		&selfJSON, &peerJSON, &parentJSON, &entry.TeacherRemarks,
		&attJSON, &entry.GeneratedPdfURL, &entry.Version,
		&entry.LockedAt, &entry.LockedBy, &entry.CreatedAt, &entry.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		// return empty entry with student info for new entry creation
		renderJSON(w, http.StatusOK, models.APIResponse{
			Data: map[string]interface{}{
				"id":             "",
				"student_id":     studentID,
				"student_name":   studentName,
				"status":         "new",
				"scholastic":     []models.ScholasticEntry{},
				"co_scholastic":  map[string]interface{}{},
				"health_pe":      map[string]interface{}{},
				"work_education": map[string]interface{}{},
				"self_assessment": map[string]interface{}{},
				"peer_assessment": map[string]interface{}{},
				"parent_feedback": map[string]interface{}{},
				"teacher_remarks": "",
				"attendance_summary": map[string]interface{}{},
				"version":        1,
			},
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("get hpc entry failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch HPC entry"},
		})
		return
	}

	json.Unmarshal(scholasticJSON, &entry.Scholastic)
	json.Unmarshal(coSchJSON, &entry.CoScholastic)
	json.Unmarshal(healthJSON, &entry.HealthPE)
	json.Unmarshal(workJSON, &entry.WorkEducation)
	json.Unmarshal(selfJSON, &entry.SelfAssessment)
	json.Unmarshal(peerJSON, &entry.PeerAssessment)
	json.Unmarshal(parentJSON, &entry.ParentFeedback)
	json.Unmarshal(attJSON, &entry.AttendanceSummary)

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"id":                 entry.ID,
			"student_id":         entry.StudentID,
			"student_name":       studentName,
			"academic_year_id":   entry.AcademicYearID,
			"term":               entry.Term,
			"status":              entry.Status,
			"scholastic":         entry.Scholastic,
			"co_scholastic":      entry.CoScholastic,
			"health_pe":          entry.HealthPE,
			"work_education":     entry.WorkEducation,
			"self_assessment":    entry.SelfAssessment,
			"peer_assessment":    entry.PeerAssessment,
			"parent_feedback":    entry.ParentFeedback,
			"teacher_remarks":    entry.TeacherRemarks,
			"attendance_summary": entry.AttendanceSummary,
			"generated_pdf_url":  entry.GeneratedPdfURL,
			"version":            entry.Version,
			"locked_at":          entry.LockedAt,
			"locked_by":          entry.LockedBy,
		},
	})
}

// PUT /api/v1/hpc/entries
func (h *HPCHandler) SaveEntry(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		ID       string           `json:"id,omitempty"`
		Entry    models.HPCEntryInput `json:"entry"`
		OverrideVersion int       `json:"version,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"},
		})
		return
	}

	scholasticJSON, _ := json.Marshal(req.Entry.Scholastic)
	coSchJSON, _ := json.Marshal(req.Entry.CoScholastic)
	healthJSON, _ := json.Marshal(req.Entry.HealthPE)
	workJSON, _ := json.Marshal(req.Entry.WorkEducation)
	selfJSON, _ := json.Marshal(req.Entry.SelfAssessment)
	peerJSON, _ := json.Marshal(req.Entry.PeerAssessment)
	parentJSON, _ := json.Marshal(req.Entry.ParentFeedback)
	attJSON, _ := json.Marshal(req.Entry.AttendanceSummary)

	if req.ID != "" {
		result, err := h.db.Exec(r.Context(),
			`UPDATE hpc_entries SET
				scholastic = $1, co_scholastic = $2, health_pe = $3, work_education = $4,
				self_assessment = $5, peer_assessment = $6, parent_feedback = $7,
				teacher_remarks = $8, attendance_summary = $9,
				version = version + 1, updated_at = NOW()
			 WHERE id = $10 AND deleted_at IS NULL AND version = $11`,
			string(scholasticJSON), string(coSchJSON), string(healthJSON), string(workJSON),
			string(selfJSON), string(peerJSON), string(parentJSON), req.Entry.TeacherRemarks,
			string(attJSON), req.ID, req.OverrideVersion,
		)
		if err != nil {
			log.Error().Err(err).Msg("update hpc entry failed")
			renderJSON(w, http.StatusInternalServerError, models.APIResponse{
				Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update HPC entry"},
			})
			return
		}
		if result.RowsAffected() == 0 {
			renderJSON(w, http.StatusConflict, models.APIResponse{
				Error: &models.APIError{Code: "VERSION_CONFLICT", Message: "entry was modified by another user. Please refresh."},
			})
			return
		}
		renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": req.ID}})
	} else {
		id := uuid.New().String()
		_, err := h.db.Exec(r.Context(),
			`INSERT INTO hpc_entries (id, student_id, academic_year_id, term, status,
				scholastic, co_scholastic, health_pe, work_education,
				self_assessment, peer_assessment, parent_feedback, teacher_remarks,
				attendance_summary, version, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,'draft',$5,$6,$7,$8,$9,$10,$11,$12,$13,1,NOW(),NOW())`,
			id, req.Entry.StudentID, req.Entry.AcademicYearID, req.Entry.Term,
			string(scholasticJSON), string(coSchJSON), string(healthJSON), string(workJSON),
			string(selfJSON), string(peerJSON), string(parentJSON), req.Entry.TeacherRemarks,
			string(attJSON),
		)
		if err != nil {
			log.Error().Err(err).Msg("create hpc entry failed")
			renderJSON(w, http.StatusInternalServerError, models.APIResponse{
				Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create HPC entry"},
			})
			return
		}
		renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
	}
}

// POST /api/v1/hpc/entries/publish
func (h *HPCHandler) PublishEntry(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		EntryID string `json:"entry_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "entry_id is required"},
		})
		return
	}

	result, err := h.db.Exec(r.Context(),
		`UPDATE hpc_entries SET status = 'published', locked_at = NOW(), locked_by = $1, updated_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL
		 AND student_id IN (SELECT id FROM students WHERE school_id = $3)`,
		claims.UserID, req.EntryID, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "HPC entry not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

// POST /api/v1/hpc/assess
func (h *HPCHandler) AssessLO(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.LOBatchAssessmentInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"},
		})
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "transaction failed"},
		})
		return
	}
	defer tx.Rollback(r.Context())

	saved := 0
	for _, a := range req.Assessments {
		_, err := tx.Exec(r.Context(),
			`INSERT INTO learning_outcome_assessments
				(id, student_id, learning_outcome_id, subject_id, term, proficiency_level, assessed_by, assessment_date, created_at, updated_at)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, CURRENT_DATE, NOW(), NOW())
			 ON CONFLICT (student_id, learning_outcome_id, term) WHERE deleted_at IS NULL DO UPDATE SET
				proficiency_level = EXCLUDED.proficiency_level,
				assessed_by = EXCLUDED.assessed_by,
				assessment_date = EXCLUDED.assessment_date,
				updated_at = NOW()`,
			req.StudentID, a.LearningOutcomeID, req.SubjectID, req.Term,
			a.ProficiencyLevel, claims.UserID,
		)
		if err != nil {
			log.Error().Err(err).Msg("save LO assessment failed")
			continue
		}
		saved++
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Error().Err(err).Msg("commit LO assessments failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save assessments"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]int{"saved": saved}})
}

// GET /api/v1/hpc/assessments?class_id=&subject_id=&term=
func (h *HPCHandler) GetLOAssessmentGrid(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	subjectID := r.URL.Query().Get("subject_id")
	term := r.URL.Query().Get("term")

	if classID == "" || subjectID == "" || term == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "class_id, subject_id, and term are required"},
		})
		return
	}

	los, err := h.db.Query(r.Context(),
		`SELECT id, code, description, domain, expected_level
		 FROM learning_outcomes
		 WHERE subject_id = $1 AND class_id = $2 AND school_id = $3 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, code ASC`,
		subjectID, classID, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("fetch LOs for grid failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch LOs"},
		})
		return
	}
	defer los.Close()

	type LOCol struct {
		ID            string `json:"id"`
		Code          string `json:"code"`
		Description   string `json:"description"`
		Domain        string `json:"domain"`
		ExpectedLevel int    `json:"expected_level"`
	}
	loCols := []LOCol{}
	for los.Next() {
		var c LOCol
		if err := los.Scan(&c.ID, &c.Code, &c.Description, &c.Domain, &c.ExpectedLevel); err != nil {
			continue
		}
		loCols = append(loCols, c)
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT s.id, s.sats_number, s.first_name, COALESCE(s.last_name, ''), s.roll_no
		 FROM students s
		 WHERE s.class_id = $1 AND s.school_id = $2 AND s.deleted_at IS NULL AND s.is_active = true
		 ORDER BY s.roll_no ASC`,
		classID, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("fetch students for LO grid failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch students"},
		})
		return
	}
	defer rows.Close()

	type StudentRow struct {
		StudentID  string `json:"student_id"`
		SATSNumber string `json:"sats_number"`
		Name       string `json:"name"`
		RollNo     int    `json:"roll_no"`
	}
	students := []StudentRow{}
	for rows.Next() {
		var s StudentRow
		var lastName string
		if err := rows.Scan(&s.StudentID, &s.SATSNumber, &s.Name, &lastName, &s.RollNo); err != nil {
			continue
		}
		if lastName != "" {
			s.Name += " " + lastName
		}
		students = append(students, s)
	}

	type CellValue struct {
		AssessmentID string `json:"assessment_id,omitempty"`
		Level        int    `json:"level,omitempty"`
	}
	type RowData struct {
		Student StudentRow            `json:"student"`
		Cells   map[string]CellValue  `json:"cells"`
	}

	gridData := []RowData{}
	for _, s := range students {
		row := RowData{Student: s, Cells: make(map[string]CellValue)}
		for _, lo := range loCols {
			row.Cells[lo.ID] = CellValue{}
		}
		gridData = append(gridData, row)
	}

	assessRows, err := h.db.Query(r.Context(),
		`SELECT loa.student_id, loa.learning_outcome_id, loa.proficiency_level, loa.id
		 FROM learning_outcome_assessments loa
		 JOIN students s ON s.id = loa.student_id
		 WHERE loa.subject_id = $1 AND loa.term = $2 AND s.class_id = $3 AND loa.deleted_at IS NULL`,
		subjectID, term, classID,
	)
	if err == nil {
		defer assessRows.Close()
		for assessRows.Next() {
			var studentID, loID, assessmentID string
			var level int
			if err := assessRows.Scan(&studentID, &loID, &level, &assessmentID); err != nil {
				continue
			}
			for i := range gridData {
				if gridData[i].Student.StudentID == studentID {
					if _, ok := gridData[i].Cells[loID]; ok {
						gridData[i].Cells[loID] = CellValue{AssessmentID: assessmentID, Level: level}
					}
					break
				}
			}
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"columns":  loCols,
		"students": students,
		"grid":     gridData,
	}})
}

// POST /api/v1/hpc/entries/generate-pdf
func (h *HPCHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EntryID string `json:"entry_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "entry_id is required"},
		})
		return
	}

	pdfURL := fmt.Sprintf("/api/v1/hpc/entries/%s/pdf", req.EntryID)

	_, err := h.db.Exec(r.Context(),
		`UPDATE hpc_entries SET generated_pdf_url = $1, updated_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL`,
		pdfURL, req.EntryID,
	)
	if err != nil {
		log.Error().Err(err).Msg("update pdf url failed")
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{
		"pdf_url": pdfURL,
		"message": "PDF generation queued. Use the URL to download when ready.",
	}})
}

// GET /api/v1/hpc/reports/class?class_id=&term=&academic_year_id=
func (h *HPCHandler) GetClassReport(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	term := r.URL.Query().Get("term")
	academicYearID := r.URL.Query().Get("academic_year_id")

	var total, published, draft int
	err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN e.status = 'published' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN e.status = 'draft' THEN 1 ELSE 0 END), 0)
		 FROM students s
		 LEFT JOIN hpc_entries e ON e.student_id = s.id AND e.term = $1
			AND e.academic_year_id = $2 AND e.deleted_at IS NULL
		 WHERE s.class_id = $3 AND s.deleted_at IS NULL AND s.is_active = true`,
		term, academicYearID, classID,
	).Scan(&total, &published, &draft)
	if err != nil {
		log.Error().Err(err).Msg("class hpc report failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate report"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: models.HPCReportSummary{
			TotalStudents:  total,
			PublishedCount: published,
			DraftCount:     draft,
		},
	})
}

// POST /api/v1/hpc/migrate-from-marks?class_id=&academic_year_id=&term=
func (h *HPCHandler) MigrateFromMarks(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	classID := r.URL.Query().Get("class_id")
	academicYearID := r.URL.Query().Get("academic_year_id")
	term := r.URL.Query().Get("term")
	if classID == "" || academicYearID == "" || term == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "class_id, academic_year_id, and term are required"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT s.id, a.subject_id,
			COALESCE(SUM(m.marks_obtained), 0),
			COALESCE(a.max_marks, 100)
		FROM students s
		JOIN assessments a ON a.class_id = s.class_id AND a.school_id = $1
			AND a.academic_year_id = $2 AND a.deleted_at IS NULL
		LEFT JOIN marks m ON m.assessment_id = a.id AND m.student_id = s.id AND m.deleted_at IS NULL
		WHERE s.class_id = $3 AND s.school_id = $1 AND s.deleted_at IS NULL AND s.is_active = true
		GROUP BY s.id, a.subject_id, a.max_marks
		ORDER BY s.id, a.subject_id`,
		claims.SchoolID, academicYearID, classID,
	)
	if err != nil {
		log.Error().Err(err).Msg("migrate marks query failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to query marks"},
		})
		return
	}
	defer rows.Close()

	type subjectMarks struct {
		Scored, Max float64
	}
	studentMarks := make(map[string]map[string]subjectMarks)

	for rows.Next() {
		var studentID, subjectID string
		var scored, max float64
		if err := rows.Scan(&studentID, &subjectID, &scored, &max); err != nil {
			continue
		}
		if studentMarks[studentID] == nil {
			studentMarks[studentID] = make(map[string]subjectMarks)
		}
		sm := studentMarks[studentID][subjectID]
		sm.Scored += scored
		sm.Max += max
		studentMarks[studentID][subjectID] = sm
	}

	var scheme []models.GradeDefinition
	err = h.db.QueryRow(r.Context(),
		`SELECT grading_scheme FROM hpc_configurations
		 WHERE school_id = $1 AND deleted_at IS NULL LIMIT 1`,
		claims.SchoolID,
	).Scan(&scheme)
	if err != nil {
		scheme = []models.GradeDefinition{
			{Grade: "A+", MinPct: 90, MaxPct: 100, Descriptor: "Outstanding"},
			{Grade: "A", MinPct: 75, MaxPct: 89, Descriptor: "Excellent"},
			{Grade: "B+", MinPct: 60, MaxPct: 74, Descriptor: "Very Good"},
			{Grade: "B", MinPct: 45, MaxPct: 59, Descriptor: "Good"},
			{Grade: "C", MinPct: 33, MaxPct: 44, Descriptor: "Satisfactory"},
			{Grade: "D", MinPct: 20, MaxPct: 32, Descriptor: "Needs Improvement"},
			{Grade: "E", MinPct: 0, MaxPct: 19, Descriptor: "Requires Remedial"},
		}
	}

	migrated := 0
	for studentID, subjects := range studentMarks {
		scholastic := hpc.ComputeHPCScholastic(subjects, scheme)
		scholasticJSON, _ := json.Marshal(scholastic)

		_, err := h.db.Exec(r.Context(),
			`INSERT INTO hpc_entries (id, student_id, academic_year_id, term, status, scholastic, version, created_at, updated_at)
			 VALUES (gen_random_uuid(), $1, $2, $3, 'draft', $4, 1, NOW(), NOW())
			 ON CONFLICT (student_id, academic_year_id, term) WHERE deleted_at IS NULL DO UPDATE SET
				scholastic = EXCLUDED.scholastic,
				updated_at = NOW()`,
			studentID, academicYearID, term, string(scholasticJSON),
		)
		if err != nil {
			log.Error().Err(err).Str("student_id", studentID).Msg("migrate hpc entry failed")
			continue
		}
		migrated++
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]int{"migrated": migrated}})
}
