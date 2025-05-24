package port

import (
	"context"
	"io"
	"time"
)

// StorageProviderType defines the type of adapters provider.
type StorageProviderType string

const (
	ProviderS3           StorageProviderType = "s3"            // Amazon S3, DigitalOcean Spaces, etc.
	ProviderCloudflareR2 StorageProviderType = "cloudflare_r2" // Cloudflare R2 is S3-compatible
	ProviderLocal        StorageProviderType = "local"
	ProviderFirebase     StorageProviderType = "firebase"  // Firebase Storage
	ProviderAzure        StorageProviderType = "azure"     // Azure Blob Storage
	ProviderDiscord      StorageProviderType = "discord"   // Discord channel storage
	ProviderScaleway     StorageProviderType = "scaleway"  // Scaleway Object Storage (S3-compatible)
	ProviderBackBlaze    StorageProviderType = "backblaze" // Backblaze B2 Cloud Storage
)

// FileObject represents a file stored in the adapters.
type FileObject struct {
	Key          string              `json:"key"`           // Unique identifier for the file in the adapters (e.g., path/to/file.jpg)
	URL          string              `json:"url"`           // Publicly accessible URL (if applicable)
	SignedURL    string              `json:"signed_url"`    // Time-limited signed URL for private access
	Size         int64               `json:"size"`          // File size in bytes
	ContentType  string              `json:"content_type"`  // MIME type of the file
	LastModified time.Time           `json:"last_modified"` // Last modified timestamp
	ETag         string              `json:"etag"`          // Entity tag, often an MD5 hash of the content
	Provider     StorageProviderType `json:"provider"`
}

// UploadOptions provides options for uploading a file.
type UploadOptions struct {
	ContentType string            // MIME type of the file
	Metadata    map[string]string // Custom metadata for the file
	ACL         string            // Access Control List (e.g., "public-read", "private") - specific to provider
}

// StorageProvider defines the interface for a adapters provider.
type StorageProvider interface {
	// CheckHealth checks if the storage provider is healthy and accessible.
	// Returns error if the provider is not healthy or cannot be accessed.
	CheckHealth(ctx context.Context) error

	// Upload uploads a file to the adapters.
	// key is the unique identifier for the file (e.g., "images/avatars/user123.jpg").
	// reader is the content of the file.
	// size is the size of the file in bytes.
	// opts provides additional upload options.
	Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *UploadOptions) (*FileObject, error)

	// GetURL returns a publicly accessible URL for the given key.
	// May return an empty string if the object is private or not directly accessible.
	GetURL(ctx context.Context, key string) (string, error)

	// GetSignedURL generates a time-limited signed URL for accessing a private object.
	// duration specifies how long the URL will be valid.
	GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error)

	// Delete removes a file from the adapters.
	Delete(ctx context.Context, key string) error

	// GetObject retrieves file information (metadata) without downloading the content.
	GetObject(ctx context.Context, key string) (*FileObject, error)

	// Download downloads a file.
	// Returns an io.ReadCloser that needs to be closed by the caller.
	Download(ctx context.Context, key string) (io.ReadCloser, *FileObject, error)

	// ProviderType returns the type of the adapters provider.
	ProviderType() StorageProviderType
}

// StorageFactory defines the interface for a factory that creates StorageProvider instances.
type StorageFactory interface {
	CreateProvider(providerType StorageProviderType) (StorageProvider, error)
}
