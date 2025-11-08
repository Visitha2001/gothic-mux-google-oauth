// handlers/auth_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"list-of-maldives/internal/auth"
	"list-of-maldives/internal/database"
	"list-of-maldives/internal/server/middleware"
	"list-of-maldives/internal/server/models"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db         database.Service
	jwtService *auth.JWTService
}

func NewAuthHandler(db database.Service, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtService: jwtService,
	}
}

// Request and Response structures
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type AuthResponse struct {
	Token string       `json:"token,omitempty"`
	User  *models.User `json:"user"`
}

// GetAuth initiates OAuth authentication flow
func (h *AuthHandler) GetAuth(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]

	// Set the provider in the context for Gothic
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	// Begin the OAuth authentication process
	gothic.BeginAuthHandler(w, r)
}

// GetAuthCallback handles OAuth callback and creates user in DB
func (h *AuthHandler) GetAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error completing authentication: %v", err), http.StatusInternalServerError)
		return
	}

	// Find or create user in database
	dbUser, err := models.FindOrCreateByProvider(h.db, provider, user.UserID, user.Email, user.NickName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("User created:", dbUser)

	// Generate JWT token for OAuth user
	token, err := h.jwtService.GenerateToken(dbUser.UUID, dbUser.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Redirect to frontend with success
	frontendURL := os.Getenv("FRONTEND_URL")
	http.Redirect(w, r, frontendURL+"?auth=success", http.StatusSeeOther)
}

// Register handles email/password registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	db := h.db.GormDB()

	// Check if user already exists
	var existingUser models.User
	err := db.Where("email = ? AND provider = 'email'", req.Email).First(&existingUser).Error
	if err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	} else if err != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}

	// Create new user
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
		NickName: req.Nickname,
		Provider: "email",
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user.UUID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	response := AuthResponse{
		Token: token,
		User:  &user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login handles email/password login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Find user by email
	db := h.db.GormDB()
	var user models.User
	err := db.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if user.Provider != "email" {
		http.Error(w, "Please use the correct login method", http.StatusUnauthorized)
		return
	}

	if !user.CheckPassword(req.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user.UUID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	response := AuthResponse{
		Token: token,
		User:  &user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-time.Hour),
	})

	// Also handle OAuth logout if provider is specified
	provider := mux.Vars(r)["provider"]
	if provider != "" {
		r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
		gothic.Logout(w, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// GetUser returns current user info
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userObj, ok := user.(*models.User)
	if !ok {
		http.Error(w, "Invalid user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userObj)
}
