package repository

import (
	"budget-tracker/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrBudgetNotFound = errors.New("budget limit not found")
	ErrBudgetExists   = errors.New("budget limit already exists for this month/year")
)

// BudgetRepository handles budget_limits database operations
type BudgetRepository struct {
	db *DB
}

// NewBudgetRepository creates a new BudgetRepository
func NewBudgetRepository(db *DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create creates a new budget limit
func (r *BudgetRepository) Create(
	req *models.CreateBudgetLimitRequest,
) (*models.BudgetLimit, error) {
	query := `
		INSERT INTO budget_limits (month, year, amount, notification_threshold)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, req.Month, req.Year, req.Amount, req.NotificationThreshold)
	if err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			return nil, ErrBudgetExists
		}
		return nil, fmt.Errorf("failed to create budget limit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(id)
}

// GetByID retrieves a budget limit by ID
func (r *BudgetRepository) GetByID(id int64) (*models.BudgetLimit, error) {
	query := `
		SELECT id, month, year, amount, notification_threshold, created_at, updated_at
		FROM budget_limits
		WHERE id = ?
	`

	var b models.BudgetLimit
	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &b.Month, &b.Year, &b.Amount,
		&b.NotificationThreshold, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, fmt.Errorf("failed to get budget limit: %w", err)
	}

	return &b, nil
}

// GetAll retrieves all budget limits
func (r *BudgetRepository) GetAll() ([]models.BudgetLimit, error) {
	query := `
		SELECT id, month, year, amount, notification_threshold, created_at, updated_at
		FROM budget_limits
		ORDER BY year DESC, month DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query budget limits: %w", err)
	}
	defer rows.Close()

	var budgets []models.BudgetLimit
	for rows.Next() {
		var b models.BudgetLimit
		if err := rows.Scan(
			&b.ID, &b.Month, &b.Year, &b.Amount,
			&b.NotificationThreshold, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan budget limit: %w", err)
		}
		budgets = append(budgets, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating budget limits: %w", err)
	}

	return budgets, nil
}

// Update updates a budget limit
func (r *BudgetRepository) Update(
	id int64,
	req *models.UpdateBudgetLimitRequest,
) (*models.BudgetLimit, error) {
	// First check if it exists
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Amount != nil {
		existing.Amount = *req.Amount
	}
	if req.NotificationThreshold != nil {
		existing.NotificationThreshold = *req.NotificationThreshold
	}

	query := `
		UPDATE budget_limits
		SET amount = ?, notification_threshold = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err = r.db.Exec(query, existing.Amount, existing.NotificationThreshold, now, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget limit: %w", err)
	}

	return r.GetByID(id)
}

// Delete deletes a budget limit
func (r *BudgetRepository) Delete(id int64) error {
	query := `DELETE FROM budget_limits WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete budget limit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrBudgetNotFound
	}

	return nil
}

// GetByMonthYear retrieves a budget limit by month and year
func (r *BudgetRepository) GetByMonthYear(month, year int) (*models.BudgetLimit, error) {
	query := `
		SELECT id, month, year, amount, notification_threshold, created_at, updated_at
		FROM budget_limits
		WHERE month = ? AND year = ?
	`

	var b models.BudgetLimit
	err := r.db.QueryRow(query, month, year).Scan(
		&b.ID, &b.Month, &b.Year, &b.Amount,
		&b.NotificationThreshold, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, fmt.Errorf("failed to get budget limit: %w", err)
	}

	return &b, nil
}

// Helper function to check for unique constraint errors
func isUniqueConstraintError(err error) bool {
	return err != nil && (
	// SQLite unique constraint error
	err.Error() == "UNIQUE constraint failed: budget_limits.month, budget_limits.year" ||
		// Generic check
		contains(err.Error(), "UNIQUE constraint failed"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpl(s, substr))
}

func containsImpl(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
