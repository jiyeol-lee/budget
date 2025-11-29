package ai

import (
	"bytes"
	"strings"
	"testing"
)

// Test PDF magic bytes for creating test data
var (
	// Valid PDF header
	validPDFData = []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\n%%EOF")
	// JPEG magic bytes (FFD8FF)
	jpegData = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	// PNG magic bytes (89 50 4E 47 0D 0A 1A 0A)
	pngData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
)

func TestPDFProcessor_AcceptsPDF(t *testing.T) {
	processor := NewPDFProcessor()

	result, err := processor.ProcessDocument(bytes.NewReader(validPDFData))
	if err != nil {
		t.Fatalf("Expected PDF to be accepted, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result for valid PDF")
	}

	if result.MimeType != "application/pdf" {
		t.Errorf("Expected MimeType 'application/pdf', got '%s'", result.MimeType)
	}

	if result.Base64Data == "" {
		t.Error("Expected non-empty Base64Data")
	}
}

func TestPDFProcessor_RejectsJPEG(t *testing.T) {
	processor := NewPDFProcessor()

	result, err := processor.ProcessDocument(bytes.NewReader(jpegData))

	if err == nil {
		t.Fatal("Expected JPEG to be rejected, got no error")
	}

	if result != nil {
		t.Error("Expected nil result for rejected JPEG")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_RejectsPNG(t *testing.T) {
	processor := NewPDFProcessor()

	result, err := processor.ProcessDocument(bytes.NewReader(pngData))

	if err == nil {
		t.Fatal("Expected PNG to be rejected, got no error")
	}

	if result != nil {
		t.Error("Expected nil result for rejected PNG")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_ErrorMessage(t *testing.T) {
	// Verify the error message is PDF-centric
	expectedSubstring := "only PDF is supported"

	if !strings.Contains(ErrUnsupportedFormat.Error(), expectedSubstring) {
		t.Errorf(
			"Expected error message to contain '%s', got: '%s'",
			expectedSubstring,
			ErrUnsupportedFormat.Error(),
		)
	}
}

func TestPDFProcessor_ValidateFormat_AcceptsPDF(t *testing.T) {
	processor := NewPDFProcessor()

	mimeType, err := processor.ValidateFormat(validPDFData)
	if err != nil {
		t.Fatalf("Expected PDF to be accepted by ValidateFormat, got error: %v", err)
	}

	if mimeType != "application/pdf" {
		t.Errorf("Expected MimeType 'application/pdf', got '%s'", mimeType)
	}
}

func TestPDFProcessor_ValidateFormat_RejectsJPEG(t *testing.T) {
	processor := NewPDFProcessor()

	mimeType, err := processor.ValidateFormat(jpegData)

	if err == nil {
		t.Fatal("Expected JPEG to be rejected by ValidateFormat")
	}

	if mimeType != "" {
		t.Errorf("Expected empty MimeType for rejected format, got '%s'", mimeType)
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_ValidateFormat_RejectsPNG(t *testing.T) {
	processor := NewPDFProcessor()

	mimeType, err := processor.ValidateFormat(pngData)

	if err == nil {
		t.Fatal("Expected PNG to be rejected by ValidateFormat")
	}

	if mimeType != "" {
		t.Errorf("Expected empty MimeType for rejected format, got '%s'", mimeType)
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_ValidateFormat_RejectsEmptyData(t *testing.T) {
	processor := NewPDFProcessor()

	mimeType, err := processor.ValidateFormat([]byte{})

	if err == nil {
		t.Fatal("Expected empty data to be rejected by ValidateFormat")
	}

	if mimeType != "" {
		t.Errorf("Expected empty MimeType for rejected format, got '%s'", mimeType)
	}
}

func TestPDFProcessor_ValidateFormat_RejectsShortData(t *testing.T) {
	processor := NewPDFProcessor()

	// Data shorter than PDF magic bytes
	shortData := []byte{0x25, 0x50} // "%P" - incomplete PDF header

	mimeType, err := processor.ValidateFormat(shortData)

	if err == nil {
		t.Fatal("Expected short data to be rejected by ValidateFormat")
	}

	if mimeType != "" {
		t.Errorf("Expected empty MimeType for rejected format, got '%s'", mimeType)
	}
}

func TestPDFProcessor_ProcessPDF_RejectsNonPDF(t *testing.T) {
	processor := NewPDFProcessor()

	// ProcessPDF expects the data to already be validated as PDF
	// but it has a secondary check
	result, err := processor.ProcessPDF(bytes.NewReader(jpegData))

	if err == nil {
		t.Fatal("Expected ProcessPDF to reject non-PDF data")
	}

	if result != nil {
		t.Error("Expected nil result for rejected data")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_ReadAndProcessReader_AcceptsPDF(t *testing.T) {
	processor := NewPDFProcessor()

	result, err := processor.ReadAndProcessReader(bytes.NewReader(validPDFData))
	if err != nil {
		t.Fatalf("Expected PDF to be accepted, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result for valid PDF")
	}

	if result.MimeType != "application/pdf" {
		t.Errorf("Expected MimeType 'application/pdf', got '%s'", result.MimeType)
	}
}

func TestPDFProcessor_ReadAndProcessReader_RejectsJPEG(t *testing.T) {
	processor := NewPDFProcessor()

	result, err := processor.ReadAndProcessReader(bytes.NewReader(jpegData))

	if err == nil {
		t.Fatal("Expected JPEG to be rejected")
	}

	if result != nil {
		t.Error("Expected nil result for rejected JPEG")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_RejectsGIF(t *testing.T) {
	processor := NewPDFProcessor()

	// GIF magic bytes (GIF89a or GIF87a)
	gifData := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}

	result, err := processor.ProcessDocument(bytes.NewReader(gifData))

	if err == nil {
		t.Fatal("Expected GIF to be rejected, got no error")
	}

	if result != nil {
		t.Error("Expected nil result for rejected GIF")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_RejectsWebP(t *testing.T) {
	processor := NewPDFProcessor()

	// WebP magic bytes (RIFF....WEBP)
	webpData := []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}

	result, err := processor.ProcessDocument(bytes.NewReader(webpData))

	if err == nil {
		t.Fatal("Expected WebP to be rejected, got no error")
	}

	if result != nil {
		t.Error("Expected nil result for rejected WebP")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}

func TestPDFProcessor_RejectsRandomBytes(t *testing.T) {
	processor := NewPDFProcessor()

	randomData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

	result, err := processor.ProcessDocument(bytes.NewReader(randomData))

	if err == nil {
		t.Fatal("Expected random bytes to be rejected, got no error")
	}

	if result != nil {
		t.Error("Expected nil result for rejected data")
	}

	if err != ErrUnsupportedFormat {
		t.Errorf("Expected ErrUnsupportedFormat, got: %v", err)
	}
}
