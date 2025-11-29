package validation

import (
	"strings"
	"testing"
)

func TestValidateMonth(t *testing.T) {
	testCases := []struct {
		name    string
		month   int
		wantErr bool
	}{
		{"valid month 1", 1, false},
		{"valid month 6", 6, false},
		{"valid month 12", 12, false},
		{"invalid month 0", 0, true},
		{"invalid month -1", -1, true},
		{"invalid month 13", 13, true},
		{"invalid month 100", 100, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMonth(tc.month)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateMonth(%d) error = %v, wantErr %v", tc.month, err, tc.wantErr)
			}
		})
	}
}

func TestValidateYear(t *testing.T) {
	testCases := []struct {
		name    string
		year    int
		wantErr bool
	}{
		{"valid year 2020", 2020, false},
		{"valid year 2024", 2024, false},
		{"valid year 2100", 2100, false},
		{"invalid year 2019", 2019, true},
		{"invalid year 1900", 1900, true},
		{"invalid year 2101", 2101, true},
		{"invalid year 3000", 3000, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateYear(tc.year)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateYear(%d) error = %v, wantErr %v", tc.year, err, tc.wantErr)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	testCases := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"valid amount 0.01", 0.01, false},
		{"valid amount 100", 100.00, false},
		{"valid amount 999999", 999999.99, false},
		{"invalid amount 0", 0, true},
		{"invalid amount -1", -1.00, true},
		{"invalid amount -100", -100.00, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAmount(tc.amount, "amount")
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateAmount(%f) error = %v, wantErr %v", tc.amount, err, tc.wantErr)
			}
		})
	}
}

func TestValidateAmountNonNegative(t *testing.T) {
	testCases := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"valid amount 0", 0, false},
		{"valid amount 0.01", 0.01, false},
		{"valid amount 100", 100.00, false},
		{"invalid amount -0.01", -0.01, true},
		{"invalid amount -100", -100.00, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAmountNonNegative(tc.amount, "expected_amount")
			if (err != nil) != tc.wantErr {
				t.Errorf(
					"ValidateAmountNonNegative(%f) error = %v, wantErr %v",
					tc.amount,
					err,
					tc.wantErr,
				)
			}
		})
	}
}

func TestValidateNotificationThreshold(t *testing.T) {
	testCases := []struct {
		name      string
		threshold float64
		wantErr   bool
	}{
		{"valid threshold 0", 0, false},
		{"valid threshold 0.5", 0.5, false},
		{"valid threshold 0.8", 0.8, false},
		{"valid threshold 1", 1, false},
		{"invalid threshold -0.1", -0.1, true},
		{"invalid threshold 1.1", 1.1, true},
		{"invalid threshold 5", 5, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateNotificationThreshold(tc.threshold)
			if (err != nil) != tc.wantErr {
				t.Errorf(
					"ValidateNotificationThreshold(%f) error = %v, wantErr %v",
					tc.threshold,
					err,
					tc.wantErr,
				)
			}
		})
	}
}

func TestValidateRequiredString(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid string", "test", false},
		{"valid string with spaces", "test string", false},
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"tabs only", "\t\t", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateRequiredString(tc.value, "field")
			if (err != nil) != tc.wantErr {
				t.Errorf(
					"ValidateRequiredString(%q) error = %v, wantErr %v",
					tc.value,
					err,
					tc.wantErr,
				)
			}
		})
	}
}

func TestValidateStringMaxLength(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		maxLength int
		wantErr   bool
	}{
		{"within limit", "test", 10, false},
		{"exact limit", "test", 4, false},
		{"empty string", "", 10, false},
		{"exceeds limit", "test string", 5, true},
		{"exceeds by one", "test", 3, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateStringMaxLength(tc.value, "field", tc.maxLength)
			if (err != nil) != tc.wantErr {
				t.Errorf(
					"ValidateStringMaxLength(%q, %d) error = %v, wantErr %v",
					tc.value,
					tc.maxLength,
					err,
					tc.wantErr,
				)
			}
		})
	}
}

