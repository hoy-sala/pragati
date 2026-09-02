package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/rs/zerolog/log"
)

type TeamQuizHandler struct {
	db *pgxpool.Pool
}

func NewTeamQuizHandler(db *pgxpool.Pool) *TeamQuizHandler {
	return &TeamQuizHandler{db: db}
}

type TeamQuizCreateInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Teams       int      `json:"teams"`
	PerTeam     int      `json:"per_team"`
	Chapters    []string `json:"chapters"`
	TimerSec    int      `json:"timer_sec"`
}

func (h *TeamQuizHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	schoolID := "00000000-0000-0000-0000-000000000001"
	userID := ""
	if claims != nil {
		schoolID = claims.SchoolID
		userID = claims.UserID
	} else {
		// allow public creation with default school
		if v := r.URL.Query().Get("school_id"); v != "" {
			schoolID = v
		}
	}
	var input TeamQuizCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error().Err(err).Msg("team quiz: decode failed")
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body: " + err.Error()})
		return
	}
	if input.Title == "" {
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "title required"})
		return
	}
	if input.Teams < 2 || input.Teams > 8 {
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "teams must be 2-8"})
		return
	}
	if input.PerTeam < 5 || input.PerTeam > 20 || input.PerTeam%5 != 0 {
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "per_team must be 5,10,15,20"})
		return
	}
	if len(input.Chapters) == 0 {
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "select at least one chapter"})
		return
	}
	if input.TimerSec == 0 {
		input.TimerSec = 30
	}
	if input.TimerSec != 30 {
		input.TimerSec = 30
	}

	// Fetch questions for selected chapters, distinct
	// Use jsonb array overlap
	rows, err := h.db.Query(r.Context(), `
		SELECT q.id, q.question_text, q.question_type, q.options, q.difficulty, q.chapters
		FROM questions q
		WHERE q.school_id = $1 AND q.deleted_at IS NULL AND q.is_active = true
		AND EXISTS (SELECT 1 FROM jsonb_array_elements_text(q.chapters) AS ch WHERE ch = ANY($2))
		AND jsonb_typeof(q.chapters) = 'array'
		ORDER BY RANDOM()
	`, schoolID, input.Chapters)
	if err != nil {
		log.Error().Err(err).Msg("team quiz: fetch questions failed")
		renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "fetch failed"})
		return
	}
	defer rows.Close()

	type Q struct {
		ID       string          `json:"id"`
		Text     string          `json:"question_text"`
		Type     string          `json:"question_type"`
		Options  json.RawMessage `json:"options"`
		Diff     string          `json:"difficulty"`
		Chapters json.RawMessage `json:"chapters"`
	}
	allQs := []Q{}
	for rows.Next() {
		var q Q
		var opts []byte
		var ch []byte
		if err := rows.Scan(&q.ID, &q.Text, &q.Type, &opts, &q.Diff, &ch); err != nil {
			continue
		}
		q.Options = opts
		q.Chapters = ch
		allQs = append(allQs, q)
	}
	needed := input.Teams * input.PerTeam
	if len(allQs) < needed {
		renderJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "not enough questions", "need": needed, "found": len(allQs), "message": fmt.Sprintf("Not enough questions for selected chapters. Need %d, found %d. Select more chapters or reduce per team.", needed, len(allQs))})
		return
	}
	// Ensure no repeats and difficulty equally distributed per team
	// Group by difficulty
	byDiff := map[string][]Q{"easy": {}, "medium": {}, "hard": {}}
	for _, q := range allQs {
		d := q.Diff
		if d != "easy" && d != "medium" && d != "hard" {
			d = "medium"
		}
		byDiff[d] = append(byDiff[d], q)
	}
	// Shuffle each difficulty pool to ensure random distribution and no repeats
	for k := range byDiff {
		rand.Shuffle(len(byDiff[k]), func(i, j int) { byDiff[k][i], byDiff[k][j] = byDiff[k][j], byDiff[k][i] })
	}
	// Calculate per-team distribution (multiples of 5, balanced)
	base := input.PerTeam / 3
	rem := input.PerTeam % 3
	easyPer, mediumPer, hardPer := base, base, base
	if rem >= 1 {
		mediumPer++
	}
	if rem >= 2 {
		easyPer++
	}
	// Verify enough per difficulty; if not, fallback to random distinct
	totalEasyNeeded := easyPer * input.Teams
	totalMediumNeeded := mediumPer * input.Teams
	totalHardNeeded := hardPer * input.Teams
	if len(byDiff["easy"]) < totalEasyNeeded || len(byDiff["medium"]) < totalMediumNeeded || len(byDiff["hard"]) < totalHardNeeded {
		// fallback: random distinct without balancing
		teamNames := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
		questionsByTeam := map[string][]Q{}
		idx := 0
		for t := 0; t < input.Teams; t++ {
			name := teamNames[t]
			questionsByTeam[name] = []Q{}
			for i := 0; i < input.PerTeam; i++ {
				questionsByTeam[name] = append(questionsByTeam[name], allQs[idx])
				idx++
			}
		}
		questionsJSON, _ := json.Marshal(questionsByTeam)
		chaptersJSON, _ := json.Marshal(input.Chapters)
		var id string
		err = h.db.QueryRow(r.Context(), `
			INSERT INTO team_quizzes (school_id, title, description, teams, per_team, timer_sec, chapters, questions, created_by)
			VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb, NULLIF($9,'')::uuid) RETURNING id
		`, schoolID, input.Title, input.Description, input.Teams, input.PerTeam, input.TimerSec, string(chaptersJSON), string(questionsJSON), userID).Scan(&id)
		if err != nil {
			log.Error().Err(err).Msg("team quiz: insert failed")
			renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
			return
		}
		renderJSON(w, http.StatusCreated, map[string]interface{}{"id": id, "questions_by_team": questionsByTeam})
		return
	}
	// Balanced distribution, no repeat
	teamNames := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	questionsByTeam := map[string][]Q{}
	eIdx, mIdx, hIdx := 0, 0, 0
	for t := 0; t < input.Teams; t++ {
		name := teamNames[t]
		teamQs := []Q{}
		for i := 0; i < easyPer; i++ {
			teamQs = append(teamQs, byDiff["easy"][eIdx])
			eIdx++
		}
		for i := 0; i < mediumPer; i++ {
			teamQs = append(teamQs, byDiff["medium"][mIdx])
			mIdx++
		}
		for i := 0; i < hardPer; i++ {
			teamQs = append(teamQs, byDiff["hard"][hIdx])
			hIdx++
		}
		rand.Shuffle(len(teamQs), func(i, j int) { teamQs[i], teamQs[j] = teamQs[j], teamQs[i] })
		questionsByTeam[name] = teamQs
	}
	questionsJSON, _ := json.Marshal(questionsByTeam)
	chaptersJSON, _ := json.Marshal(input.Chapters)

	var id string
	err = h.db.QueryRow(r.Context(), `
		INSERT INTO team_quizzes (school_id, title, description, teams, per_team, timer_sec, chapters, questions, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb, NULLIF($9,'')::uuid) RETURNING id
	`, schoolID, input.Title, input.Description, input.Teams, input.PerTeam, input.TimerSec, string(chaptersJSON), string(questionsJSON), userID).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("team quiz: insert failed")
		renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}
	renderJSON(w, http.StatusCreated, map[string]interface{}{"id": id, "questions_by_team": questionsByTeam})
}

