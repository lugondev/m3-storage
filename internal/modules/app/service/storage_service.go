package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	ports "github.com/lugondev/m3-storage/internal/modules/app/port"

	"github.com/google/uuid"
)

// StorageServiceImpl implements the ports.StorageService interface
type StorageServiceImpl struct {
	basePath string
}

// NewStorageService creates a new adapters service
func NewStorageService() ports.StorageService {
	// In a production environment, this would be configured from environment variables
	// For now, we'll use a local directory
	return &StorageServiceImpl{
		basePath: "./uploads",
	}
}

// UploadFile uploads a file to adapters
func (s *StorageServiceImpl) UploadFile(ctx context.Context, bucket string, filename string, file io.Reader) (string, error) {
	// Create directory if it doesn't exist
	uploadPath := filepath.Join(s.basePath, bucket)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate a unique filename if not provided
	if filename == "" {
		filename = uuid.New().String()
	}

	filePath := filepath.Join(uploadPath, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the file contents
	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return the relative path to the file
	return filepath.Join(bucket, filename), nil
}

// DeleteFile deletes a file from adapters
func (s *StorageServiceImpl) DeleteFile(ctx context.Context, filePath string) error {
	// Ensure the path is within our base path
	if !strings.HasPrefix(filePath, s.basePath) {
		filePath = filepath.Join(s.basePath, filePath)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFileURL returns the URL for a file (helper method, not part of the interface)
func (s *StorageServiceImpl) GetFileURL(filePath string) string {
	// In a real implementation, this would return a proper URL
	// For now, we'll just return the relative path
	return filePath
}
