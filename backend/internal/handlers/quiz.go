package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type QuizHandler struct {
	db *pgxpool.Pool
}

func NewQuizHandler(db *pgxpool.Pool) *QuizHandler {
	return &QuizHandler{db: db}
}

// POST /api/v1/quizzes
func (h *QuizHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.QuizCreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"}})
		return
	}
	if req.Title == "" || req.TargetType == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "title and target_type are required"}})
		return
	}
	if req.TargetType != "student" && req.TargetType != "staff" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "target_type must be 'student' or 'staff'"}})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO quiz_assignments (id, school_id, title, description, target_type, target_id,
			pass_pct, max_attempts, duration_min, shuffle_questions, shuffle_options, show_result,
			start_at, end_at, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW(),NOW())`,
		id, claims.SchoolID, req.Title, req.Description, req.TargetType, nullIfEmpty(req.TargetID),
		req.PassPct, req.MaxAttempts, req.DurationMin, req.ShuffleQuestions, req.ShuffleOptions, req.ShowResult,
		parseNullableTime(req.StartAt), parseNullableTime(req.EndAt), claims.UserID,
	)
	if err != nil {
		log.Error().Err(err).Msg("create quiz failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create quiz"}})
		return
	}
	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

// GET /api/v1/quizzes?target_type=&page=&limit=
func (h *QuizHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query := `FROM quiz_assignments qa
		JOIN users u ON u.id = qa.created_by
		WHERE qa.school_id = $1 AND qa.deleted_at IS NULL`
	args := []interface{}{claims.SchoolID}
	n := 2

	if tt := r.URL.Query().Get("target_type"); tt != "" {
		query += fmt.Sprintf(" AND qa.target_type = $%d", n)
		args = append(args, tt)
		n++
	}

	var total int
	h.db.QueryRow(r.Context(), "SELECT COUNT(*) "+query, args...).Scan(&total)

	dataQuery := `SELECT qa.id, qa.school_id, qa.title, qa.description, qa.target_type,
		COALESCE(qa.target_id, ''), qa.pass_pct, qa.max_attempts, qa.duration_min,
		qa.shuffle_questions, qa.shuffle_options, qa.show_result,
		qa.start_at, qa.end_at, qa.is_published, qa.created_by, qa.is_active,
		qa.created_at, qa.updated_at,
		(SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = qa.id),
		(SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = qa.id),
		u.name ` + query + ` ORDER BY qa.created_at DESC LIMIT $` + fmt.Sprintf("%d", n) + ` OFFSET $` + fmt.Sprintf("%d", n+1)

	args = append(args, limit, offset)
	rows, err := h.db.Query(r.Context(), dataQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("list quizzes failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch quizzes"}})
		return
	}
	defer rows.Close()

	quizzes := []models.QuizListItem{}
	for rows.Next() {
		var q models.QuizListItem
		if err := rows.Scan(&q.ID, &q.SchoolID, &q.Title, &q.Description, &q.TargetType,
			&q.TargetID, &q.PassPct, &q.MaxAttempts, &q.DurationMin,
			&q.ShuffleQuestions, &q.ShuffleOptions, &q.ShowResult,
			&q.StartAt, &q.EndAt, &q.IsPublished, &q.CreatedBy, &q.IsActive,
			&q.CreatedAt, &q.UpdatedAt,
			&q.QuestionCount, &q.AttemptCount, &q.CreatedByName); err != nil {
			log.Error().Err(err).Msg("scan quiz row failed")
			continue
		}
		quizzes = append(quizzes, q)
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: quizzes, Meta: models.Pagination{Offset: offset, Limit: limit, Total: total}})
}

// GET /api/v1/quizzes/{id}
func (h *QuizHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	var q models.QuizAssignment
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, title, description, target_type, COALESCE(target_id, ''),
			pass_pct, max_attempts, duration_min, shuffle_questions, shuffle_options, show_result,
			start_at, end_at, is_published, created_by, is_active, created_at, updated_at
		 FROM quiz_assignments WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID,
	).Scan(&q.ID, &q.SchoolID, &q.Title, &q.Description, &q.TargetType, &q.TargetID,
		&q.PassPct, &q.MaxAttempts, &q.DurationMin, &q.ShuffleQuestions, &q.ShuffleOptions, &q.ShowResult,
		&q.StartAt, &q.EndAt, &q.IsPublished, &q.CreatedBy, &q.IsActive, &q.CreatedAt, &q.UpdatedAt)
	if err == pgx.ErrNoRows {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "quiz not found"}})
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("get quiz failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch quiz"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: q})
}

// PUT /api/v1/quizzes/{id}
func (h *QuizHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())
	var req models.QuizCreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}
	result, err := h.db.Exec(r.Context(),
		`UPDATE quiz_assignments SET title=$1, description=$2, target_type=$3, target_id=$4,
			pass_pct=$5, max_attempts=$6, duration_min=$7, shuffle_questions=$8, shuffle_options=$9,
			show_result=$10, start_at=$11, end_at=$12, updated_at=NOW()
		 WHERE id=$13 AND school_id=$14 AND deleted_at IS NULL AND NOT is_published`,
		req.Title, req.Description, req.TargetType, nullIfEmpty(req.TargetID),
		req.PassPct, req.MaxAttempts, req.DurationMin, req.ShuffleQuestions, req.ShuffleOptions,
		req.ShowResult, parseNullableTime(req.StartAt), parseNullableTime(req.EndAt),
		id, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("update quiz failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update quiz"}})
		return
	}
	if result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "quiz not found or already published"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": id}})
}

// DELETE /api/v1/quizzes/{id}
func (h *QuizHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())
	_, err := h.db.Exec(r.Context(),
		`UPDATE quiz_assignments SET deleted_at = NOW() WHERE id = $1 AND school_id = $2`,
		id, claims.SchoolID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to delete quiz"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

// POST /api/v1/quizzes/{id}/publish
func (h *QuizHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	var qCount int
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = $1`, id).Scan(&qCount)
	if qCount == 0 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "add at least one question before publishing"}})
		return
	}

	result, err := h.db.Exec(r.Context(),
		`UPDATE quiz_assignments SET is_published = true, updated_at = NOW()
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID)
	if err != nil || result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "quiz not found"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": id}})
}

// GET /api/v1/quizzes/{id}/questions
func (h *QuizHandler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	rows, err := h.db.Query(r.Context(),
		`SELECT qq.question_id, qq.sort_order, qq.marks,
			q.question_text, q.question_type, q.options, q.answer, q.marks, q.difficulty
		 FROM quiz_questions qq
		 JOIN questions q ON q.id = qq.question_id AND q.school_id = $1
		 WHERE qq.quiz_id = $2
		 ORDER BY qq.sort_order ASC`,
		claims.SchoolID, id)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch questions"}})
		return
	}
	defer rows.Close()

	type QuizQuestionDetail struct {
		QuestionID   string          `json:"question_id"`
		SortOrder    int             `json:"sort_order"`
		Marks        int             `json:"marks"`
		QuestionText string          `json:"question_text"`
		QuestionType string          `json:"question_type"`
		Options      json.RawMessage `json:"options"`
		Answer       string          `json:"answer"`
		Difficulty   string          `json:"difficulty"`
	}

	questions := []QuizQuestionDetail{}
	for rows.Next() {
		var q QuizQuestionDetail
		if err := rows.Scan(&q.QuestionID, &q.SortOrder, &q.Marks,
			&q.QuestionText, &q.QuestionType, &q.Options, &q.Answer, &q.Marks, &q.Difficulty); err != nil {
			continue
		}
		questions = append(questions, q)
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: questions})
}

// POST /api/v1/quizzes/{id}/questions
func (h *QuizHandler) AddQuestions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())
	var req struct {
		QuestionIDs []string `json:"question_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}

	added := 0
	for i, qid := range req.QuestionIDs {
		_, err := h.db.Exec(r.Context(),
			`INSERT INTO quiz_questions (quiz_id, question_id, sort_order, marks)
			 VALUES ($1,$2,$3, (SELECT marks FROM questions WHERE id = $2))
			 ON CONFLICT DO NOTHING`,
			id, qid, i)
		if err == nil {
			added++
		}
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]int{"added": added}})
}

