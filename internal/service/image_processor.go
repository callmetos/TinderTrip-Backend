package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

// ImageProcessor handles image compression and resizing
type ImageProcessor struct {
	maxWidth    int
	maxHeight   int
	jpegQuality int
}

// NewImageProcessor creates a new image processor with configurable settings
func NewImageProcessor() *ImageProcessor {
	maxWidth := 1920  // Max width in pixels
	maxHeight := 1920 // Max height in pixels
	jpegQuality := 85 // JPEG quality (1-100)

	// Allow configuration via environment variables
	if v := os.Getenv("IMAGE_MAX_WIDTH"); v != "" {
		if w := parseInt(v); w > 0 {
			maxWidth = w
		}
	}
	if v := os.Getenv("IMAGE_MAX_HEIGHT"); v != "" {
		if h := parseInt(v); h > 0 {
			maxHeight = h
		}
	}
	if v := os.Getenv("IMAGE_JPEG_QUALITY"); v != "" {
		if q := parseInt(v); q > 0 && q <= 100 {
			jpegQuality = q
		}
	}

	return &ImageProcessor{
		maxWidth:    maxWidth,
		maxHeight:   maxHeight,
		jpegQuality: jpegQuality,
	}
}

// ProcessImage compresses and/or resizes an image
// Returns processed image data, content type, and error
func (p *ImageProcessor) ProcessImage(imageData []byte, contentType string) ([]byte, string, error) {
	// Parse the image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Get original bounds
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions if resizing is needed
	newWidth, newHeight := width, height
	if width > p.maxWidth || height > p.maxHeight {
		ratio := float64(width) / float64(height)
		if width > height {
			newWidth = p.maxWidth
			newHeight = int(float64(p.maxWidth) / ratio)
			if newHeight > p.maxHeight {
				newHeight = p.maxHeight
				newWidth = int(float64(p.maxHeight) * ratio)
			}
		} else {
			newHeight = p.maxHeight
			newWidth = int(float64(p.maxHeight) * ratio)
			if newWidth > p.maxWidth {
				newWidth = p.maxWidth
				newHeight = int(float64(p.maxWidth) / ratio)
			}
		}
	}

	// Resize if needed
	var processedImg image.Image = img
	if newWidth != width || newHeight != height {
		// Create a new RGBA image for the resized version
		dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
		processedImg = dst
	}

	// Encode the processed image
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		// Convert to JPEG with quality setting
		if contentType != "image/jpeg" {
			contentType = "image/jpeg"
		}
		err = jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: p.jpegQuality})
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode JPEG: %w", err)
		}
	case "png":
		// Keep PNG format
		err = png.Encode(&buf, processedImg)
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode PNG: %w", err)
		}
	default:
		// For other formats, convert to JPEG
		contentType = "image/jpeg"
		err = jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: p.jpegQuality})
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode as JPEG: %w", err)
		}
	}

	return buf.Bytes(), contentType, nil
}

// ShouldProcess determines if an image should be processed
func (p *ImageProcessor) ShouldProcess(imageData []byte, contentType string) bool {
	// Only process image/jpeg and image/png
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
		return false
	}

	// Check file size - only process if larger than 100KB
	if len(imageData) < 100*1024 {
		return false
	}

	return true
}

// Helper function to parse integer from string
func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

