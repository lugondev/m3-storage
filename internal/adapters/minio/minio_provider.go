package minio

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// normalizeMinIOEndpoint validates and normalizes MinIO endpoint URL
func normalizeMinIOEndpoint(endpoint string, useSSL bool) (string, error) {
	if endpoint == "" {
		return "", fmt.Errorf("endpoint cannot be empty")
	}

	// Parse the endpoint URL
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint URL: %w", err)
	}

	// If no scheme is provided, add the appropriate one based on useSSL
	if parsedURL.Scheme == "" {
		if useSSL {
			endpoint = "https://" + endpoint
		} else {
			endpoint = "http://" + endpoint
		}
		// Re-parse with scheme
		parsedURL, err = url.Parse(endpoint)
		if err != nil {
			return "", fmt.Errorf("invalid endpoint URL after adding scheme: %w", err)
		}
	}

	// Validate scheme matches useSSL setting
	if useSSL && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("useSSL is true but endpoint scheme is %s", parsedURL.Scheme)
	}
	if !useSSL && parsedURL.Scheme != "http" {
		return "", fmt.Errorf("useSSL is false but endpoint scheme is %s", parsedURL.Scheme)
	}

	// Ensure host is not empty
	if parsedURL.Host == "" {
		return "", fmt.Errorf("endpoint host cannot be empty")
	}

	return parsedURL.String(), nil
}

// isMinIOEndpoint attempts to detect if an endpoint is likely a MinIO server
// This is used for validation and logging purposes
func isMinIOEndpoint(endpoint string) bool {
	if endpoint == "" {
		return false
	}

	// Parse the endpoint
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return false
	}

	host := strings.ToLower(parsedURL.Host)

	// Common MinIO patterns
	minioPatterns := []string{
		"min.io",         // Official MinIO domains
		"minio",          // Common in hostnames/subdomains
		":9000",          // Default MinIO port
		":9001",          // Default MinIO console port
		"localhost:9000", // Local development
		"127.0.0.1:9000", // Local development
	}

	for _, pattern := range minioPatterns {
		if strings.Contains(host, pattern) {
			return true
		}
	}

	// Check port specifically for common MinIO ports
	if parsedURL.Port() == "9000" || parsedURL.Port() == "9001" {
		return true
	}

	return false
}

// minioProvider implements the port.StorageProvider interface for MinIO.
type minioProvider struct {
	client      *minio.Client
	bucketName  string
	region      string
	endpointURL string
	useSSL      bool
	logger      logger.Logger
}

