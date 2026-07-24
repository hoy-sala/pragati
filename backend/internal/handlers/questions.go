package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/pragati/backend/internal/services/question_import"
	"github.com/rs/zerolog/log"
)

type QuestionHandler struct {
	db *pgxpool.Pool
}

func NewQuestionHandler(db *pgxpool.Pool) *QuestionHandler {
	return &QuestionHandler{db: db}
}

func (h *QuestionHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	baseWhere := "WHERE school_id = $1 AND deleted_at IS NULL"
	args := []interface{}{claims.SchoolID}
	n := 2

	if subj := r.URL.Query().Get("subject_id"); subj != "" {
		baseWhere += fmt.Sprintf(" AND subject_id = $%d", n)
		args = append(args, subj)
		n++
	}
	if qtype := r.URL.Query().Get("type"); qtype != "" {
		baseWhere += fmt.Sprintf(" AND question_type = $%d", n)
		args = append(args, qtype)
		n++
	}
	if search := r.URL.Query().Get("search"); search != "" {
		baseWhere += fmt.Sprintf(" AND question_text ILIKE $%d", n)
		args = append(args, "%"+search+"%")
		n++
	}

	whereArgs := args[:]

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM questions %s", baseWhere)
	if err := h.db.QueryRow(r.Context(), countQuery, args...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("count questions failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count questions"}})
		return
	}

	dataQuery := `SELECT id, school_id, subject_id, COALESCE(teacher_id, ''), question_type, question_text,
		COALESCE(question_image, ''), options, answer, marks, difficulty, chapters, tags,
		COALESCE(explanation, ''), is_active, created_at, updated_at
		FROM questions ` + baseWhere + ` ORDER BY created_at DESC`
	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", n, n+1)
	args = append(whereArgs, limit, offset)

	rows, err := h.db.Query(r.Context(), dataQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("list questions failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch questions"}})
		return
	}
	defer rows.Close()

	questions := []models.Question{}
	for rows.Next() {
		var q models.Question
		var optsJSON, chaptersJSON, tagsJSON []byte
		if err := rows.Scan(&q.ID, &q.SchoolID, &q.SubjectID, &q.TeacherID, &q.QuestionType, &q.QuestionText,
			&q.QuestionImage, &optsJSON, &q.Answer, &q.Marks, &q.Difficulty, &chaptersJSON, &tagsJSON,
			&q.Explanation, &q.IsActive, &q.CreatedAt, &q.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scan question row failed")
			continue
		}
		json.Unmarshal(optsJSON, &q.Options)
		json.Unmarshal(chaptersJSON, &q.Chapters)
		json.Unmarshal(tagsJSON, &q.Tags)
		questions = append(questions, q)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: questions, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}
	if req.SubjectID == "" || req.QuestionText == "" || req.Answer == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "subject_id, question_text, and answer are required"}})
		return
	}
	if req.Marks <= 0 {
		req.Marks = 1
	}
	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}

	optsJSON, _ := json.Marshal(req.Options)
	chaptersJSON, _ := json.Marshal(req.Chapters)
	tagsJSON, _ := json.Marshal(req.Tags)

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO questions (id, school_id, subject_id, teacher_id, question_type, question_text,
			question_image, options, answer, marks, difficulty, chapters, tags, explanation, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,true,NOW(),NOW())`,
		id, claims.SchoolID, req.SubjectID, claims.UserID, req.QuestionType, req.QuestionText,
		req.QuestionImage, string(optsJSON), req.Answer, req.Marks, req.Difficulty,
		string(chaptersJSON), string(tagsJSON), req.Explanation,
	)
	if err != nil {
		log.Error().Err(err).Msg("create question failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create question"}})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *QuestionHandler) ImportGIFT(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	subjectID := r.URL.Query().Get("subject_id")
	if subjectID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "subject_id is required"}})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "failed to read body"}})
		return
	}

	result, questions := question_import.ParseGIFT(string(body))

	for i := range questions {
		q := questions[i]
		optsJSON, _ := json.Marshal(q.Options)
		chaptersJSON, _ := json.Marshal(q.Chapters)
		tagsJSON, _ := json.Marshal(q.Tags)

		id := uuid.New().String()
		_, err := h.db.Exec(r.Context(),
			`INSERT INTO questions (id, school_id, subject_id, teacher_id, question_type, question_text,
				options, answer, marks, difficulty, chapters, tags, is_active, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,true,NOW(),NOW())`,
			id, claims.SchoolID, subjectID, claims.UserID, q.QuestionType, q.QuestionText,
			string(optsJSON), q.Answer, q.Marks, q.Difficulty,
			string(chaptersJSON), string(tagsJSON),
		)
		if err != nil {
			result.Errors = append(result.Errors, question_import.ImportRowError{Line: q.Line, Message: err.Error()})
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: result})
}

func (h *QuestionHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	subjectID := r.URL.Query().Get("subject_id")
	if subjectID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "subject_id is required"}})
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "file is required"}})
		return
	}
	defer file.Close()

	result, questions := question_import.ParseCSV(file)

	for i := range questions {
		q := questions[i]
		optsJSON, _ := json.Marshal(q.Options)
		chaptersJSON, _ := json.Marshal(q.Chapters)
		tagsJSON, _ := json.Marshal(q.Tags)

		id := uuid.New().String()
		_, err := h.db.Exec(r.Context(),
			`INSERT INTO questions (id, school_id, subject_id, teacher_id, question_type, question_text,
				options, answer, marks, difficulty, chapters, tags, is_active, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,true,NOW(),NOW())`,
			id, claims.SchoolID, subjectID, claims.UserID, q.QuestionType, q.QuestionText,
			string(optsJSON), q.Answer, q.Marks, q.Difficulty,
			string(chaptersJSON), string(tagsJSON),
		)
		if err != nil {
			result.Errors = append(result.Errors, question_import.ImportRowError{Line: q.Line, Message: err.Error()})
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: result})
}
