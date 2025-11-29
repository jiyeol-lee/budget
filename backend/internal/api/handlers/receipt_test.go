package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/services/ai"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test data for PDF and image formats
var (
	// Valid PDF header - minimal valid PDF structure
	testValidPDFData = []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\n%%EOF")
	// JPEG magic bytes (FFD8FF)
	testJPEGData = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00}
	// PNG magic bytes (89 50 4E 47 0D 0A 1A 0A)
	testPNGData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
)

// createTestReceiptMux creates a router with the receipt handler for testing
func createTestReceiptMux(receiptHandler *ReceiptHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/receipts/process", receiptHandler.Process)
	return mux
}

// createMultipartRequest creates a multipart form request with the given file data
func createMultipartRequest(t *testing.T, fieldName string, fileName string, fileData []byte) (*http.Request, error) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/api/receipts/process", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

// TestReceiptHandler_FormFieldKeyConstant verifies the form field key is "document"
func TestReceiptHandler_FormFieldKeyConstant(t *testing.T) {
	if FormFileKey != "document" {
		t.Errorf("Expected FormFileKey to be 'document', got '%s'", FormFileKey)
	}
}

// TestReceiptHandler_MaxUploadSizeConstant verifies MaxUploadSize is 10MB
func TestReceiptHandler_MaxUploadSizeConstant(t *testing.T) {
	expectedSize := int64(10 << 20) // 10 MB
	if MaxUploadSize != expectedSize {
		t.Errorf("Expected MaxUploadSize to be %d bytes, got %d", expectedSize, MaxUploadSize)
	}
}

// TestReceiptHandler_ErrorCodes verifies all expected error codes are defined correctly
func TestReceiptHandler_ErrorCodes(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected string
	}{
		{"INVALID_DOCUMENT", models.ErrCodeInvalidDocument, "INVALID_DOCUMENT"},
		{"TIMEOUT", models.ErrCodeTimeout, "TIMEOUT"},
		{"RATE_LIMIT", models.ErrCodeRateLimit, "RATE_LIMIT"},
		{"API_ERROR", models.ErrCodeAPIError, "API_ERROR"},
		{"INTERNAL_ERROR", models.ErrCodeInternalError, "INTERNAL_ERROR"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.code != tc.expected {
				t.Errorf("Expected error code '%s', got '%s'", tc.expected, tc.code)
			}
		})
	}
}

// TestReceiptHandler_ServiceUnavailableWithoutAIClient verifies the handler
// returns 503 Service Unavailable when no AI client is configured
func TestReceiptHandler_ServiceUnavailableWithoutAIClient(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Handler without AI client
	handler := NewReceiptHandler(nil, nil, nil)
	mux := createTestReceiptMux(handler)

	// Upload valid PDF
	req, err := createMultipartRequest(t, FormFileKey, "test.pdf", testValidPDFData)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Without AI client, should return service unavailable
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d without AI client, got %d", http.StatusServiceUnavailable, rec.Code)
	}

	var errResp models.ProcessReceiptError
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Success {
		t.Error("Expected success to be false")
	}

	if errResp.Code != models.ErrCodeInternalError {
		t.Errorf("Expected error code '%s', got '%s'", models.ErrCodeInternalError, errResp.Code)
	}
}

// TestReceiptHandler_ErrorResponseStructure verifies the error response has the correct structure
func TestReceiptHandler_ErrorResponseStructure(t *testing.T) {
	handler := NewReceiptHandler(nil, nil, nil)
	mux := createTestReceiptMux(handler)

	// Create request with no file to trigger error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/receipts/process", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Verify the error response structure
	var errResp models.ProcessReceiptError
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	// Check all required fields are present
	if errResp.Error == "" {
		t.Error("Expected non-empty error message")
	}

	if errResp.Code == "" {
		t.Error("Expected non-empty error code")
	}

	// Success should be false for error responses
	if errResp.Success {
		t.Error("Expected success to be false for error response")
	}
}

// TestReceiptHandler_DocumentProcessorIntegration tests that the document
// processor correctly handles PDF and rejects images
// This tests the integration between the handler and the PDFProcessor
func TestReceiptHandler_DocumentProcessorIntegration(t *testing.T) {
	processor := ai.NewPDFProcessor()

	t.Run("accepts PDF", func(t *testing.T) {
		result, err := processor.ProcessDocument(bytes.NewReader(testValidPDFData))
		if err != nil {
			t.Fatalf("Expected PDF to be accepted, got error: %v", err)
		}
		if result.MimeType != "application/pdf" {
			t.Errorf("Expected MimeType 'application/pdf', got '%s'", result.MimeType)
		}
	})

	t.Run("rejects JPEG with INVALID_DOCUMENT", func(t *testing.T) {
		_, err := processor.ProcessDocument(bytes.NewReader(testJPEGData))
		if err == nil {
			t.Fatal("Expected JPEG to be rejected")
		}
		if err != ai.ErrUnsupportedFormat {
			t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
		}
	})

	t.Run("rejects PNG with INVALID_DOCUMENT", func(t *testing.T) {
		_, err := processor.ProcessDocument(bytes.NewReader(testPNGData))
		if err == nil {
			t.Fatal("Expected PNG to be rejected")
		}
		if err != ai.ErrUnsupportedFormat {
			t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
		}
	})
}

