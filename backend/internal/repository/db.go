package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	Path string
}

// NewDB creates a new database connection
func NewDB(cfg Config) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(cfg.Path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Build connection string with SQLite optimizations
	// WAL mode: better concurrent read performance
	// Foreign keys: enforce referential integrity
	// Busy timeout: wait up to 5 seconds when database is locked
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON&_busy_timeout=5000", cfg.Path)
	fmt.Printf("Connecting to database with DSN: %s\n", cfg.Path)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLite best practice: limit to 1 connection to avoid "database is locked" errors
	db.SetMaxOpenConns(1)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected: %s", cfg.Path)

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	log.Println("Closing database connection")
	return db.DB.Close()
}
