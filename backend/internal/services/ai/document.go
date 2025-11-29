package ai

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// PDF processing errors
var (
	ErrUnsupportedFormat = errors.New("unsupported format: only PDF is supported")
	ErrReadFile          = errors.New("failed to read file")
)

// PDFProcessor handles PDF processing operations
type PDFProcessor struct{}

// NewPDFProcessor creates a new PDFProcessor
func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{}
}

// ProcessedDocument represents a document ready for AI processing
type ProcessedDocument struct {
	Base64Data string
	MimeType   string
}

// ReadAndProcessFile reads a PDF file and processes it for AI analysis
func (p *PDFProcessor) ReadAndProcessFile(filePath string) (*ProcessedDocument, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFile, err)
	}
	defer f.Close()

	return p.ProcessDocument(f)
}

// ReadAndProcessReader reads a PDF from an io.ReadSeeker and processes it for AI analysis
func (p *PDFProcessor) ReadAndProcessReader(r io.ReadSeeker) (*ProcessedDocument, error) {
	return p.ProcessDocument(r)
}

// ProcessDocument processes a PDF for AI analysis
func (p *PDFProcessor) ProcessDocument(r io.ReadSeeker) (*ProcessedDocument, error) {
	// Read first bytes to detect format
	header := make([]byte, 8)
	if _, err := r.Read(header); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFile, err)
	}

	// Rewind
	if _, err := r.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek stream: %w", err)
	}

	// Check if it's a PDF
	if string(header[:4]) != "%PDF" {
		return nil, ErrUnsupportedFormat
	}

	return p.ProcessPDF(r)
}

// ProcessPDF reads a PDF and returns base64 encoded data
func (p *PDFProcessor) ProcessPDF(r io.Reader) (*ProcessedDocument, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFile, err)
	}

	// Validate it's a PDF
	if len(data) < 4 || string(data[:4]) != "%PDF" {
		return nil, ErrUnsupportedFormat
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	return &ProcessedDocument{
		Base64Data: base64Data,
		MimeType:   "application/pdf",
	}, nil
}

// ValidateFormat validates the format of file data (PDF only)
func (p *PDFProcessor) ValidateFormat(data []byte) (string, error) {
	// Check PDF magic bytes
	if len(data) >= 4 && string(data[:4]) == "%PDF" {
		return "application/pdf", nil
	}

	return "", ErrUnsupportedFormat
}
