package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type MapQuizHandler struct {
	db *pgxpool.Pool
}

func NewMapQuizHandler(db *pgxpool.Pool) *MapQuizHandler {
	return &MapQuizHandler{db: db}
}

type MapPlace struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Kind      string  `json:"kind"`
	Category  string  `json:"category"`
	Map       string  `json:"map"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	SvgX      float64 `json:"svg_x"`
	SvgY      float64 `json:"svg_y"`
	State     string  `json:"state"`
	District  string  `json:"district"`
	WhyInNews string  `json:"why_in_news,omitempty"`
}

type MapOption struct {
	Key     string  `json:"key"`
	Label   string  `json:"label"`
	PlaceID string  `json:"place_id"`
	SvgX    float64 `json:"svg_x"`
	SvgY    float64 `json:"svg_y"`
	Correct bool    `json:"correct"`
}

type MapQuestion struct {
	QuestionText string      `json:"question_text"`
	PlaceID      string      `json:"place_id"`
	Category     string      `json:"category"`
	Map          string      `json:"map"`
	Explanation  string      `json:"explanation,omitempty"`
	Options      []MapOption `json:"options"`
}

func (h *MapQuizHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT category, map, COUNT(*)::int
		FROM map_places
		WHERE deleted_at IS NULL AND is_active = true
		GROUP BY category, map ORDER BY category
	`)
	if err != nil {
		log.Error().Err(err).Msg("map: list categories failed")
		renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed"})
		return
	}
	defer rows.Close()
	type Cat struct {
		Category string `json:"category"`
		Label    string `json:"label"`
		Map      string `json:"map"`
		Count    int    `json:"count"`
	}
	cats := []Cat{}
	for rows.Next() {
		var c Cat
		if err := rows.Scan(&c.Category, &c.Map, &c.Count); err != nil {
			continue
		}
		label := c.Category
		if len(label) > 5 && label[:5] == "Maps:" {
			label = label[5:]
		}
		c.Label = label
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []Cat{}
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: cats})
}

// GenerateQuiz builds 4-pin map questions: 1 correct place + 3 same-category
// distractors. Recent (news_date) places are weighted first.
func (h *MapQuizHandler) GenerateQuiz(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	mapName := r.URL.Query().Get("map")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	if category == "" {
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "category is required"})
		return
	}

	query := `
		SELECT id, name, kind, category, map, COALESCE(lat,0), COALESCE(lng,0),
		       COALESCE(svg_x,0), COALESCE(svg_y,0), COALESCE(state,''), COALESCE(district,''),
		       COALESCE(why_in_news,'')
		FROM map_places
		WHERE deleted_at IS NULL AND is_active = true AND category = $1
	`
	args := []interface{}{category}
	if mapName != "" {
		query += ` AND map = $2`
		args = append(args, mapName)
	}
	// Weight recently-in-news places first, random within recency bands
	query += ` ORDER BY (news_date IS NULL), news_date DESC NULLS LAST, RANDOM() LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("map: generate quiz failed")
		renderJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed"})
		return
	}
	defer rows.Close()

	corrects := []MapPlace{}
	for rows.Next() {
		var p MapPlace
		if err := rows.Scan(&p.ID, &p.Name, &p.Kind, &p.Category, &p.Map, &p.Lat, &p.Lng, &p.SvgX, &p.SvgY, &p.State, &p.District, &p.WhyInNews); err != nil {
			continue
		}
		corrects = append(corrects, p)
	}
	if len(corrects) == 0 {
		renderJSON(w, http.StatusNotFound, map[string]string{"error": "no places for category"})
		return
	}

	questions := []MapQuestion{}
	for _, c := range corrects {
		// 3 distractors, same category (+map), different place
		drows, err := h.db.Query(r.Context(), `
			SELECT id, name, COALESCE(svg_x,0), COALESCE(svg_y,0)
			FROM map_places
			WHERE deleted_at IS NULL AND is_active = true
			AND category = $1 AND map = $2 AND id != $3
			ORDER BY RANDOM() LIMIT 3
		`, c.Category, c.Map, c.ID)
		if err != nil {
			continue
		}
		type D struct {
			ID   string
			Name string
			X, Y float64
		}
		ds := []D{}
		for drows.Next() {
			var d D
			if err := drows.Scan(&d.ID, &d.Name, &d.X, &d.Y); err != nil {
				continue
			}
			ds = append(ds, d)
		}
		drows.Close()
		if len(ds) < 3 {
			continue
		}
		pool := []MapOption{{Key: "", Label: c.Name, PlaceID: c.ID, SvgX: c.SvgX, SvgY: c.SvgY, Correct: true}}
		for _, d := range ds {
			pool = append(pool, MapOption{Label: d.Name, PlaceID: d.ID, SvgX: d.X, SvgY: d.Y})
		}
		rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		for i := range pool {
			pool[i].Key = string(rune('A' + i))
		}
		qtext := "Click the location of " + c.Name
		if c.WhyInNews != "" {
			qtext += " (" + c.WhyInNews + ")"
		}
		questions = append(questions, MapQuestion{
			QuestionText: qtext,
			PlaceID:      c.ID,
			Category:     c.Category,
			Map:          c.Map,
			Explanation:  c.WhyInNews,
			Options:      pool,
		})
	}

	if len(questions) == 0 {
		renderJSON(w, http.StatusNotFound, map[string]string{"error": "not enough places"})
		return
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: questions})
}

// CheckAnswer verifies a tapped pin key for a generated question.
// The client sends question options + selected key; server re-checks the
// correct flag against the DB place to keep scoring authoritative.
func (h *MapQuizHandler) CheckAnswer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PlaceID     string `json:"place_id"`
		SelectedKey string `json:"selected_key"`
		Options     []struct {
			Key     string `json:"key"`
			PlaceID string `json:"place_id"`
			Correct bool   `json:"correct"`
		} `json:"options"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	correct := false
	for _, o := range input.Options {
		if o.Key == input.SelectedKey && o.PlaceID == input.PlaceID && o.Correct {
			correct = true
		}
	}
	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{"correct": correct}})
}
