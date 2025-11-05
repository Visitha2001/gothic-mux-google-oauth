package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"list-of-maldives/internal/database"
	"list-of-maldives/internal/server/models"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	// Migrate models
	User := models.User{}
	if err := NewServer.db.GormDB().AutoMigrate(&User); err != nil {
		log.Fatalf("failed to migrate User: %v", err)
	} else {
		log.Println("User migrated successfully")
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
