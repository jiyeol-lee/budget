package models

import "time"

// BudgetLimit represents a monthly budget limit
type BudgetLimit struct {
	ID                    int64     `json:"id"`
	Month                 int       `json:"month"`
	Year                  int       `json:"year"`
	Amount                float64   `json:"amount"`
	NotificationThreshold float64   `json:"notification_threshold"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// CreateBudgetLimitRequest represents the request body for creating a budget limit
type CreateBudgetLimitRequest struct {
	Month                 int     `json:"month"`
	Year                  int     `json:"year"`
	Amount                float64 `json:"amount"`
	NotificationThreshold float64 `json:"notification_threshold,omitempty"`
}

// UpdateBudgetLimitRequest represents the request body for updating a budget limit
type UpdateBudgetLimitRequest struct {
	Amount                *float64 `json:"amount,omitempty"`
	NotificationThreshold *float64 `json:"notification_threshold,omitempty"`
}

// Validate validates the CreateBudgetLimitRequest
func (r *CreateBudgetLimitRequest) Validate() error {
	if r.Month < 1 || r.Month > 12 {
		return ErrInvalidMonth
	}
	if r.Year < 2020 || r.Year > 2100 {
		return ErrInvalidYear
	}
	if r.Amount <= 0 {
		return ErrInvalidAmount
	}
	if r.NotificationThreshold == 0 {
		r.NotificationThreshold = 0.8 // Default value
	}
	if r.NotificationThreshold < 0 || r.NotificationThreshold > 1 {
		return ErrInvalidThreshold
	}
	return nil
}

// Validate validates the UpdateBudgetLimitRequest
func (r *UpdateBudgetLimitRequest) Validate() error {
	if r.Amount != nil && *r.Amount <= 0 {
		return ErrInvalidAmount
	}
	if r.NotificationThreshold != nil &&
		(*r.NotificationThreshold < 0 || *r.NotificationThreshold > 1) {
		return ErrInvalidThreshold
	}
	return nil
}
