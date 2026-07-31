package handlers

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type DashboardHandler struct {
	db *pgxpool.Pool
}

func NewDashboardHandler(db *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) StudentInsights(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	var student struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ClassName string `json:"class_name"`
	}
	err := h.db.QueryRow(r.Context(),
		`SELECT s.first_name, COALESCE(s.last_name, ''), COALESCE(c.name, '')
		 FROM students s
		 LEFT JOIN classes c ON c.id = s.class_id
		 WHERE s.id = $1 AND s.deleted_at IS NULL`,
		claims.UserID,
	).Scan(&student.FirstName, &student.LastName, &student.ClassName)
	if err != nil {
		log.Error().Err(err).Msg("fetch student for insights failed")
		renderJSON(w, http.StatusNotFound, models.APIResponse{Error: &models.APIError{Code: "NOT_FOUND", Message: "student not found"}})
		return
	}

	var quizzesTaken, quizzesPassed int
	var bestPct, avgPct float64
	if err := h.db.QueryRow(r.Context(),
		`SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE passed = true),
			COALESCE(MAX(percentage), 0),
			COALESCE(AVG(percentage), 0)
		 FROM quiz_attempts
		 WHERE user_id = $1 AND status = 'submitted' AND percentage IS NOT NULL`,
		claims.UserID,
	).Scan(&quizzesTaken, &quizzesPassed, &bestPct, &avgPct); err != nil {
		log.Error().Err(err).Msg("quiz stats for insights failed")
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT qa.title, a.percentage, a.passed, a.submitted_at
		 FROM quiz_attempts a
		 JOIN quiz_assignments qa ON qa.id = a.quiz_id
		 WHERE a.user_id = $1 AND a.status = 'submitted' AND a.percentage IS NOT NULL
		 ORDER BY a.submitted_at DESC
		 LIMIT 10`,
		claims.UserID,
	)
	recent := []map[string]interface{}{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var title string
			var pct float64
			var passed *bool
			var submittedAt *time.Time
			if err := rows.Scan(&title, &pct, &passed, &submittedAt); err != nil {
				continue
			}
			ts := ""
			if submittedAt != nil {
				ts = submittedAt.Format("2006-01-02")
			}
			isPassed := false
			if passed != nil {
				isPassed = *passed
			}
			recent = append(recent, map[string]interface{}{
				"quiz_title":   title,
				"percentage":   pct,
				"passed":       isPassed,
				"submitted_at": ts,
			})
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"student":       student,
		"quizzes_taken": quizzesTaken,
		"quizzes_passed": quizzesPassed,
		"best_percentage": bestPct,
		"average_percentage": avgPct,
		"recent_attempts": recent,
	}})
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	var totalStudents int
	if err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM students WHERE school_id = $1 AND deleted_at IS NULL AND is_active = true`,
		claims.SchoolID,
	).Scan(&totalStudents); err != nil {
		log.Error().Err(err).Msg("count students failed")
	}

	var totalTeachers int
	if err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM users WHERE school_id = $1 AND role = 'teacher' AND deleted_at IS NULL AND is_active = true`,
		claims.SchoolID,
	).Scan(&totalTeachers); err != nil {
		log.Error().Err(err).Msg("count teachers failed")
	}

	var totalClasses int
	if err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM classes WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&totalClasses); err != nil {
		log.Error().Err(err).Msg("count classes failed")
	}

	var totalAssessments int
	if err := h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM assessments WHERE school_id = $1 AND deleted_at IS NULL`,
		claims.SchoolID,
	).Scan(&totalAssessments); err != nil {
		log.Error().Err(err).Msg("count assessments failed")
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT c.name, COUNT(s.id)
		 FROM classes c
		 LEFT JOIN students s ON s.class_id = c.id AND s.deleted_at IS NULL AND s.is_active = true
		 WHERE c.school_id = $1 AND c.deleted_at IS NULL
		 GROUP BY c.id, c.name
		 ORDER BY c.sort_order ASC`,
		claims.SchoolID,
	)
	studentsByClass := []map[string]interface{}{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var className string
			var count int
			if err := rows.Scan(&className, &count); err != nil {
				log.Error().Err(err).Msg("scan class stats row failed")
				continue
			}
			studentsByClass = append(studentsByClass, map[string]interface{}{
				"class":  className,
				"count":  count,
			})
		}
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"total_students":    totalStudents,
		"total_teachers":    totalTeachers,
		"total_classes":     totalClasses,
		"total_assessments": totalAssessments,
		"students_by_class": studentsByClass,
	}})
}
