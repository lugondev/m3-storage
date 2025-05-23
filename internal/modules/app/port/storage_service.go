package port

import (
	"context"
	"io"
)

// StorageService defines the interface for file adapters operations
type StorageService interface {
	// UploadFile uploads a file to the specified bucket and returns its public URL
	UploadFile(ctx context.Context, bucket string, filename string, file io.Reader) (string, error)

	// DeleteFile deletes a file from adapters
	DeleteFile(ctx context.Context, filepath string) error
}