// NewMinIOProvider creates a new instance of minioProvider.
func NewMinIOProvider(cfg config.MinIOConfig, log logger.Logger) (port.StorageProvider, error) {
	log = log.WithFields(map[string]any{"component": "MinIOProvider"})

	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for MinIOProvider")
	}
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("bucket_name is required for MinIOProvider")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("access_key_id is required for MinIOProvider")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("secret_access_key is required for MinIOProvider")
	}

	// Validate and normalize endpoint URL
	normalizedEndpoint, err := normalizeMinIOEndpoint(cfg.Endpoint, cfg.UseSSL)
	if err != nil {
		log.Errorf(context.Background(), "Invalid MinIO endpoint", map[string]any{"endpoint": cfg.Endpoint, "error": err})
		return nil, fmt.Errorf("invalid MinIO endpoint: %w", err)
	}

	// Log warning if endpoint doesn't look like MinIO
	if !isMinIOEndpoint(cfg.Endpoint) {
		log.Warnf(context.Background(), "Endpoint does not appear to be a MinIO server", map[string]any{
			"endpoint": cfg.Endpoint,
			"hint":     "Make sure this is a MinIO server endpoint. If it's a different S3-compatible service, consider using the S3 provider instead.",
		})
	}

	// Parse the endpoint to get host and port
	parsedURL, err := url.Parse(normalizedEndpoint)
	if err != nil {
		log.Errorf(context.Background(), "Failed to parse MinIO endpoint", map[string]any{"endpoint": normalizedEndpoint, "error": err})
		return nil, fmt.Errorf("failed to parse MinIO endpoint: %w", err)
	}

	// Create MinIO client
	minioClient, err := minio.New(parsedURL.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		log.Errorf(context.Background(), "Failed to create MinIO client", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	log.Infof(context.Background(), "MinIOProvider initialized", map[string]any{
		"bucket":             cfg.BucketName,
		"endpoint":           cfg.Endpoint,
		"normalizedEndpoint": normalizedEndpoint,
		"region":             cfg.Region,
		"useSSL":             cfg.UseSSL,
	})

	return &minioProvider{
		client:      minioClient,
		bucketName:  cfg.BucketName,
		region:      cfg.Region,
		endpointURL: normalizedEndpoint,
		useSSL:      cfg.UseSSL,
		logger:      log,
	}, nil
}

// Upload uploads a file to MinIO.
func (p *minioProvider) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *port.UploadOptions) (*port.FileObject, error) {
	if key == "" {
		return nil, fmt.Errorf("upload key cannot be empty")
	}

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

	putObjectOpts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	if opts != nil {
		if opts.Metadata != nil {
			putObjectOpts.UserMetadata = opts.Metadata
		}
	}

	p.logger.Infof(ctx, "Attempting to upload file to MinIO", map[string]any{"key": key, "bucket": p.bucketName, "contentType": contentType})

	info, err := p.client.PutObject(ctx, p.bucketName, key, reader, size, putObjectOpts)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to upload file to MinIO", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to upload to MinIO key %s: %w", key, err)
	}

	fileURL := p.generateObjectURL(ctx, key)

	p.logger.Infof(ctx, "File uploaded successfully to MinIO", map[string]any{"key": key, "size": info.Size, "etag": info.ETag})
	return &port.FileObject{
		Key:          key,
		URL:          fileURL,
		Size:         info.Size,
		ContentType:  contentType,
		LastModified: info.LastModified,
		ETag:         strings.Trim(info.ETag, "\""),
		Provider:     p.ProviderType(),
	}, nil
}

// generateObjectURL generates the public URL for accessing the object.
func (p *minioProvider) generateObjectURL(ctx context.Context, key string) string {
	// MinIO always uses path-style URL: http(s)://endpoint/bucket/key
	// Since endpointURL is already normalized with proper scheme, we can use it directly
	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(p.endpointURL, "/"), p.bucketName, strings.TrimPrefix(key, "/"))
}

// GetURL returns a publicly accessible URL for the given key.
func (p *minioProvider) GetURL(ctx context.Context, key string) (string, error) {
	// First check if object exists
	_, err := p.client.StatObject(ctx, p.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" || errResponse.Code == "NotFound" {
			p.logger.Warnf(ctx, "Object not found, cannot get URL", map[string]any{"key": key})
			return "", fmt.Errorf("object %s not found: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to check object existence for GetURL", map[string]any{"key": key, "error": err})
		return "", fmt.Errorf("failed to check object %s: %w", key, err)
	}
	return p.generateObjectURL(ctx, key), nil
}

// GetSignedURL generates a time-limited signed URL for accessing a private object.
func (p *minioProvider) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := p.client.PresignedGetObject(ctx, p.bucketName, key, duration, reqParams)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to generate MinIO signed URL", map[string]any{"key": key, "error": err})
		return "", fmt.Errorf("failed to generate MinIO signed URL for key %s: %w", key, err)
	}
	p.logger.Infof(ctx, "Generated MinIO signed URL", map[string]any{"key": key, "duration": duration})
	return presignedURL.String(), nil
}

// Delete removes a file from MinIO.
func (p *minioProvider) Delete(ctx context.Context, key string) error {
	err := p.client.RemoveObject(ctx, p.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		p.logger.Errorf(ctx, "Failed to delete MinIO object", map[string]any{"key": key, "error": err})
		return fmt.Errorf("failed to delete MinIO object %s: %w", key, err)
	}
	p.logger.Infof(ctx, "MinIO object deleted successfully", map[string]any{"key": key})
	return nil
}