func TestValidateExpenseType(t *testing.T) {
	testCases := []struct {
		name        string
		expenseType string
		wantErr     bool
	}{
		{"valid WEEKLY", "WEEKLY", false},
		{"valid MONTHLY", "MONTHLY", false},
		{"invalid lowercase weekly", "weekly", true},
		{"invalid lowercase monthly", "monthly", true},
		{"invalid DAILY", "DAILY", true},
		{"invalid empty", "", true},
		{"invalid random", "RANDOM", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateExpenseType(tc.expenseType)
			if (err != nil) != tc.wantErr {
				t.Errorf(
					"ValidateExpenseType(%q) error = %v, wantErr %v",
					tc.expenseType,
					err,
					tc.wantErr,
				)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	testCases := []struct {
		name    string
		size    int64
		wantErr bool
	}{
		{"valid 1 byte", 1, false},
		{"valid 1MB", 1 << 20, false},
		{"valid 10MB", 10 << 20, false},
		{"invalid 0 bytes", 0, true},
		{"invalid over 10MB", 10<<20 + 1, true},
		{"invalid 20MB", 20 << 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFileSize(tc.size)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateFileSize(%d) error = %v, wantErr %v", tc.size, err, tc.wantErr)
			}
		})
	}
}

func TestValidateBudgetCreate(t *testing.T) {
	testCases := []struct {
		name      string
		month     int
		year      int
		amount    float64
		threshold float64
		wantErr   bool
	}{
		{"valid all fields", 6, 2024, 1000.00, 0.8, false},
		{"valid edge month 1", 1, 2024, 1000.00, 0.8, false},
		{"valid edge month 12", 12, 2024, 1000.00, 0.8, false},
		{"valid edge year 2020", 6, 2020, 1000.00, 0.8, false},
		{"valid edge year 2100", 6, 2100, 1000.00, 0.8, false},
		{"valid threshold 0", 6, 2024, 1000.00, 0, false},
		{"valid threshold 1", 6, 2024, 1000.00, 1, false},
		{"invalid month 0", 0, 2024, 1000.00, 0.8, true},
		{"invalid month 13", 13, 2024, 1000.00, 0.8, true},
		{"invalid year 2019", 6, 2019, 1000.00, 0.8, true},
		{"invalid year 2101", 6, 2101, 1000.00, 0.8, true},
		{"invalid amount 0", 6, 2024, 0, 0.8, true},
		{"invalid amount negative", 6, 2024, -100.00, 0.8, true},
		{"invalid threshold negative", 6, 2024, 1000.00, -0.1, true},
		{"invalid threshold over 1", 6, 2024, 1000.00, 1.1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateBudgetCreate(tc.month, tc.year, tc.amount, tc.threshold)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateBudgetCreate(%d, %d, %f, %f) error = %v, wantErr %v",
					tc.month, tc.year, tc.amount, tc.threshold, err, tc.wantErr)
			}
		})
	}
}

func TestValidateExpenseCreate(t *testing.T) {
	itemCode := "CODE123"
	longItemCode := strings.Repeat("a", 51)

	testCases := []struct {
		name           string
		itemName       string
		source         string
		expenseType    string
		expectedAmount float64
		itemCode       *string
		wantErr        bool
	}{
		{"valid all fields", "Test Item", "Test Source", "WEEKLY", 100.00, &itemCode, false},
		{"valid without item code", "Test Item", "Test Source", "MONTHLY", 100.00, nil, false},
		{"valid zero amount", "Test Item", "Test Source", "WEEKLY", 0, nil, false},
		{"invalid empty item name", "", "Test Source", "WEEKLY", 100.00, nil, true},
		{"invalid empty source", "Test Item", "", "WEEKLY", 100.00, nil, true},
		{"invalid expense type", "Test Item", "Test Source", "DAILY", 100.00, nil, true},
		{"invalid negative amount", "Test Item", "Test Source", "WEEKLY", -100.00, nil, true},
		{
			"invalid long item code",
			"Test Item",
			"Test Source",
			"WEEKLY",
			100.00,
			&longItemCode,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateExpenseCreate(
				tc.itemName,
				tc.source,
				tc.expenseType,
				tc.expectedAmount,
				tc.itemCode,
			)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateExpenseCreate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "month", Message: "must be between 1 and 12"}
	expected := "month: must be between 1 and 12"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestValidationErrors(t *testing.T) {
	errs := &ValidationErrors{}
	errs.Add("month", "must be between 1 and 12")
	errs.Add("year", "must be between 2020 and 2100")

	if !errs.HasErrors() {
		t.Error("Expected HasErrors() to return true")
	}

	if len(errs.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errs.Errors))
	}

	errStr := errs.Error()
	if !strings.Contains(errStr, "month") || !strings.Contains(errStr, "year") {
		t.Errorf("Error string should contain both field names: %s", errStr)
	}
}

func TestValidationErrors_Empty(t *testing.T) {
	errs := &ValidationErrors{}

	if errs.HasErrors() {
		t.Error("Expected HasErrors() to return false for empty errors")
	}

	if errs.Error() != "validation failed" {
		t.Errorf("Expected 'validation failed', got %q", errs.Error())
	}
}
