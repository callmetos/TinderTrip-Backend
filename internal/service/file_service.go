package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"TinderTrip-Backend/internal/service/storage"
	"TinderTrip-Backend/internal/utils"

	"github.com/google/uuid"
)

type FileService struct {
	uploader       storage.Uploader
	maxBytes       int64
	allow          map[string]bool
	imageProcessor *ImageProcessor
}

func NewFileService() (*FileService, error) {
	up, err := storage.NewUploader()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage uploader: %w", err)
	}

	maxMB := int64(10)
	if v := os.Getenv("MAX_UPLOAD_MB"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &maxMB); err != nil {
			return nil, fmt.Errorf("invalid MAX_UPLOAD_MB value: %s", v)
		}
		if maxMB <= 0 || maxMB > 100 {
			return nil, fmt.Errorf("MAX_UPLOAD_MB must be between 1 and 100, got: %d", maxMB)
		}
	}

	allow := map[string]bool{}
	allowed := os.Getenv("ALLOWED_IMAGE_TYPES")
	if allowed == "" {
		allowed = "image/jpeg,image/png,image/webp"
	}
	for _, t := range strings.Split(allowed, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			allow[t] = true
		}
	}

	if len(allow) == 0 {
		return nil, fmt.Errorf("no allowed image types configured")
	}

	return &FileService{
		uploader:       up,
		maxBytes:       maxMB * 1024 * 1024,
		allow:          allow,
		imageProcessor: NewImageProcessor(),
	}, nil
}

func detectContentType(head []byte) string {
	return http.DetectContentType(head)
}

func extFromCT(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	}
	if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
		return exts[0]
	}
	return ""
}

// UploadImage streams an image to storage and returns (key, url, size, checksum, contentType).
func (s *FileService) UploadImage(ctx context.Context, folder, filename string, body io.Reader) (key, url string, size int64, checksum, contentType string, err error) {
	// Validate inputs
	if folder == "" {
		return "", "", 0, "", "", fmt.Errorf("folder cannot be empty")
	}
	if filename == "" {
		return "", "", 0, "", "", fmt.Errorf("filename cannot be empty")
	}
	if body == nil {
		return "", "", 0, "", "", fmt.Errorf("body cannot be nil")
	}

	lr := &io.LimitedReader{R: body, N: s.maxBytes + 1}

	head := make([]byte, 512)
	n, err := io.ReadFull(lr, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", "", 0, "", "", fmt.Errorf("failed to read file header: %w", err)
	}
	head = head[:n]

	// Check if we have any data
	if n == 0 {
		return "", "", 0, "", "", fmt.Errorf("file is empty")
	}

	ct := detectContentType(head)
	if !s.allow[ct] {
		allowedTypes := make([]string, 0, len(s.allow))
		for t := range s.allow {
			allowedTypes = append(allowedTypes, t)
		}
		return "", "", 0, "", "", fmt.Errorf("unsupported content type: %s. Allowed types: %v", ct, allowedTypes)
	}

	reader := io.MultiReader(bytes.NewReader(head), lr)

	buf := new(bytes.Buffer)
	written, err := io.Copy(buf, reader)
	if err != nil {
		return "", "", 0, "", "", fmt.Errorf("failed to read file content: %w", err)
	}
	if written > s.maxBytes {
		return "", "", 0, "", "", fmt.Errorf("file too large: %d bytes exceeds limit of %d bytes", written, s.maxBytes)
	}

	// Validate minimum file size
	if written < 100 {
		return "", "", 0, "", "", fmt.Errorf("file too small: minimum 100 bytes required")
	}

	// Process image if it's an image type that should be optimized
	processedData := buf.Bytes()
	processedContentType := ct
	if s.imageProcessor.ShouldProcess(processedData, ct) {
		processed, newContentType, err := s.imageProcessor.ProcessImage(processedData, ct)
		if err != nil {
			// If processing fails, use original image
			utils.Logger().WithField("error", err).Warn("Failed to process image, using original")
		} else {
			processedData = processed
			processedContentType = newContentType
			written = int64(len(processedData))
		}
	}

	sum := sha256.Sum256(processedData)
	checksum = fmt.Sprintf("sha256:%x", sum[:])

	day := time.Now().Format("2006/01/02")
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = extFromCT(processedContentType)
	}
	if ext == "" {
		ext = ".bin"
	}

	id := uuid.New().String()
	key = fmt.Sprintf("tindertrip/%s/%s/%s%s", strings.Trim(folder, "/"), day, id, ext)
	key = strings.ReplaceAll(key, "//", "/")

	url, err = s.uploader.Upload(ctx, key, bytes.NewReader(processedData), processedContentType)
	if err != nil {
		return "", "", 0, "", "", fmt.Errorf("failed to upload file to storage: %w", err)
	}
	return key, url, written, checksum, processedContentType, nil
}
