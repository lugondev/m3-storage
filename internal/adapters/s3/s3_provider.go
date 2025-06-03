package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// s3Provider implements the port.StorageProvider interface for AWS S3.
type s3Provider struct {
	client         *s3.Client
	presignClient  *s3.PresignClient
	uploader       *manager.Uploader
	bucketName     string
	region         string
	endpointURL    string // Optional: for S3-compatible services like MinIO or Cloudflare R2
	forcePathStyle bool   // Optional: for S3-compatible services
	logger         logger.Logger
}

// NewS3Provider creates a new instance of s3Provider.
func NewS3Provider(cfg config.S3Config, log logger.Logger) (port.StorageProvider, error) {
	log = log.WithFields(map[string]any{"component": "S3Provider"})

	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for S3Provider")
	}
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("bucket_name is required for S3Provider")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("access_key_id is required for S3Provider")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("secret_access_key is required for S3Provider")
	}

	// Use config struct fields
	bucketName := cfg.BucketName
	region := cfg.Region
	accessKeyID := cfg.AccessKeyID
	secretAccessKey := cfg.SecretAccessKey
	endpointURL := cfg.Endpoint
	forcePathStyle := cfg.ForcePathStyle

	if region == "" && endpointURL == "" {
		// Region is typically required for AWS S3.
		// For S3-compatible services, if endpoint is given, region might be dummy or specific to that service.
		log.Warn(context.Background(), "S3 region is not set and no custom endpoint provided. This might lead to issues.", nil)
		// return nil, fmt.Errorf("region is required for S3Provider if endpoint is not specified")
	}

	cfgLoadOpts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}
	if region != "" {
		cfgLoadOpts = append(cfgLoadOpts, awsConfig.WithRegion(region))
	}

	if endpointURL != "" {
		// Use the modern BaseEndpoint approach instead of deprecated EndpointResolver
		cfgLoadOpts = append(cfgLoadOpts, awsConfig.WithBaseEndpoint(endpointURL))
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), cfgLoadOpts...)
	if err != nil {
		log.Errorf(context.Background(), "Failed to load AWS SDK config", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	if strings.Contains(endpointURL, "backblaze") {
		// Backblaze B2 S3-compatible storage requires specific settings
		awsCfg.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		awsCfg.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if forcePathStyle {
			o.UsePathStyle = true
		}
	})
	presignClient := s3.NewPresignClient(s3Client)
	uploader := manager.NewUploader(s3Client)

	log.Infof(context.Background(), "S3Provider initialized", map[string]any{"bucket": bucketName, "region": region, "endpoint": endpointURL})
	return &s3Provider{
		client:         s3Client,
		presignClient:  presignClient,
		uploader:       uploader,
		bucketName:     bucketName,
		region:         region,
		endpointURL:    endpointURL,
		forcePathStyle: forcePathStyle,
		logger:         log,
	}, nil
}

// Upload uploads a file to S3.
func (p *s3Provider) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *port.UploadOptions) (*port.FileObject, error) {
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

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
		// ContentLength: aws.Int64(size), // manager.Uploader handles this, but can be set.
	}

	if opts != nil {
		if opts.ACL != "" {
			uploadInput.ACL = types.ObjectCannedACL(opts.ACL)
		}
		if opts.Metadata != nil {
			uploadInput.Metadata = opts.Metadata
		}
	}

	p.logger.Infof(ctx, "Attempting to upload file to S3", map[string]any{"key": key, "bucket": p.bucketName, "contentType": contentType})
	result, err := p.uploader.Upload(ctx, uploadInput)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to upload file to S3", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to upload to S3 key %s: %w", key, err)
	}

	// After upload, get object attributes to populate FileObject
	headObjectOutput, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		p.logger.Errorf(ctx, "Failed to get object metadata after S3 upload", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to get metadata for S3 key %s: %w", key, err)
	}

	fileURL := p.generateObjectURL(ctx, key)
	if p.endpointURL != "" && !strings.HasPrefix(p.endpointURL, "https://s3.") && !strings.HasSuffix(p.endpointURL, ".amazonaws.com") {
		// For non-AWS S3-compatible services, the URL might be different
		fileURL = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(p.endpointURL, "/"), p.bucketName, strings.TrimPrefix(key, "/"))
	}

	p.logger.Infof(ctx, "File uploaded successfully to S3", map[string]any{"key": key, "location": result.Location, "versionId": result.VersionID})
	return &port.FileObject{
		Key:          key,
		URL:          fileURL,
		Size:         aws.ToInt64(headObjectOutput.ContentLength),
		ContentType:  aws.ToString(headObjectOutput.ContentType),
		LastModified: aws.ToTime(headObjectOutput.LastModified),
		ETag:         strings.Trim(aws.ToString(headObjectOutput.ETag), "\""), // ETag often comes with quotes
		Provider:     p.ProviderType(),
	}, nil
}

