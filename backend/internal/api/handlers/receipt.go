package handlers

import (
	"budget-tracker/internal/models"
	"budget-tracker/internal/repository"
	"budget-tracker/internal/services/ai"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// MaxUploadSize is the maximum file size for receipt uploads (10MB)
	MaxUploadSize = 10 << 20 // 10 MB
	// FormFileKey is the key for the document file in the multipart form
	FormFileKey = "document"
)

// ReceiptHandler handles receipt-related HTTP requests
type ReceiptHandler struct {
	aiClient            *ai.Client
	documentProcessor   *ai.ImageProcessor
	expectedExpenseRepo *repository.ExpectedExpenseRepository
	actualExpenseRepo   *repository.ActualExpenseRepository
}

// NewReceiptHandler creates a new ReceiptHandler
func NewReceiptHandler(
	aiClient *ai.Client,
	expectedExpenseRepo *repository.ExpectedExpenseRepository,
	actualExpenseRepo *repository.ActualExpenseRepository,
) *ReceiptHandler {
	return &ReceiptHandler{
		aiClient:            aiClient,
		documentProcessor:   ai.NewImageProcessor(),
		expectedExpenseRepo: expectedExpenseRepo,
		actualExpenseRepo:   actualExpenseRepo,
	}
}

// Process handles POST /api/receipts/process
// Accepts multipart form data with a PDF document and returns extracted receipt items
func (h *ReceiptHandler) Process(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in ReceiptHandler: %v\n", r)
			h.respondReceiptError(
				w,
				http.StatusInternalServerError,
				"Internal server error during processing",
				models.ErrCodeInternalError,
			)
		}
	}()

	startTime := time.Now()
	fmt.Printf("[Receipt] Starting receipt processing\n")

	// Check if AI client is configured
	if h.aiClient == nil {
		h.respondReceiptError(
			w,
			http.StatusServiceUnavailable,
			"AI service not configured",
			models.ErrCodeInternalError,
		)
		return
	}

	// Limit the request body size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	// Parse the multipart form
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			h.respondReceiptError(
				w,
				http.StatusRequestEntityTooLarge,
				"PDF file too large (max 10MB)",
				models.ErrCodeInvalidDocument,
			)
			return
		}
		h.respondReceiptError(
			w,
			http.StatusBadRequest,
			"Failed to parse form data",
			models.ErrCodeInvalidDocument,
		)
		return
	}

	// Get the uploaded file
	file, header, err := r.FormFile(FormFileKey)
	if err != nil {
		h.respondReceiptError(
			w,
			http.StatusBadRequest,
			"No document file provided. Use form field 'image'",
			models.ErrCodeInvalidDocument,
		)
		return
	}
	defer file.Close()
	fmt.Printf("[Receipt] File received: name=%s, size=%d bytes\n", header.Filename, header.Size)

	// Validate file size
	if header.Size == 0 {
		h.respondReceiptError(
			w,
			http.StatusBadRequest,
			"Empty document file",
			models.ErrCodeInvalidDocument,
		)
		return
	}

	// Process the document
	processedDocument, err := h.documentProcessor.ReadAndProcessReader(file)
	if err != nil {
		if errors.Is(err, ai.ErrUnsupportedFormat) {
			h.respondReceiptError(
				w,
				http.StatusBadRequest,
				"Unsupported format. Only PDF is supported",
				models.ErrCodeInvalidDocument,
			)
			return
		}
		h.respondReceiptError(
			w,
			http.StatusBadRequest,
			"Failed to process document",
			models.ErrCodeInvalidDocument,
		)
		return
	}

	fmt.Printf("[Receipt] Document processed: mimeType=%s, dataLength=%d\n", processedDocument.MimeType, len(processedDocument.Base64Data))

	// Call the AI service with context timeout
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	// Fetch existing expected expenses to build budget categories for AI categorization
	var budgetCategories []string
	if h.expectedExpenseRepo != nil {
		expenses, err := h.expectedExpenseRepo.GetAll()
		if err == nil {
			// Build unique category list from expense item names
			categoryMap := make(map[string]bool)
			for _, expense := range expenses {
				if !categoryMap[expense.ItemName] {
					categoryMap[expense.ItemName] = true
					// Include the type information for better AI categorization
					categoryInfo := expense.ItemName + " (" + string(expense.ExpenseType) + ")"
					budgetCategories = append(budgetCategories, categoryInfo)
				}
			}
		}
	}

	fmt.Printf("[Receipt] Calling AI service with %d budget categories\n", len(budgetCategories))

	// Process receipt: OCR extraction + categorization in one request
	result, err := h.aiClient.ProcessReceiptImage(
		ctx,
		processedDocument.Base64Data,
		processedDocument.MimeType,
		budgetCategories,
	)
	if err != nil {
		h.handleAIError(w, err)
		return
	}

	// Calculate processing time
	processingTimeMs := time.Since(startTime).Milliseconds()

	// Get source from result
	source := result.Source
	if source == "" {
		source = "Unknown"
	}

	// Prepare the response items from result
	responseItems := make([]models.ReceiptItem, len(result.Items))
	for i, item := range result.Items {
		itemType := item.ItemType
		if itemType == "" {
			itemType = "misc"
		}
		responseItems[i] = models.ReceiptItem{
			Source:    source,
			Type:      itemType,
			ItemCode:  item.ItemCode,
			ItemPrice: item.ItemPrice,
			ItemName:  item.ItemName,
		}
	}

	fmt.Printf("[Receipt] Success: extracted %d items in %dms\n", len(responseItems), processingTimeMs)

	// Return the response
	respondJSON(w, http.StatusOK, models.ProcessReceiptResponse{
		Success:          true,
		Items:            responseItems,
		ProcessingTimeMs: processingTimeMs,
	})
}

