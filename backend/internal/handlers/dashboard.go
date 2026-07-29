package handlers

import (
	"net/http"

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
