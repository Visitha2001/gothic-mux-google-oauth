package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------
// Service – unchanged public contract
// ---------------------------------------------------------------------
type Service interface {
	Health() map[string]string
	Close() error
	GormDB() *gorm.DB
}

// ---------------------------------------------------------------------
// service – now wraps a *gorm.DB
// ---------------------------------------------------------------------
type service struct {
	db *gorm.DB
}

func (s *service) GormDB() *gorm.DB {
	return s.db
}

// ---------------------------------------------------------------------
// Environment variables (required)
// ---------------------------------------------------------------------
var (
	dbName   = mustEnv("BLUEPRINT_DB_DATABASE")
	password = mustEnv("BLUEPRINT_DB_PASSWORD")
	username = mustEnv("BLUEPRINT_DB_USERNAME")
	host     = mustEnv("BLUEPRINT_DB_HOST")
	portStr  = mustEnv("BLUEPRINT_DB_PORT")
	schema   = os.Getenv("BLUEPRINT_DB_SCHEMA")
)

// singleton
var dbInstance *service

func mustEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("required env var %s is missing", name)
	}
	return v
}

// ---------------------------------------------------------------------
// New – returns the singleton Service
// ---------------------------------------------------------------------
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	// ---- validate port -------------------------------------------------
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		log.Fatalf("invalid BLUEPRINT_DB_PORT: %s", portStr)
	}

	// ---- build DSN ------------------------------------------------------
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, username, password, dbName, port,
	)

	// optional schema (search_path)
	if schema != "" {
		dsn += " search_path=" + schema
	}

	// ---- open GORM connection -------------------------------------------
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// You can tune logger, naming strategy, etc. here.
		DisableAutomaticPing: false,
	})
	if err != nil {
		log.Fatalf("failed to connect with gorm: %v", err)
	} else {
		log.Printf("Connected to PostgreSQL database: %s@%s/%s", username, host, dbName)
	}

	// ---- obtain the underlying sql.DB for pool tuning -------------------
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get *sql.DB from gorm: %v", err)
	}

	// sensible defaults – feel free to adjust
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbInstance = &service{db: gormDB}
	return dbInstance
}

// ---------------------------------------------------------------------
// Health – uses GORM's DB.PingContext + sql.DB.Stats()
// ---------------------------------------------------------------------
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// ---- ping the database ---------------------------------------------
	sqlDB, _ := s.db.DB()
	if err := sqlDB.PingContext(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	// ---- DB is alive ---------------------------------------------------
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// ---- pool statistics -----------------------------------------------
	p := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(p.OpenConnections)
	stats["in_use"] = strconv.Itoa(p.InUse)
	stats["idle"] = strconv.Itoa(p.Idle)
	stats["wait_count"] = strconv.FormatInt(p.WaitCount, 10)
	stats["wait_duration"] = p.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(p.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(p.MaxLifetimeClosed, 10)

	// ---- simple heuristics (same as before) ----------------------------
	if p.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if p.WaitCount > 1000 {
		stats["message"] = "High wait events – possible bottleneck."
	}
	if p.MaxIdleClosed > int64(p.OpenConnections)/2 {
		stats["message"] = "Many idle connections closed – review pool settings."
	}
	if p.MaxLifetimeClosed > int64(p.OpenConnections)/2 {
		stats["message"] = "Connections hitting max lifetime – consider increasing it."
	}

	return stats
}

// ---------------------------------------------------------------------
// Close – shuts down the underlying pool
// ---------------------------------------------------------------------
func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	log.Printf("Disconnected from PostgreSQL database: %s@%s/%s", username, host, dbName)
	return sqlDB.Close()
}
