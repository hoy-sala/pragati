package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type CertificateHandler struct {
	db        *pgxpool.Pool
	uploadDir string
}

func NewCertificateHandler(db *pgxpool.Pool, uploadDir string) *CertificateHandler {
	return &CertificateHandler{db: db, uploadDir: uploadDir}
}

func parseNullableDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

// POST /certificates/events
func (h *CertificateHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.CertificateCreateInput
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
	if req.Category == "" {
		req.Category = "other"
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO certificate_events (id, school_id, academic_year_id, name, category, held_date, venue, description, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())`,
		id, claims.SchoolID, nullIfEmpty(req.AcademicYearID), req.Name, req.Category,
		parseNullableDate(req.HeldDate), req.Venue, req.Description,
	)
	if err != nil {
		log.Error().Err(err).Msg("create certificate event failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create event"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

// GET /certificates/events
func (h *CertificateHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
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
		`SELECT COUNT(*) FROM certificate_events WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count certificate events failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count events"},
		})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, school_id, academic_year_id, name, category, held_date, venue, description,
		        created_at, updated_at
		 FROM certificate_events
		 WHERE school_id = $1 AND deleted_at IS NULL
		 ORDER BY held_date DESC NULLS LAST, created_at DESC
		 LIMIT $2 OFFSET $3`,
		claims.SchoolID, limit, offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("list certificate events failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch events"},
		})
		return
	}
	defer rows.Close()

	events := []models.CertificateEvent{}
	for rows.Next() {
		var e models.CertificateEvent
		if err := rows.Scan(&e.ID, &e.SchoolID, &e.AcademicYearID, &e.Name, &e.Category, &e.HeldDate, &e.Venue, &e.Description, &e.CreatedAt, &e.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan certificate event row failed")
			continue
		}
		events = append(events, e)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: events, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

// GET /certificates/events/{id}
func (h *CertificateHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	eventID := chi.URLParam(r, "id")

	var event models.CertificateEvent
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, academic_year_id, name, category, held_date, venue, description,
		        created_at, updated_at
		 FROM certificate_events
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		eventID, claims.SchoolID,
	).Scan(&event.ID, &event.SchoolID, &event.AcademicYearID, &event.Name, &event.Category, &event.HeldDate, &event.Venue, &event.Description, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "event not found"},
		})
		return
	}

	participants, err := h.listParticipants(r, eventID, claims.SchoolID)
	if err != nil {
		log.Error().Err(err).Msg("list certificate participants failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch participants"},
		})
		return
	}

	signatories, err := h.listSignatories(r, eventID)
	if err != nil {
		log.Error().Err(err).Msg("list certificate signatories failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch signatories"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"event":        event,
		"participants": participants,
		"signatories":  signatories,
	}})
}

func (h *CertificateHandler) listParticipants(r *http.Request, eventID, schoolID string) ([]models.CertificateParticipant, error) {
	rows, err := h.db.Query(r.Context(),
		`SELECT c.id, c.student_id, s.first_name, s.last_name, s.sats_number, cl.name, c.position, c.prize_title, c.issue_date
		 FROM certificates c
		 JOIN students s ON s.id = c.student_id
		 LEFT JOIN classes cl ON cl.id = s.class_id
		 WHERE c.event_id = $1 AND c.school_id = $2 AND c.deleted_at IS NULL
		 ORDER BY c.created_at ASC`,
		eventID, schoolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	participants := []models.CertificateParticipant{}
	for rows.Next() {
		var p models.CertificateParticipant
		var firstName, lastName string
		if err := rows.Scan(&p.ID, &p.StudentID, &firstName, &lastName, &p.SATSNumber, &p.ClassName, &p.Position, &p.PrizeTitle, &p.IssueDate); err != nil {
			log.Error().Err(err).Msg("scan certificate participant row failed")
			continue
		}
		p.StudentName = strings.TrimSpace(firstName + " " + lastName)
		participants = append(participants, p)
	}
	return participants, nil
}

func (h *CertificateHandler) listSignatories(r *http.Request, eventID string) ([]models.CertificateSignatory, error) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, event_id, name, role, title, signature_url, sort_order, created_at, updated_at
		 FROM certificate_signatories
		 WHERE event_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, created_at ASC`,
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	signatories := []models.CertificateSignatory{}
	for rows.Next() {
		var s models.CertificateSignatory
		if err := rows.Scan(&s.ID, &s.EventID, &s.Name, &s.Role, &s.Title, &s.SignatureURL, &s.SortOrder, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan certificate signatory row failed")
			continue
		}
		signatories = append(signatories, s)
	}
	return signatories, nil
}

