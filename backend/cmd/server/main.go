package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/auth"
	"github.com/pragati/backend/internal/config"
	"github.com/pragati/backend/internal/database"
	"github.com/pragati/backend/internal/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	cfg := config.Load()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		runSeed(cfg)
		return
	}

	log.Info().Str("port", cfg.Port).Msg("starting pragati api server")

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := database.RunMigrations(ctx, db, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	jwtService := auth.NewJWTService(cfg)
	router := handlers.NewRouter(db, jwtService, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("server shutdown failed")
	}
	log.Info().Msg("server stopped")
}

func runSeed(cfg *config.Config) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to database")
	}
	defer pool.Close()

	var count int
	pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if count > 0 {
		fmt.Println("Users already exist. Skipping seed.")
		return
	}

	hash, err := auth.HashPassword("pragati123")
	if err != nil {
		log.Fatal().Err(err).Msg("hash password")
	}

	adminID := uuid.New().String()
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, school_id, email, password_hash, name, role, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, true, NOW(), NOW())`,
		adminID, cfg.DefaultSchoolID, "admin@pragati.edu", hash, "School Admin", "admin",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("create admin user")
	}

	yearID := uuid.New().String()
	_, err = pool.Exec(ctx,
		`INSERT INTO academic_years (id, school_id, name, start_date, end_date, is_current)
		 VALUES ($1, $2, $3, $4, $5, true)`,
		yearID, cfg.DefaultSchoolID, "2025-26", "2025-06-01", "2026-03-31",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("create academic year")
	}

	fmt.Println("Seed complete!")
	fmt.Println("  Admin: admin@pragati.edu / pragati123")
	fmt.Println("  Academic Year: 2025-26 (current)")
}
