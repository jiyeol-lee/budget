package models

// ReceiptItem represents an item extracted from a receipt
type ReceiptItem struct {
	Source    string  `json:"source"`
	Type      string  `json:"type"`
	ItemCode  string  `json:"item_code"`
	ItemPrice float64 `json:"item_price"`
	ItemName  string  `json:"item_name"`
}

// ProcessReceiptResponse represents the response for receipt processing
type ProcessReceiptResponse struct {
	Success          bool          `json:"success"`
	Items            []ReceiptItem `json:"items"`
	ProcessingTimeMs int64         `json:"processing_time_ms"`
}

// ProcessReceiptError represents an error response for receipt processing
type ProcessReceiptError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
}

// Error codes for receipt processing
const (
	ErrCodeTimeout         = "TIMEOUT"
	ErrCodeRateLimit       = "RATE_LIMIT"
	ErrCodeInvalidDocument = "INVALID_DOCUMENT"
	ErrCodeParseError      = "PARSE_ERROR"
	ErrCodeAPIError        = "API_ERROR"
	ErrCodeInternalError   = "INTERNAL_ERROR"
)
