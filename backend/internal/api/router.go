package api

import (
	"budget-tracker/internal/api/handlers"
	"encoding/json"
	"net/http"
)

// Handlers holds all API handlers
type Handlers struct {
	Budget          *handlers.BudgetHandler
	ExpectedExpense *handlers.ExpectedExpenseHandler
	ActualExpense   *handlers.ActualExpenseHandler
	Receipt         *handlers.ReceiptHandler
	Notification    *handlers.NotificationHandler
}

// NewRouter creates a new HTTP router with all routes configured
// Uses Go 1.22+ enhanced routing patterns
func NewRouter(h *Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", healthCheck)

	// Budget routes
	mux.HandleFunc("GET /api/budgets", h.Budget.List)
	mux.HandleFunc("POST /api/budgets", h.Budget.Create)
	mux.HandleFunc("GET /api/budgets/{id}", h.Budget.Get)
	mux.HandleFunc("PUT /api/budgets/{id}", h.Budget.Update)
	mux.HandleFunc("DELETE /api/budgets/{id}", h.Budget.Delete)

	// Expected Expenses routes
	mux.HandleFunc("GET /api/expected-expenses", h.ExpectedExpense.List)
	mux.HandleFunc("POST /api/expected-expenses", h.ExpectedExpense.Create)
	mux.HandleFunc("GET /api/expected-expenses/{id}", h.ExpectedExpense.Get)
	mux.HandleFunc("PUT /api/expected-expenses/{id}", h.ExpectedExpense.Update)
	mux.HandleFunc("DELETE /api/expected-expenses/{id}", h.ExpectedExpense.Delete)

	// Actual Expenses routes
	mux.HandleFunc("GET /api/actual-expenses", h.ActualExpense.List)
	mux.HandleFunc("POST /api/actual-expenses", h.ActualExpense.Create)
	mux.HandleFunc(
		"GET /api/actual-expenses/next-receipt-number",
		h.ActualExpense.GetNextReceiptNumber,
	)
	mux.HandleFunc("GET /api/actual-expenses/summary", h.ActualExpense.GetSummary)
	mux.HandleFunc("GET /api/actual-expenses/{id}", h.ActualExpense.Get)
	mux.HandleFunc("PUT /api/actual-expenses/{id}", h.ActualExpense.Update)
	mux.HandleFunc("DELETE /api/actual-expenses/{id}", h.ActualExpense.Delete)

	// Receipt processing route
	mux.HandleFunc("POST /api/receipts/process", h.Receipt.Process)

	// Notification routes
	mux.HandleFunc("GET /api/notifications/budget-status", h.Notification.BudgetStatus)

	return mux
}

// healthCheck handles the health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
