package repository

import (
	"budget-tracker/internal/models"
	"database/sql"
	"time"
)

type ActualExpenseRepository struct {
	db *DB
}

func NewActualExpenseRepository(db *DB) *ActualExpenseRepository {
	return &ActualExpenseRepository{db: db}
}

func (r *ActualExpenseRepository) Create(
	req *models.CreateActualExpenseRequest,
) (*models.ActualExpense, error) {
	receiptDate := time.Now()
	if req.ReceiptDate != nil {
		receiptDate = *req.ReceiptDate
	}
	month := int(receiptDate.Month())
	year := receiptDate.Year()

	result, err := r.db.Exec(`
		INSERT INTO actual_expenses (item_name, source, actual_amount, expense_type, item_code, expected_expense_id, receipt_date, receipt_number, month, year)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.ItemName, req.Source, req.ActualAmount, req.ExpenseType, req.ItemCode, req.ExpectedExpenseID, receiptDate, req.ReceiptNumber, month, year)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *ActualExpenseRepository) GetByID(id int64) (*models.ActualExpense, error) {
	var expense models.ActualExpense
	var itemCode sql.NullString
	var expectedExpenseID sql.NullInt64

	err := r.db.QueryRow(`
		SELECT id, item_name, source, actual_amount, expense_type, item_code, expected_expense_id, receipt_date, receipt_number, month, year, created_at, updated_at
		FROM actual_expenses WHERE id = ?
	`, id).Scan(
		&expense.ID, &expense.ItemName, &expense.Source, &expense.ActualAmount,
		&expense.ExpenseType, &itemCode, &expectedExpenseID, &expense.ReceiptDate,
		&expense.ReceiptNumber, &expense.Month, &expense.Year, &expense.CreatedAt, &expense.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, models.ErrExpenseNotFound
	}
	if err != nil {
		return nil, err
	}

	if itemCode.Valid {
		expense.ItemCode = &itemCode.String
	}
	if expectedExpenseID.Valid {
		expense.ExpectedExpenseID = &expectedExpenseID.Int64
	}

	return &expense, nil
}

func (r *ActualExpenseRepository) GetAll() ([]models.ActualExpense, error) {
	rows, err := r.db.Query(`
		SELECT id, item_name, source, actual_amount, expense_type, item_code, expected_expense_id, receipt_date, receipt_number, month, year, created_at, updated_at
		FROM actual_expenses ORDER BY receipt_date DESC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *ActualExpenseRepository) GetByMonthYear(month, year int) ([]models.ActualExpense, error) {
	rows, err := r.db.Query(`
		SELECT id, item_name, source, actual_amount, expense_type, item_code, expected_expense_id, receipt_date, receipt_number, month, year, created_at, updated_at
		FROM actual_expenses WHERE month = ? AND year = ? ORDER BY receipt_date DESC, created_at DESC
	`, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *ActualExpenseRepository) GetByType(
	expenseType models.ExpenseType,
) ([]models.ActualExpense, error) {
	rows, err := r.db.Query(`
		SELECT id, item_name, source, actual_amount, expense_type, item_code, expected_expense_id, receipt_date, receipt_number, month, year, created_at, updated_at
		FROM actual_expenses WHERE expense_type = ? ORDER BY receipt_date DESC, created_at DESC
	`, expenseType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *ActualExpenseRepository) GetByTypeAndMonthYear(
	expenseType models.ExpenseType,
	month, year int,
) ([]models.ActualExpense, error) {
	rows, err := r.db.Query(`
		SELECT id, item_name, source, actual_amount, expense_type, item_code, expected_expense_id, receipt_date, receipt_number, month, year, created_at, updated_at
		FROM actual_expenses WHERE expense_type = ? AND month = ? AND year = ? ORDER BY receipt_date DESC, created_at DESC
	`, expenseType, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *ActualExpenseRepository) GetMonthlyTotal(month, year int) (float64, error) {
	var total sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(actual_amount), 0) FROM actual_expenses WHERE month = ? AND year = ?
	`, month, year).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total.Float64, nil
}

func (r *ActualExpenseRepository) GetMonthlySummary(
	month, year int,
) (*models.ActualExpenseSummary, error) {
	summary := &models.ActualExpenseSummary{Month: month, Year: year}

	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN expense_type = 'weekly' THEN actual_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN expense_type = 'monthly' THEN actual_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN expense_type = 'misc' THEN actual_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN expense_type = 'tax' THEN actual_amount ELSE 0 END), 0),
			COALESCE(SUM(actual_amount), 0)
		FROM actual_expenses WHERE month = ? AND year = ?
	`, month, year).Scan(&summary.TotalWeekly, &summary.TotalMonthly, &summary.TotalMisc, &summary.TotalTax, &summary.TotalActual)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *ActualExpenseRepository) Update(
	id int64,
	req *models.UpdateActualExpenseRequest,
) (*models.ActualExpense, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ItemName != nil {
		existing.ItemName = *req.ItemName
	}
	if req.Source != nil {
		existing.Source = *req.Source
	}
	if req.ActualAmount != nil {
		existing.ActualAmount = *req.ActualAmount
	}
	if req.ExpenseType != nil {
		existing.ExpenseType = *req.ExpenseType
	}
	if req.ItemCode != nil {
		existing.ItemCode = req.ItemCode
	}
	if req.ExpectedExpenseID != nil {
		existing.ExpectedExpenseID = req.ExpectedExpenseID
	}

	_, err = r.db.Exec(`
		UPDATE actual_expenses SET item_name = ?, source = ?, actual_amount = ?, expense_type = ?, item_code = ?, expected_expense_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, existing.ItemName, existing.Source, existing.ActualAmount, existing.ExpenseType, existing.ItemCode, existing.ExpectedExpenseID, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *ActualExpenseRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM actual_expenses WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrExpenseNotFound
	}

	return nil
}

func (r *ActualExpenseRepository) GetNextReceiptNumber() (int64, error) {
	var maxReceiptNumber sql.NullInt64
	err := r.db.QueryRow(`
		SELECT MAX(receipt_number) FROM actual_expenses
	`).Scan(&maxReceiptNumber)
	if err != nil {
		return 0, err
	}

	if maxReceiptNumber.Valid {
		return maxReceiptNumber.Int64 + 1, nil
	}
	return 1, nil
}

func (r *ActualExpenseRepository) scanRows(rows *sql.Rows) ([]models.ActualExpense, error) {
	var expenses []models.ActualExpense

	for rows.Next() {
		var expense models.ActualExpense
		var itemCode sql.NullString
		var expectedExpenseID sql.NullInt64

		err := rows.Scan(
			&expense.ID, &expense.ItemName, &expense.Source, &expense.ActualAmount,
			&expense.ExpenseType, &itemCode, &expectedExpenseID, &expense.ReceiptDate,
			&expense.ReceiptNumber, &expense.Month, &expense.Year, &expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if itemCode.Valid {
			expense.ItemCode = &itemCode.String
		}
		if expectedExpenseID.Valid {
			expense.ExpectedExpenseID = &expectedExpenseID.Int64
		}

		expenses = append(expenses, expense)
	}

	return expenses, rows.Err()
}
