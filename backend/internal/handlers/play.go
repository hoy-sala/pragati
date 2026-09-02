package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type PlayHandler struct {
	db *pgxpool.Pool
}

func NewPlayHandler(db *pgxpool.Pool) *PlayHandler {
	return &PlayHandler{db: db}
}

type PlayClass struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	QuestionCount int    `json:"question_count"`
}

type PlaySubject struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	QuestionCount int    `json:"question_count"`
}

type PlayTopic struct {
	Name string `json:"name"`
}

type PlayQuestion struct {
	ID           string          `json:"id"`
	QuestionText string          `json:"question_text"`
	QuestionType string          `json:"question_type"`
	Options      []models.Option `json:"options"`
	Difficulty   string          `json:"difficulty"`
}

type PlayScoreInput struct {
	PlayerName     string `json:"player_name"`
	ClassID        string `json:"class_id"`
	SubjectID      string `json:"subject_id"`
	Topic          string `json:"topic"`
	Difficulty     string `json:"difficulty"`
	Score          int    `json:"score"`
	TotalQuestions int    `json:"total_questions"`
	CorrectCount   int    `json:"correct_count"`
	BestStreak     int    `json:"best_streak"`
	TimeTakenMs    int    `json:"time_taken_ms"`
}

type HighScore struct {
	Rank           int    `json:"rank"`
	PlayerName     string `json:"player_name"`
	Score          int    `json:"score"`
	CorrectCount   int    `json:"correct_count"`
	TotalQuestions int    `json:"total_questions"`
	BestStreak     int    `json:"best_streak"`
	TimeTakenMs    int    `json:"time_taken_ms"`
	Difficulty     string `json:"difficulty"`
	CreatedAt      string `json:"created_at"`
}

