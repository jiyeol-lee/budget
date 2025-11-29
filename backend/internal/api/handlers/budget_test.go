package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBudgetList_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	req := httptest.NewRequest("GET", "/api/budgets", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var budgets []models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budgets); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(budgets) != 0 {
		t.Errorf("Expected empty list, got %d budgets", len(budgets))
	}
}

func TestBudgetList_WithData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create some test budgets
	testBudgets := []models.CreateBudgetLimitRequest{
		{Month: 1, Year: 2024, Amount: 1000.00, NotificationThreshold: 0.8},
		{Month: 2, Year: 2024, Amount: 1500.00, NotificationThreshold: 0.75},
		{Month: 3, Year: 2024, Amount: 2000.00, NotificationThreshold: 0.9},
	}

	for _, b := range testBudgets {
		_, err := repo.Create(&b)
		if err != nil {
			t.Fatalf("Failed to create test budget: %v", err)
		}
	}

	req := httptest.NewRequest("GET", "/api/budgets", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var budgets []models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budgets); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(budgets) != 3 {
		t.Errorf("Expected 3 budgets, got %d", len(budgets))
	}
}

func TestBudgetCreate_Valid(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	reqBody := models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                2500.00,
		NotificationThreshold: 0.85,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf(
			"Expected status %d, got %d. Body: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var budget models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budget); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if budget.Month != 6 {
		t.Errorf("Expected month 6, got %d", budget.Month)
	}
	if budget.Year != 2024 {
		t.Errorf("Expected year 2024, got %d", budget.Year)
	}
	if budget.Amount != 2500.00 {
		t.Errorf("Expected amount 2500.00, got %f", budget.Amount)
	}
	if budget.NotificationThreshold != 0.85 {
		t.Errorf("Expected threshold 0.85, got %f", budget.NotificationThreshold)
	}
	if budget.ID == 0 {
		t.Error("Expected non-zero ID")
	}
}