// handleAIError handles errors from the AI service and returns appropriate responses
func (h *ReceiptHandler) handleAIError(w http.ResponseWriter, err error) {
	fmt.Printf("[Receipt] AI Error: %v\n", err)
	switch {
	case errors.Is(err, ai.ErrTimeout):
		h.respondReceiptError(
			w,
			http.StatusGatewayTimeout,
			"Receipt processing timed out. Please try again",
			models.ErrCodeTimeout,
		)
	case errors.Is(err, ai.ErrRateLimit):
		h.respondReceiptError(
			w,
			http.StatusTooManyRequests,
			"Service is busy. Please try again in a moment",
			models.ErrCodeRateLimit,
		)
	case errors.Is(err, ai.ErrOverloaded):
		h.respondReceiptError(
			w,
			http.StatusServiceUnavailable,
			"AI service is temporarily overloaded. Please try again in a few moments",
			models.ErrCodeAPIError,
		)
	case errors.Is(err, ai.ErrAPIKeyNotSet):
		h.respondReceiptError(
			w,
			http.StatusServiceUnavailable,
			"AI service not configured",
			models.ErrCodeInternalError,
		)
	case errors.Is(err, ai.ErrMaxRetries):
		h.respondReceiptError(
			w,
			http.StatusServiceUnavailable,
			"Failed to process receipt after multiple attempts",
			models.ErrCodeAPIError,
		)
	case errors.Is(err, ai.ErrAPIError):
		h.respondReceiptError(
			w,
			http.StatusBadGateway,
			"AI service error. Please try again",
			models.ErrCodeAPIError,
		)
	case errors.Is(err, context.DeadlineExceeded):
		h.respondReceiptError(
			w,
			http.StatusGatewayTimeout,
			"Request timed out",
			models.ErrCodeTimeout,
		)
	case errors.Is(err, context.Canceled):
		h.respondReceiptError(
			w,
			http.StatusRequestTimeout,
			"Request was canceled",
			models.ErrCodeTimeout,
		)
	default:
		h.respondReceiptError(
			w,
			http.StatusInternalServerError,
			"Failed to process receipt",
			models.ErrCodeInternalError,
		)
	}
}

// respondReceiptError sends an error response for receipt processing
func (h *ReceiptHandler) respondReceiptError(
	w http.ResponseWriter,
	status int,
	message string,
	code string,
) {
	fmt.Printf("[Receipt] Error Response: status=%d, code=%s, message=%s\n", status, code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ProcessReceiptError{
		Success: false,
		Error:   message,
		Code:    code,
	})
}
