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
	"list-of-maldives/internal/server/routes"

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

	// Auth routes
	auth := r.PathPrefix("/auth").Subrouter()
	routes.RegisterAuthRoutes(auth, handlers.NewAuthHandler(s.db, jwtService))

	// Protected routes example
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.RequireAuth)
	protected.HandleFunc("/protected", s.protectedHandler).Methods("GET")

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
		w.Header().Set("Access-Control-Allow-Origin", "*") // Wildcard allows all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Credentials not allowed with wildcard origins

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
