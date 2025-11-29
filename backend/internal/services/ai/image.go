package ai

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

// Image processing errors
var (
	ErrUnsupportedFormat = errors.New("unsupported format: only JPEG, PNG, and PDF are supported")
	ErrImageTooSmall     = errors.New("image is too small to process")
	ErrImageTooLarge     = errors.New("image is too large to process")
	ErrReadImage         = errors.New("failed to read image file")
	ErrDecodeImage       = errors.New("failed to decode image")
	ErrEncodeImage       = errors.New("failed to encode image")
)

const (
	// MaxImageDimension is the maximum size for the longest side of an image
	MaxImageDimension = 1024
	// MaxInputDimension is the maximum allowed input dimension
	MaxInputDimension = 4096
	// JPEGQuality is the quality setting for JPEG encoding (1-100)
	JPEGQuality = 85
)

// ImageProcessor handles image processing operations
type ImageProcessor struct {
	maxDimension int
	jpegQuality  int
}

// NewImageProcessor creates a new ImageProcessor with default settings
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		maxDimension: MaxImageDimension,
		jpegQuality:  JPEGQuality,
	}
}

// ProcessedImage represents an image ready for AI processing
type ProcessedImage struct {
	Base64Data string
	MimeType   string
	Width      int
	Height     int
}

// ReadAndProcessFile reads an image file and processes it for AI analysis
func (p *ImageProcessor) ReadAndProcessFile(filePath string) (*ProcessedImage, error) {
	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadImage, err)
	}
	defer f.Close()

	return p.ProcessImage(f)
}

// ReadAndProcessReader reads an image from an io.ReadSeeker and processes it for AI analysis
func (p *ImageProcessor) ReadAndProcessReader(r io.ReadSeeker) (*ProcessedImage, error) {
	return p.ProcessImage(r)
}

// ProcessPDF reads a PDF and returns base64 encoded data
func (p *ImageProcessor) ProcessPDF(r io.Reader) (*ProcessedImage, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadImage, err)
	}

	// Validate it's a PDF
	if len(data) < 4 || string(data[:4]) != "%PDF" {
		return nil, ErrUnsupportedFormat
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	return &ProcessedImage{
		Base64Data: base64Data,
		MimeType:   "application/pdf",
		Width:      0,
		Height:     0,
	}, nil
}

// ProcessImage processes image or PDF for AI analysis
func (p *ImageProcessor) ProcessImage(r io.ReadSeeker) (*ProcessedImage, error) {
	// Read first bytes to detect format
	header := make([]byte, 8)
	if _, err := r.Read(header); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadImage, err)
	}

	// Rewind
	if _, err := r.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek stream: %w", err)
	}

	// Check if it's a PDF
	if string(header[:4]) == "%PDF" {
		return p.ProcessPDF(r)
	}

	// Continue with existing image processing logic...
	// Decode config to check dimensions and format first
	config, format, err := image.DecodeConfig(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodeImage, err)
	}

	// Validate format
	if format != "jpeg" && format != "png" {
		return nil, ErrUnsupportedFormat
	}

	// Check dimensions
	if config.Width > MaxInputDimension || config.Height > MaxInputDimension {
		return nil, ErrImageTooLarge
	}

	// Rewind the reader
	if _, err := r.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek stream: %w", err)
	}

	// Decode the image
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodeImage, err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resizing is needed
	var processedImg image.Image
	var finalWidth, finalHeight int

	if width > p.maxDimension || height > p.maxDimension {
		processedImg, finalWidth, finalHeight = p.resizeImage(img, width, height)
	} else {
		processedImg = img
		finalWidth = width
		finalHeight = height
	}

	// Encode the image
	var buf bytes.Buffer
	var outputMimeType string

	// Use the original format if possible
	switch format {
	case "png":
		if err := png.Encode(&buf, processedImg); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrEncodeImage, err)
		}
		outputMimeType = "image/png"
	default:
		// Default to JPEG for other formats
		if err := jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: p.jpegQuality}); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrEncodeImage, err)
		}
		outputMimeType = "image/jpeg"
	}

	// Encode to base64
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	return &ProcessedImage{
		Base64Data: base64Data,
		MimeType:   outputMimeType,
		Width:      finalWidth,
		Height:     finalHeight,
	}, nil
}

// isValidFormat checks if the MIME type is a supported format
func (p *ImageProcessor) isValidFormat(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	return mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/jpg" || mimeType == "application/pdf"
}

// ValidateFormat validates the format of file data (images or PDF)
func (p *ImageProcessor) ValidateFormat(data []byte) (string, error) {
	mimeType := http.DetectContentType(data)

	// http.DetectContentType may not detect PDF correctly, check magic bytes
	if len(data) >= 4 && string(data[:4]) == "%PDF" {
		mimeType = "application/pdf"
	}

	if !p.isValidFormat(mimeType) {
		return "", ErrUnsupportedFormat
	}
	return mimeType, nil
}

// ValidateImageFormat validates the format of image data
// Deprecated: Use ValidateFormat instead
func (p *ImageProcessor) ValidateImageFormat(data []byte) error {
	_, err := p.ValidateFormat(data)
	return err
}

// resizeImage resizes an image to fit within the max dimension while maintaining aspect ratio
func (p *ImageProcessor) resizeImage(img image.Image, width, height int) (image.Image, int, int) {
	var newWidth, newHeight int

	if width > height {
		newWidth = p.maxDimension
		newHeight = int(float64(height) * float64(p.maxDimension) / float64(width))
	} else {
		newHeight = p.maxDimension
		newWidth = int(float64(width) * float64(p.maxDimension) / float64(height))
	}

	// Ensure dimensions are at least 1x1
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	// Create a new image with the target size
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Use high-quality resampling
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	return dst, newWidth, newHeight
}

// GetBase64FromFile is a convenience function to read a file and return base64-encoded data
func GetBase64FromFile(filePath string) (string, string, error) {
	processor := NewImageProcessor()
	result, err := processor.ReadAndProcessFile(filePath)
	if err != nil {
		return "", "", err
	}
	return result.Base64Data, result.MimeType, nil
}

// GetBase64FromReader is a convenience function to read from a reader and return base64-encoded data
func GetBase64FromReader(r io.Reader) (string, string, error) {
	processor := NewImageProcessor()

	// If reader is a seeker, use it directly
	if seeker, ok := r.(io.ReadSeeker); ok {
		result, err := processor.ReadAndProcessReader(seeker)
		if err != nil {
			return "", "", err
		}
		return result.Base64Data, result.MimeType, nil
	}

	// Otherwise, we must buffer it to make it seekable
	data, err := io.ReadAll(r)
	if err != nil {
		return "", "", err
	}

	// Recursive call with a seeker
	return GetBase64FromReader(bytes.NewReader(data))
}
