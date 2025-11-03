package handlers

import (
	"context"
	"fmt"
	"list-of-maldives/internal/database"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	db database.Service
}

func NewAuthHandler(db database.Service) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) GetAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error completing authentication: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Successfully authenticated user: %+v\n", user)
	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusSeeOther)
}

func (h *AuthHandler) GetAuth(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during logout: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusSeeOther)
}
