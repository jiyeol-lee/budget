package models

import "errors"

// Common validation errors
var (
	ErrInvalidMonth       = errors.New("month must be between 1 and 12")
	ErrInvalidYear        = errors.New("year must be between 2020 and 2100")
	ErrInvalidAmount      = errors.New("amount must be greater than 0")
	ErrInvalidThreshold   = errors.New("notification threshold must be between 0 and 1")
	ErrInvalidItemName    = errors.New("item name is required")
	ErrInvalidSource      = errors.New("source is required")
	ErrInvalidExpenseType = errors.New("expense type must be weekly, monthly, misc, or tax")
	ErrInvalidItemNameLen = errors.New("item name must not exceed 200 characters")
	ErrInvalidSourceLen   = errors.New("source must not exceed 100 characters")
	ErrInvalidItemCodeLen = errors.New("item code must not exceed 50 characters")
	ErrInvalidExpectedAmt = errors.New("expected amount must be greater than or equal to 0")
	ErrExpenseNotFound    = errors.New("expense not found")

	// Actual expense validation errors
	ErrItemNameRequired = errors.New("item name is required")
	ErrItemNameTooLong  = errors.New("item name must not exceed 255 characters")
	ErrSourceRequired   = errors.New("source is required")
	ErrSourceTooLong    = errors.New("source must not exceed 255 characters")
)
