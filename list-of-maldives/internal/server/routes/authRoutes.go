// routes/auth.go
package routes

import (
	"list-of-maldives/internal/server/handlers"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes mounts auth-related endpoints
func RegisterAuthRoutes(r *mux.Router, h *handlers.AuthHandler) {
	// OAuth routes
	r.HandleFunc("/{provider}/callback", h.GetAuthCallback).Methods("GET")
	r.HandleFunc("/{provider}", h.GetAuth).Methods("GET")

	// Email/Password routes
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")

	// User info
	r.HandleFunc("/me", h.GetUser).Methods("GET")
}
