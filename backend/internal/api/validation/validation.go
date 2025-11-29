package validation

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error with field context
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// Add adds a validation error
func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are validation errors
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// Budget validation constants
const (
	MinMonth = 1
	MaxMonth = 12
	MinYear  = 2020
	MaxYear  = 2100
)

// Expense validation constants
const (
	MaxItemNameLength = 200
	MaxSourceLength   = 100
	MaxItemCodeLength = 50
)

// File validation constants
const (
	MaxFileSize = 10 << 20 // 10 MB
)

// ValidateMonth validates a month value
func ValidateMonth(month int) error {
	if month < MinMonth || month > MaxMonth {
		return &ValidationError{
			Field:   "month",
			Message: fmt.Sprintf("must be between %d and %d", MinMonth, MaxMonth),
		}
	}
	return nil
}

// ValidateYear validates a year value
func ValidateYear(year int) error {
	if year < MinYear || year > MaxYear {
		return &ValidationError{
			Field:   "year",
			Message: fmt.Sprintf("must be between %d and %d", MinYear, MaxYear),
		}
	}
	return nil
}

// ValidateAmount validates an amount value (must be positive)
func ValidateAmount(amount float64, field string) error {
	if amount <= 0 {
		return &ValidationError{
			Field:   field,
			Message: "must be greater than 0",
		}
	}
	return nil
}

// ValidateAmountNonNegative validates an amount value (must be >= 0)
func ValidateAmountNonNegative(amount float64, field string) error {
	if amount < 0 {
		return &ValidationError{
			Field:   field,
			Message: "must be greater than or equal to 0",
		}
	}
	return nil
}

// ValidateNotificationThreshold validates a notification threshold value
func ValidateNotificationThreshold(threshold float64) error {
	if threshold < 0 || threshold > 1 {
		return &ValidationError{
			Field:   "notification_threshold",
			Message: "must be between 0 and 1",
		}
	}
	return nil
}

// ValidateRequiredString validates a required string field
func ValidateRequiredString(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   field,
			Message: "is required",
		}
	}
	return nil
}

// ValidateStringMaxLength validates string maximum length
func ValidateStringMaxLength(value, field string, maxLength int) error {
	if len(value) > maxLength {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must not exceed %d characters", maxLength),
		}
	}
	return nil
}

// ValidateExpenseType validates expense type value
func ValidateExpenseType(expenseType string) error {
	if expenseType != "WEEKLY" && expenseType != "MONTHLY" {
		return &ValidationError{
			Field:   "expense_type",
			Message: "must be WEEKLY or MONTHLY",
		}
	}
	return nil
}

// ValidateFileSize validates file size
func ValidateFileSize(size int64) error {
	if size == 0 {
		return &ValidationError{
			Field:   "file",
			Message: "file is empty",
		}
	}
	if size > MaxFileSize {
		return &ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("file size must not exceed %d MB", MaxFileSize/(1<<20)),
		}
	}
	return nil
}

// ValidateBudgetCreate validates budget creation request
func ValidateBudgetCreate(month, year int, amount, threshold float64) error {
	errs := &ValidationErrors{}

	if err := ValidateMonth(month); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if err := ValidateYear(year); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if err := ValidateAmount(amount, "amount"); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if err := ValidateNotificationThreshold(threshold); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ValidateBudgetUpdate validates budget update request
func ValidateBudgetUpdate(amount, threshold *float64) error {
	errs := &ValidationErrors{}

	if amount != nil {
		if err := ValidateAmount(*amount, "amount"); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if threshold != nil {
		if err := ValidateNotificationThreshold(*threshold); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ValidateExpenseCreate validates expense creation request
func ValidateExpenseCreate(
	itemName, source, expenseType string,
	expectedAmount float64,
	itemCode *string,
) error {
	errs := &ValidationErrors{}

	if err := ValidateRequiredString(itemName, "item_name"); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	} else if err := ValidateStringMaxLength(itemName, "item_name", MaxItemNameLength); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if err := ValidateRequiredString(source, "source"); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	} else if err := ValidateStringMaxLength(source, "source", MaxSourceLength); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if err := ValidateAmountNonNegative(expectedAmount, "expected_amount"); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if err := ValidateExpenseType(expenseType); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	if itemCode != nil {
		if err := ValidateStringMaxLength(*itemCode, "item_code", MaxItemCodeLength); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ValidateExpenseUpdate validates expense update request
func ValidateExpenseUpdate(
	itemName, source *string,
	expectedAmount *float64,
	expenseType *string,
	itemCode *string,
) error {
	errs := &ValidationErrors{}

	if itemName != nil {
		if err := ValidateRequiredString(*itemName, "item_name"); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		} else if err := ValidateStringMaxLength(*itemName, "item_name", MaxItemNameLength); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if source != nil {
		if err := ValidateRequiredString(*source, "source"); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		} else if err := ValidateStringMaxLength(*source, "source", MaxSourceLength); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if expectedAmount != nil {
		if err := ValidateAmountNonNegative(*expectedAmount, "expected_amount"); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if expenseType != nil {
		if err := ValidateExpenseType(*expenseType); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if itemCode != nil {
		if err := ValidateStringMaxLength(*itemCode, "item_code", MaxItemCodeLength); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
