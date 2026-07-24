package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type ClassHandler struct {
	db *pgxpool.Pool
}

func NewClassHandler(db *pgxpool.Pool) *ClassHandler {
	return &ClassHandler{db: db}
}

func (h *ClassHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var total int
	err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM classes WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count classes failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count classes"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, academic_year_id, name, code, sort_order,
		        created_at, updated_at
		 FROM classes
		 WHERE school_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, name ASC
		 LIMIT $2 OFFSET $3`,
		claims.SchoolID, limit, offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("list classes failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch classes"},
		})
		return
	}
	defer rows.Close()

	classes := []models.Class{}
	for rows.Next() {
		var c models.Class
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.AcademicYearID, &c.Name, &c.Code, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan class row failed")
			continue
		}
		classes = append(classes, c)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: classes, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		Name           string `json:"name"`
		Code           string `json:"code"`
		SortOrder      int    `json:"sort_order"`
		AcademicYearID string `json:"academic_year_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}
	if req.Name == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "name is required"},
		})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO classes (id, school_id, academic_year_id, name, code, sort_order, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`,
		id, claims.SchoolID, req.AcademicYearID, req.Name, req.Code, req.SortOrder,
	)
	if err != nil {
		log.Error().Err(err).Msg("create class failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create class"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

type SubjectHandler struct {
	db *pgxpool.Pool
}

func NewSubjectHandler(db *pgxpool.Pool) *SubjectHandler {
	return &SubjectHandler{db: db}
}

func (h *SubjectHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var total int
	err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM subjects WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count subjects"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, name, code, is_language, is_core,
		        created_at, updated_at
		 FROM subjects
		 WHERE school_id = $1 AND deleted_at IS NULL
		 ORDER BY name ASC
		 LIMIT $2 OFFSET $3`,
		claims.SchoolID, limit, offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("list subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch subjects"},
		})
		return
	}
	defer rows.Close()

	subjects := []models.Subject{}
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.SchoolID, &s.Name, &s.Code, &s.IsLanguage, &s.IsCore, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan subject row failed")
			continue
		}
		subjects = append(subjects, s)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: subjects, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *SubjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		Name       string `json:"name"`
		Code       string `json:"code"`
		IsLanguage bool   `json:"is_language"`
		IsCore     bool   `json:"is_core"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}
	if req.Name == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "name is required"},
		})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO subjects (id, school_id, name, code, is_language, is_core, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`,
		id, claims.SchoolID, req.Name, req.Code, req.IsLanguage, req.IsCore,
	)
	if err != nil {
		log.Error().Err(err).Msg("create subject failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create subject"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

type AcademicYearHandler struct {
	db *pgxpool.Pool
}

func NewAcademicYearHandler(db *pgxpool.Pool) *AcademicYearHandler {
	return &AcademicYearHandler{db: db}
}

func (h *AcademicYearHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var total int
	err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM academic_years WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count academic years failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count academic years"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, name, start_date, end_date, is_current,
		        created_at, updated_at
		 FROM academic_years
		 WHERE school_id = $1 AND deleted_at IS NULL
		 ORDER BY start_date DESC
		 LIMIT $2 OFFSET $3`,
		claims.SchoolID, limit, offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("list academic years failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch academic years"},
		})
		return
	}
	defer rows.Close()

	years := []models.AcademicYear{}
	for rows.Next() {
		var y models.AcademicYear
		if err := rows.Scan(&y.ID, &y.SchoolID, &y.Name, &y.StartDate, &y.EndDate, &y.IsCurrent, &y.CreatedAt, &y.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan academic year row failed")
			continue
		}
		years = append(years, y)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: years, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *AcademicYearHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		IsCurrent bool   `json:"is_current"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}
	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "name, start_date, and end_date are required"},
		})
		return
	}

	if req.IsCurrent {
		if _, err := h.db.Exec(r.Context(),
			`UPDATE academic_years SET is_current = false WHERE school_id = $1`,
			claims.SchoolID,
		); err != nil {
			log.Error().Err(err).Msg("reset current academic year failed")
		}
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO academic_years (id, school_id, name, start_date, end_date, is_current, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`,
		id, claims.SchoolID, req.Name, req.StartDate, req.EndDate, req.IsCurrent,
	)
	if err != nil {
		log.Error().Err(err).Msg("create academic year failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create academic year"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *AcademicYearHandler) SetCurrent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update academic year"},
		})
		return
	}
	defer tx.Rollback(r.Context())

	if _, err := tx.Exec(r.Context(), `UPDATE academic_years SET is_current = false WHERE school_id = $1`, claims.SchoolID); err != nil {
		log.Error().Err(err).Msg("reset current academic year in transaction failed")
	}
	result, err := tx.Exec(r.Context(),
		`UPDATE academic_years SET is_current = true WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		tx.Rollback(r.Context())
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "academic year not found"},
		})
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		log.Error().Err(err).Msg("commit set current year failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update academic year"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}
