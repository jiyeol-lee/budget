package repository

import (
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	SQL         string
}

// migrations holds all database migrations in order
var migrations = []Migration{
	{
		Version:     1,
		Description: "Create budget_limits table",
		SQL: `
			CREATE TABLE IF NOT EXISTS budget_limits (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				month INTEGER NOT NULL,
				year INTEGER NOT NULL,
				amount REAL NOT NULL,
				notification_threshold REAL DEFAULT 0.8,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(month, year)
			);
		`,
	},
	{
		Version:     2,
		Description: "Create expected_expenses table",
		SQL: `
			CREATE TABLE IF NOT EXISTS expected_expenses (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				item_name TEXT NOT NULL,
				source TEXT NOT NULL,
				expected_amount REAL NOT NULL,
				expense_type TEXT NOT NULL CHECK(expense_type IN ('weekly', 'monthly')),
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);
		`,
	},
	{
		Version:     3,
		Description: "Create actual_expenses table",
		SQL: `
			CREATE TABLE IF NOT EXISTS actual_expenses (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				item_name TEXT NOT NULL,
				source TEXT NOT NULL,
				actual_amount REAL NOT NULL,
				expense_type TEXT NOT NULL CHECK(expense_type IN ('weekly', 'monthly', 'misc', 'tax')),
				item_code TEXT,
				expected_expense_id INTEGER,
				receipt_date DATE DEFAULT (DATE('now')),
				month INTEGER NOT NULL,
				year INTEGER NOT NULL,
				receipt_number INTEGER NOT NULL DEFAULT 0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (expected_expense_id) REFERENCES expected_expenses(id) ON DELETE SET NULL
			);
			CREATE INDEX IF NOT EXISTS idx_actual_expenses_month_year ON actual_expenses(year, month);
			CREATE INDEX IF NOT EXISTS idx_actual_expenses_expected ON actual_expenses(expected_expense_id);
			CREATE INDEX IF NOT EXISTS idx_actual_expenses_receipt_number ON actual_expenses(receipt_number);
		`,
	},
}

// RunMigrations executes all pending database migrations
func (db *DB) RunMigrations() error {
	log.Println("Running database migrations...")

	// First, ensure schema_migrations table exists
	_, err := db.Exec(`
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

	// Run pending migrations
	for _, m := range migrations {
		if applied[m.Version] {
			log.Printf("Migration %d already applied: %s", m.Version, m.Description)
			continue
		}

		log.Printf("Applying migration %d: %s", m.Version, m.Description)

		// Execute migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.Version, err)
		}

		if _, err := tx.Exec(m.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", m.Version, err)
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

		log.Printf("Migration %d applied successfully", m.Version)
	}

	log.Println("All migrations completed")
	return nil
}
