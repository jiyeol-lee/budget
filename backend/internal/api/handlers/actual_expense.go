package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ActualExpenseHandler struct {
	repo *repository.ActualExpenseRepository
}

func NewActualExpenseHandler(repo *repository.ActualExpenseRepository) *ActualExpenseHandler {
	return &ActualExpenseHandler{repo: repo}
}

type ActualExpenseListResponse struct {
	Expenses []models.ActualExpense `json:"expenses"`
	Total    int                    `json:"total"`
}

func (h *ActualExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query params: month, year, type
	query := r.URL.Query()
	monthStr := query.Get("month")
	yearStr := query.Get("year")
	expenseType := query.Get("type")

	var expenses []models.ActualExpense
	var err error

	// Default to current month/year if provided
	if monthStr != "" && yearStr != "" {
		month, _ := strconv.Atoi(monthStr)
		year, _ := strconv.Atoi(yearStr)

		if expenseType != "" && expenseType != "ALL" {
			expenses, err = h.repo.GetByTypeAndMonthYear(
				models.ExpenseType(strings.ToLower(expenseType)),
				month,
				year,
			)
		} else {
			expenses, err = h.repo.GetByMonthYear(month, year)
		}
	} else if expenseType != "" && expenseType != "ALL" {
		expenses, err = h.repo.GetByType(models.ExpenseType(strings.ToLower(expenseType)))
	} else {
		expenses, err = h.repo.GetAll()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if expenses == nil {
		expenses = []models.ActualExpense{}
	}

	response := ActualExpenseListResponse{
		Expenses: expenses,
		Total:    len(expenses),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ActualExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateActualExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expense, err := h.repo.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

func (h *ActualExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	expense, err := h.repo.GetByID(id)
	if err != nil {
		if err == models.ErrExpenseNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ActualExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateActualExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expense, err := h.repo.Update(id, &req)
	if err != nil {
		if err == models.ErrExpenseNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ActualExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err == models.ErrExpenseNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ActualExpenseHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	monthStr := query.Get("month")
	yearStr := query.Get("year")

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			month = m
		}
	}
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	summary, err := h.repo.GetMonthlySummary(month, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (h *ActualExpenseHandler) GetNextReceiptNumber(w http.ResponseWriter, r *http.Request) {
	nextNumber, err := h.repo.GetNextReceiptNumber()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{
		"next_receipt_number": nextNumber,
	})
}
