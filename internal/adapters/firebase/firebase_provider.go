package firebase

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// firebaseProvider implements the port.StorageProvider interface for Firebase Cloud Storage.
type firebaseProvider struct {
	bucket     *storage.BucketHandle
	bucketName string
	logger     logger.Logger
}

// NewFirebaseService creates a new instance of FirebaseStorageService.
func NewFirebaseService(cfg config.FireStoreConfig, log logger.Logger) (port.StorageProvider, error) {
	app, _, err := InitializeFirebase(cfg, log)
	if err != nil {
		log.Errorf(context.Background(), "Firebase initialization failed: %v", err)
		return nil, fmt.Errorf("firebase initialization failed: %w", err)
	}

	firebaseSDKClient, err := app.Storage(context.Background())
	if err != nil {
		log.Errorf(context.Background(), "Error getting Firebase Storage client (SDK): %v", map[string]any{"error": err})
		return nil, fmt.Errorf("error getting Firebase Storage client (SDK): %w", err)
	}

	bucketName := cfg.BucketName
	if bucketName == "" {
		log.Error(context.Background(), "Firebase Storage bucket name not configured", nil)
		return nil, fmt.Errorf("firebase StorageBucket is not configured")
	}

	bucketHandle, err := firebaseSDKClient.Bucket(bucketName)
	if err != nil {
		log.Errorf(context.Background(), "Error getting bucket handle: %v", map[string]any{"bucket": bucketName, "error": err})
		return nil, fmt.Errorf("error getting bucket handle for %s: %w", bucketName, err)
	}

	log.Infof(context.Background(), "FirebaseStorageService initialized", map[string]any{"bucket": bucketName})

	return &firebaseProvider{
		bucket:     bucketHandle, // Store the bucket handle
		bucketName: bucketName,
		logger:     log.WithFields(map[string]any{"component": "FirebaseStorageService"}),
	}, nil
}

// Upload uploads a file to Firebase Cloud Storage.
func (p *firebaseProvider) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *port.UploadOptions) (*port.FileObject, error) {
	if key == "" {
		p.logger.Error(ctx, "Upload key cannot be empty", nil)
		return nil, fmt.Errorf("upload key cannot be empty")
	} else {
		key = fmt.Sprintf("%s/%s", "m3-storage", key)
	}

	// If key doesn't have an extension, try to infer or use a default.
	// For Firebase, often the key is the full path including a generated filename.
	// If opts.ContentType is not set, try to determine from key's extension.
	contentType := ""
	if opts != nil && opts.ContentType != "" {
		contentType = opts.ContentType
	} else {
		ext := filepath.Ext(key)
		if ext != "" {
			contentType = mime.TypeByExtension(ext)
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream" // Default
	}

	// Generate a unique ID if key is meant to be a folder path or if we want to ensure uniqueness.
	// For this provider, 'key' is treated as the final object path.
	// If key needs to be dynamic:
	// finalKey := filepath.Join(key, uuid.New().String()+filepath.Ext(originalFileNameFromOptsOrContext))
	finalKey := key

	obj := p.bucket.Object(finalKey)
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType
	wc.Size = size // Set the size for resumable uploads or progress tracking

	if opts != nil && opts.Metadata != nil {
		wc.Metadata = opts.Metadata
	}
	// ACL handling for Firebase/GCS is typically done via bucket/object IAM policies
	// or predefined ACLs like "publicRead".
	// wc.ACL = ... (if direct ACL setting is needed and supported by the writer)

	p.logger.Infof(ctx, "Attempting to upload file", map[string]any{"key": finalKey, "contentType": contentType, "size": size})

	if _, err := io.Copy(wc, reader); err != nil {
		p.logger.Errorf(ctx, "Failed to copy file to Firebase Storage", map[string]any{"key": finalKey, "error": err})
		return nil, fmt.Errorf("failed to copy file to Firebase Storage for key %s: %w", finalKey, err)
	}
	if err := wc.Close(); err != nil {
		p.logger.Errorf(ctx, "Failed to close Firebase Storage writer", map[string]any{"key": finalKey, "error": err})
		return nil, fmt.Errorf("failed to close Firebase Storage writer for key %s: %w", finalKey, err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to get attributes after upload", map[string]any{"key": finalKey, "error": err})
		return nil, fmt.Errorf("failed to get attributes for key %s after upload: %w", finalKey, err)
	}

	fileURL := p.generatePublicURL(finalKey)
	p.logger.Infof(ctx, "File uploaded successfully", map[string]any{"key": finalKey, "url": fileURL})

	return &port.FileObject{
		Key:          finalKey,
		URL:          fileURL, // This URL might not be publicly accessible by default
		Size:         attrs.Size,
		ContentType:  attrs.ContentType,
		LastModified: attrs.Updated,
		ETag:         attrs.Etag,
		Provider:     p.ProviderType(),
	}, nil
}

// generatePublicURL creates a standard public URL for Firebase Storage objects.
// Note: This URL is only accessible if the object has public read permissions.
func (p *firebaseProvider) generatePublicURL(key string) string {
	obj := p.bucket.Object(key)
	publicURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media",
		p.bucketName, url.PathEscape(obj.ObjectName()),
	)
	return publicURL
}

// GetURL returns a publicly accessible URL for the given key.
// For Firebase, this usually means the object must be publicly readable.
func (p *firebaseProvider) GetURL(ctx context.Context, key string) (string, error) {
	// Check if object exists first, though not strictly necessary for URL generation
	// _, err := p.bucket.Object(key).Attrs(ctx)
	// if err != nil {
	// 	if err == adapters.ErrObjectNotExist {
	// 		p.logger.Warnf(ctx, "Object does not exist, cannot get URL", map[string]any{"key": key})
	// 		return "", fmt.Errorf("object %s not found: %w", key, err)
	// 	}
	// 	p.logger.Errorf(ctx, "Failed to get object attributes for URL", map[string]any{"key": key, "error": err})
	// 	return "", fmt.Errorf("failed to get attributes for %s: %w", key, err)
	// }
	return p.generatePublicURL(key), nil
}

// GetSignedURL generates a time-limited signed URL for accessing a private object.
func (p *firebaseProvider) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(duration),
		// Add other options like Headers, QueryParams if needed
	}

	signedURL, err := p.bucket.SignedURL(key, opts)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to generate signed URL", map[string]any{"key": key, "duration": duration, "error": err})
		return "", fmt.Errorf("failed to generate signed URL for key %s: %w", key, err)
	}
	p.logger.Infof(ctx, "Generated signed URL", map[string]any{"key": key, "duration": duration})
	return signedURL, nil
}

