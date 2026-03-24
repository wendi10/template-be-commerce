package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	jwtpkg "github.com/template-be-commerce/pkg/jwt"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/apperrors"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUserRole contextKey = "user_role"
	ContextKeyEmail    contextKey = "email"
)

// Authenticate validates the Bearer JWT token on the Authorization header.
func Authenticate(jwtManager *jwtpkg.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, apperrors.ErrUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Error(w, apperrors.ErrInvalidToken)
				return
			}

			claims, err := jwtManager.ValidateAccessToken(parts[1])
			if err != nil {
				response.Error(w, apperrors.ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)
			ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin ensures the authenticated user has the admin role.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextKeyUserRole).(domain.UserRole)
		if !ok || role != domain.RoleAdmin {
			response.Error(w, apperrors.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireCustomer ensures the authenticated user has the customer role.
func RequireCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextKeyUserRole).(domain.UserRole)
		if !ok || role != domain.RoleCustomer {
			response.Error(w, apperrors.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserIDFromContext extracts the user UUID from the request context.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return id, ok
}

// RoleFromContext extracts the user role from the request context.
func RoleFromContext(ctx context.Context) (domain.UserRole, bool) {
	role, ok := ctx.Value(ContextKeyUserRole).(domain.UserRole)
	return role, ok
}
