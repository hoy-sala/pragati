package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/auth"
	"github.com/pragati/backend/internal/config"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	db         *pgxpool.Pool
	jwtService *auth.JWTService
	cfg        *config.Config
}

func NewAuthHandler(db *pgxpool.Pool, jwtService *auth.JWTService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, jwtService: jwtService, cfg: cfg}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	if req.Email == "" || req.Password == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "email and password are required"},
		})
		return
	}

	var user models.User
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, email, password_hash, name, role,
		        COALESCE(phone, ''), COALESCE(avatar_url, ''), is_active
		 FROM users WHERE email = $1 AND is_active = true AND deleted_at IS NULL`,
		req.Email,
	).Scan(&user.ID, &user.SchoolID, &user.Email, &user.PasswordHash, &user.Name,
		&user.Role, &user.Phone, &user.AvatarURL, &user.IsActive)
	if err != nil {
		log.Debug().Err(err).Str("email", req.Email).Msg("login failed: user not found")
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_CREDENTIALS", Message: "invalid email or password"},
		})
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_CREDENTIALS", Message: "invalid email or password"},
		})
		return
	}

	accessToken, expiresAt, err := h.jwtService.GenerateAccessToken(&user)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate token"},
		})
		return
	}

	rawToken, hashedToken, lookupHash, err := h.jwtService.GenerateRefreshToken()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate refresh token")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate refresh token"},
		})
		return
	}

	if _, err := h.db.Exec(r.Context(),
		`INSERT INTO refresh_tokens (id, user_id, token_hash, lookup_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New().String(), user.ID, hashedToken, lookupHash, time.Now().Add(h.jwtService.RefreshTokenExpiry()), time.Now(),
	); err != nil {
		log.Error().Err(err).Msg("failed to store refresh token")
	}

	if _, err := h.db.Exec(r.Context(),
		`UPDATE users SET last_login_at = NOW() WHERE id = $1`, user.ID); err != nil {
		log.Error().Err(err).Msg("failed to update last login")
	}

	user.PasswordHash = ""

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: models.LoginResponse{
			User:         &user,
			AccessToken:  accessToken,
			RefreshToken: rawToken,
			ExpiresIn:    expiresAt - time.Now().Unix(),
		},
	})
}

func (h *AuthHandler) StaffLogin(w http.ResponseWriter, r *http.Request) {
	var req models.StaffLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	if req.Mobile == "" || req.Password == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "mobile and password are required"},
		})
		return
	}

	if len(req.Mobile) != 10 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "mobile must be 10 digits"},
		})
		return
	}

	var user models.User
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, email, password_hash, name, role,
		        COALESCE(phone, ''), COALESCE(avatar_url, ''), is_active
		 FROM users WHERE mobile = $1 AND is_active = true AND deleted_at IS NULL`,
		req.Mobile,
	).Scan(&user.ID, &user.SchoolID, &user.Email, &user.PasswordHash, &user.Name,
		&user.Role, &user.Phone, &user.AvatarURL, &user.IsActive)
	if err != nil {
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_CREDENTIALS", Message: "invalid mobile or password"},
		})
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_CREDENTIALS", Message: "invalid mobile or password"},
		})
		return
	}

	accessToken, expiresAt, err := h.jwtService.GenerateAccessToken(&user)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate token"},
		})
		return
	}

	rawToken, hashedToken, lookupHash, err := h.jwtService.GenerateRefreshToken()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate refresh token")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate refresh token"},
		})
		return
	}

	if _, err := h.db.Exec(r.Context(),
		`INSERT INTO refresh_tokens (id, user_id, token_hash, lookup_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New().String(), user.ID, hashedToken, lookupHash, time.Now().Add(h.jwtService.RefreshTokenExpiry()), time.Now(),
	); err != nil {
		log.Error().Err(err).Msg("failed to store refresh token")
	}

	if _, err := h.db.Exec(r.Context(),
		`UPDATE users SET last_login_at = NOW() WHERE id = $1`, user.ID); err != nil {
		log.Error().Err(err).Msg("failed to update last login")
	}

	user.PasswordHash = ""

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: models.LoginResponse{
			User:         &user,
			AccessToken:  accessToken,
			RefreshToken: rawToken,
			ExpiresIn:    expiresAt - time.Now().Unix(),
		},
	})
}