func TestBudgetCreate_InvalidMonth(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	testCases := []struct {
		name  string
		month int
	}{
		{"month zero", 0},
		{"month negative", -1},
		{"month too high", 13},
		{"month way too high", 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.CreateBudgetLimitRequest{
				Month:  tc.month,
				Year:   2024,
				Amount: 1000.00,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestBudgetCreate_InvalidYear(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	testCases := []struct {
		name string
		year int
	}{
		{"year too low", 2019},
		{"year way too low", 1900},
		{"year too high", 2101},
		{"year way too high", 3000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.CreateBudgetLimitRequest{
				Month:  6,
				Year:   tc.year,
				Amount: 1000.00,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestBudgetCreate_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	testCases := []struct {
		name   string
		amount float64
	}{
		{"zero amount", 0},
		{"negative amount", -100.00},
		{"very negative amount", -999999.99},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.CreateBudgetLimitRequest{
				Month:  6,
				Year:   2024,
				Amount: tc.amount,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestBudgetCreate_InvalidThreshold(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	testCases := []struct {
		name      string
		threshold float64
	}{
		{"negative threshold", -0.1},
		{"threshold too high", 1.1},
		{"threshold way too high", 5.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.CreateBudgetLimitRequest{
				Month:                 6,
				Year:                  2024,
				Amount:                1000.00,
				NotificationThreshold: tc.threshold,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestBudgetCreate_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	reqBody := models.CreateBudgetLimitRequest{
		Month:  6,
		Year:   2024,
		Amount: 1000.00,
	}

	// Create first budget
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"Expected first create to succeed with status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}

	// Try to create duplicate
	body, _ = json.Marshal(reqBody)
	req = httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Expected status %d for duplicate, got %d", http.StatusConflict, rec.Code)
	}
}

func TestBudgetCreate_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestBudgetGet_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create a budget first
	created, err := repo.Create(&models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                1500.00,
		NotificationThreshold: 0.8,
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/budgets/"+itoa(created.ID), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var budget models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budget); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if budget.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, budget.ID)
	}
	if budget.Month != 6 {
		t.Errorf("Expected month 6, got %d", budget.Month)
	}
}

func TestBudgetGet_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	req := httptest.NewRequest("GET", "/api/budgets/99999", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestBudgetGet_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	req := httptest.NewRequest("GET", "/api/budgets/invalid", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestBudgetUpdate_Valid(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create a budget first
	created, err := repo.Create(&models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                1000.00,
		NotificationThreshold: 0.8,
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	newAmount := 2000.00
	newThreshold := 0.9
	updateReq := models.UpdateBudgetLimitRequest{
		Amount:                &newAmount,
		NotificationThreshold: &newThreshold,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/budgets/"+itoa(created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var budget models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budget); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if budget.Amount != 2000.00 {
		t.Errorf("Expected amount 2000.00, got %f", budget.Amount)
	}
	if budget.NotificationThreshold != 0.9 {
		t.Errorf("Expected threshold 0.9, got %f", budget.NotificationThreshold)
	}
}

func TestBudgetUpdate_PartialUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create a budget first
	created, err := repo.Create(&models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                1000.00,
		NotificationThreshold: 0.8,
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	// Only update amount
	newAmount := 1500.00
	updateReq := models.UpdateBudgetLimitRequest{
		Amount: &newAmount,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/budgets/"+itoa(created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var budget models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budget); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if budget.Amount != 1500.00 {
		t.Errorf("Expected amount 1500.00, got %f", budget.Amount)
	}
	// Threshold should remain unchanged
	if budget.NotificationThreshold != 0.8 {
		t.Errorf("Expected threshold 0.8 (unchanged), got %f", budget.NotificationThreshold)
	}
}

func TestBudgetUpdate_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	newAmount := 2000.00
	updateReq := models.UpdateBudgetLimitRequest{
		Amount: &newAmount,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/budgets/99999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestBudgetUpdate_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create a budget first
	created, err := repo.Create(&models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                1000.00,
		NotificationThreshold: 0.8,
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	invalidAmount := -100.00
	updateReq := models.UpdateBudgetLimitRequest{
		Amount: &invalidAmount,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/budgets/"+itoa(created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestBudgetUpdate_InvalidThreshold(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create a budget first
	created, err := repo.Create(&models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                1000.00,
		NotificationThreshold: 0.8,
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	invalidThreshold := 1.5
	updateReq := models.UpdateBudgetLimitRequest{
		NotificationThreshold: &invalidThreshold,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/budgets/"+itoa(created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestBudgetDelete_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create a budget first
	created, err := repo.Create(&models.CreateBudgetLimitRequest{
		Month:                 6,
		Year:                  2024,
		Amount:                1000.00,
		NotificationThreshold: 0.8,
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/api/budgets/"+itoa(created.ID), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Verify it's deleted
	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected budget to be deleted")
	}
}

func TestBudgetDelete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	req := httptest.NewRequest("DELETE", "/api/budgets/99999", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestBudgetDelete_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	req := httptest.NewRequest("DELETE", "/api/budgets/invalid", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestBudgetCreate_DefaultThreshold(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBudgetRepository(db)
	handler := NewBudgetHandler(repo)
	mux := createTestMux(handler, nil)

	// Create budget without threshold (should default to 0.8)
	reqBody := models.CreateBudgetLimitRequest{
		Month:  7,
		Year:   2024,
		Amount: 1000.00,
		// NotificationThreshold not set, should default
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf(
			"Expected status %d, got %d. Body: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var budget models.BudgetLimit
	if err := json.NewDecoder(rec.Body).Decode(&budget); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if budget.NotificationThreshold != 0.8 {
		t.Errorf("Expected default threshold 0.8, got %f", budget.NotificationThreshold)
	}
}

// Helper function to convert int64 to string
func itoa(i int64) string {
	if i == 0 {
		return "0"
	}

	var result []byte
	negative := i < 0
	if negative {
		i = -i
	}

	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}

	if negative {
		result = append([]byte{'-'}, result...)
	}

	return string(result)
}
