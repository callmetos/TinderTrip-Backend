package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"TinderTrip-Backend/internal/service/storage"
)

type ImageService struct {
	uploader storage.Uploader
	client   *http.Client
}

func NewImageService() (*ImageService, error) {
	uploader, err := storage.NewUploader()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage uploader: %w", err)
	}

	return &ImageService{
		uploader: uploader,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetImageFromKey downloads image from storage using key and returns the image data
func (s *ImageService) GetImageFromKey(ctx context.Context, key string) ([]byte, string, error) {
	if key == "" {
		return nil, "", fmt.Errorf("image key cannot be empty")
	}

	// Build the full URL for the image
	baseURL := strings.TrimRight(getEnvOrDefault("NEXTCLOUD_BASE_URL", ""), "/")
	username := getEnvOrDefault("NEXTCLOUD_USERNAME", "")

	if baseURL == "" || username == "" {
		return nil, "", fmt.Errorf("storage configuration missing")
	}

	// Check if key is already a full Nextcloud URL
	if strings.HasPrefix(key, "https://") {
		return s.GetImageFromURL(ctx, key)
	}

	// Construct the image URL using WebDAV path
	imageURL := fmt.Sprintf("%s/remote.php/dav/files/%s/%s", baseURL, username, strings.TrimLeft(key, "/"))

	return s.GetImageFromURL(ctx, imageURL)
}

// GetImageFromURL downloads image from Nextcloud URL and returns the image data
func (s *ImageService) GetImageFromURL(ctx context.Context, imageURL string) ([]byte, string, error) {
	if imageURL == "" {
		return nil, "", fmt.Errorf("image URL cannot be empty")
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add basic auth for Nextcloud
	username := getEnvOrDefault("NEXTCLOUD_USERNAME", "")
	password := getEnvOrDefault("NEXTCLOUD_PASSWORD", "")
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	req.Header.Set("User-Agent", "TinderTrip-Backend/1.0")

	// Make request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Get content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return imageData, contentType, nil
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
