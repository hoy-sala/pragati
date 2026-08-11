package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/auth"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type TeacherHandler struct {
	db *pgxpool.Pool
}

func NewTeacherHandler(db *pgxpool.Pool) *TeacherHandler {
	return &TeacherHandler{db: db}
}

func (h *TeacherHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	var total int
	err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM users WHERE school_id = $1 AND role = 'teacher' AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count teachers failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count teachers"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, email, name, role, COALESCE(phone, ''), is_active,
		        created_at, updated_at
		 FROM users
		 WHERE school_id = $1 AND role = 'teacher' AND deleted_at IS NULL
		 ORDER BY name ASC
		 LIMIT $2 OFFSET $3`,
		claims.SchoolID, limit, offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("list teachers failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch teachers"},
		})
		return
	}
	defer rows.Close()

	type TeacherItem struct {
		ID        string    `json:"id"`
		SchoolID  string    `json:"school_id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		Role      string    `json:"role"`
		Phone     string    `json:"phone"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	teachers := []TeacherItem{}
	for rows.Next() {
		var t TeacherItem
		if err := rows.Scan(&t.ID, &t.SchoolID, &t.Email, &t.Name, &t.Role, &t.Phone, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan teacher row failed")
			continue
		}
		teachers = append(teachers, t)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: teachers, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *TeacherHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	type TeacherDetail struct {
		ID        string    `json:"id"`
		SchoolID  string    `json:"school_id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		Role      string    `json:"role"`
		Phone     string    `json:"phone"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var t TeacherDetail
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, email, name, role, COALESCE(phone, ''), is_active,
		        created_at, updated_at
		 FROM users WHERE id = $1 AND school_id = $2 AND role = 'teacher' AND deleted_at IS NULL`,
		id, claims.SchoolID,
	).Scan(&t.ID, &t.SchoolID, &t.Email, &t.Name, &t.Role, &t.Phone, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "teacher not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: t})
}

func (h *TeacherHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "name, email, and password are required"},
		})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("hash teacher password failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create teacher"},
		})
		return
	}

	id := uuid.New().String()
	_, err = h.db.Exec(r.Context(),
		`INSERT INTO users (id, school_id, email, password_hash, name, role, phone, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,'teacher',$6,true,NOW(),NOW())`,
		id, claims.SchoolID, req.Email, hash, req.Name, req.Phone,
	)
	if err != nil {
		log.Error().Err(err).Msg("create teacher failed")
		renderJSON(w, http.StatusConflict, models.APIResponse{
			Error: &models.APIError{Code: "CONFLICT", Message: "email already exists"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *TeacherHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	result, err := h.db.Exec(r.Context(),
		`UPDATE users SET name = COALESCE(NULLIF($1,''), name),
		                  phone = COALESCE(NULLIF($2,''), phone),
		                  updated_at = NOW()
		 WHERE id = $3 AND school_id = $4 AND role = 'teacher' AND deleted_at IS NULL`,
		req.Name, req.Phone, id, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("update teacher failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update teacher"},
		})
		return
	}
	if result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "teacher not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

func (h *TeacherHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	result, err := h.db.Exec(r.Context(),
		`UPDATE users SET deleted_at = NOW(), is_active = false, updated_at = NOW()
		 WHERE id = $1 AND school_id = $2 AND role = 'teacher' AND deleted_at IS NULL`,
		id, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("delete teacher failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to delete teacher"},
		})
		return
	}
	if result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "teacher not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

func (h *TeacherHandler) ListSubjects(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	rows, err := h.db.Query(r.Context(),
		`SELECT ts.teacher_id, ts.subject_id, ts.class_id, s.name as subject_name, c.name as class_name
		 FROM teacher_subjects ts
		 JOIN subjects s ON s.id = ts.subject_id
		 JOIN classes c ON c.id = ts.class_id
		 WHERE ts.teacher_id = $1 AND c.school_id = $2
		 ORDER BY c.sort_order, s.name`,
		id, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("list teacher subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch teacher subjects"},
		})
		return
	}
	defer rows.Close()

	type SubjectAssignment struct {
		TeacherID   string `json:"teacher_id"`
		SubjectID   string `json:"subject_id"`
		ClassID     string `json:"class_id"`
		SubjectName string `json:"subject_name"`
		ClassName   string `json:"class_name"`
	}

	assignments := []SubjectAssignment{}
	for rows.Next() {
		var a SubjectAssignment
		if err := rows.Scan(&a.TeacherID, &a.SubjectID, &a.ClassID, &a.SubjectName, &a.ClassName); err != nil {
			log.Error().Err(err).Msg("scan teacher subject row failed")
			continue
		}
		assignments = append(assignments, a)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: assignments})
}

func (h *TeacherHandler) SetSubjects(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Subjects []struct {
			SubjectID string `json:"subject_id"`
			ClassID   string `json:"class_id"`
		} `json:"subjects"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("begin tx for teacher subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update teacher subjects"},
		})
		return
	}
	defer tx.Rollback(r.Context())

	if _, err := tx.Exec(r.Context(),
		`DELETE FROM teacher_subjects WHERE teacher_id = $1`, id); err != nil {
		log.Error().Err(err).Msg("delete teacher subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update teacher subjects"},
		})
		return
	}

	for _, s := range req.Subjects {
		if _, err := tx.Exec(r.Context(),
			`INSERT INTO teacher_subjects (teacher_id, subject_id, class_id)
			 VALUES ($1,$2,$3)`,
			id, s.SubjectID, s.ClassID,
		); err != nil {
			log.Error().Err(err).Msg("insert teacher subject failed")
			renderJSON(w, http.StatusInternalServerError, models.APIResponse{
				Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update teacher subjects"},
			})
			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Error().Err(err).Msg("commit teacher subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update teacher subjects"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}
