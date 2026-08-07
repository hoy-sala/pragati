package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

// GET /api/v1/users — list all users (admin only)
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, email, name, role, phone, is_active, created_at
		 FROM users WHERE school_id = $1 AND deleted_at IS NULL
		 ORDER BY role, name`,
		claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("list users failed")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch users"))
		return
	}
	defer rows.Close()

	type userInfo struct {
		ID        string    `json:"id"`
		Email     *string   `json:"email,omitempty"`
		Name      string    `json:"name"`
		Role      string    `json:"role"`
		Phone     *string   `json:"phone,omitempty"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
	}

	users := []userInfo{}
	for rows.Next() {
		var u userInfo
		var schoolID string
		if err := rows.Scan(&u.ID, &schoolID, &u.Email, &u.Name, &u.Role, &u.Phone, &u.IsActive, &u.CreatedAt); err != nil {
			log.Error().Err(err).Msg("scan user row failed")
			continue
		}
		users = append(users, u)
	}
	log.Info().Int("count", len(users)).Msg("UserList result")
	renderJSON(w, http.StatusOK, apiOK(users))
}

// POST /api/v1/users — create user (admin only)
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Role     string `json:"role"`
		Phone    string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, apiErr("INVALID_INPUT", "invalid request"))
		return
	}
	if req.Name == "" || req.Role == "" {
		renderJSON(w, http.StatusBadRequest, apiErr("VALIDATION_ERROR", "name and role are required"))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to hash password"))
		return
	}

	id := uuid.New().String()
	_, err = h.db.Exec(r.Context(),
		`INSERT INTO users (id, school_id, email, password_hash, name, role, phone, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,true,NOW(),NOW())`,
		id, claims.SchoolID, req.Email, string(hashed), req.Name, req.Role, req.Phone,
	)
	if err != nil {
		log.Error().Err(err).Msg("create user failed")
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to create user"))
		return
	}
	renderJSON(w, http.StatusCreated, apiOK(map[string]string{"id": id}))
}

// PATCH /api/v1/users/{id}/toggle — activate/deactivate user
func (h *UserHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	userID := chi.URLParam(r, "id")

	result, err := h.db.Exec(r.Context(),
		`UPDATE users SET is_active = NOT is_active, updated_at = NOW()
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		userID, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, apiErr("NOT_FOUND", "user not found"))
		return
	}
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}

// POST /api/v1/users/{id}/reset-password — reset user password
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	userID := chi.URLParam(r, "id")

	var req struct {
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to hash password"))
		return
	}

	result, err := h.db.Exec(r.Context(),
		`UPDATE users SET password_hash = $1, updated_at = NOW()
		 WHERE id = $2 AND school_id = $3 AND deleted_at IS NULL`,
		string(hashed), userID, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, apiErr("NOT_FOUND", "user not found"))
		return
	}
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}

// GET /api/v1/users/{id}/teacher-detail — teacher's subjects and class
func (h *UserHandler) TeacherDetail(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	rows, err := h.db.Query(r.Context(), `SELECT subject_id FROM teacher_subjects WHERE teacher_id = $1`, userID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to fetch teacher subjects"))
		return
	}
	defer rows.Close()
	type subj struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var subjects []subj
	for rows.Next() {
		var sid string
		if err := rows.Scan(&sid); err != nil { continue }
		var sname string
		h.db.QueryRow(r.Context(), `SELECT name FROM subjects WHERE id = $1 AND deleted_at IS NULL`, sid).Scan(&sname)
		subjects = append(subjects, subj{ID: sid, Name: sname})
	}

	var classID string
	h.db.QueryRow(r.Context(), `SELECT class_id FROM teacher_classes WHERE teacher_id = $1 LIMIT 1`, userID).Scan(&classID)

	renderJSON(w, http.StatusOK, apiOK(map[string]interface{}{"subjects": subjects, "class_id": classID}))
}

// PUT /api/v1/users/{id}/teacher-detail — update teacher's subjects and class
func (h *UserHandler) UpdateTeacherDetail(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	var req struct {
		SubjectIDs []string `json:"subject_ids"`
		ClassID    string   `json:"class_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, apiErr("INVALID_INPUT", "invalid request"))
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to start transaction"))
		return
	}
	defer tx.Rollback(r.Context())

	tx.Exec(r.Context(), `DELETE FROM teacher_subjects WHERE teacher_id = $1`, userID)
	for _, sid := range req.SubjectIDs {
		tx.Exec(r.Context(), `INSERT INTO teacher_subjects (teacher_id, subject_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, sid)
	}

	tx.Exec(r.Context(), `DELETE FROM teacher_classes WHERE teacher_id = $1`, userID)
	if req.ClassID != "" {
		tx.Exec(r.Context(), `INSERT INTO teacher_classes (teacher_id, class_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, req.ClassID)
	}

	if err := tx.Commit(r.Context()); err != nil {
		renderJSON(w, http.StatusInternalServerError, apiErr("INTERNAL_ERROR", "failed to update teacher detail"))
		return
	}
	renderJSON(w, http.StatusOK, apiOK(map[string]bool{"success": true}))
}
