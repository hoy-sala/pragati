package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID    string   `json:"uid"`
	SchoolID  string   `json:"sid"`
	Role      string   `json:"rol"`
	Email     string   `json:"eml"`
	TokenType string   `json:"ttp"`
	Permissions []string `json:"prm,omitempty"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	TokenHash string    `json:"token_hash" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}