func (h *AuthHandler) StudentLogin(w http.ResponseWriter, r *http.Request) {
	var req models.StudentLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	if req.SATSNumber == "" || req.DateOfBirth == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "sats_number and date_of_birth are required"},
		})
		return
	}

	var student models.Student
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, sats_number, first_name, COALESCE(last_name, ''), class_id,
		        COALESCE(section_id::text, ''), academic_year_id, is_active
		 FROM students
		 WHERE sats_number = $1 AND date_of_birth = $2 AND is_active = true AND deleted_at IS NULL`,
		req.SATSNumber, parsedDOB(req.DateOfBirth),
	).Scan(&student.ID, &student.SchoolID, &student.SATSNumber, &student.FirstName, &student.LastName,
		&student.ClassID, &student.SectionID, &student.AcademicYearID, &student.IsActive)
	if err != nil {
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_CREDENTIALS", Message: "invalid SATS number or date of birth"},
		})
		return
	}

	accessToken, expiresAt, err := h.jwtService.GenerateStudentToken(student.ID, student.SchoolID, student.SATSNumber)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate student access token")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate token"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: models.StudentLoginResponse{
			Student:     &student,
			AccessToken: accessToken,
			ExpiresIn:   expiresAt - time.Now().Unix(),
		},
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "refresh_token is required"},
		})
		return
	}

	lookupBytes := sha256.Sum256([]byte(req.RefreshToken))
	lookupHash := hex.EncodeToString(lookupBytes[:])

	var tokenID, userID, tokenHash, schoolID, name, email, role string
	var expiresAt time.Time
	err := h.db.QueryRow(r.Context(),
		`SELECT rt.id, rt.user_id, rt.token_hash, rt.expires_at,
		        u.school_id, u.name, u.email, u.role
		 FROM refresh_tokens rt
		 JOIN users u ON u.id = rt.user_id
		 WHERE rt.lookup_hash = $1 AND rt.revoked_at IS NULL AND rt.expires_at > NOW()`,
		lookupHash,
	).Scan(&tokenID, &userID, &tokenHash, &expiresAt, &schoolID, &name, &email, &role)
	if err != nil || !auth.CheckPassword(req.RefreshToken, tokenHash) {
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_TOKEN", Message: "invalid refresh token"},
		})
		return
	}

	if _, err := h.db.Exec(r.Context(),
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`, tokenID); err != nil {
		log.Error().Err(err).Msg("failed to revoke old refresh token")
	}

	accessToken, expiresAtUnix, err := h.jwtService.GenerateAccessToken(&models.User{
		ID: userID, SchoolID: schoolID, Name: name, Email: email, Role: role,
	})
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate token"},
		})
		return
	}

	rawNew, hashedNew, lookupHashNew, err := h.jwtService.GenerateRefreshToken()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate new refresh token")
	} else if _, err := h.db.Exec(r.Context(),
		`INSERT INTO refresh_tokens (id, user_id, token_hash, lookup_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New().String(), userID, hashedNew, lookupHashNew, time.Now().Add(h.jwtService.RefreshTokenExpiry()), time.Now(),
	); err != nil {
		log.Error().Err(err).Msg("failed to store new refresh token")
	}

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": rawNew,
			"expires_in":    expiresAtUnix - time.Now().Unix(),
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if _, err := h.db.Exec(r.Context(),
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	); err != nil {
		log.Error().Err(err).Msg("failed to revoke tokens on logout")
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		renderJSON(w, http.StatusUnauthorized, models.APIResponse{
			Error: &models.APIError{Code: "UNAUTHORIZED", Message: "not authenticated"},
		})
		return
	}

	if claims.Role == "student" {
		var s struct {
			models.Student
			DateOfBirthStr string `json:"date_of_birth"`
		}
		err := h.db.QueryRow(r.Context(),
			`SELECT id, school_id, sats_number, first_name, COALESCE(last_name, ''),
			        COALESCE(date_of_birth::text, ''), class_id, COALESCE(section_id::text, ''),
			        academic_year_id, is_active, created_at, updated_at
			 FROM students WHERE id = $1 AND deleted_at IS NULL`,
			claims.UserID,
		).Scan(&s.ID, &s.SchoolID, &s.SATSNumber, &s.FirstName, &s.LastName,
			&s.DateOfBirthStr, &s.ClassID, &s.SectionID,
			&s.AcademicYearID, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			renderJSON(w, http.StatusNotFound, models.APIResponse{
				Error: &models.APIError{Code: "NOT_FOUND", Message: "student not found"},
			})
			return
		}
		renderJSON(w, http.StatusOK, models.APIResponse{Data: s})
		return
	}

	var user models.User
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, email, name, role,
		        COALESCE(phone, ''), COALESCE(avatar_url, ''), is_active,
		        created_at, updated_at
		 FROM users WHERE id = $1 AND deleted_at IS NULL`,
		claims.UserID,
	).Scan(&user.ID, &user.SchoolID, &user.Email, &user.Name, &user.Role,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "user not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: user})
}

func parsedDOB(s string) interface{} {
	formats := []string{"2006-01-02", "02-01-2006", "02/01/2006", "01/02/2006", "2/1/2006", "2-1-2006"}
	for _, f := range formats {
		if t, err := time.Parse(f, strings.TrimSpace(s)); err == nil {
			return t
		}
	}
	return strings.TrimSpace(s)
}


