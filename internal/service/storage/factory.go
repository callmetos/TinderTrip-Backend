package storage

import (
	"fmt"
	"os"
)

func NewUploader() (Uploader, error) {
	switch os.Getenv("STORAGE_PROVIDER") {
	case "webdav", "":
		return NewWebDAVUploader()
	default:
		return nil, fmt.Errorf("unsupported STORAGE_PROVIDER: %s", os.Getenv("STORAGE_PROVIDER"))
	}
}
