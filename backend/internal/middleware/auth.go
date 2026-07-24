package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pragati/backend/internal/auth"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type contextKey string

const (
	UserClaimsKey contextKey = "user_claims"
	UserIDKey     contextKey = "user_id"
)

type roleMiddleware struct {
	jwtService *auth.JWTService
}

func NewRoleMiddleware(jwtService *auth.JWTService) *roleMiddleware {
	return &roleMiddleware{jwtService: jwtService}
}

func (m *roleMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"missing authorization header"}}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid authorization format"}}`, http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			log.Debug().Err(err).Msg("authentication failed")
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid or expired token"}}`, http.StatusUnauthorized)
			return
		}

		if claims.TokenType != "access" {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid token type"}}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *roleMiddleware) RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsKey).(*models.TokenClaims)
			if !ok {
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"not authenticated"}}`, http.StatusUnauthorized)
				return
			}

			for _, role := range roles {
				if claims.Role == role {
					next(w, r)
					return
				}
			}

			http.Error(w, `{"error":{"code":"FORBIDDEN","message":"insufficient permissions"}}`, http.StatusForbidden)
		}
	}
}

func GetUserClaims(ctx context.Context) *models.TokenClaims {
	claims, _ := ctx.Value(UserClaimsKey).(*models.TokenClaims)
	return claims
}

func GetUserID(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}