// Delete removes a file from Firebase Cloud Storage.
func (p *firebaseProvider) Delete(ctx context.Context, key string) error {
	obj := p.bucket.Object(key)
	if err := obj.Delete(ctx); err != nil {
		if err == storage.ErrObjectNotExist {
			p.logger.Warnf(ctx, "Attempted to delete non-existent object", map[string]any{"key": key})
			return nil // Or a specific "not found" error if preferred by the interface contract
		}
		p.logger.Errorf(ctx, "Failed to delete file from Firebase Storage", map[string]any{"key": key, "error": err})
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}
	p.logger.Infof(ctx, "File deleted successfully", map[string]any{"key": key})
	return nil
}

// GetObject retrieves file information (metadata) without downloading the content.
func (p *firebaseProvider) GetObject(ctx context.Context, key string) (*port.FileObject, error) {
	attrs, err := p.bucket.Object(key).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			p.logger.Warnf(ctx, "Object not found for GetObject", map[string]any{"key": key})
			return nil, fmt.Errorf("object %s not found: %w", key, err) // Consider a custom error like port.ErrObjectNotFound
		}
		p.logger.Errorf(ctx, "Failed to get object attributes", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to get attributes for object %s: %w", key, err)
	}

	fileURL := p.generatePublicURL(key)
	// A signed URL is not typically generated here unless specifically requested by the use case for FileObject.
	// If a short-lived signed URL is always needed with GetObject, generate it here.

	return &port.FileObject{
		Key:          attrs.Name, // Attrs.Name is the full path
		URL:          fileURL,
		Size:         attrs.Size,
		ContentType:  attrs.ContentType,
		LastModified: attrs.Updated,
		ETag:         attrs.Etag,
		Provider:     p.ProviderType(),
		// SignedURL: could be generated here if needed, e.g. p.GetSignedURL(ctx, key, 5*time.Minute)
	}, nil
}

// Download downloads a file. Returns an io.ReadCloser that needs to be closed by the caller.
func (p *firebaseProvider) Download(ctx context.Context, key string) (io.ReadCloser, *port.FileObject, error) {
	objHandle := p.bucket.Object(key)
	attrs, err := objHandle.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			p.logger.Warnf(ctx, "Object not found for Download", map[string]any{"key": key})
			return nil, nil, fmt.Errorf("object %s not found for download: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to get object attributes before download", map[string]any{"key": key, "error": err})
		return nil, nil, fmt.Errorf("failed to get attributes for %s before download: %w", key, err)
	}

	reader, err := objHandle.NewReader(ctx)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to create reader for object", map[string]any{"key": key, "error": err})
		return nil, nil, fmt.Errorf("failed to create reader for object %s: %w", key, err)
	}

	fileObject := &port.FileObject{
		Key:          attrs.Name,
		URL:          p.generatePublicURL(key),
		Size:         attrs.Size,
		ContentType:  attrs.ContentType,
		LastModified: attrs.Updated,
		ETag:         attrs.Etag,
		Provider:     p.ProviderType(),
	}
	p.logger.Infof(ctx, "Prepared file for download", map[string]any{"key": key})
	return reader, fileObject, nil
}

// ProviderType returns the type of the adapters provider.
func (p *firebaseProvider) ProviderType() port.StorageProviderType {
	return port.ProviderFirebase
}
