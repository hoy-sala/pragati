package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type AssessmentHandler struct {
	db *pgxpool.Pool
}

func NewAssessmentHandler(db *pgxpool.Pool) *AssessmentHandler {
	return &AssessmentHandler{db: db}
}

func (h *AssessmentHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	baseWhere := "WHERE a.school_id = $1 AND a.deleted_at IS NULL"
	args := []interface{}{claims.SchoolID}
	n := 2

	if classID := r.URL.Query().Get("class_id"); classID != "" {
		baseWhere += fmt.Sprintf(" AND a.class_id = $%d", n)
		args = append(args, classID)
		n++
	}
	if subjectID := r.URL.Query().Get("subject_id"); subjectID != "" {
		baseWhere += fmt.Sprintf(" AND a.subject_id = $%d", n)
		args = append(args, subjectID)
		n++
	}
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		baseWhere += fmt.Sprintf(" AND a.category_id = $%d", n)
		args = append(args, categoryID)
		n++
	}

	whereArgs := args[:]

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM assessments a %s", baseWhere)
	if err := h.db.QueryRow(r.Context(), countQuery, args...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("count assessments failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count assessments"}})
		return
	}

	dataQuery := `SELECT a.id, a.school_id, a.category_id, a.subject_id, a.teacher_id,
		a.class_id, a.section_id, COALESCE(a.name, ''), a.max_marks, a.weightage,
		COALESCE(a.date::text, ''), a.academic_year_id, a.is_published, a.is_locked, a.version,
		a.created_at, a.updated_at
		FROM assessments a ` + baseWhere + ` ORDER BY a.date DESC NULLS LAST, a.created_at DESC`
	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", n, n+1)
	args = append(whereArgs, limit, offset)

	rows, err := h.db.Query(r.Context(), dataQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("list assessments failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch assessments"}})
		return
	}
	defer rows.Close()

	assessments := []models.Assessment{}
	for rows.Next() {
		var a models.Assessment
		var dateStr string
		if err := rows.Scan(&a.ID, &a.SchoolID, &a.CategoryID, &a.SubjectID, &a.TeacherID,
			&a.ClassID, &a.SectionID, &a.Name, &a.MaxMarks, &a.Weightage,
			&dateStr, &a.AcademicYearID, &a.IsPublished, &a.IsLocked, &a.Version,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan assessment row failed")
			continue
		}
		a.Date = dateStr
		assessments = append(assessments, a)
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: assessments, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *AssessmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		CategoryID    string   `json:"category_id"`
		SubjectID     string   `json:"subject_id"`
		ClassID       string   `json:"class_id"`
		SectionID     string   `json:"section_id"`
		Name          string   `json:"name"`
		MaxMarks      float64  `json:"max_marks"`
		Weightage     int      `json:"weightage"`
		Date          string   `json:"date"`
		Chapters      []string `json:"chapters"`
		AcademicYearID string  `json:"academic_year_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}
	if req.CategoryID == "" || req.SubjectID == "" || req.ClassID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "category_id, subject_id, and class_id are required"}})
		return
	}
	if req.MaxMarks <= 0 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "max_marks must be positive"}})
		return
	}

	var date *time.Time
	if req.Date != "" {
		if t, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = &t
		}
	}

	chaptersJSON, _ := json.Marshal(req.Chapters)

	teacherID := claims.UserID
	if rRole := claims.Role; rRole == "admin" || rRole == "principal" {
		if t := r.URL.Query().Get("teacher_id"); t != "" {
			teacherID = t
		}
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO assessments (id, school_id, category_id, subject_id, teacher_id,
			class_id, section_id, name, max_marks, weightage, date, chapters,
			academic_year_id, is_published, is_locked, version, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,false,false,1,NOW(),NOW())`,
		id, claims.SchoolID, req.CategoryID, req.SubjectID, teacherID,
		req.ClassID, nullIfEmpty(req.SectionID), req.Name, req.MaxMarks, req.Weightage,
		date, string(chaptersJSON), req.AcademicYearID,
	)
	if err != nil {
		log.Error().Err(err).Msg("create assessment failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create assessment"}})
		return
	}
	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *AssessmentHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	result, err := h.db.Exec(r.Context(),
		`UPDATE assessments SET is_published = true, is_locked = true, updated_at = NOW()
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "assessment not found"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": id}})
}