// DELETE /api/v1/quizzes/{id}/questions/{qid}
func (h *QuizHandler) RemoveQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	qid := chi.URLParam(r, "qid")
	_, err := h.db.Exec(r.Context(), `DELETE FROM quiz_questions WHERE quiz_id = $1 AND question_id = $2`, id, qid)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to remove question"}})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

// GET /api/v1/quizzes/available
func (h *QuizHandler) GetAvailable(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	rows, err := h.db.Query(r.Context(),
		`SELECT qa.id, qa.title, qa.description, qa.duration_min, qa.pass_pct, qa.max_attempts,
			qa.start_at, qa.end_at,
			(SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = qa.id),
			COALESCE((SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = qa.id AND user_id = $1), 0),
			COALESCE((SELECT status FROM quiz_attempts WHERE quiz_id = qa.id AND user_id = $1 ORDER BY attempt_no DESC LIMIT 1), ''),
			COALESCE((SELECT score FROM quiz_attempts WHERE quiz_id = qa.id AND user_id = $1 ORDER BY attempt_no DESC LIMIT 1), 0),
			COALESCE((SELECT passed FROM quiz_attempts WHERE quiz_id = qa.id AND user_id = $1 ORDER BY attempt_no DESC LIMIT 1), false)
		 FROM quiz_assignments qa
		 WHERE qa.school_id = $2 AND qa.is_published = true AND qa.deleted_at IS NULL
		 AND (qa.start_at IS NULL OR qa.start_at <= NOW())
		 AND (qa.end_at IS NULL OR qa.end_at >= NOW())
		 AND (qa.target_type = 'staff' OR (qa.target_type = 'student' AND qa.target_id IN (
			SELECT class_id FROM students WHERE user_id = $1
		 )))
		 ORDER BY qa.end_at ASC NULLS LAST, qa.created_at DESC`,
		claims.UserID, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("fetch available quizzes failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch quizzes"}})
		return
	}
	defer rows.Close()

	quizzes := []models.AvailableQuizItem{}
	for rows.Next() {
		var q models.AvailableQuizItem
		var lastStatus string
		var lastScore float64
		var lastPassed bool
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.DurationMin, &q.PassPct, &q.MaxAttempts,
			&q.StartAt, &q.EndAt, &q.QuestionCount, &q.AttemptsUsed, &lastStatus, &lastScore, &lastPassed); err != nil {
			continue
		}
		if lastStatus != "" {
			q.LastStatus = lastStatus
			q.LastScore = &lastScore
			q.LastPassed = &lastPassed
		}
		quizzes = append(quizzes, q)
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: quizzes})
}

