package routes

import (
	"list-of-maldives/internal/server/handlers"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes mounts auth-related endpoints under the provided router
func RegisterAuthRoutes(r *mux.Router, h *handlers.AuthHandler) {
	r.HandleFunc("/{provider}/callback", h.GetAuthCallback).Methods("GET")
	r.HandleFunc("/{provider}", h.GetAuth).Methods("GET")
	r.HandleFunc("/{provider}/logout", h.Logout).Methods("GET")
}
