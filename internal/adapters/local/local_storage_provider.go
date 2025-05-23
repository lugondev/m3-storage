package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// LocalStorageProvider implements the StorageProvider interface for local file system.
type LocalStorageProvider struct {
	config config.LocalStorageConfig
}

// NewLocalStorageProvider creates a new LocalStorageProvider.
// It expects a config map that can be unmarshalled into LocalStorageConfig.
func NewLocalStorageProvider(cfg config.LocalStorageConfig) (port.StorageProvider, error) {
	// A more robust way would be to use a library like mapstructure to convert map to struct
	// For simplicity, we'll do direct type assertion here, but this is not production-ready.
	basePath := cfg.Path
	if basePath == "" {
		return nil, errors.New("local adapters: basePath is required in config")
	}

	// Ensure basePath exists and is a directory
	info, err := os.Stat(basePath)
	if os.IsNotExist(err) {
		// Attempt to create the directory
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return nil, fmt.Errorf("local adapters: failed to create basePath '%s': %w", basePath, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("local adapters: error accessing basePath '%s': %w", basePath, err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("local adapters: basePath '%s' is not a directory", basePath)
	}
	if cfg.SignedURLExpiry == 0 {
		cfg.SignedURLExpiry = 30 * time.Minute // Default expiry
	}
	if cfg.SignedURLSecret == "" {
		return nil, errors.New("local adapters: signedURLSecret is required in config")
	}

	return &LocalStorageProvider{
		config: cfg,
	}, nil
}

func (p *LocalStorageProvider) resolvePath(key string) string {
	// Sanitize key to prevent directory traversal issues
	cleanKey := filepath.Clean(key)
	// Ensure the key is relative and doesn't try to escape the base path
	if strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		// Handle potentially malicious key, e.g., return an error or a default safe path
		// For now, we'll just join it, but this needs careful consideration for security.
		// A better approach might be to ensure the key doesn't contain ".."
		// or to use a more robust path joining/cleaning mechanism.
		cleanKey = strings.ReplaceAll(cleanKey, "..", "") // Basic sanitization
	}
	return filepath.Join(p.config.Path, cleanKey)
}

// Upload uploads a file to the local file system.
func (p *LocalStorageProvider) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *port.UploadOptions) (*port.FileObject, error) {
	filePath := p.resolvePath(key)
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	if written != size && size != -1 { // size == -1 can mean unknown size for chunked transfer
		return nil, fmt.Errorf("file size mismatch: expected %d, wrote %d", size, written)
	}

	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
	}

	contentType := ""
	if opts != nil && opts.ContentType != "" {
		contentType = opts.ContentType
	}

	return &port.FileObject{
		Key:          key,
		URL:          p.buildPublicURL(key),
		Size:         fileInfo.Size(),
		ContentType:  contentType,
		LastModified: fileInfo.ModTime(),
		Provider:     p.ProviderType(),
	}, nil
}

func (p *LocalStorageProvider) buildPublicURL(key string) string {
	if p.config.BaseURL == "" {
		return "" // Not publicly accessible via a direct URL from this provider
	}
	// Ensure no double slashes and proper URL joining
	base := strings.TrimSuffix(p.config.BaseURL, "/")
	cleanKey := strings.TrimPrefix(filepath.ToSlash(key), "/") // Use URL slashes
	return base + "/" + cleanKey
}

// GetURL returns a publicly accessible URL for the given key.
func (p *LocalStorageProvider) GetURL(ctx context.Context, key string) (string, error) {
	filePath := p.resolvePath(key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("file not found")
	}
	return p.buildPublicURL(key), nil
}

// GetSignedURL generates a time-limited "signed" URL (simulated for local).
// This is a basic simulation and not cryptographically secure for production without more work.
func (p *LocalStorageProvider) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	filePath := p.resolvePath(key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("file not found")
	}

	if p.config.BaseURL == "" {
		return "", errors.New("base_url not configured, cannot generate signed URL")
	}

	// Basic "signing": append expiry and a simple hash/token (not secure for production)
	// In a real scenario, you'd use HMAC or a proper JWT-like mechanism.
	expiry := time.Now().Add(duration).Unix()
	tokenPayload := fmt.Sprintf("%s:%d:%s", key, expiry, p.config.SignedURLSecret)
	// In a real app, use crypto/hmac and sha256
	simulatedSignature := fmt.Sprintf("%x", []byte(tokenPayload)) // Very basic, NOT SECURE

	signedURL, err := url.Parse(p.buildPublicURL(key))
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := signedURL.Query()
	q.Set("expires", fmt.Sprintf("%d", expiry))
	q.Set("signature", simulatedSignature) // This signature is trivial to forge
	signedURL.RawQuery = q.Encode()

	return signedURL.String(), nil
}

// Delete removes a file from the local file system.
func (p *LocalStorageProvider) Delete(ctx context.Context, key string) error {
	filePath := p.resolvePath(key)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted or never existed, treat as success
		}
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}
	return nil
}

// GetObject retrieves file information.
func (p *LocalStorageProvider) GetObject(ctx context.Context, key string) (*port.FileObject, error) {
	filePath := p.resolvePath(key)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file not found")
		}
		return nil, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
	}

	// Determine content type (can be improved, e.g., by http.DetectContentType on a small portion of the file)
	contentType := "" // Placeholder, ideally detect or retrieve from metadata if stored

	return &port.FileObject{
		Key:          key,
		URL:          p.buildPublicURL(key),
		Size:         fileInfo.Size(),
		ContentType:  contentType,
		LastModified: fileInfo.ModTime(),
		Provider:     p.ProviderType(),
	}, nil
}

// Download retrieves a file from local adapters.
func (p *LocalStorageProvider) Download(ctx context.Context, key string) (io.ReadCloser, *port.FileObject, error) {
	filePath := p.resolvePath(key)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("file not found")
		}
		return nil, nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
	}

	objInfo := &port.FileObject{
		Key:          key,
		URL:          p.buildPublicURL(key),
		Size:         fileInfo.Size(),
		LastModified: fileInfo.ModTime(),
		Provider:     p.ProviderType(),
		// ContentType should be determined here as well
	}

	return file, objInfo, nil
}

// ProviderType returns the type of the adapters provider.
func (p *LocalStorageProvider) ProviderType() port.StorageProviderType {
	return port.ProviderLocal
}