// GetObject retrieves file information (metadata) from MinIO.
func (p *minioProvider) GetObject(ctx context.Context, key string) (*port.FileObject, error) {
	objectInfo, err := p.client.StatObject(ctx, p.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" || errResponse.Code == "NotFound" {
			p.logger.Warnf(ctx, "MinIO object not found for GetObject", map[string]any{"key": key})
			return nil, fmt.Errorf("minio object %s not found: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to get MinIO object info for GetObject", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to get MinIO object metadata for %s: %w", key, err)
	}

	return &port.FileObject{
		Key:          key,
		URL:          p.generateObjectURL(ctx, key),
		Size:         objectInfo.Size,
		ContentType:  objectInfo.ContentType,
		LastModified: objectInfo.LastModified,
		ETag:         strings.Trim(objectInfo.ETag, "\""),
		Provider:     p.ProviderType(),
	}, nil
}

// Download downloads a file from MinIO.
func (p *minioProvider) Download(ctx context.Context, key string) (io.ReadCloser, *port.FileObject, error) {
	object, err := p.client.GetObject(ctx, p.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" || errResponse.Code == "NotFound" {
			p.logger.Warnf(ctx, "MinIO object not found for Download", map[string]any{"key": key})
			return nil, nil, fmt.Errorf("minio object %s not found for download: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to get MinIO object for Download", map[string]any{"key": key, "error": err})
		return nil, nil, fmt.Errorf("failed to get MinIO object %s for download: %w", key, err)
	}

	// Get object info for metadata
	objectInfo, err := object.Stat()
	if err != nil {
		object.Close()
		p.logger.Errorf(ctx, "Failed to get MinIO object stats for Download", map[string]any{"key": key, "error": err})
		return nil, nil, fmt.Errorf("failed to get MinIO object stats for %s: %w", key, err)
	}

	fileObject := &port.FileObject{
		Key:          key,
		URL:          p.generateObjectURL(ctx, key),
		Size:         objectInfo.Size,
		ContentType:  objectInfo.ContentType,
		LastModified: objectInfo.LastModified,
		ETag:         strings.Trim(objectInfo.ETag, "\""),
		Provider:     p.ProviderType(),
	}
	p.logger.Infof(ctx, "Prepared MinIO file for download", map[string]any{"key": key})
	return object, fileObject, nil
}

// CheckHealth checks if the MinIO storage provider is healthy and accessible.
func (p *minioProvider) CheckHealth(ctx context.Context) error {
	// First, check if we can list buckets (basic connectivity test)
	buckets, err := p.client.ListBuckets(ctx)
	if err != nil {
		p.logger.Errorf(ctx, "MinIO health check failed - cannot list buckets", map[string]any{"error": err, "endpoint": p.endpointURL})
		return fmt.Errorf("minio health check failed - cannot connect to server: %w", err)
	}

	p.logger.Infof(ctx, "MinIO server connectivity verified", map[string]any{"bucketsCount": len(buckets)})

	// Then, check if our specific bucket exists
	exists, err := p.client.BucketExists(ctx, p.bucketName)
	if err != nil {
		p.logger.Errorf(ctx, "MinIO health check failed - cannot check bucket existence", map[string]any{"error": err, "bucket": p.bucketName})
		return fmt.Errorf("minio health check failed - cannot check bucket '%s' existence: %w", p.bucketName, err)
	}
	if !exists {
		p.logger.Warnf(ctx, "MinIO bucket does not exist", map[string]any{"bucket": p.bucketName})
		return fmt.Errorf("minio bucket '%s' does not exist", p.bucketName)
	}

	p.logger.Infof(ctx, "MinIO health check passed", map[string]any{"bucket": p.bucketName, "endpoint": p.endpointURL})
	return nil
}

// ProviderType returns the type of the storage provider.
func (p *minioProvider) ProviderType() port.StorageProviderType {
	return port.ProviderMinIO
}
