package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"list-of-maldives/internal/auth"
	"list-of-maldives/internal/server/handlers"
	"list-of-maldives/internal/server/middleware"
	"list-of-maldives/internal/server/models"

	"github.com/gorilla/mux"
)

// server/server.go (update RegisterRoutes method)
func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(s.corsMiddleware)
	r.Use(s.requestLogger)

	r.HandleFunc("/", s.HelloWorldHandler)
	r.HandleFunc("/health", s.healthHandler)

	// Initialize JWT service
	jwtService := auth.NewJWTService()

	// Apply auth middleware (sets user in context if authenticated)
	r.Use(middleware.AuthMiddleware(jwtService, s.db))

	// Auth routes (UNPROTECTED: register, login, oauth)
	authHandler := handlers.NewAuthHandler(s.db, jwtService)

	// User Info/Protected Auth Routes (PROTECTED: /auth/me)
	userAuth := r.PathPrefix("/auth/me").Subrouter()
	userAuth.Use(middleware.RequireAuth)
	userAuth.HandleFunc("", authHandler.GetUser).Methods("GET")

	auth := r.PathPrefix("/auth").Subrouter()
	// Register all routes EXCEPT /me here
	auth.HandleFunc("/{provider}/callback", authHandler.GetAuthCallback).Methods("GET", "OPTIONS")
	auth.HandleFunc("/{provider}", authHandler.GetAuth).Methods("GET", "OPTIONS")
	auth.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/log-out", authHandler.Logout).Methods("POST", "OPTIONS")

	// Protected API routes example (already correctly protected)
	protectedAPI := r.PathPrefix("/api").Subrouter()
	protectedAPI.Use(middleware.RequireAuth)
	protectedAPI.HandleFunc("/protected", s.protectedHandler).Methods("GET", "OPTIONS")

	// Protected routes example
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.RequireAuth)
	protected.HandleFunc("/protected", s.protectedHandler).Methods("GET", "OPTIONS")

	return r
}

// Example protected handler
func (s *Server) protectedHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*models.User)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "This is a protected route",
		"user":    user,
	})
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// You also need to add the Max-Age header for preflight caching
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingResponseWriter captures status code and bytes written
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}

// requestLogger logs method, path, status, bytes and duration for each request
func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)

		status := lrw.status
		if status == 0 {
			status = http.StatusOK
		}

		duration := time.Since(start)
		log.Printf("%s %s %d %dB %s", r.Method, r.URL.RequestURI(), status, lrw.bytes, duration)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
