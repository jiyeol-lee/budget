package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// ExpectedExpenseListResponse represents the response for listing expected expenses with filter info
type ExpectedExpenseListResponse struct {
	Expenses []models.ExpectedExpense `json:"expenses"`
	Filter   string                   `json:"filter"`
	Count    int                      `json:"count"`
}

// ExpectedExpenseHandler handles expected expense-related HTTP requests
type ExpectedExpenseHandler struct {
	repo *repository.ExpectedExpenseRepository
}

// NewExpectedExpenseHandler creates a new ExpectedExpenseHandler
func NewExpectedExpenseHandler(repo *repository.ExpectedExpenseRepository) *ExpectedExpenseHandler {
	return &ExpectedExpenseHandler{repo: repo}
}

// List handles GET /api/expected-expenses
// Supports optional query parameter: ?type=WEEKLY or ?type=MONTHLY (no MISC for expected expenses)
func (h *ExpectedExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check for type filter query parameter
	typeFilter := r.URL.Query().Get("type")

	var expenses []models.ExpectedExpense
	var err error
	var filterLabel string

	if typeFilter != "" {
		// Normalize to lowercase
		typeFilter = strings.ToLower(typeFilter)

		// Validate the filter value - only weekly and monthly allowed for expected expenses
		if typeFilter != string(models.ExpenseTypeWeekly) &&
			typeFilter != string(models.ExpenseTypeMonthly) {
			respondError(w, http.StatusBadRequest, "Invalid type filter. Must be weekly or monthly")
			return
		}

		expenses, err = h.repo.GetByType(models.ExpenseType(typeFilter))
		filterLabel = strings.ToUpper(typeFilter)
	} else {
		expenses, err = h.repo.GetAll()
		filterLabel = "ALL"
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch expected expenses")
		return
	}

	// Ensure we return an empty array instead of null
	if expenses == nil {
		expenses = []models.ExpectedExpense{}
	}

	response := ExpectedExpenseListResponse{
		Expenses: expenses,
		Filter:   filterLabel,
		Count:    len(expenses),
	}

	respondJSON(w, http.StatusOK, response)
}

// Create handles POST /api/expected-expenses
func (h *ExpectedExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateExpectedExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	expense, err := h.repo.Create(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create expected expense")
		return
	}

	respondJSON(w, http.StatusCreated, expense)
}

// Get handles GET /api/expected-expenses/{id}
func (h *ExpectedExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	expense, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrExpenseNotFound) {
			respondError(w, http.StatusNotFound, "Expense not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch expected expense")
		return
	}

	respondJSON(w, http.StatusOK, expense)
}

// Update handles PUT /api/expected-expenses/{id}
func (h *ExpectedExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	var req models.UpdateExpectedExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	expense, err := h.repo.Update(id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrExpenseNotFound) {
			respondError(w, http.StatusNotFound, "Expense not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update expected expense")
		return
	}

	respondJSON(w, http.StatusOK, expense)
}

// Delete handles DELETE /api/expected-expenses/{id}
func (h *ExpectedExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrExpenseNotFound) {
			respondError(w, http.StatusNotFound, "Expense not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete expected expense")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
