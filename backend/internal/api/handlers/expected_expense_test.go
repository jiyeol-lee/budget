package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExpenseList_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest("GET", "/api/expected-expenses", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response ExpectedExpenseListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Expenses) != 0 {
		t.Errorf("Expected empty list, got %d expenses", len(response.Expenses))
	}
	if response.Filter != "ALL" {
		t.Errorf("Expected filter 'ALL', got '%s'", response.Filter)
	}
	if response.Count != 0 {
		t.Errorf("Expected count 0, got %d", response.Count)
	}
}

func TestExpenseList_WithData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create test expenses (only WEEKLY and MONTHLY allowed for expected expenses)
	testExpenses := []models.CreateExpectedExpenseRequest{
		{
			ItemName:       "Rent",
			Source:         "Landlord",
			ExpectedAmount: 1200.00,
			ExpenseType:    models.ExpenseTypeMonthly,
		},
		{
			ItemName:       "Groceries",
			Source:         "Supermarket",
			ExpectedAmount: 150.00,
			ExpenseType:    models.ExpenseTypeWeekly,
		},
		{
			ItemName:       "Internet",
			Source:         "ISP",
			ExpectedAmount: 50.00,
			ExpenseType:    models.ExpenseTypeMonthly,
		},
	}

	for _, e := range testExpenses {
		_, err := repo.Create(&e)
		if err != nil {
			t.Fatalf("Failed to create test expense: %v", err)
		}
	}

	req := httptest.NewRequest("GET", "/api/expected-expenses", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response ExpectedExpenseListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Expenses) != 3 {
		t.Errorf("Expected 3 expenses, got %d", len(response.Expenses))
	}
	if response.Count != 3 {
		t.Errorf("Expected count 3, got %d", response.Count)
	}
}

func TestExpenseList_FilteredByType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create test expenses (only WEEKLY and MONTHLY allowed for expected expenses)
	testExpenses := []models.CreateExpectedExpenseRequest{
		{
			ItemName:       "Rent",
			Source:         "Landlord",
			ExpectedAmount: 1200.00,
			ExpenseType:    models.ExpenseTypeMonthly,
		},
		{
			ItemName:       "Groceries",
			Source:         "Supermarket",
			ExpectedAmount: 150.00,
			ExpenseType:    models.ExpenseTypeWeekly,
		},
		{
			ItemName:       "Internet",
			Source:         "ISP",
			ExpectedAmount: 50.00,
			ExpenseType:    models.ExpenseTypeMonthly,
		},
		{
			ItemName:       "Gas",
			Source:         "Gas Station",
			ExpectedAmount: 40.00,
			ExpenseType:    models.ExpenseTypeWeekly,
		},
	}

	for _, e := range testExpenses {
		_, err := repo.Create(&e)
		if err != nil {
			t.Fatalf("Failed to create test expense: %v", err)
		}
	}

	// Test filter by WEEKLY
	t.Run("filter by WEEKLY", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/expected-expenses?type=WEEKLY", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response ExpectedExpenseListResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(response.Expenses) != 2 {
			t.Errorf("Expected 2 weekly expenses, got %d", len(response.Expenses))
		}
		if response.Filter != "WEEKLY" {
			t.Errorf("Expected filter 'WEEKLY', got '%s'", response.Filter)
		}
	})

	// Test filter by MONTHLY
	t.Run("filter by MONTHLY", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/expected-expenses?type=MONTHLY", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response ExpectedExpenseListResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(response.Expenses) != 2 {
			t.Errorf("Expected 2 monthly expenses, got %d", len(response.Expenses))
		}
		if response.Filter != "MONTHLY" {
			t.Errorf("Expected filter 'MONTHLY', got '%s'", response.Filter)
		}
	})

	// Test lowercase filter (should be normalized)
	t.Run("filter lowercase", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/expected-expenses?type=weekly", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response ExpectedExpenseListResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(response.Expenses) != 2 {
			t.Errorf(
				"Expected 2 weekly expenses with lowercase filter, got %d",
				len(response.Expenses),
			)
		}
	})
}

