package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// NextcloudConfig represents the Nextcloud configuration
type NextcloudConfig struct {
	BaseURL     string
	Username    string
	Password    string
	AppPassword string
	WebDAVPath  string
}

// NextcloudClient represents a Nextcloud WebDAV client
type NextcloudClient struct {
	config     *NextcloudConfig
	httpClient *http.Client
}

// FileInfo represents file information from Nextcloud
type FileInfo struct {
	Path         string    `json:"path"`
	DisplayName  string    `json:"displayname"`
	ContentType  string    `json:"contenttype"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastmodified"`
	ETag         string    `json:"etag"`
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	ETag     string `json:"etag"`
	MimeType string `json:"mime_type"`
}

// NewNextcloudClient creates a new Nextcloud client
func NewNextcloudClient(baseURL, username, appPassword string) *NextcloudClient {
	config := &NextcloudConfig{
		BaseURL:     strings.TrimSuffix(baseURL, "/"),
		Username:    username,
		AppPassword: appPassword,
		WebDAVPath:  "/remote.php/dav/files/" + username,
	}

	return &NextcloudClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadFile uploads a file to Nextcloud
func (c *NextcloudClient) UploadFile(filePath string, fileData []byte, mimeType string) (*UploadResult, error) {
	// Create the full URL
	url := c.config.BaseURL + c.config.WebDAVPath + "/" + filePath

	// Create the request
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	req.SetBasicAuth(c.config.Username, c.config.AppPassword)

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Get file info
	fileInfo, err := c.GetFileInfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &UploadResult{
		URL:      c.config.BaseURL + "/index.php/apps/files/?dir=/" + filePath,
		Path:     filePath,
		Size:     fileInfo.Size,
		ETag:     fileInfo.ETag,
		MimeType: mimeType,
	}, nil
}

// UploadMultipartFile uploads a multipart file to Nextcloud
func (c *NextcloudClient) UploadMultipartFile(filePath string, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Get MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return c.UploadFile(filePath, fileData, mimeType)
}

// GetFileInfo gets information about a file
func (c *NextcloudClient) GetFileInfo(filePath string) (*FileInfo, error) {
	url := c.config.BaseURL + c.config.WebDAVPath + "/" + filePath

	req, err := http.NewRequest("PROPFIND", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Depth", "0")
	req.SetBasicAuth(c.config.Username, c.config.AppPassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("failed to get file info: status %d", resp.StatusCode)
	}

	// Parse the response (simplified - in production you'd use a proper XML parser)
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// This is a simplified parsing - in production you'd use xml.Unmarshal
	// For now, we'll return a basic file info
	return &FileInfo{
		Path:         filePath,
		DisplayName:  filepath.Base(filePath),
		ContentType:  "application/octet-stream",
		Size:         0, // Would be parsed from XML response
		LastModified: time.Now(),
		ETag:         "", // Would be parsed from XML response
	}, nil
}

// DeleteFile deletes a file from Nextcloud
func (c *NextcloudClient) DeleteFile(filePath string) error {
	url := c.config.BaseURL + c.config.WebDAVPath + "/" + filePath

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.config.Username, c.config.AppPassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateFolder creates a folder in Nextcloud
func (c *NextcloudClient) CreateFolder(folderPath string) error {
	url := c.config.BaseURL + c.config.WebDAVPath + "/" + folderPath + "/"

	req, err := http.NewRequest("MKCOL", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.config.Username, c.config.AppPassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMethodNotAllowed {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create folder failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListFiles lists files in a directory
func (c *NextcloudClient) ListFiles(directoryPath string) ([]FileInfo, error) {
	url := c.config.BaseURL + c.config.WebDAVPath + "/" + directoryPath + "/"

	req, err := http.NewRequest("PROPFIND", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Depth", "1")
	req.SetBasicAuth(c.config.Username, c.config.AppPassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("failed to list files: status %d", resp.StatusCode)
	}

	// Parse the response (simplified - in production you'd use a proper XML parser)
	// For now, return empty slice
	return []FileInfo{}, nil
}

// GetPublicURL generates a public share URL for a file
func (c *NextcloudClient) GetPublicURL(filePath string) (string, error) {
	// This would typically involve creating a public share via the Nextcloud API
	// For now, return a basic URL
	return c.config.BaseURL + "/index.php/apps/files/?dir=/" + filePath, nil
}

// GenerateFilePath generates a unique file path for upload
func (c *NextcloudClient) GenerateFilePath(originalName string, folder string) string {
	// Generate a unique filename with timestamp
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()

	uniqueName := fmt.Sprintf("%s_%d%s", name, timestamp, ext)

	if folder != "" {
		return folder + "/" + uniqueName
	}

	return uniqueName
}

// ValidateFileType validates if the file type is allowed
func (c *NextcloudClient) ValidateFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowedType := range allowedTypes {
		if ext == "."+strings.ToLower(allowedType) {
			return true
		}
	}

	return false
}

// GetFileSizeLimit returns the maximum file size allowed
func (c *NextcloudClient) GetFileSizeLimit() int64 {
	// Return 10MB as default limit
	return 10 * 1024 * 1024
}
