package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// handleAPIError processes Anthropic SDK errors and maps them to appropriate error types
func handleAPIError(err error) error {
	var apiErr *anthropic.Error
	if errors.As(err, &apiErr) {
		// Log detailed error information
		fmt.Printf("Anthropic API Error: Status=%d\n", apiErr.StatusCode)
		fmt.Printf("Request: %s\n", string(apiErr.DumpRequest(true)))
		fmt.Printf("Response: %s\n", string(apiErr.DumpResponse(true)))

		// Map to appropriate error type based on status code
		switch apiErr.StatusCode {
		case 401:
			return fmt.Errorf(
				"%w: authentication failed - check ANTHROPIC_API_KEY",
				ErrAPIKeyNotSet,
			)
		case 429:
			return ErrRateLimit
		case 408, 504:
			return ErrTimeout
		case 503, 529:
			return ErrOverloaded
		default:
			return fmt.Errorf("%w: status %d - %v", ErrAPIError, apiErr.StatusCode, err)
		}
	}
	// For non-API errors (network issues, etc.)
	fmt.Printf("Non-API Error: %v\n", err)
	return fmt.Errorf("%w: %v", ErrAPIError, err)
}

// stripMarkdownCodeBlock removes markdown code block formatting from AI responses
// Handles responses wrapped in ```json ... ``` or ``` ... ```
func stripMarkdownCodeBlock(s string) string {
	s = strings.TrimSpace(s)

	// Check for ```json or ``` at the start
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}

	// Remove trailing ```
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}

	return strings.TrimSpace(s)
}

// Common errors
var (
	ErrAPIKeyNotSet    = errors.New("ANTHROPIC_API_KEY environment variable not set")
	ErrTimeout         = errors.New("API request timed out")
	ErrRateLimit       = errors.New("API rate limit exceeded")
	ErrInvalidDocument = errors.New("invalid document provided")
	ErrParseResponse   = errors.New("failed to parse AI response")
	ErrAPIError        = errors.New("API returned an error")
	ErrMaxRetries      = errors.New("max retries exceeded")
	ErrOverloaded      = errors.New("AI service is temporarily overloaded")
)

const (
	defaultMaxTokens = 8192
)

// Client represents the AI service client for receipt processing
type Client struct {
	client    anthropic.Client
	model     anthropic.Model
	maxTokens int
}

// Config holds AI client configuration
type Config struct {
	APIKey    string
	Model     string
	MaxTokens int
}

// RawReceiptItem represents an item extracted from OCR (uncategorized)
type RawReceiptItem struct {
	ItemCode  string  `json:"item_code"`
	ItemPrice float64 `json:"item_price"`
	ItemName  string  `json:"item_name"`
}

// OCRExtractionResult represents the output of OCR extraction
type OCRExtractionResult struct {
	Source    string           `json:"source"`
	Items     []RawReceiptItem `json:"items"`
	Total     float64          `json:"total"`
	Tax       float64          `json:"tax"`
	ItemCount int              `json:"item_count"`
}

// CategorizedItem represents an item with budget category assigned
type CategorizedItem struct {
	ItemCode  string  `json:"item_code"`
	ItemPrice float64 `json:"item_price"`
	ItemName  string  `json:"item_name"`
	ItemType  string  `json:"item_type"`
}

// CategorizationResult represents the output of categorization
type CategorizationResult struct {
	Items []CategorizedItem `json:"items"`
}

// ReceiptProcessingResult represents the combined OCR + categorization result
type ReceiptProcessingResult struct {
	Source    string            `json:"source"`
	Items     []CategorizedItem `json:"items"`
	Total     float64           `json:"total"`
	Tax       float64           `json:"tax"`
	ItemCount int               `json:"item_count"`
}

// NewClient creates a new AI service client
func NewClient(cfg Config) (*Client, error) {
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return nil, ErrAPIKeyNotSet
	}

	var model anthropic.Model
	if cfg.Model != "" {
		model = anthropic.Model(cfg.Model)
	} else {
		model = anthropic.ModelClaudeSonnet4_5
	}

	maxTokens := cfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	// Create the Anthropic client with the API key
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &Client{
		client:    client,
		model:     model,
		maxTokens: maxTokens,
	}, nil
}

