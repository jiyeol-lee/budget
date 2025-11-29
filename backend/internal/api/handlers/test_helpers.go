package handlers

import (
	"budget-tracker/internal/repository"
	"database/sql"
	"net/http"
	"testing"

	_ "github.com/tursodatabase/go-libsql"
)

// TestDB wraps the database for testing
type TestDB struct {
	*repository.DB
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *repository.DB {
	t.Helper()

	// Create in-memory database
	sqlDB, err := sql.Open("libsql", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	db := &repository.DB{DB: sqlDB}

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// createTestMux creates a router with the given handlers for testing
// Note: This is a simplified router for testing individual handlers
func createTestMux(
	budgetHandler *BudgetHandler,
	expectedExpenseHandler *ExpectedExpenseHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	if budgetHandler != nil {
		mux.HandleFunc("GET /api/budgets", budgetHandler.List)
		mux.HandleFunc("POST /api/budgets", budgetHandler.Create)
		mux.HandleFunc("GET /api/budgets/{id}", budgetHandler.Get)
		mux.HandleFunc("PUT /api/budgets/{id}", budgetHandler.Update)
		mux.HandleFunc("DELETE /api/budgets/{id}", budgetHandler.Delete)
	}

	if expectedExpenseHandler != nil {
		mux.HandleFunc("GET /api/expected-expenses", expectedExpenseHandler.List)
		mux.HandleFunc("POST /api/expected-expenses", expectedExpenseHandler.Create)
		mux.HandleFunc("GET /api/expected-expenses/{id}", expectedExpenseHandler.Get)
		mux.HandleFunc("PUT /api/expected-expenses/{id}", expectedExpenseHandler.Update)
		mux.HandleFunc("DELETE /api/expected-expenses/{id}", expectedExpenseHandler.Delete)
	}

	return mux
}