func (h *PlayHandler) ListClasses(w http.ResponseWriter, r *http.Request) {
	schoolID := r.URL.Query().Get("school_id")
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT c.id, c.name, COUNT(DISTINCT q.id)::int
		FROM classes c
		JOIN class_subjects cs ON cs.class_id = c.id
		JOIN questions q ON q.subject_id = cs.subject_id
		JOIN subjects s ON s.id = q.subject_id
		WHERE c.school_id = $1 AND c.deleted_at IS NULL
		AND q.deleted_at IS NULL AND q.is_active = true
		AND s.code != 'GK'
		GROUP BY c.id, c.name
		ORDER BY c.sort_order, c.name
	`, schoolID)
	if err != nil {
		log.Error().Err(err).Msg("play: list classes failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch classes"}})
		return
	}
	defer rows.Close()

	classes := []PlayClass{}
	for rows.Next() {
		var c PlayClass
		if err := rows.Scan(&c.ID, &c.Name, &c.QuestionCount); err != nil {
			continue
		}
		if c.QuestionCount > 0 {
			classes = append(classes, c)
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: classes})
}

func (h *PlayHandler) ListSubjects(w http.ResponseWriter, r *http.Request) {
	classID := r.URL.Query().Get("class_id")

	query := `
		SELECT s.id, s.name, COUNT(DISTINCT q.id)::int
		FROM subjects s
		JOIN questions q ON q.subject_id = s.id
		WHERE q.deleted_at IS NULL AND q.is_active = true
	`
	args := []interface{}{}
	if classID != "" {
		query += ` AND q.id IN (SELECT q2.id FROM questions q2 JOIN class_subjects cs2 ON cs2.subject_id = q2.subject_id WHERE cs2.class_id = $1)`
		args = append(args, classID)
	}
	query += ` GROUP BY s.id, s.name ORDER BY s.name`

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("play: list subjects failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch subjects"}})
		return
	}
	defer rows.Close()

	subjects := []PlaySubject{}
	for rows.Next() {
		var s PlaySubject
		if err := rows.Scan(&s.ID, &s.Name, &s.QuestionCount); err != nil {
			continue
		}
		if s.QuestionCount > 0 {
			subjects = append(subjects, s)
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: subjects})
}

func (h *PlayHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	classID := r.URL.Query().Get("class_id")
	subjectID := r.URL.Query().Get("subject_id")
	if subjectID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "subject_id is required"}})
		return
	}

	query := `
		SELECT DISTINCT jsonb_array_elements_text(q.chapters) AS topic
		FROM questions q
		WHERE q.subject_id = $1
		AND q.deleted_at IS NULL AND q.is_active = true
		AND q.chapters != '[]'::jsonb AND q.chapters IS NOT NULL
	`
	args := []interface{}{subjectID}
	if classID != "" {
		query += ` AND q.id IN (SELECT q2.id FROM questions q2 JOIN class_subjects cs2 ON cs2.subject_id = q2.subject_id WHERE cs2.class_id = $2)`
		args = append(args, classID)
	}
	query += ` ORDER BY topic`

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("play: list topics failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch topics"}})
		return
	}
	defer rows.Close()

	topics := []PlayTopic{}
	seen := map[string]bool{}
	for rows.Next() {
		var t PlayTopic
		if err := rows.Scan(&t.Name); err != nil {
			continue
		}
		if !seen[t.Name] {
			seen[t.Name] = true
			topics = append(topics, t)
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: topics})
}

func (h *PlayHandler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	classID := r.URL.Query().Get("class_id")
	subjectID := r.URL.Query().Get("subject_id")
	topic := r.URL.Query().Get("topic")
	difficulty := r.URL.Query().Get("difficulty")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	if subjectID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "subject_id is required"}})
		return
	}

	query := `
		SELECT q.id, q.question_text, q.question_type, q.options, q.difficulty
		FROM questions q
		WHERE q.subject_id = $1
		AND q.deleted_at IS NULL AND q.is_active = true
		AND q.question_type IN ('mcq', 'true_false')
	`
	args := []interface{}{subjectID}
	n := 2

	if classID != "" {
		query += ` AND q.id IN (SELECT q2.id FROM questions q2 JOIN class_subjects cs2 ON cs2.subject_id = q2.subject_id WHERE cs2.class_id = $` + strconv.Itoa(n) + `)`
		args = append(args, classID)
		n++
	}

	if topic != "" {
		query += ` AND q.chapters @> $` + strconv.Itoa(n) + `::jsonb`
		args = append(args, `["`+topic+`"]`)
		n++
	}
	if difficulty != "" {
		query += ` AND q.difficulty = $` + strconv.Itoa(n)
		args = append(args, difficulty)
		n++
	}

	query += ` ORDER BY RANDOM() LIMIT $` + strconv.Itoa(n)
	args = append(args, limit)

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("play: get quiz failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch questions"}})
		return
	}
	defer rows.Close()

	questions := []PlayQuestion{}
	for rows.Next() {
		var q PlayQuestion
		var optsJSON []byte
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType, &optsJSON, &q.Difficulty); err != nil {
			continue
		}
		json.Unmarshal(optsJSON, &q.Options)
		if q.QuestionType == "mcq" || q.QuestionType == "true_false" {
			questions = append(questions, q)
		}
	}

	if len(questions) == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NO_QUESTIONS", Message: "no questions found for the selected criteria"}})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: questions})
}

func (h *PlayHandler) SaveScore(w http.ResponseWriter, r *http.Request) {
	var input PlayScoreInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"}})
		return
	}

	if input.PlayerName == "" || input.ClassID == "" || input.SubjectID == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "player_name, class_id, and subject_id are required"}})
		return
	}
	if input.Score < 0 {
		input.Score = 0
	}
	if input.TotalQuestions <= 0 {
		input.TotalQuestions = 10
	}
	if input.Difficulty == "" {
		input.Difficulty = "medium"
	}

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO play_high_scores (id, player_name, class_id, subject_id, topic, difficulty, score, total_questions, correct_count, best_streak, time_taken_ms)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		id, input.PlayerName, input.ClassID, input.SubjectID, input.Topic, input.Difficulty,
		input.Score, input.TotalQuestions, input.CorrectCount, input.BestStreak, input.TimeTakenMs,
	)
	if err != nil {
		log.Error().Err(err).Msg("play: save score failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to save score"}})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *PlayHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	classID := r.URL.Query().Get("class_id")
	subjectID := r.URL.Query().Get("subject_id")
	difficulty := r.URL.Query().Get("difficulty")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	query := `
		SELECT player_name, score, correct_count, total_questions, best_streak, time_taken_ms, difficulty,
		       to_char(created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM play_high_scores
		WHERE 1=1
	`
	args := []interface{}{}
	n := 1

	if classID != "" {
		query += ` AND class_id = $` + strconv.Itoa(n)
		args = append(args, classID)
		n++
	}
	if subjectID != "" {
		query += ` AND subject_id = $` + strconv.Itoa(n)
		args = append(args, subjectID)
		n++
	}
	if difficulty != "" {
		query += ` AND difficulty = $` + strconv.Itoa(n)
		args = append(args, difficulty)
		n++
	}

	query += ` ORDER BY score DESC, time_taken_ms ASC LIMIT $` + strconv.Itoa(n)
	args = append(args, limit)

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("play: leaderboard failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch leaderboard"}})
		return
	}
	defer rows.Close()

	scores := []HighScore{}
	rank := 1
	for rows.Next() {
		var s HighScore
		s.Rank = rank
		if err := rows.Scan(&s.PlayerName, &s.Score, &s.CorrectCount, &s.TotalQuestions, &s.BestStreak, &s.TimeTakenMs, &s.Difficulty, &s.CreatedAt); err != nil {
			continue
		}
		scores = append(scores, s)
		rank++
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: scores})
}

// Shuffle options for each question (client-side can also shuffle)
func init() {
	rand.Seed(0)
}
