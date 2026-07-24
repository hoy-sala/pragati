package auth

import (
	"testing"
	"time"

	"github.com/pragati/backend/internal/config"
	"github.com/pragati/backend/internal/models"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "test-password-123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword returned false for correct password")
	}
	if CheckPassword("wrong-password", hash) {
		t.Fatal("CheckPassword returned true for wrong password")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:       "test-secret-key-for-unit-tests",
		JWTIssuer:       "test-issuer",
		JWTAccessExpiry: time.Hour,
	}
	srv := NewJWTService(cfg)

	user := &models.User{
		ID:       "user-123",
		SchoolID: "school-456",
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "admin",
	}

	token, expiresAt, err := srv.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateAccessToken returned empty token")
	}
	if expiresAt <= time.Now().Unix() {
		t.Fatal("GenerateAccessToken returned expired timestamp")
	}
}

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:       "test-secret-key-for-unit-tests",
		JWTIssuer:       "test-issuer",
		JWTAccessExpiry: time.Hour,
	}
	srv := NewJWTService(cfg)

	user := &models.User{
		ID:       "user-123",
		SchoolID: "school-456",
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "admin",
	}

	token, _, err := srv.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := srv.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got '%s'", claims.UserID)
	}
	if claims.SchoolID != "school-456" {
		t.Errorf("expected SchoolID 'school-456', got '%s'", claims.SchoolID)
	}
	if claims.Role != "admin" {
		t.Errorf("expected Role 'admin', got '%s'", claims.Role)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got '%s'", claims.Email)
	}
	if claims.TokenType != "access" {
		t.Errorf("expected TokenType 'access', got '%s'", claims.TokenType)
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:       "test-secret-key-for-unit-tests",
		JWTIssuer:       "test-issuer",
		JWTAccessExpiry: time.Hour,
	}
	srv := NewJWTService(cfg)
	srvWrong := NewJWTService(&config.Config{
		JWTSecret:       "different-secret-key",
		JWTIssuer:       "test-issuer",
		JWTAccessExpiry: time.Hour,
	})

	user := &models.User{ID: "user-123", SchoolID: "school-456", Role: "admin"}
	token, _, _ := srv.GenerateAccessToken(user)

	_, err := srvWrong.ValidateToken(token)
	if err == nil {
		t.Fatal("ValidateToken should fail with wrong secret")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:       "test-secret-key-for-unit-tests",
		JWTIssuer:       "test-issuer",
		JWTAccessExpiry: -time.Hour,
	}
	srv := NewJWTService(cfg)

	user := &models.User{ID: "user-123", SchoolID: "school-456", Role: "admin"}
	token, _, err := srv.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	_, err = srv.ValidateToken(token)
	if err == nil {
		t.Fatal("ValidateToken should fail for expired token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	srv := NewJWTService(cfg)

	raw, hashed, err := srv.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if len(raw) != 64 {
		t.Errorf("expected raw token length 64, got %d", len(raw))
	}
	if hashed == "" {
		t.Fatal("GenerateRefreshToken returned empty hash")
	}
	if !CheckPassword(raw, hashed) {
		t.Fatal("bcrypt check of raw token against hash should succeed")
	}
}

func TestRefreshTokenExpiry(t *testing.T) {
	cfg := &config.Config{JWTRefreshExpiry: 7 * 24 * time.Hour}
	srv := NewJWTService(cfg)

	expiry := srv.RefreshTokenExpiry()
	if expiry != 7*24*time.Hour {
		t.Errorf("expected 7d expiry, got %v", expiry)
	}
}
