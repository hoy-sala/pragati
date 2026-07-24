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

type CategoryHandler struct {
	db *pgxpool.Pool
}

func NewCategoryHandler(db *pgxpool.Pool) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var total int
	err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM assessment_categories WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count categories failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count categories"}})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, name, code, weightage, sort_order, is_active, created_at, updated_at
		 FROM assessment_categories
		 WHERE school_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC
		 LIMIT $2 OFFSET $3`,
		claims.SchoolID, limit, offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("list categories failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch categories"}})
		return
	}
	defer rows.Close()

	categories := []models.AssessmentCategory{}
	for rows.Next() {
		var c models.AssessmentCategory
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &c.Code, &c.Weightage, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan category row failed")
			continue
		}
		categories = append(categories, c)
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: categories, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		Name      string `json:"name"`
		Code      string `json:"code"`
		Weightage int    `json:"weightage"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}
	if req.Name == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "name is required"}})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO assessment_categories (id, school_id, name, code, weightage, sort_order, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,true,NOW(),NOW())`,
		id, claims.SchoolID, req.Name, req.Code, req.Weightage, req.SortOrder,
	)
	if err != nil {
		log.Error().Err(err).Msg("create category failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create category"}})
		return
	}
	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		Name      *string `json:"name"`
		Code      *string `json:"code"`
		Weightage *int    `json:"weightage"`
		SortOrder *int    `json:"sort_order"`
		IsActive  *bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}

	fields := []string{}
	args := []interface{}{}
	n := 1

	add := func(f string, v interface{}) {
		n++
		fields = append(fields, fmt.Sprintf("%s = $%d", f, n))
		args = append(args, v)
	}
	if req.Name != nil { add("name", *req.Name) }
	if req.Code != nil { add("code", *req.Code) }
	if req.Weightage != nil { add("weightage", *req.Weightage) }
	if req.SortOrder != nil { add("sort_order", *req.SortOrder) }
	if req.IsActive != nil { add("is_active", *req.IsActive) }

	if len(fields) == 0 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "no fields to update"}})
		return
	}

	fields = append(fields, "updated_at = NOW()")
	setSQL := ""
	for i, f := range fields {
		if i > 0 { setSQL += ", " }
		setSQL += f
	}

	args = append([]interface{}{id, claims.SchoolID}, args...)
	sql := fmt.Sprintf("UPDATE assessment_categories SET %s WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL", setSQL)

	result, err := h.db.Exec(r.Context(), sql, args...)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "category not found"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": id}})
}
