package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL       string
	RedisURL          string
	JWTSecret         string
	JWTIssuer         string
	JWTAccessExpiry   time.Duration
	JWTRefreshExpiry  time.Duration
	Port              string
	LogLevel          string
	CORSOrigins       []string
	UploadDir         string
	DefaultSchoolID   string
}

func Load() *Config {
	cfg := &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://pragati:pragati@localhost:5432/pragati?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:        getEnv("JWT_SECRET", "pragati-dev-secret-change-in-production"),
		JWTIssuer:        getEnv("JWT_ISSUER", "pragati"),
		JWTAccessExpiry:  getDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: getDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		Port:             getEnv("PORT", "8080"),
		LogLevel:         getEnv("LOG_LEVEL", "debug"),
		CORSOrigins:      splitEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:5050"),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
		DefaultSchoolID:  getEnv("DEFAULT_SCHOOL_ID", "00000000-0000-0000-0000-000000000001"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return fallback
}

func splitEnv(key, fallback string) []string {
	v := getEnv(key, fallback)
	var result []string
	for _, s := range split(v, ",") {
		if trimmed := trim(s); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	if start <= len(s) {
		result = append(result, s[start:])
	}
	return result
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
