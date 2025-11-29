package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// BudgetStatusType represents the status of budget usage
type BudgetStatusType string

const (
	BudgetStatusSafe    BudgetStatusType = "safe"
	BudgetStatusWarning BudgetStatusType = "warning"
	BudgetStatusDanger  BudgetStatusType = "danger"
	BudgetStatusOver    BudgetStatusType = "over"
)

// BudgetStatusResponse represents the budget status response
type BudgetStatusResponse struct {
	CurrentBudget  *models.BudgetLimit `json:"current_budget"`
	TotalSpent     float64             `json:"total_spent"`
	ExpectedTotal  float64             `json:"expected_total"`
	PercentageUsed float64             `json:"percentage_used"`
	Status         BudgetStatusType    `json:"status"`
	Message        string              `json:"message"`
}

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	budgetRepo          *repository.BudgetRepository
	expectedExpenseRepo *repository.ExpectedExpenseRepository
	actualExpenseRepo   *repository.ActualExpenseRepository
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(
	budgetRepo *repository.BudgetRepository,
	expectedExpenseRepo *repository.ExpectedExpenseRepository,
	actualExpenseRepo *repository.ActualExpenseRepository,
) *NotificationHandler {
	return &NotificationHandler{
		budgetRepo:          budgetRepo,
		expectedExpenseRepo: expectedExpenseRepo,
		actualExpenseRepo:   actualExpenseRepo,
	}
}

// BudgetStatus handles GET /api/notifications/budget-status
// Returns the current month's budget status with spending calculations
func (h *NotificationHandler) BudgetStatus(w http.ResponseWriter, r *http.Request) {
	// Get current month and year
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	// Parse month and year from query params if provided
	if m := r.URL.Query().Get("month"); m != "" {
		if val, err := strconv.Atoi(m); err == nil && val >= 1 && val <= 12 {
			currentMonth = val
		}
	}
	if y := r.URL.Query().Get("year"); y != "" {
		if val, err := strconv.Atoi(y); err == nil && val > 2000 {
			currentYear = val
		}
	}

	// Get budget for current month
	budget, err := h.budgetRepo.GetByMonthYear(currentMonth, currentYear)
	if err != nil {
		if errors.Is(err, repository.ErrBudgetNotFound) {
			respondJSON(w, http.StatusOK, BudgetStatusResponse{
				CurrentBudget:  nil,
				TotalSpent:     0,
				ExpectedTotal:  0,
				PercentageUsed: 0,
				Status:         BudgetStatusSafe,
				Message: fmt.Sprintf(
					"No budget set for %s %d",
					now.Month().String(),
					currentYear,
				),
			})
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch budget")
		return
	}

	// Calculate actual spending from actual_expenses table using the same summary logic
	summary, err := h.actualExpenseRepo.GetMonthlySummary(currentMonth, currentYear)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to calculate spending")
		return
	}
	totalSpent := summary.TotalActual

	// Calculate expected total from expected_expenses
	expectedTotal, err := h.expectedExpenseRepo.GetMonthlyExpectedTotal()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to calculate expected spending")
		return
	}

	// Calculate percentage used
	percentageUsed := 0.0
	if budget.Amount > 0 {
		percentageUsed = (totalSpent / budget.Amount) * 100
	}

	// Determine status and message
	status, message := h.determineStatus(
		percentageUsed,
		budget.NotificationThreshold,
		totalSpent,
		budget.Amount,
	)

	response := BudgetStatusResponse{
		CurrentBudget:  budget,
		TotalSpent:     totalSpent,
		ExpectedTotal:  expectedTotal,
		PercentageUsed: percentageUsed,
		Status:         status,
		Message:        message,
	}

	respondJSON(w, http.StatusOK, response)
}

// determineStatus determines the budget status based on percentage used
func (h *NotificationHandler) determineStatus(
	percentageUsed, threshold float64,
	spent, budget float64,
) (BudgetStatusType, string) {
	thresholdPercent := threshold * 100

	switch {
	case percentageUsed > 100:
		return BudgetStatusOver, fmt.Sprintf(
			"You've exceeded your monthly budget by $%.2f",
			spent-budget,
		)
	case percentageUsed >= 90:
		return BudgetStatusDanger, fmt.Sprintf(
			"You've used %.0f%% of your monthly budget - approaching limit!",
			percentageUsed,
		)
	case percentageUsed >= thresholdPercent:
		return BudgetStatusWarning, fmt.Sprintf(
			"You've used %.0f%% of your monthly budget",
			percentageUsed,
		)
	default:
		return BudgetStatusSafe, fmt.Sprintf(
			"You've used %.0f%% of your monthly budget - on track!",
			percentageUsed,
		)
	}
}