func (p *s3Provider) generateObjectURL(ctx context.Context, key string) string {
	// Standard S3 URL format: https://<bucket-name>.s3.<region>.amazonaws.com/<key>
	// Or path-style: https://s3.<region>.amazonaws.com/<bucket-name>/<key>
	// If a custom endpoint is used, it might be different.
	if p.endpointURL != "" {
		if p.forcePathStyle {
			return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(p.endpointURL, "/"), p.bucketName, strings.TrimPrefix(key, "/"))
		}
		// Attempt to construct a virtual-hosted-style URL if endpoint is like a base S3 domain
		// This is a simplification; complex endpoint configurations might need more robust logic.
		parsedEndpoint, err := url.Parse(p.endpointURL)
		if err == nil {
			return fmt.Sprintf("%s://%s.%s/%s", parsedEndpoint.Scheme, p.bucketName, parsedEndpoint.Host, strings.TrimPrefix(key, "/"))
		}
		// Fallback for custom endpoint if parsing fails or it's not a typical S3 domain structure
		return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(p.endpointURL, "/"), p.bucketName, strings.TrimPrefix(key, "/"))
	}

	// Default AWS S3 URL
	if p.region == "" { // Should not happen if configured correctly for AWS
		// Using Warn without context as this is an internal utility method
		p.logger.Warn(nil, "S3 region is empty, cannot construct standard AWS S3 URL accurately.", map[string]any{"key": key})
		return "" // Or some other placeholder
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucketName, p.region, strings.TrimPrefix(key, "/"))
}

// GetURL returns a publicly accessible URL for the given key.
func (p *s3Provider) GetURL(ctx context.Context, key string) (string, error) {
	// This typically returns the same as generateObjectURL if the object is public.
	// For S3, ACLs determine public accessibility.
	// We can check if object exists first.
	_, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		var nf *types.NotFound
		if errors.As(err, &nsk) || errors.As(err, &nf) {
			p.logger.Warnf(ctx, "Object not found, cannot get URL", map[string]any{"key": key})
			return "", fmt.Errorf("object %s not found: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to HeadObject for GetURL", map[string]any{"key": key, "error": err})
		return "", fmt.Errorf("failed to check object %s: %w", key, err)
	}
	return p.generateObjectURL(ctx, key), nil
}

// GetSignedURL generates a time-limited signed URL for accessing a private object.
func (p *s3Provider) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	request, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		p.logger.Errorf(ctx, "Failed to generate S3 signed URL", map[string]any{"key": key, "error": err})
		return "", fmt.Errorf("failed to generate S3 signed URL for key %s: %w", key, err)
	}
	p.logger.Infof(ctx, "Generated S3 signed URL", map[string]any{"key": key, "duration": duration})
	return request.URL, nil
}