func TestExpenseList_InvalidFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest("GET", "/api/expected-expenses?type=INVALID", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseCreate_Valid(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	reqBody := models.CreateExpectedExpenseRequest{
		ItemName:       "Monthly Rent",
		Source:         "Landlord ABC",
		ExpectedAmount: 1500.00,
		ExpenseType:    models.ExpenseTypeMonthly,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
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

	var expense models.ExpectedExpense
	if err := json.NewDecoder(rec.Body).Decode(&expense); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if expense.ItemName != "Monthly Rent" {
		t.Errorf("Expected item name 'Monthly Rent', got '%s'", expense.ItemName)
	}
	if expense.Source != "Landlord ABC" {
		t.Errorf("Expected source 'Landlord ABC', got '%s'", expense.Source)
	}
	if expense.ExpectedAmount != 1500.00 {
		t.Errorf("Expected amount 1500.00, got %f", expense.ExpectedAmount)
	}
	if expense.ExpenseType != models.ExpenseTypeMonthly {
		t.Errorf("Expected type MONTHLY, got %s", expense.ExpenseType)
	}
	if expense.ID == 0 {
		t.Error("Expected non-zero ID")
	}
}

func TestExpenseCreate_InvalidItemName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	testCases := []struct {
		name     string
		itemName string
	}{
		{"empty", ""},
		{"whitespace only", "   "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.CreateExpectedExpenseRequest{
				ItemName:       tc.itemName,
				Source:         "Test Source",
				ExpectedAmount: 100.00,
				ExpenseType:    models.ExpenseTypeWeekly,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestExpenseCreate_ItemNameTooLong(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create item name with more than 200 characters
	longName := strings.Repeat("a", 201)

	reqBody := models.CreateExpectedExpenseRequest{
		ItemName:       longName,
		Source:         "Test Source",
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseCreate_InvalidSource(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	testCases := []struct {
		name   string
		source string
	}{
		{"empty", ""},
		{"whitespace only", "   "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.CreateExpectedExpenseRequest{
				ItemName:       "Test Item",
				Source:         tc.source,
				ExpectedAmount: 100.00,
				ExpenseType:    models.ExpenseTypeWeekly,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestExpenseCreate_SourceTooLong(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create source with more than 100 characters
	longSource := strings.Repeat("a", 101)

	reqBody := models.CreateExpectedExpenseRequest{
		ItemName:       "Test Item",
		Source:         longSource,
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseCreate_InvalidExpectedAmount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Note: expected_amount >= 0 is valid, only negative is invalid
	reqBody := models.CreateExpectedExpenseRequest{
		ItemName:       "Test Item",
		Source:         "Test Source",
		ExpectedAmount: -100.00, // Negative amount
		ExpenseType:    models.ExpenseTypeWeekly,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseCreate_ZeroExpectedAmount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Zero amount should be valid
	reqBody := models.CreateExpectedExpenseRequest{
		ItemName:       "Free Service",
		Source:         "Test Source",
		ExpectedAmount: 0.00,
		ExpenseType:    models.ExpenseTypeMonthly,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf(
			"Expected status %d for zero amount, got %d. Body: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestExpenseCreate_InvalidExpenseType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	testCases := []struct {
		name        string
		expenseType string
	}{
		{"empty", ""},
		{"invalid", "DAILY"},
		{"uppercase", "WEEKLY"},
		{"mixed case", "Weekly"},
		{"misc not allowed", "misc"}, // misc is not allowed for expected expenses
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"item_name":       "Test Item",
				"source":          "Test Source",
				"expected_amount": 100.00,
				"expense_type":    tc.expenseType,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/expected-expenses", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf(
					"Expected status %d for expense type '%s', got %d",
					http.StatusBadRequest,
					tc.expenseType,
					rec.Code,
				)
			}
		})
	}
}

func TestExpenseCreate_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest(
		"POST",
		"/api/expected-expenses",
		bytes.NewReader([]byte("invalid json")),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseGet_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create an expense first
	created, err := repo.Create(&models.CreateExpectedExpenseRequest{
		ItemName:       "Test Expense",
		Source:         "Test Source",
		ExpectedAmount: 200.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	})
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/expected-expenses/"+itoa(created.ID), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var expense models.ExpectedExpense
	if err := json.NewDecoder(rec.Body).Decode(&expense); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if expense.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, expense.ID)
	}
	if expense.ItemName != "Test Expense" {
		t.Errorf("Expected item name 'Test Expense', got '%s'", expense.ItemName)
	}
}

func TestExpenseGet_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest("GET", "/api/expected-expenses/99999", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestExpenseGet_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest("GET", "/api/expected-expenses/invalid", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseUpdate_Valid(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create an expense first
	created, err := repo.Create(&models.CreateExpectedExpenseRequest{
		ItemName:       "Original Name",
		Source:         "Original Source",
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	})
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	newName := "Updated Name"
	newSource := "Updated Source"
	newAmount := 200.00
	newType := models.ExpenseTypeMonthly

	updateReq := models.UpdateExpectedExpenseRequest{
		ItemName:       &newName,
		Source:         &newSource,
		ExpectedAmount: &newAmount,
		ExpenseType:    &newType,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(
		"PUT",
		"/api/expected-expenses/"+itoa(created.ID),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var expense models.ExpectedExpense
	if err := json.NewDecoder(rec.Body).Decode(&expense); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if expense.ItemName != "Updated Name" {
		t.Errorf("Expected item name 'Updated Name', got '%s'", expense.ItemName)
	}
	if expense.Source != "Updated Source" {
		t.Errorf("Expected source 'Updated Source', got '%s'", expense.Source)
	}
	if expense.ExpectedAmount != 200.00 {
		t.Errorf("Expected amount 200.00, got %f", expense.ExpectedAmount)
	}
	if expense.ExpenseType != models.ExpenseTypeMonthly {
		t.Errorf("Expected type MONTHLY, got %s", expense.ExpenseType)
	}
}

func TestExpenseUpdate_PartialUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create an expense first
	created, err := repo.Create(&models.CreateExpectedExpenseRequest{
		ItemName:       "Original Name",
		Source:         "Original Source",
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	})
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	// Only update item name
	newName := "Updated Name"
	updateReq := models.UpdateExpectedExpenseRequest{
		ItemName: &newName,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(
		"PUT",
		"/api/expected-expenses/"+itoa(created.ID),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var expense models.ExpectedExpense
	if err := json.NewDecoder(rec.Body).Decode(&expense); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if expense.ItemName != "Updated Name" {
		t.Errorf("Expected item name 'Updated Name', got '%s'", expense.ItemName)
	}
	// Other fields should remain unchanged
	if expense.Source != "Original Source" {
		t.Errorf("Expected source 'Original Source' (unchanged), got '%s'", expense.Source)
	}
	if expense.ExpectedAmount != 100.00 {
		t.Errorf("Expected amount 100.00 (unchanged), got %f", expense.ExpectedAmount)
	}
}

func TestExpenseUpdate_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	newName := "Updated Name"
	updateReq := models.UpdateExpectedExpenseRequest{
		ItemName: &newName,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/expected-expenses/99999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestExpenseUpdate_InvalidItemName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create an expense first
	created, err := repo.Create(&models.CreateExpectedExpenseRequest{
		ItemName:       "Original Name",
		Source:         "Original Source",
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	})
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	emptyName := ""
	updateReq := models.UpdateExpectedExpenseRequest{
		ItemName: &emptyName,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(
		"PUT",
		"/api/expected-expenses/"+itoa(created.ID),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseUpdate_InvalidExpenseType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create an expense first
	created, err := repo.Create(&models.CreateExpectedExpenseRequest{
		ItemName:       "Original Name",
		Source:         "Original Source",
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	})
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	invalidType := models.ExpenseType("INVALID")
	updateReq := models.UpdateExpectedExpenseRequest{
		ExpenseType: &invalidType,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(
		"PUT",
		"/api/expected-expenses/"+itoa(created.ID),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExpenseDelete_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	// Create an expense first
	created, err := repo.Create(&models.CreateExpectedExpenseRequest{
		ItemName:       "To Delete",
		Source:         "Test Source",
		ExpectedAmount: 100.00,
		ExpenseType:    models.ExpenseTypeWeekly,
	})
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/api/expected-expenses/"+itoa(created.ID), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Verify it's deleted
	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected expense to be deleted")
	}
}

func TestExpenseDelete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest("DELETE", "/api/expected-expenses/99999", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestExpenseDelete_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewExpectedExpenseRepository(db)
	handler := NewExpectedExpenseHandler(repo)
	mux := createTestMux(nil, handler)

	req := httptest.NewRequest("DELETE", "/api/expected-expenses/invalid", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
