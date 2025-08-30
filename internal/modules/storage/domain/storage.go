package domain

import (
	"context"
	"io"
	"time"
)

// StorageProviderType represents the type of storage provider
type StorageProviderType string

const (
	ProviderS3           StorageProviderType = "s3"
	ProviderCloudflareR2 StorageProviderType = "cloudflare_r2"
	ProviderLocal        StorageProviderType = "local"
	ProviderFirebase     StorageProviderType = "firebase"
	ProviderAzure        StorageProviderType = "azure"
	ProviderDiscord      StorageProviderType = "discord"
	ProviderScaleway     StorageProviderType = "scaleway"
	ProviderBackBlaze    StorageProviderType = "backblaze"
	ProviderMinIO        StorageProviderType = "minio"
)

// FileObject represents a file stored in the storage system
type FileObject struct {
	Key          string
	URL          string
	SignedURL    string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
	Provider     StorageProviderType
}

// UploadOptions provides options for uploading a file
type UploadOptions struct {
	ContentType string
	Metadata    map[string]string
	ACL         string
}

// StorageProvider defines the domain interface for storage operations
type StorageProvider interface {
	CheckHealth(ctx context.Context) error
	Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *UploadOptions) (*FileObject, error)
	GetURL(ctx context.Context, key string) (string, error)
	GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
	GetObject(ctx context.Context, key string) (*FileObject, error)
	Download(ctx context.Context, key string) (io.ReadCloser, *FileObject, error)
	ProviderType() StorageProviderType
}

// StorageFactory defines the domain interface for creating storage providers
type StorageFactory interface {
	CreateProvider(providerType StorageProviderType) (StorageProvider, error)
}

// HealthStatus represents the health status of a storage provider
type HealthStatus struct {
	Provider StorageProviderType
	Status   string
	Message  string
}

// IsHealthy returns true if the status indicates the provider is healthy
func (h *HealthStatus) IsHealthy() bool {
	return h.Status == "healthy"
}

// NewHealthStatus creates a new HealthStatus instance
func NewHealthStatus(provider StorageProviderType, status, message string) *HealthStatus {
	return &HealthStatus{
		Provider: provider,
		Status:   status,
		Message:  message,
	}
}
