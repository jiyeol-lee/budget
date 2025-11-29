package repository

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// legacyVersionMapping maps old integer versions to new date-based versions.
// This ensures databases with existing migrations (1, 2, 3) don't re-run
// the initial migration file (2025-11-29-001.sql = version 20251129001).
var legacyVersionMapping = map[int]int{
	1: 20251129001, // budget_limits
	2: 20251129001, // expected_expenses
	3: 20251129001, // actual_expenses
}

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	SQL         string
}

// splitSQLStatements splits SQL content into individual statements.
// It handles semicolons inside single-quoted string literals by tracking
// quote state. Each statement is trimmed of whitespace, and empty statements
// are excluded from the result.
//
// Limitations:
// - Does not handle escaped quotes within strings (e.g., 'it”s')
// - Does not handle double-quoted identifiers
// - Does not strip SQL comments (-- or /* */)
// These edge cases are rare in typical migration files.
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inString := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]

		if ch == '\'' {
			// Toggle string state on single quote
			inString = !inString
			current.WriteByte(ch)
		} else if ch == ';' && !inString {
			// Statement delimiter found outside of string
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}

	// Handle final statement (may not end with semicolon)
	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// parseFilename extracts version number from a migration filename.
// Format: YYYY-MM-DD-NNN.sql -> YYYYMMDDNNN (e.g., "2025-11-29-001.sql" -> 20251129001)
func parseFilename(filename string) (int, error) {
	// Remove .sql extension
	name := strings.TrimSuffix(filename, ".sql")

	// Remove dashes to form version number: 2025-11-29-001 -> 20251129001
	versionStr := strings.ReplaceAll(name, "-", "")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("invalid migration filename %q: %w", filename, err)
	}

	return version, nil
}

// loadMigrations reads all SQL migration files from the embedded filesystem,
// parses their filenames to extract versions, and returns them sorted by version.
func loadMigrations() ([]Migration, error) {
	// Read directory entries from migrationsFS
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration

	for _, entry := range entries {
		// Skip directories and non-.sql files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Parse filename to get version
		version, err := parseFilename(entry.Name())
		if err != nil {
			return nil, err
		}

		// Read file content
		content, err := fs.ReadFile(migrationsFS, "migrations/"+entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		// Create Migration struct with description as filename without .sql extension
		description := strings.TrimSuffix(entry.Name(), ".sql")
		migrations = append(migrations, Migration{
			Version:     version,
			Description: description,
			SQL:         string(content),
		})
	}

	// Sort migrations by version ascending
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// RunMigrations executes all pending database migrations
func (db *DB) RunMigrations() error {
	log.Println("Running database migrations...")

	// Load migrations from embedded files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// First, ensure schema_migrations table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get applied migrations
	applied := make(map[int]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query schema_migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating migrations: %w", err)
	}

	// Handle backward compatibility: mark new versions as applied if their legacy equivalents exist
	for legacyVer, newVer := range legacyVersionMapping {
		if applied[legacyVer] {
			applied[newVer] = true
		}
	}

	// Run pending migrations
	for _, m := range migrations {
		if applied[m.Version] {
			log.Printf("Migration %d (%s) already applied", m.Version, m.Description)
			continue
		}

		log.Printf("Applying migration %d: %s", m.Version, m.Description)

		// Execute migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.Version, err)
		}

		// Split migration SQL into individual statements and execute each
		statements := splitSQLStatements(m.SQL)
		for i, stmt := range statements {
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %d (statement %d): %w", m.Version, i+1, err)
			}
		}

		// Record the migration
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, description) VALUES (?, ?)",
			m.Version, m.Description,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.Version, err)
		}

		log.Printf("Migration %d (%s) applied successfully", m.Version, m.Description)
	}

	log.Println("All migrations completed")
	return nil
}
