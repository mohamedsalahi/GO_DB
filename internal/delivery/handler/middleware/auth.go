package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/infra/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "user_role"

func AuthMiddleware(tokenManager *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := tokenManager.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(UserIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(r.Context())
		if role != "admin" {
			http.Error(w, `{"error":"forbidden: admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