// Delete removes a file from S3.
func (p *s3Provider) Delete(ctx context.Context, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// S3 DeleteObject doesn't typically error if the object doesn't exist,
		// but checking HeadObject first or inspecting error details might be needed
		// if strict "not found" behavior is required.
		p.logger.Errorf(ctx, "Failed to delete S3 object", map[string]any{"key": key, "error": err})
		return fmt.Errorf("failed to delete S3 object %s: %w", key, err)
	}
	p.logger.Infof(ctx, "S3 object deleted successfully", map[string]any{"key": key})
	return nil
}

// GetObject retrieves file information (metadata) from S3.
func (p *s3Provider) GetObject(ctx context.Context, key string) (*port.FileObject, error) {
	headOutput, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		var nf *types.NotFound
		if errors.As(err, &nsk) || errors.As(err, &nf) {
			p.logger.Warnf(ctx, "S3 object not found for GetObject", map[string]any{"key": key})
			return nil, fmt.Errorf("s3 object %s not found: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to HeadObject for S3 GetObject", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to get S3 object metadata for %s: %w", key, err)
	}

	return &port.FileObject{
		Key:          key,
		URL:          p.generateObjectURL(ctx, key),
		Size:         aws.ToInt64(headOutput.ContentLength),
		ContentType:  aws.ToString(headOutput.ContentType),
		LastModified: aws.ToTime(headOutput.LastModified),
		ETag:         strings.Trim(aws.ToString(headOutput.ETag), "\""),
		Provider:     p.ProviderType(),
	}, nil
}

// Download downloads a file from S3.
func (p *s3Provider) Download(ctx context.Context, key string) (io.ReadCloser, *port.FileObject, error) {
	getObjectOutput, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		var nf *types.NotFound
		if errors.As(err, &nsk) || errors.As(err, &nf) {
			p.logger.Warnf(ctx, "S3 object not found for Download", map[string]any{"key": key})
			return nil, nil, fmt.Errorf("s3 object %s not found for download: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to GetObject for S3 Download", map[string]any{"key": key, "error": err})
		return nil, nil, fmt.Errorf("failed to get S3 object %s for download: %w", key, err)
	}

	fileObject := &port.FileObject{
		Key:          key,
		URL:          p.generateObjectURL(ctx, key),
		Size:         aws.ToInt64(getObjectOutput.ContentLength),
		ContentType:  aws.ToString(getObjectOutput.ContentType),
		LastModified: aws.ToTime(getObjectOutput.LastModified),
		ETag:         strings.Trim(aws.ToString(getObjectOutput.ETag), "\""),
		Provider:     p.ProviderType(),
	}
	p.logger.Infof(ctx, "Prepared S3 file for download", map[string]any{"key": key})
	return getObjectOutput.Body, fileObject, nil
}

// CheckHealth checks if the storage provider is healthy and accessible.
func (p *s3Provider) CheckHealth(ctx context.Context) error {
	var err error
	if p.ProviderType() != port.ProviderScaleway {
		_, err = p.client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(p.bucketName),
		})
	} else {
		_, err = p.client.ListBuckets(ctx, &s3.ListBucketsInput{
			MaxBuckets: aws.Int32(1),
		})
	}
	if err != nil {
		p.logger.Errorf(ctx, "S3 health check failed", map[string]any{"error": err})
		return fmt.Errorf("s3 health check failed: %w", err)
	}
	return nil
}

// ProviderType returns the type of the adapters provider.
func (p *s3Provider) ProviderType() port.StorageProviderType {
	// If this provider is also used for Cloudflare R2, this might need adjustment
	// or Cloudflare R2 could have its own type.
	if strings.Contains(p.endpointURL, "r2.cloudflarestorage.com") {
		return port.ProviderCloudflareR2
	}
	// If using a custom endpoint, we can check for specific providers
	if strings.Contains(p.endpointURL, "backblaze.com") {
		return port.ProviderBackBlaze
	}
	if strings.Contains(p.endpointURL, "scw.cloud") || strings.Contains(p.endpointURL, "scaleway.com") {
		return port.ProviderScaleway
	}
	return port.ProviderS3
}