// NewClientFromEnv creates a new AI service client using environment variables
func NewClientFromEnv() (*Client, error) {
	return NewClient(Config{})
}

// AnalyzeDocument sends a PDF document with a prompt to the AI and returns the response
// Only PDF format (application/pdf) is supported
func (c *Client) AnalyzeDocument(
	ctx context.Context,
	base64Data, mimeType, prompt string,
) (string, error) {
	// Only PDF is supported
	if mimeType != "application/pdf" {
		return "", fmt.Errorf("%w: unsupported mime type: %s (only application/pdf is supported)", ErrInvalidDocument, mimeType)
	}

	contentBlock := anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{
		Type:      "base64",
		MediaType: "application/pdf",
		Data:      base64Data,
	})

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: int64(c.maxTokens),
		Model:     c.model,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					contentBlock,
					anthropic.NewTextBlock(prompt),
				},
			},
		},
	})
	if err != nil {
		return "", handleAPIError(err)
	}

	// Extract response text from content
	for _, block := range message.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("%w: no text in response content", ErrParseResponse)
}

// SendTextPrompt sends a text-only prompt to the AI and returns the response
func (c *Client) SendTextPrompt(ctx context.Context, prompt string) (string, error) {
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: int64(c.maxTokens),
		Model:     c.model,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					anthropic.NewTextBlock(prompt),
				},
			},
		},
	})
	if err != nil {
		return "", handleAPIError(err)
	}

	// Extract response text from content
	for _, block := range message.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("%w: no text in response content", ErrParseResponse)
}

// ReceiptProcessingPrompt returns the prompt for combined OCR extraction and categorization
func ReceiptProcessingPrompt(budgets []string) string {
	budgetList := "None"
	if len(budgets) > 0 {
		budgetList = strings.Join(budgets, ", ")
	}

	return fmt.Sprintf(
		`You are a precise receipt OCR and categorization system. Extract ALL items from the receipt and categorize them.

=== CRITICAL REQUIREMENTS ===
*** EXTRACT EVERY SINGLE ITEM - No omissions allowed ***
*** COPY ITEM CODES EXACTLY AS PRINTED - Do not modify or abbreviate ***
*** COPY PRICES EXACTLY AS SHOWN - Do not round or modify ***
*** PRESERVE ORDER - Extract items from top to bottom ***

=== EXTRACTION RULES ===
For EACH line item on the receipt, extract:
1. item_code: The abbreviated item name or SKU code as printed on the receipt (e.g., "ORG BANAN", "MLK 2%%", "BROC CRN"). Copy it EXACTLY as shown. If not visible, use "N/A"
2. item_price: The EXACT price as a decimal number. If the price has a minus sign (-) or is marked as a refund/discount/credit, use a NEGATIVE number (e.g., -5.99)
3. item_name: Your best guess of the full item name based on the item_code abbreviation (e.g., "ORG BANAN" → "Organic Bananas", "MLK 2%%" → "2%% Milk", "BROC CRN" → "Broccoli Crown")
4. item_type: Categorize based on rules below

Also extract:
- source: Store name from receipt header (use "Unknown" if not visible)
- total: The total amount shown on receipt
- tax: The tax amount (0 if not shown)
- item_count: Total number of items extracted

=== CATEGORIZATION RULES ===
Budget Categories: %s

1. Compare each item against the Budget Categories list
2. If item matches a category, assign the type in parentheses (e.g., "Apple (monthly)" → "monthly")
3. Types must be lowercase: "weekly", "monthly", "misc", or "tax"
4. If item does NOT match any category, assign "misc"
5. If item_name contains "tax", "TAX", "HST", "GST", "VAT", assign "tax"
6. Do NOT guess - only match against provided categories
7. *** TAX LINE ITEMS ARE MANDATORY ***: If you see ANY tax line (sales tax, VAT, HST, GST, PST, etc.) on the receipt, you MUST extract it as a separate item with item_type "tax", item_code "TAX", and item_name "Tax"

=== OUTPUT FORMAT ===
IMPORTANT: Return ONLY the raw JSON object, nothing else.
- NO markdown formatting
- NO code blocks (no `+"`"+``+"`"+``+"`"+` before or after)
- NO explanatory text
- Start your response with { and end with }

{
  "source": "Store Name",
  "item_count": 0,
  "total": 0.00,
  "tax": 0.00,
  "items": [
    {
      "item_code": "EXACT_CODE",
      "item_price": EXACT_PRICE,
      "item_name": "Item Name",
      "item_type": "weekly|monthly|misc|tax"
    }
  ]
}

=== WARNINGS ===
- EVERY line item must be extracted
- Item codes must be EXACTLY as printed
- Prices must be EXACTLY as shown
- Items must be in receipt order (top to bottom)
- Return ONLY raw JSON, absolutely NO markdown formatting or code blocks
- *** ALWAYS EXTRACT TAX *** - Tax line items must be included with item_type "tax"
- *** PRESERVE NEGATIVE PRICES *** - Refunds, discounts, and credits should be negative numbers (e.g., -5.99)`,
		budgetList,
	)
}

