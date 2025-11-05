// middleware/auth.go
package middleware

import (
	"context"
	"list-of-maldives/internal/auth"
	"list-of-maldives/internal/database"
	"list-of-maldives/internal/server/models"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware validates JWT token and sets user in context
func AuthMiddleware(jwtService *auth.JWTService, db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Validate token
			claims, err := jwtService.ValidateToken(cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Find user
			var user models.User
			gormDB := db.GormDB()
			if err := gormDB.Where("uuid = ?", claims.UserID).First(&user).Error; err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth middleware protects routes that require authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserContextKey)
		if user == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
