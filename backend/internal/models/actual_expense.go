package models

import (
	"strings"
	"time"
)

// ActualExpense represents real spending tracked from receipts
type ActualExpense struct {
	ID                int64       `json:"id"`
	ItemName          string      `json:"item_name"`
	Source            string      `json:"source"`
	ActualAmount      float64     `json:"actual_amount"`
	ExpenseType       ExpenseType `json:"expense_type"`
	ItemCode          *string     `json:"item_code,omitempty"`
	ExpectedExpenseID *int64      `json:"expected_expense_id,omitempty"`
	ReceiptDate       time.Time   `json:"receipt_date"`
	ReceiptNumber     int64       `json:"receipt_number"`
	Month             int         `json:"month"`
	Year              int         `json:"year"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

// CreateActualExpenseRequest for creating actual expenses
type CreateActualExpenseRequest struct {
	ItemName          string      `json:"item_name"`
	Source            string      `json:"source"`
	ActualAmount      float64     `json:"actual_amount"`
	ExpenseType       ExpenseType `json:"expense_type"`
	ItemCode          *string     `json:"item_code,omitempty"`
	ExpectedExpenseID *int64      `json:"expected_expense_id,omitempty"`
	ReceiptDate       *time.Time  `json:"receipt_date,omitempty"`
	ReceiptNumber     int64       `json:"receipt_number"`
}

func (r *CreateActualExpenseRequest) Validate() error {
	r.ItemName = strings.TrimSpace(r.ItemName)
	r.Source = strings.TrimSpace(r.Source)

	if r.ItemName == "" {
		return ErrItemNameRequired
	}
	if len(r.ItemName) > 255 {
		return ErrItemNameTooLong
	}
	if r.Source == "" {
		return ErrSourceRequired
	}
	if len(r.Source) > 255 {
		return ErrSourceTooLong
	}
	if r.ActualAmount <= 0 {
		return ErrInvalidAmount
	}
	// Allow WEEKLY, MONTHLY, MISC, and TAX for actual expenses
	if r.ExpenseType != ExpenseTypeWeekly && r.ExpenseType != ExpenseTypeMonthly &&
		r.ExpenseType != ExpenseTypeMisc && r.ExpenseType != ExpenseTypeTax {
		return ErrInvalidExpenseType
	}
	return nil
}

// UpdateActualExpenseRequest for updating actual expenses
type UpdateActualExpenseRequest struct {
	ItemName          *string      `json:"item_name,omitempty"`
	Source            *string      `json:"source,omitempty"`
	ActualAmount      *float64     `json:"actual_amount,omitempty"`
	ExpenseType       *ExpenseType `json:"expense_type,omitempty"`
	ItemCode          *string      `json:"item_code,omitempty"`
	ExpectedExpenseID *int64       `json:"expected_expense_id,omitempty"`
}

func (r *UpdateActualExpenseRequest) Validate() error {
	if r.ItemName != nil {
		*r.ItemName = strings.TrimSpace(*r.ItemName)
		if *r.ItemName == "" {
			return ErrItemNameRequired
		}
		if len(*r.ItemName) > 255 {
			return ErrItemNameTooLong
		}
	}
	if r.Source != nil {
		*r.Source = strings.TrimSpace(*r.Source)
		if *r.Source == "" {
			return ErrSourceRequired
		}
		if len(*r.Source) > 255 {
			return ErrSourceTooLong
		}
	}
	if r.ActualAmount != nil && *r.ActualAmount <= 0 {
		return ErrInvalidAmount
	}
	if r.ExpenseType != nil {
		if *r.ExpenseType != ExpenseTypeWeekly && *r.ExpenseType != ExpenseTypeMonthly &&
			*r.ExpenseType != ExpenseTypeMisc && *r.ExpenseType != ExpenseTypeTax {
			return ErrInvalidExpenseType
		}
	}
	return nil
}

// ActualExpenseSummary for aggregated data
type ActualExpenseSummary struct {
	Month        int     `json:"month"`
	Year         int     `json:"year"`
	TotalWeekly  float64 `json:"total_weekly"`
	TotalMonthly float64 `json:"total_monthly"`
	TotalMisc    float64 `json:"total_misc"`
	TotalTax     float64 `json:"total_tax"`
	TotalActual  float64 `json:"total_actual"`
}
