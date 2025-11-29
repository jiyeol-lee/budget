package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

// Mode represents the database connection mode
type Mode string

const (
	ModeLocal  Mode = "local"  // Local file database for development
	ModeRemote Mode = "remote" // Turso cloud database for production
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	Mode        Mode   // Connection mode: "local" or "remote"
	LocalPath   string // Path for local mode (e.g., "./data/budget.db")
	DatabaseURL string // Turso URL for remote mode (e.g., "libsql://xxx.turso.io")
	AuthToken   string // Turso auth token for remote mode
}

// NewConfigFromEnv creates a Config from environment variables
func NewConfigFromEnv() Config {
	mode := os.Getenv("TURSO_MODE")
	if mode == "" {
		mode = "local" // Default to local mode
	}

	return Config{
		Mode:        Mode(mode),
		LocalPath:   getEnvOrDefault("TURSO_LOCAL_PATH", "./data/budget.db"),
		DatabaseURL: os.Getenv("TURSO_DATABASE_URL"),
		AuthToken:   os.Getenv("TURSO_AUTH_TOKEN"),
	}
}

// getEnvOrDefault returns the environment variable value or a default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewDB creates a new database connection
func NewDB(cfg Config) (*DB, error) {
	var dsn string

	switch cfg.Mode {
	case ModeLocal:
		// Ensure the directory exists
		dir := filepath.Dir(cfg.LocalPath)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}
		}
		// Local mode: use file path with pragmas
		// WAL mode: better concurrent read performance
		// Foreign keys: enforce referential integrity
		// Busy timeout: wait up to 5 seconds when database is locked
		dsn = fmt.Sprintf(
			"file:%s?_journal_mode=WAL&_foreign_keys=ON&_busy_timeout=5000",
			cfg.LocalPath,
		)
		log.Printf("Connecting to local database: %s", cfg.LocalPath)

	case ModeRemote:
		// Validate required fields for remote mode
		if cfg.DatabaseURL == "" {
			return nil, fmt.Errorf("DatabaseURL is required for remote mode")
		}
		if cfg.AuthToken == "" {
			return nil, fmt.Errorf("AuthToken is required for remote mode")
		}
		// Remote mode: use Turso URL with auth token
		dsn = fmt.Sprintf("%s?authToken=%s", cfg.DatabaseURL, cfg.AuthToken)
		log.Printf("Connecting to remote database: %s", cfg.DatabaseURL)

	default:
		return nil, fmt.Errorf("invalid database mode: %s", cfg.Mode)
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool settings
	if cfg.Mode == ModeLocal {
		// SQLite best practice: limit to 1 connection to avoid "database is locked" errors
		db.SetMaxOpenConns(1)
	} else {
		// Allow more connections for remote Turso database
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(5 * time.Minute)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully (mode: %s)", cfg.Mode)

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	log.Println("Closing database connection")
	return db.DB.Close()
}
