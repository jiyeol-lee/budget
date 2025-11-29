package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// BudgetHandler handles budget-related HTTP requests
type BudgetHandler struct {
	repo *repository.BudgetRepository
}

// NewBudgetHandler creates a new BudgetHandler
func NewBudgetHandler(repo *repository.BudgetRepository) *BudgetHandler {
	return &BudgetHandler{repo: repo}
}

// List handles GET /api/budgets
func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	budgets, err := h.repo.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch budgets")
		return
	}

	// Ensure we return an empty array instead of null
	if budgets == nil {
		budgets = []models.BudgetLimit{}
	}

	respondJSON(w, http.StatusOK, budgets)
}

// Create handles POST /api/budgets
func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBudgetLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	budget, err := h.repo.Create(&req)
	if err != nil {
		if errors.Is(err, repository.ErrBudgetExists) {
			respondError(w, http.StatusConflict, "Budget for this month/year already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create budget")
		return
	}

	respondJSON(w, http.StatusCreated, budget)
}

// Get handles GET /api/budgets/{id}
func (h *BudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	budget, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrBudgetNotFound) {
			respondError(w, http.StatusNotFound, "Budget not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch budget")
		return
	}

	respondJSON(w, http.StatusOK, budget)
}

// Update handles PUT /api/budgets/{id}
func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	var req models.UpdateBudgetLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	budget, err := h.repo.Update(id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrBudgetNotFound) {
			respondError(w, http.StatusNotFound, "Budget not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update budget")
		return
	}

	respondJSON(w, http.StatusOK, budget)
}

// Delete handles DELETE /api/budgets/{id}
func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrBudgetNotFound) {
			respondError(w, http.StatusNotFound, "Budget not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete budget")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseIDFromPath extracts the ID from the URL path using Go 1.22+ PathValue
func parseIDFromPath(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errors.New("missing id")
	}
	return strconv.ParseInt(idStr, 10, 64)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
