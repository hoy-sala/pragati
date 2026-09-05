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
	dashH := NewDashboardHandler(db)
	questionH := NewQuestionHandler(db)
 	quizH := NewQuizHandler(db)
 	hpcH := NewHPCHandler(db)
 	mentorH := NewMentorHandler(db)
 	reportsH := NewReportsHandler(db)
	certH := NewCertificateHandler(db, cfg.UploadDir)
	playH := NewPlayHandler(db)
	teamQuizH := NewTeamQuizHandler(db)
	mapH := NewMapQuizHandler(db)

	roleMw := middleware.NewRoleMiddleware(jwtService)
	loginLimiter := middleware.NewRateLimiter(10, time.Minute)

	userH := NewUserHandler(db)

	r.Route("/api/v1", func(r chi.Router) {
		r.With(loginLimiter.Limit).Post("/auth/login", authH.Login)
		r.With(loginLimiter.Limit).Post("/auth/staff-login", authH.StaffLogin)
		r.With(loginLimiter.Limit).Post("/auth/student-login", authH.StudentLogin)
		r.Post("/auth/refresh", authH.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Post("/auth/logout", authH.Logout)
			r.Get("/auth/me", authH.Me)
		})

		r.Route("/play", func(r chi.Router) {
			r.Get("/classes", playH.ListClasses)
			r.Get("/subjects", playH.ListSubjects)
			r.Get("/topics", playH.ListTopics)
			r.Get("/quiz", playH.GetQuiz)
			r.Post("/score", playH.SaveScore)
			r.Get("/leaderboard", playH.GetLeaderboard)
			r.Get("/map-categories", mapH.ListCategories)
			r.Get("/map-quiz", mapH.GenerateQuiz)
			r.Post("/map-answer", mapH.CheckAnswer)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/", roleMw.RequireRole("admin")(http.HandlerFunc(userH.List)))
			r.Post("/", roleMw.RequireRole("admin")(http.HandlerFunc(userH.Create)))
			r.Patch("/{id}/toggle", roleMw.RequireRole("admin")(http.HandlerFunc(userH.ToggleActive)))
			r.Post("/{id}/reset-password", roleMw.RequireRole("admin")(http.HandlerFunc(userH.ResetPassword)))
			r.Get("/{id}/teacher-detail", roleMw.RequireRole("admin")(http.HandlerFunc(userH.TeacherDetail)))
			r.Put("/{id}/teacher-detail", roleMw.RequireRole("admin")(http.HandlerFunc(userH.UpdateTeacherDetail)))
		})

		r.Route("/students", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/search", studentH.Search)
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

		r.With(roleMw.Authenticate).Get("/dashboard/stats", dashH.Stats)
		r.With(roleMw.Authenticate).Get("/dashboard/staff", dashH.StaffDashboard)
		r.With(roleMw.Authenticate).Get("/dashboard/student", dashH.StudentInsights)

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

 		r.Route("/mentors", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/assignments", mentorH.ListAssignments)
			r.Post("/assignments", roleMw.RequireRole("admin")(http.HandlerFunc(mentorH.CreateAssignment)))
			r.Delete("/assignments/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(mentorH.DeleteAssignment)))
			r.Get("/stats", mentorH.Stats)
			r.Get("/roster", mentorH.Roster)
			r.Get("/attendance", mentorH.GetAttendance)
			r.Put("/attendance", mentorH.SaveAttendance)
			r.Post("/contact-parent", mentorH.ContactParent)
			r.Get("/logs", mentorH.ListLogs)
			r.Post("/logs", mentorH.CreateLog)
			r.Get("/principal/alerts", roleMw.RequireRole("admin","principal")(http.HandlerFunc(mentorH.PrincipalAlerts)))
			r.Put("/logs/{id}/review", roleMw.RequireRole("admin","principal")(http.HandlerFunc(mentorH.ReviewLog)))
			r.Get("/principal/summary", roleMw.RequireRole("admin","principal")(http.HandlerFunc(mentorH.MonthlySummary)))
		})

		r.Route("/reports", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/mark-sheet", reportsH.MarkSheet)
			r.Get("/student", reportsH.StudentReport)
			r.Get("/student-me", reportsH.StudentSelf)
			r.Get("/mentors", reportsH.MentorReport)
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

		r.Route("/team-quizzes", func(r chi.Router) {
			// Public for creation/listing under quizzes — school_id defaults if unauthenticated
			r.Post("/", teamQuizH.Create)
			r.Get("/", teamQuizH.List)
			r.Get("/{id}", teamQuizH.Get)
			r.Delete("/{id}", teamQuizH.Delete)
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
		r.Route("/certificates", func(r chi.Router) {
			r.Use(roleMw.Authenticate)
			r.Get("/events", roleMw.RequireRole("admin")(http.HandlerFunc(certH.ListEvents)))
			r.Post("/events", roleMw.RequireRole("admin")(http.HandlerFunc(certH.CreateEvent)))
			r.Get("/events/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(certH.GetEvent)))
			r.Post("/events/{id}/participants", roleMw.RequireRole("admin")(http.HandlerFunc(certH.AddParticipant)))
			r.Post("/events/{id}/signatories", roleMw.RequireRole("admin")(http.HandlerFunc(certH.AddSignatory)))
			r.Post("/signatures", roleMw.RequireRole("admin")(http.HandlerFunc(certH.UploadSignature)))
			r.Get("/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(certH.GetCertificate)))
			r.Put("/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(certH.UpdateCertificate)))
			r.Delete("/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(certH.DeleteCertificate)))
			r.Delete("/signatories/{id}", roleMw.RequireRole("admin")(http.HandlerFunc(certH.DeleteSignatory)))
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Handle("/api/v1/uploads/*", http.StripPrefix("/api/v1/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))

	return r
}
