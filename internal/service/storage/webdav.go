package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Uploader interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) (string, error)
}

type WebDAVUploader struct {
	base     string
	user     string
	pass     string
	withUser bool
	client   *http.Client
	username string
}

func NewWebDAVUploader() (*WebDAVUploader, error) {
	base := strings.TrimRight(os.Getenv("NEXTCLOUD_BASE_URL"), "/")
	user := os.Getenv("NEXTCLOUD_USERNAME")
	pass := os.Getenv("NEXTCLOUD_PASSWORD")

	if base == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("webdav env missing: NEXTCLOUD_BASE_URL/NEXTCLOUD_USERNAME/NEXTCLOUD_PASSWORD")
	}

	withUser := true
	if strings.Contains(base, "/files/"+user) {
		withUser = false
	}

	return &WebDAVUploader{
		base:     base,
		user:     user,
		pass:     pass,
		withUser: withUser,
		client:   &http.Client{Timeout: 30 * time.Second},
		username: user,
	}, nil
}

func (w *WebDAVUploader) fullURL(key string) (string, error) {
	key = strings.TrimLeft(key, "/")
	base := w.base
	if w.withUser {
		base = base + "/" + url.PathEscape(w.username)
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = path.Clean(u.Path + "/" + key)
	return u.String(), nil
}

func (w *WebDAVUploader) Upload(ctx context.Context, key string, r io.Reader, contentType string) (string, error) {
	// Validate inputs
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}
	if r == nil {
		return "", fmt.Errorf("reader cannot be nil")
	}

	// Ensure directory exists before uploading
	if err := w.ensureDirectoryExists(ctx, key); err != nil {
		return "", fmt.Errorf("failed to ensure directory exists: %w", err)
	}

	target, err := w.fullURL(key)
	if err != nil {
		return "", fmt.Errorf("failed to build target URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, target, r)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(w.user, w.pass)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "TinderTrip-Backend/1.0")

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for better error messages
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// If we can't read the body, still return the status error
		body = []byte("(unable to read response body)")
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "..."
		}
		return "", fmt.Errorf("webdav put failed: %s. Response body: %s", resp.Status, bodyStr)
	}
	return target, nil
}

// ensureDirectoryExists creates the directory structure for the given key
func (w *WebDAVUploader) ensureDirectoryExists(ctx context.Context, key string) error {
	// Extract directory path from key (e.g., "tindertrip/avatars/2024/01/02/file.png" -> "tindertrip/avatars/2024/01/02")
	dirPath := path.Dir(key)
	if dirPath == "." || dirPath == "/" {
		return nil // No directory to create
	}

	// Split path into components and create each level
	pathParts := strings.Split(strings.Trim(dirPath, "/"), "/")
	currentPath := ""

	for _, part := range pathParts {
		if part == "" {
			continue
		}

		if currentPath == "" {
			currentPath = part
		} else {
			currentPath = currentPath + "/" + part
		}

		// Create this directory level
		if err := w.createDirectory(ctx, currentPath); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", currentPath, err)
		}
	}

	return nil
}

// createDirectory creates a single directory level
func (w *WebDAVUploader) createDirectory(ctx context.Context, dirPath string) error {
	// Build directory URL
	dirURL, err := w.fullURL(dirPath + "/")
	if err != nil {
		return fmt.Errorf("failed to build directory URL: %w", err)
	}

	// Create MKCOL request to create directory
	req, err := http.NewRequestWithContext(ctx, "MKCOL", dirURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create MKCOL request: %w", err)
	}
	req.SetBasicAuth(w.user, w.pass)
	req.Header.Set("User-Agent", "TinderTrip-Backend/1.0")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute MKCOL request: %w", err)
	}
	defer resp.Body.Close()

	// MKCOL returns 201 (Created) for new directories, 405 (Method Not Allowed) if already exists
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMethodNotAllowed {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return fmt.Errorf("failed to create directory: %s. Response: %s", resp.Status, bodyStr)
	}

	return nil
}