// POST /certificates/events/{id}/participants
func (h *CertificateHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	eventID := chi.URLParam(r, "id")

	var req models.CertificateParticipantInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}
	if req.StudentID == "" || req.Position == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "student_id and position are required"},
		})
		return
	}

	var exists bool
	h.db.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM certificate_events WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL)`,
		eventID, claims.SchoolID,
	).Scan(&exists)
	if !exists {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "event not found"},
		})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO certificates (id, school_id, event_id, student_id, position, prize_title, issue_date, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
		 ON CONFLICT (event_id, student_id) DO UPDATE SET position = EXCLUDED.position, prize_title = EXCLUDED.prize_title, issue_date = EXCLUDED.issue_date, updated_at = NOW()`,
		id, claims.SchoolID, eventID, req.StudentID, req.Position, req.PrizeTitle, parseNullableDate(req.IssueDate),
	)
	if err != nil {
		log.Error().Err(err).Msg("add certificate participant failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to add participant"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

// PUT /certificates/{id}
func (h *CertificateHandler) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	var req models.CertificateParticipantInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	result, err := h.db.Exec(r.Context(),
		`UPDATE certificates SET position = $1, prize_title = $2, issue_date = $3, updated_at = NOW()
		 WHERE id = $4 AND school_id = $5 AND deleted_at IS NULL`,
		req.Position, req.PrizeTitle, parseNullableDate(req.IssueDate), id, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "certificate not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

// DELETE /certificates/{id}
func (h *CertificateHandler) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	result, err := h.db.Exec(r.Context(),
		`UPDATE certificates SET deleted_at = NOW(), updated_at = NOW()
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "certificate not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

// POST /certificates/events/{id}/signatories
func (h *CertificateHandler) AddSignatory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	eventID := chi.URLParam(r, "id")

	var req models.CertificateSignatoryInput
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
	if req.Role == "" {
		req.Role = "judge"
	}

	var exists bool
	h.db.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM certificate_events WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL)`,
		eventID, claims.SchoolID,
	).Scan(&exists)
	if !exists {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "event not found"},
		})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO certificate_signatories (id, event_id, name, role, title, signature_url, sort_order, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`,
		id, eventID, req.Name, req.Role, req.Title, req.SignatureURL, req.SortOrder,
	)
	if err != nil {
		log.Error().Err(err).Msg("add certificate signatory failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to add signatory"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

// DELETE /certificates/signatories/{id}
func (h *CertificateHandler) DeleteSignatory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.db.Exec(r.Context(),
		`UPDATE certificate_signatories SET deleted_at = NOW(), updated_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "signatory not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

// POST /certificates/signatures — multipart signature image upload
func (h *CertificateHandler) UploadSignature(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "file is required"},
		})
		return
	}
	defer file.Close()

	dir := filepath.Join(h.uploadDir, "signatures")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Error().Err(err).Msg("create signature upload dir failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save signature"},
		})
		return
	}

	filename := uuid.New().String() + ".png"
	path := filepath.Join(dir, filename)
	out, err := os.Create(path)
	if err != nil {
		log.Error().Err(err).Msg("create signature file failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save signature"},
		})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		log.Error().Err(err).Msg("write signature file failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save signature"},
		})
		return
	}

	url := fmt.Sprintf("/api/v1/uploads/signatures/%s", filename)
	log.Info().Str("school", claims.SchoolID).Str("file", filename).Msg("signature uploaded")
	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"url": url}})
}

// GET /certificates/{id} — full print payload
func (h *CertificateHandler) GetCertificate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	id := chi.URLParam(r, "id")

	var cert models.Certificate
	var firstName, lastName, satsNumber, className string
	err := h.db.QueryRow(r.Context(),
		`SELECT c.id, c.school_id, c.event_id, c.student_id, c.position, c.prize_title, c.issue_date,
		        c.created_at, c.updated_at, s.first_name, s.last_name, s.sats_number, COALESCE(cl.name, '')
		 FROM certificates c
		 JOIN students s ON s.id = c.student_id
		 LEFT JOIN classes cl ON cl.id = s.class_id
		 WHERE c.id = $1 AND c.school_id = $2 AND c.deleted_at IS NULL`,
		id, claims.SchoolID,
	).Scan(&cert.ID, &cert.SchoolID, &cert.EventID, &cert.StudentID, &cert.Position, &cert.PrizeTitle, &cert.IssueDate,
		&cert.CreatedAt, &cert.UpdatedAt, &firstName, &lastName, &satsNumber, &className)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "certificate not found"},
		})
		return
	}

	detail := models.CertificateDetail{
		Certificate: cert,
		StudentName: strings.TrimSpace(firstName + " " + lastName),
		SATSNumber:  satsNumber,
		ClassName:   className,
	}

	// fetch event
	var event models.CertificateEvent
	err = h.db.QueryRow(r.Context(),
		`SELECT id, school_id, academic_year_id, name, category, held_date, venue, description, created_at, updated_at
		 FROM certificate_events WHERE id = $1 AND deleted_at IS NULL`,
		cert.EventID,
	).Scan(&event.ID, &event.SchoolID, &event.AcademicYearID, &event.Name, &event.Category, &event.HeldDate, &event.Venue, &event.Description, &event.CreatedAt, &event.UpdatedAt)
	if err == nil {
		detail.Event = &event
	}

	signatories, err := h.listSignatories(r, cert.EventID)
	if err == nil {
		detail.Signatories = signatories
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: detail})
}