// POST /api/v1/quizzes/{id}/start
func (h *QuizHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	var maxAttempts int
	h.db.QueryRow(r.Context(),
		`SELECT max_attempts FROM quiz_assignments WHERE id = $1 AND is_published = true AND deleted_at IS NULL`,
		id).Scan(&maxAttempts)
	if maxAttempts == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "quiz not found or not published"}})
		return
	}

	var used int
	h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = $1 AND user_id = $2`, id, claims.UserID).Scan(&used)
	if used >= maxAttempts {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "MAX_ATTEMPTS", Message: "maximum attempts reached"}})
		return
	}

	var existingID string
	err := h.db.QueryRow(r.Context(),
		`SELECT id FROM quiz_attempts WHERE quiz_id = $1 AND user_id = $2 AND status = 'in_progress' ORDER BY started_at DESC LIMIT 1`,
		id, claims.UserID).Scan(&existingID)
	if err == nil {
		renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
			"id":     existingID,
			"status": "in_progress",
		}})
		return
	}

	attemptID := uuid.New().String()
	_, err = h.db.Exec(r.Context(),
		`INSERT INTO quiz_attempts (id, quiz_id, user_id, attempt_no, status, started_at)
		 VALUES ($1, $2, $3, $4, 'in_progress', NOW())`,
		attemptID, id, claims.UserID, used+1)
	if err != nil {
		log.Error().Err(err).Msg("start quiz attempt failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to start attempt"}})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]interface{}{
		"id":     attemptID,
		"status": "in_progress",
	}})
}

// GET /api/v1/quizzes/attempts/{attemptId}
func (h *QuizHandler) GetAttempt(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "attemptId")
	claims := middleware.GetUserClaims(r.Context())

	var attempt models.QuizAttempt
	err := h.db.QueryRow(r.Context(),
		`SELECT id, quiz_id, user_id, attempt_no, status, score, percentage, passed,
			started_at, submitted_at
		 FROM quiz_attempts WHERE id = $1 AND user_id = $2`,
		attemptID, claims.UserID,
	).Scan(&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.AttemptNo, &attempt.Status,
		&attempt.Score, &attempt.Percentage, &attempt.Passed, &attempt.StartedAt, &attempt.SubmittedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "attempt not found"}})
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT qq.question_id, qq.sort_order, qq.marks,
			q.question_text, q.question_type, q.options
		 FROM quiz_questions qq
		 JOIN questions q ON q.id = qq.question_id
		 WHERE qq.quiz_id = $1
		 ORDER BY qq.sort_order ASC`,
		attempt.QuizID,
	)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch questions"}})
		return
	}
	defer rows.Close()

	type AttemptQuestion struct {
		QuestionID   string          `json:"question_id"`
		SortOrder    int             `json:"sort_order"`
		Marks        int             `json:"marks"`
		QuestionText string          `json:"question_text"`
		QuestionType string          `json:"question_type"`
		Options      json.RawMessage `json:"options,omitempty"`
	}

	questions := []AttemptQuestion{}
	for rows.Next() {
		var q AttemptQuestion
		if err := rows.Scan(&q.QuestionID, &q.SortOrder, &q.Marks, &q.QuestionText, &q.QuestionType, &q.Options); err != nil {
			continue
		}
		questions = append(questions, q)
	}

	// also fetch any existing responses for this attempt
	respRows, err := h.db.Query(r.Context(),
		`SELECT question_id, selected_options, text_answer
		 FROM quiz_responses WHERE attempt_id = $1`, attemptID)
	if err == nil {
		defer respRows.Close()
		type savedAnswer struct {
			QuestionID      string   `json:"question_id"`
			SelectedOptions []string `json:"selected_options"`
			TextAnswer      string   `json:"text_answer"`
		}
		saved := []savedAnswer{}
		for respRows.Next() {
			var s savedAnswer
			if err := respRows.Scan(&s.QuestionID, &s.SelectedOptions, &s.TextAnswer); err == nil {
				saved = append(saved, s)
			}
		}
		renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
			"attempt":   attempt,
			"questions": questions,
			"saved":     saved,
		}})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"attempt":   attempt,
		"questions": questions,
		"saved":     []interface{}{},
	}})
}

