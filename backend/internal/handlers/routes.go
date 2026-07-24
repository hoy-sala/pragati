package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/auth"
	"github.com/pragati/backend/internal/config"
	"github.com/pragati/backend/internal/middleware"
)

func NewRouter(db *pgxpool.Pool, jwtService *auth.JWTService, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authH := NewAuthHandler(db, jwtService, cfg)
	studentH := NewStudentHandler(db)
	classH := NewClassHandler(db)
	subjectH := NewSubjectHandler(db)
	academicYearH := NewAcademicYearHandler(db)
	catH := NewCategoryHandler(db)
	assessH := NewAssessmentHandler(db)
	markH := NewMarkHandler(db)
	questionH := NewQuestionHandler(db)
	quizH := NewQuizHandler(db)
	hpcH := NewHPCHandler(db)

	roleMw := middleware.NewRoleMiddleware(jwtService)
	loginLimiter := middleware.NewRateLimiter(10, time.Minute)

	r.Route("/api/v1", func(r chi.Router) {
		r.With(loginLimiter.Limit).Post("/auth/login", authH.Login)
		r.Post("/auth/refresh", authH.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Post("/auth/logout", authH.Logout)
			r.Get("/auth/me", authH.Me)
		})

		r.Route("/students", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", studentH.List)
			r.Post("/", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(studentH.Create)))
			r.Post("/import", roleMw.RequireRole("admin")(http.HandlerFunc(studentH.ImportCSV)))
			r.Post("/bulk-update", roleMw.RequireRole("admin", "principal")(http.HandlerFunc(studentH.BulkUpdate)))
			r.Get("/{id}", studentH.Get)
			r.Put("/{id}", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(studentH.Update)))
			r.Delete("/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(studentH.Delete)))
		})

		r.Route("/classes", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", classH.List)
			r.Post("/", roleMw.RequireRole("admin")(http.HandlerFunc(classH.Create)))
		})

		r.Route("/subjects", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", subjectH.List)
			r.Post("/", roleMw.RequireRole("admin")(http.HandlerFunc(subjectH.Create)))
		})

		r.Route("/academic-years", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", academicYearH.List)
			r.Post("/", roleMw.RequireRole("admin")(http.HandlerFunc(academicYearH.Create)))
			r.Post("/{id}/set-current", roleMw.RequireRole("admin")(http.HandlerFunc(academicYearH.SetCurrent)))
		})

		r.Route("/assessment-categories", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", catH.List)
			r.Post("/", roleMw.RequireRole("admin")(http.HandlerFunc(catH.Create)))
			r.Put("/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(catH.Update)))
		})

		r.Route("/assessments", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", assessH.List)
			r.Post("/", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(assessH.Create)))
			r.Post("/{id}/publish", roleMw.RequireRole("admin", "principal")(http.HandlerFunc(assessH.Publish)))
		})

		r.Route("/questions", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", questionH.List)
			r.Post("/", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(questionH.Create)))
			r.Post("/import/gift", roleMw.RequireRole("admin", "teacher")(http.HandlerFunc(questionH.ImportGIFT)))
			r.Post("/import/csv", roleMw.RequireRole("admin", "teacher")(http.HandlerFunc(questionH.ImportCSV)))
		})

		r.Route("/marks", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/grid", markH.GetGrid)
			r.Put("/batch", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(markH.BatchSave)))
			r.Post("/import/{id}", roleMw.RequireRole("admin", "teacher")(http.HandlerFunc(markH.ImportExcel)))
		})

		r.Route("/quizzes", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/available", quizH.GetAvailable)
			r.Post("/", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.Create)))
			r.Get("/", quizH.List)
			r.Get("/{id}", quizH.Get)
			r.Put("/{id}", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.Update)))
			r.Delete("/{id}", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.Delete)))
			r.Post("/{id}/publish", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.Publish)))
			r.Post("/{id}/attempts", quizH.StartAttempt)
			r.Post("/{id}/questions", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.AddQuestions)))
			r.Get("/{id}/questions", quizH.ListQuestions)
			r.Delete("/{id}/questions/{questionId}", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.RemoveQuestion)))
			r.Get("/attempts/{attemptId}", quizH.GetAttempt)
			r.Put("/attempts/{attemptId}/answers", quizH.SaveAnswer)
			r.Post("/attempts/{attemptId}/submit", quizH.SubmitAttempt)
			r.Get("/attempts/{attemptId}/result", quizH.GetResult)
			r.Post("/attempts/{attemptId}/grade", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(quizH.GradeShortAnswer)))
		})

		r.Route("/hpc", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/config", roleMw.RequireRole("admin")(http.HandlerFunc(hpcH.GetConfig)))
			r.Put("/config", roleMw.RequireRole("admin")(http.HandlerFunc(hpcH.SaveConfig)))
			r.Get("/learning-outcomes", hpcH.ListLearningOutcomes)
			r.Post("/learning-outcomes/import", roleMw.RequireRole("admin")(http.HandlerFunc(hpcH.ImportLearningOutcomes)))
			r.Get("/grid", hpcH.GetGrid)
			r.Get("/entries", hpcH.GetEntry)
			r.Put("/entries", hpcH.SaveEntry)
			r.Post("/entries/publish", hpcH.PublishEntry)
			r.Post("/entries/generate-pdf", hpcH.GeneratePDF)
			r.Post("/assess", roleMw.RequireRole("admin", "principal", "teacher")(http.HandlerFunc(hpcH.AssessLO)))
			r.Get("/assessments", hpcH.GetLOAssessmentGrid)
			r.Get("/reports/class", hpcH.GetClassReport)
			r.Post("/migrate-from-marks", roleMw.RequireRole("admin")(http.HandlerFunc(hpcH.MigrateFromMarks)))
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