// TestReceiptHandler_NewReceiptHandler verifies the handler is created correctly
func TestReceiptHandler_NewReceiptHandler(t *testing.T) {
	handler := NewReceiptHandler(nil, nil, nil)

	if handler == nil {
		t.Fatal("Expected non-nil handler")
	}

	// Verify the document processor is initialized
	if handler.documentProcessor == nil {
		t.Error("Expected documentProcessor to be initialized")
	}
}

// TestReceiptHandler_ProcessReceiptResponseStructure verifies response models are correct
func TestReceiptHandler_ProcessReceiptResponseStructure(t *testing.T) {
	// Test ProcessReceiptResponse
	response := models.ProcessReceiptResponse{
		Success: true,
		Items: []models.ReceiptItem{
			{
				Source:    "Test Store",
				Type:      "weekly",
				ItemCode:  "ABC123",
				ItemPrice: 19.99,
				ItemName:  "Test Item",
			},
		},
		ProcessingTimeMs: 150,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var decoded models.ProcessReceiptResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if decoded.Success != true {
		t.Error("Expected Success to be true")
	}

	if len(decoded.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(decoded.Items))
	}

	if decoded.Items[0].ItemName != "Test Item" {
		t.Errorf("Expected ItemName 'Test Item', got '%s'", decoded.Items[0].ItemName)
	}
}

// TestReceiptHandler_ProcessReceiptErrorStructure verifies error response models are correct
func TestReceiptHandler_ProcessReceiptErrorStructure(t *testing.T) {
	// Test ProcessReceiptError
	errResponse := models.ProcessReceiptError{
		Success: false,
		Error:   "Test error message",
		Code:    models.ErrCodeInvalidDocument,
	}

	data, err := json.Marshal(errResponse)
	if err != nil {
		t.Fatalf("Failed to marshal error response: %v", err)
	}

	var decoded models.ProcessReceiptError
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if decoded.Success != false {
		t.Error("Expected Success to be false")
	}

	if decoded.Error != "Test error message" {
		t.Errorf("Expected Error 'Test error message', got '%s'", decoded.Error)
	}

	if decoded.Code != models.ErrCodeInvalidDocument {
		t.Errorf("Expected Code '%s', got '%s'", models.ErrCodeInvalidDocument, decoded.Code)
	}
}

// TestReceiptHandler_ReceiptItemStructure verifies the ReceiptItem model structure
func TestReceiptHandler_ReceiptItemStructure(t *testing.T) {
	item := models.ReceiptItem{
		Source:    "Costco",
		Type:      "monthly",
		ItemCode:  "1234567890",
		ItemPrice: 25.99,
		ItemName:  "Bulk Paper Towels",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal ReceiptItem: %v", err)
	}

	// Verify JSON keys are snake_case
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	expectedKeys := []string{"source", "type", "item_code", "item_price", "item_name"}
	for _, key := range expectedKeys {
		if _, ok := m[key]; !ok {
			t.Errorf("Expected JSON key '%s' to be present", key)
		}
	}
}

// TestReceiptHandler_UnsupportedFormatErrorMessage verifies the error message is PDF-centric
func TestReceiptHandler_UnsupportedFormatErrorMessage(t *testing.T) {
	expectedSubstring := "only PDF is supported"

	if errMsg := ai.ErrUnsupportedFormat.Error(); errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	if !bytes.Contains([]byte(ai.ErrUnsupportedFormat.Error()), []byte(expectedSubstring)) {
		t.Errorf("Expected error message to contain '%s', got: '%s'", expectedSubstring, ai.ErrUnsupportedFormat.Error())
	}
}

// TestReceiptHandler_FormFileKeyIsDocument verifies the correct form field name
// This is a critical test because the frontend must use this exact field name
func TestReceiptHandler_FormFileKeyIsDocument(t *testing.T) {
	// This test documents the API contract: the form field must be "document"
	if FormFileKey != "document" {
		t.Errorf("API contract violation: FormFileKey must be 'document', got '%s'", FormFileKey)
	}
}