func (h *TeamQuizHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	schoolID := "00000000-0000-0000-0000-000000000001"
	if claims != nil {
		schoolID = claims.SchoolID
	} else if v := r.URL.Query().Get("school_id"); v != "" {
		schoolID = v
	}
	rows, err := h.db.Query(r.Context(), `SELECT id, title, description, teams, per_team, timer_sec, chapters, created_at FROM team_quizzes WHERE school_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC`, schoolID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "list failed"})
		return
	}
	defer rows.Close()
	list := []map[string]interface{}{}
	for rows.Next() {
		var id, title, desc string
		var teams, perTeam, timer int
		var chapters []byte
		var createdAt time.Time
		rows.Scan(&id, &title, &desc, &teams, &perTeam, &timer, &chapters, &createdAt)
		var ch []string
		json.Unmarshal(chapters, &ch)
		list = append(list, map[string]interface{}{"id": id, "title": title, "description": desc, "teams": teams, "per_team": perTeam, "timer_sec": timer, "chapters": ch, "created_at": createdAt})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	renderJSON(w, http.StatusOK, list)
}

func (h *TeamQuizHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Info().Str("team_quiz_get_id", id).Msg("team quiz get")
	var title, desc, chaptersJSON, questionsJSON string
	var teams, perTeam, timer int
	var createdAt time.Time
	err := h.db.QueryRow(r.Context(), `SELECT title, description, teams, per_team, timer_sec, chapters::text, questions::text, created_at FROM team_quizzes WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&title, &desc, &teams, &perTeam, &timer, &chaptersJSON, &questionsJSON, &createdAt)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("team quiz get failed")
		renderJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	var chapters []string
	var questions map[string]interface{}
	json.Unmarshal([]byte(chaptersJSON), &chapters)
	json.Unmarshal([]byte(questionsJSON), &questions)
	renderJSON(w, http.StatusOK, map[string]interface{}{
		"id": id, "title": title, "description": desc, "teams": teams, "per_team": perTeam, "timer_sec": timer, "chapters": chapters, "questions_by_team": questions, "created_at": createdAt,
	})
}

func (h *TeamQuizHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE team_quizzes SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete failed"})
		return
	}
	renderJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}


