package repository

import (
	"budget-tracker/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrExpenseNotFound = errors.New("expense not found")

// ExpectedExpenseRepository handles expected_expenses database operations
type ExpectedExpenseRepository struct {
	db *DB
}

// NewExpectedExpenseRepository creates a new ExpectedExpenseRepository
func NewExpectedExpenseRepository(db *DB) *ExpectedExpenseRepository {
	return &ExpectedExpenseRepository{db: db}
}

// Create creates a new expected expense
func (r *ExpectedExpenseRepository) Create(
	req *models.CreateExpectedExpenseRequest,
) (*models.ExpectedExpense, error) {
	query := `
		INSERT INTO expected_expenses (item_name, source, expected_amount, expense_type)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		req.ItemName,
		req.Source,
		req.ExpectedAmount,
		req.ExpenseType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create expected expense: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(id)
}

// GetByID retrieves an expected expense by ID
func (r *ExpectedExpenseRepository) GetByID(id int64) (*models.ExpectedExpense, error) {
	query := `
		SELECT id, item_name, source, expected_amount, expense_type, created_at, updated_at
		FROM expected_expenses
		WHERE id = ?
	`

	var e models.ExpectedExpense
	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.ItemName, &e.Source, &e.ExpectedAmount,
		&e.ExpenseType, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrExpenseNotFound
		}
		return nil, fmt.Errorf("failed to get expected expense: %w", err)
	}

	return &e, nil
}

// GetAll retrieves all expected expenses
func (r *ExpectedExpenseRepository) GetAll() ([]models.ExpectedExpense, error) {
	query := `
		SELECT id, item_name, source, expected_amount, expense_type, created_at, updated_at
		FROM expected_expenses
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query expected expenses: %w", err)
	}
	defer rows.Close()

	var expenses []models.ExpectedExpense
	for rows.Next() {
		var e models.ExpectedExpense
		if err := rows.Scan(
			&e.ID, &e.ItemName, &e.Source, &e.ExpectedAmount,
			&e.ExpenseType, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan expected expense: %w", err)
		}
		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expected expenses: %w", err)
	}

	return expenses, nil
}

// Update updates an expected expense
func (r *ExpectedExpenseRepository) Update(
	id int64,
	req *models.UpdateExpectedExpenseRequest,
) (*models.ExpectedExpense, error) {
	// First check if it exists
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.ItemName != nil {
		existing.ItemName = *req.ItemName
	}
	if req.Source != nil {
		existing.Source = *req.Source
	}
	if req.ExpectedAmount != nil {
		existing.ExpectedAmount = *req.ExpectedAmount
	}
	if req.ExpenseType != nil {
		existing.ExpenseType = *req.ExpenseType
	}

	query := `
		UPDATE expected_expenses
		SET item_name = ?, source = ?, expected_amount = ?, expense_type = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err = r.db.Exec(query, existing.ItemName, existing.Source, existing.ExpectedAmount,
		existing.ExpenseType, now, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update expected expense: %w", err)
	}

	return r.GetByID(id)
}

// Delete deletes an expected expense
func (r *ExpectedExpenseRepository) Delete(id int64) error {
	query := `DELETE FROM expected_expenses WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expected expense: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrExpenseNotFound
	}

	return nil
}

// GetByType retrieves expected expenses by type
func (r *ExpectedExpenseRepository) GetByType(
	expenseType models.ExpenseType,
) ([]models.ExpectedExpense, error) {
	query := `
		SELECT id, item_name, source, expected_amount, expense_type, created_at, updated_at
		FROM expected_expenses
		WHERE expense_type = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, expenseType)
	if err != nil {
		return nil, fmt.Errorf("failed to query expected expenses by type: %w", err)
	}
	defer rows.Close()

	var expenses []models.ExpectedExpense
	for rows.Next() {
		var e models.ExpectedExpense
		if err := rows.Scan(
			&e.ID, &e.ItemName, &e.Source, &e.ExpectedAmount,
			&e.ExpenseType, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan expected expense: %w", err)
		}
		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expected expenses: %w", err)
	}

	return expenses, nil
}

// GetMonthlyExpectedTotal calculates the expected monthly total
// Weekly expenses are multiplied by 4 for monthly estimate
func (r *ExpectedExpenseRepository) GetMonthlyExpectedTotal() (float64, error) {
	expenses, err := r.GetAll()
	if err != nil {
		return 0, err
	}

	var totalMonthly float64
	for _, expense := range expenses {
		if expense.ExpenseType == models.ExpenseTypeWeekly {
			// Weekly expenses: multiply by 4 for monthly estimate
			totalMonthly += expense.ExpectedAmount * 4
		} else if expense.ExpenseType == models.ExpenseTypeMonthly {
			// Monthly expenses: add directly
			totalMonthly += expense.ExpectedAmount
		}
	}

	return totalMonthly, nil
}