// Deprecated: Use ReceiptProcessingPrompt and ProcessReceiptDocument instead
// OCRExtractionPrompt returns the prompt for pure OCR extraction (no categorization)
func OCRExtractionPrompt() string {
	return `You are a precise receipt OCR system. Your ONLY task is to extract data exactly as printed on the receipt.

=== CRITICAL REQUIREMENTS ===
*** EXTRACT EVERY SINGLE ITEM - No omissions allowed ***
*** COPY ITEM CODES EXACTLY AS PRINTED - Do not modify or abbreviate ***
*** COPY PRICES EXACTLY AS SHOWN - Do not round or modify ***
*** PRESERVE ORDER - Extract items from top to bottom ***
*** DO NOT CATEGORIZE - This is extraction only ***

=== EXTRACTION RULES ===
For EACH line item on the receipt, extract:
1. item_code: The EXACT code/SKU as printed (if not visible, use "N/A")
2. item_price: The EXACT price as a decimal number
3. item_name: Your best interpretation of the item name

Also extract:
- source: Store name from receipt header (use "Unknown" if not visible)
- total: The total amount shown on receipt
- tax: The tax amount (0 if not shown)
- item_count: Total number of items extracted

=== OUTPUT FORMAT ===
CRITICAL: Return ONLY raw JSON. Do NOT wrap in markdown code blocks.
Do NOT use ` + "`" + `` + "`" + `` + "`" + `json or ` + "`" + `` + "`" + `` + "`" + ` - just return the raw JSON object.

{
  "source": "Store Name",
  "item_count": 0,
  "total": 0.00,
  "tax": 0.00,
  "items": [
    {
      "item_code": "EXACT_CODE",
      "item_price": 0.00,
      "item_name": "Item Name"
    }
  ]
}

=== WARNINGS ===
- EVERY line item must be extracted
- Item codes must be EXACTLY as printed
- Prices must be EXACTLY as shown
- Items must be in receipt order (top to bottom)
- Return ONLY raw JSON, absolutely NO markdown formatting or code blocks`
}