// PUT /api/v1/quizzes/attempts/{attemptId}/answers
func (h *QuizHandler) SaveAnswer(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "attemptId")
	claims := middleware.GetUserClaims(r.Context())

	var req struct {
		Responses []models.QuizAnswerInput `json:"responses"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}

	// verify attempt belongs to user and is in progress
	var status string
	h.db.QueryRow(r.Context(),
		`SELECT status FROM quiz_attempts WHERE id = $1 AND user_id = $2`, attemptID, claims.UserID,
	).Scan(&status)
	if status != "in_progress" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "NOT_IN_PROGRESS", Message: "attempt is not in progress"}})
		return
	}

	for _, ans := range req.Responses {
		_, err := h.db.Exec(r.Context(),
			`INSERT INTO quiz_responses (id, attempt_id, question_id, selected_options, text_answer)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4)
			 ON CONFLICT (attempt_id, question_id) DO UPDATE SET
				selected_options = EXCLUDED.selected_options,
				text_answer = EXCLUDED.text_answer,
				updated_at = NOW()`,
			attemptID, ans.QuestionID, ans.SelectedOptions, ans.TextAnswer,
		)
		if err != nil {
			log.Error().Err(err).Msg("save answer failed")
			renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save answer"}})
			return
		}
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"saved": true}})
}

// POST /api/v1/quizzes/attempts/{attemptId}/submit
func (h *QuizHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "attemptId")
	claims := middleware.GetUserClaims(r.Context())

	var attempt models.QuizAttempt
	err := h.db.QueryRow(r.Context(),
		`SELECT id, quiz_id, user_id, status, started_at FROM quiz_attempts
		 WHERE id = $1 AND user_id = $2`, attemptID, claims.UserID,
	).Scan(&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.Status, &attempt.StartedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "attempt not found"}})
		return
	}
	if attempt.Status != "in_progress" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "ALREADY_SUBMITTED", Message: "attempt already submitted"}})
		return
	}

	h.db.Exec(r.Context(),
		`UPDATE quiz_attempts SET status = 'submitted', submitted_at = NOW(), updated_at = NOW()
		 WHERE id = $1`, attemptID)

	// auto-grade MCQ, true_false, fill_blank
	h.gradeAttempt(attemptID)

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"attempt_id": attemptID, "status": "submitted"}})
}

func (h *QuizHandler) gradeAttempt(attemptID string) {
	var quizID string
	var totalMarks float64
	h.db.QueryRow(context.Background(),
		`SELECT qa.id, COALESCE(SUM(qq.marks), 0)::numeric(5,2)
		 FROM quiz_attempts qa
		 JOIN quiz_questions qq ON qq.quiz_id = qa.quiz_id
		 WHERE qa.id = $1
		 GROUP BY qa.id`, attemptID).Scan(&quizID, &totalMarks)

	if totalMarks == 0 {
		log.Warn().Str("attempt_id", attemptID).Msg("grade attempt: no questions found")
		return
	}

	awardedMarks := 0.0

	rows, err := h.db.Query(context.Background(),
		`SELECT qr.id, qr.selected_options, qr.text_answer,
			q.question_type, q.answer, qq.marks
		 FROM quiz_responses qr
		 JOIN questions q ON q.id = qr.question_id
		 JOIN quiz_questions qq ON qq.question_id = qr.question_id AND qq.quiz_id = $2
		 WHERE qr.attempt_id = $1`,
		attemptID, quizID,
	)
	if err != nil {
		log.Error().Err(err).Msg("grade attempt: fetch responses failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var respID string
		var selectedOpts []string
		var textAnswer, qType, correctAnswer string
		var marksTotal float64

		if err := rows.Scan(&respID, &selectedOpts, &textAnswer,
			&qType, &correctAnswer, &marksTotal); err != nil {
			continue
		}

		isCorrect := false

		switch qType {
		case "mcq":
			if len(selectedOpts) == 1 && selectedOpts[0] == correctAnswer {
				isCorrect = true
			}
		case "true_false":
			if strings.EqualFold(textAnswer, correctAnswer) {
				isCorrect = true
			}
		case "fill_blank":
			if strings.EqualFold(strings.TrimSpace(textAnswer), strings.TrimSpace(correctAnswer)) {
				isCorrect = true
			}
		case "short_answer":
			if strings.EqualFold(strings.TrimSpace(textAnswer), strings.TrimSpace(correctAnswer)) {
				isCorrect = true
			}
		}

		var aw float64
		if isCorrect {
			aw = marksTotal
		}
		awardedMarks += aw

		h.db.Exec(context.Background(),
			`UPDATE quiz_responses SET is_correct = $1, marks_awarded = $2, marks_total = $3, updated_at = NOW()
			 WHERE id = $4`,
			isCorrect, aw, marksTotal, respID)
	}

	pct := (awardedMarks / totalMarks) * 100
	h.db.Exec(context.Background(),
		`UPDATE quiz_attempts SET score = $1, percentage = $2, passed = (percentage >= (SELECT pass_pct FROM quiz_assignments WHERE id = $4)),
			graded_at = NOW(), updated_at = NOW()
		 WHERE id = $3`,
		awardedMarks, pct, attemptID, quizID)
}

// GET /api/v1/quizzes/attempts/{attemptId}/result
func (h *QuizHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "attemptId")
	claims := middleware.GetUserClaims(r.Context())

	var attempt models.QuizAttempt
	err := h.db.QueryRow(r.Context(),
		`SELECT id, quiz_id, user_id, attempt_no, status, score, percentage, passed,
			started_at, submitted_at, graded_at
		 FROM quiz_attempts WHERE id = $1 AND user_id = $2`,
		attemptID, claims.UserID,
	).Scan(&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.AttemptNo, &attempt.Status,
		&attempt.Score, &attempt.Percentage, &attempt.Passed, &attempt.StartedAt, &attempt.SubmittedAt, &attempt.GradedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "attempt not found"}})
		return
	}

	var quiz models.QuizAssignment
	h.db.QueryRow(r.Context(),
		`SELECT id, title, description, pass_pct, show_result FROM quiz_assignments WHERE id = $1`,
		attempt.QuizID,
	).Scan(&quiz.ID, &quiz.Title, &quiz.Description, &quiz.PassPct, &quiz.ShowResult)

	var totalMarks float64
	h.db.QueryRow(r.Context(),
		`SELECT COALESCE(SUM(marks), 0)::numeric(5,2) FROM quiz_questions WHERE quiz_id = $1`,
		attempt.QuizID).Scan(&totalMarks)

	respRows, err := h.db.Query(r.Context(),
		`SELECT qr.id, qr.attempt_id, qr.question_id, qr.selected_options, qr.text_answer,
			qr.is_correct, qr.marks_awarded, qr.marks_total, qr.graded_at, COALESCE(qr.graded_by, ''),
			q.question_text, q.question_type, q.options, q.answer,
			COALESCE(qr.marks_total, q.marks)
		 FROM quiz_responses qr
		 JOIN questions q ON q.id = qr.question_id
		 WHERE qr.attempt_id = $1`,
		attemptID,
	)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch results"}})
		return
	}
	defer respRows.Close()

	responses := []models.QuizResponseDetail{}
	totalAwarded := 0.0

	for respRows.Next() {
		var res models.QuizResponseDetail
		var optsJSON []byte
		if err := respRows.Scan(&res.ID, &res.AttemptID, &res.QuestionID, &res.SelectedOptions, &res.TextAnswer,
			&res.IsCorrect, &res.MarksAwarded, &res.MarksTotal, &res.GradedAt, &res.GradedBy,
			&res.QuestionText, &res.QuestionType, &optsJSON, &res.CorrectAnswer, &res.MarksTotal); err != nil {
			continue
		}
		json.Unmarshal(optsJSON, &res.Options)
		if !quiz.ShowResult {
			res.CorrectAnswer = ""
		}
		totalAwarded += res.MarksAwarded
		responses = append(responses, res)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: models.QuizResultData{
		Attempt:      attempt,
		Quiz:         quiz,
		Responses:    responses,
		TotalMarks:   totalMarks,
		TotalAwarded: totalAwarded,
	}})
}

// POST /api/v1/quizzes/attempts/{attemptId}/grade
func (h *QuizHandler) GradeShortAnswer(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "attemptId")
	claims := middleware.GetUserClaims(r.Context())

	var req struct {
		Grading []struct {
			ResponseID    string  `json:"response_id"`
			MarksAwarded  float64 `json:"marks_awarded"`
			IsCorrect     bool    `json:"is_correct"`
		} `json:"grading"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request"}})
		return
	}

	for _, g := range req.Grading {
		h.db.Exec(r.Context(),
			`UPDATE quiz_responses SET marks_awarded = $1, is_correct = $2, graded_at = NOW(), graded_by = $3, updated_at = NOW()
			 WHERE id = $4`,
			g.MarksAwarded, g.IsCorrect, claims.UserID, g.ResponseID)
	}

	// recalc total using quiz-level marks
	var quizID string
	var totalMarks float64
	h.db.QueryRow(r.Context(),
		`SELECT qa.quiz_id, COALESCE(SUM(qq.marks), 0)::numeric(5,2)
		 FROM quiz_attempts qa
		 JOIN quiz_questions qq ON qq.quiz_id = qa.quiz_id
		 WHERE qa.id = $1
		 GROUP BY qa.quiz_id`, attemptID).Scan(&quizID, &totalMarks)

	var totalAwarded float64
	h.db.QueryRow(r.Context(),
		`SELECT COALESCE(SUM(marks_awarded), 0) FROM quiz_responses WHERE attempt_id = $1`,
		attemptID).Scan(&totalAwarded)
	if totalMarks > 0 {
		pct := (totalAwarded / totalMarks) * 100
		h.db.Exec(r.Context(),
			`UPDATE quiz_attempts SET score = $1, percentage = $2, passed = (percentage >= (SELECT pass_pct FROM quiz_assignments WHERE id = $4)),
				status = 'graded', graded_at = NOW(), updated_at = NOW()
			 WHERE id = $3`,
			totalAwarded, pct, attemptID, quizID)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

func parseNullableTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return nil
		}
	}
	return &t
}
