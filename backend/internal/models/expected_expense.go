package models

import (
	"strings"
	"time"
)

// ExpenseType represents the type of expense
type ExpenseType string

const (
	ExpenseTypeWeekly  ExpenseType = "weekly"
	ExpenseTypeMonthly ExpenseType = "monthly"
	ExpenseTypeMisc    ExpenseType = "misc"
	ExpenseTypeTax     ExpenseType = "tax"
)

// ExpectedExpense represents a planned recurring expense
type ExpectedExpense struct {
	ID             int64       `json:"id"`
	ItemName       string      `json:"item_name"`
	Source         string      `json:"source"`
	ExpectedAmount float64     `json:"expected_amount"`
	ExpenseType    ExpenseType `json:"expense_type"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// CreateExpectedExpenseRequest represents the request body for creating an expected expense
type CreateExpectedExpenseRequest struct {
	ItemName       string      `json:"item_name"`
	Source         string      `json:"source"`
	ExpectedAmount float64     `json:"expected_amount"`
	ExpenseType    ExpenseType `json:"expense_type"`
}

// UpdateExpectedExpenseRequest represents the request body for updating an expected expense
type UpdateExpectedExpenseRequest struct {
	ItemName       *string      `json:"item_name,omitempty"`
	Source         *string      `json:"source,omitempty"`
	ExpectedAmount *float64     `json:"expected_amount,omitempty"`
	ExpenseType    *ExpenseType `json:"expense_type,omitempty"`
}

// Validate validates the CreateExpectedExpenseRequest
// Only allows WEEKLY and MONTHLY expense types (no MISC for expected expenses)
func (r *CreateExpectedExpenseRequest) Validate() error {
	if strings.TrimSpace(r.ItemName) == "" {
		return ErrInvalidItemName
	}
	if len(r.ItemName) > 200 {
		return ErrInvalidItemNameLen
	}
	if strings.TrimSpace(r.Source) == "" {
		return ErrInvalidSource
	}
	if len(r.Source) > 100 {
		return ErrInvalidSourceLen
	}
	if r.ExpectedAmount < 0 {
		return ErrInvalidExpectedAmt
	}
	// Expected expenses only allow WEEKLY and MONTHLY (no MISC)
	if r.ExpenseType != ExpenseTypeWeekly && r.ExpenseType != ExpenseTypeMonthly {
		return ErrInvalidExpenseType
	}
	return nil
}

// Validate validates the UpdateExpectedExpenseRequest
// Only allows WEEKLY and MONTHLY expense types (no MISC for expected expenses)
func (r *UpdateExpectedExpenseRequest) Validate() error {
	if r.ItemName != nil {
		if strings.TrimSpace(*r.ItemName) == "" {
			return ErrInvalidItemName
		}
		if len(*r.ItemName) > 200 {
			return ErrInvalidItemNameLen
		}
	}
	if r.Source != nil {
		if strings.TrimSpace(*r.Source) == "" {
			return ErrInvalidSource
		}
		if len(*r.Source) > 100 {
			return ErrInvalidSourceLen
		}
	}
	if r.ExpectedAmount != nil && *r.ExpectedAmount < 0 {
		return ErrInvalidExpectedAmt
	}
	// Expected expenses only allow WEEKLY and MONTHLY (no MISC)
	if r.ExpenseType != nil && *r.ExpenseType != ExpenseTypeWeekly &&
		*r.ExpenseType != ExpenseTypeMonthly {
		return ErrInvalidExpenseType
	}
	return nil
}