// Deprecated: Use ReceiptProcessingPrompt and ProcessReceiptDocument instead
// CategorizationPrompt returns the prompt for categorizing extracted items
func CategorizationPrompt(itemsJSON string, budgets []string) string {
	budgetList := "None"
	if len(budgets) > 0 {
		budgetList = strings.Join(budgets, ", ")
	}

	return fmt.Sprintf(
		`You are a budget categorization system. Categorize each item based on the budget categories provided.

=== INPUT ===
Extracted items: %s

Budget Categories: %s

=== CATEGORIZATION RULES ===
1. Compare each item against the Budget Categories list
2. If item matches a category, assign the type in parentheses (e.g., "Apple (monthly)" → "monthly")
3. Types must be lowercase: "weekly", "monthly", "misc", or "tax"
4. If item does NOT match any category, assign "misc"
5. If item_code contains "tax", "TAX", "HST", "GST", "VAT", assign "tax"
6. Do NOT guess - only match against provided categories

=== OUTPUT FORMAT ===
CRITICAL: Return ONLY raw JSON. Do NOT wrap in markdown code blocks.
Do NOT use `+"`"+``+"`"+``+"`"+`json or `+"`"+``+"`"+``+"`"+` - just return the raw JSON object:
{
  "items": [
    {
      "item_code": "BB GARLIC HOT",
      "item_price": 0.00,
      "item_name": "Name",
      "item_type": "weekly|monthly|misc|tax"
    }
  ]
}

=== WARNINGS ===
- Preserve the exact item_code, item_price from input
- Only add item_type based on categorization rules
- Maintain the same order as input
- Return ONLY raw JSON, absolutely NO markdown formatting or code blocks`,
		itemsJSON,
		budgetList,
	)
}

// Legacy types for backward compatibility
type ReceiptData struct {
	Items []ReceiptItem `json:"items"`
	Total float64       `json:"total"`
	Date  string        `json:"date"`
	Store string        `json:"store"`
}

type ReceiptItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// ProcessReceipt is deprecated - use AnalyzeDocument instead
func (c *Client) ProcessReceipt(imageData []byte) (*ReceiptData, error) {
	return nil, errors.New("use AnalyzeDocument instead")
}

// ProcessReceiptDocument performs OCR extraction and categorization on a PDF receipt in a single AI request
// Only PDF format (application/pdf) is supported
func (c *Client) ProcessReceiptDocument(
	ctx context.Context,
	base64Data, mimeType string,
	budgets []string,
) (*ReceiptProcessingResult, error) {
	prompt := ReceiptProcessingPrompt(budgets)

	responseText, err := c.AnalyzeDocument(ctx, base64Data, mimeType, prompt)
	if err != nil {
		return nil, fmt.Errorf("receipt processing failed: %w", err)
	}

	// Strip any markdown code block formatting from the response
	responseText = stripMarkdownCodeBlock(responseText)

	var result ReceiptProcessingResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf(
			"%w: failed to parse result: %v\nResponse was: %s",
			ErrParseResponse,
			err,
			responseText,
		)
	}

	return &result, nil
}

// ProcessReceiptImage is deprecated, use ProcessReceiptDocument instead
func (c *Client) ProcessReceiptImage(
	ctx context.Context,
	base64Data, mimeType string,
	budgets []string,
) (*ReceiptProcessingResult, error) {
	return c.ProcessReceiptDocument(ctx, base64Data, mimeType, budgets)
}

// Deprecated: Use ProcessReceiptDocument instead
// ExtractReceiptItems performs OCR extraction on a PDF receipt document (Step 1)
// Only PDF format (application/pdf) is supported
func (c *Client) ExtractReceiptItems(
	ctx context.Context,
	base64Data, mimeType string,
) (*OCRExtractionResult, error) {
	prompt := OCRExtractionPrompt()

	responseText, err := c.AnalyzeDocument(ctx, base64Data, mimeType, prompt)
	if err != nil {
		return nil, fmt.Errorf("document extraction failed: %w", err)
	}

	var result OCRExtractionResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf(
			"%w: failed to parse OCR result: %v\nResponse was: %s",
			ErrParseResponse,
			err,
			responseText,
		)
	}

	return &result, nil
}

// Deprecated: Use ProcessReceiptDocument instead
// CategorizeItems categorizes extracted items against budget categories (Step 2)
func (c *Client) CategorizeItems(
	ctx context.Context,
	items []RawReceiptItem,
	budgets []string,
) (*CategorizationResult, error) {
	// Convert items to JSON for the prompt
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}

	prompt := CategorizationPrompt(string(itemsJSON), budgets)

	responseText, err := c.SendTextPrompt(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("categorization failed: %w", err)
	}

	var result CategorizationResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf(
			"%w: failed to parse categorization result: %v",
			ErrParseResponse,
			err,
		)
	}

	return &result, nil
}
