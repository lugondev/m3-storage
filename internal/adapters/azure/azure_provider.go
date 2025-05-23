package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"

	"github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// azureProvider implements the port.StorageProvider interface for Azure Blob Storage.
type azureProvider struct {
	client        *azblob.Client  // Client for service, container, and blob operations
	serviceClient *service.Client // More specific client for service-level operations
	containerName string
	accountName   string
	logger        logger.Logger
}

// NewAzureProvider creates a new instance of azureProvider.
// Config map is expected to contain:
// - "account_name": (string) Azure Storage account name (required)
// - "account_key": (string) Azure Storage account key (required)
// - "container_name": (string) Azure Blob Storage container name (required)
// - "logger": (logger.Logger) A logger instance (required)
// - "service_url": (string) Optional: custom service URL (e.g., for Azurite emulator or specific clouds)
// If service_url is not provided, it defaults to "https://<account_name>.blob.core.windows.net/"
func NewAzureProvider(config map[string]interface{}) (port.StorageProvider, error) {
	log, ok := config["logger"].(logger.Logger)
	if !ok || log == nil {
		return nil, fmt.Errorf("logger is required in config for AzureProvider")
	}
	log = log.WithFields(map[string]any{"component": "AzureProvider"})

	accountName, ok := config["account_name"].(string)
	if !ok || accountName == "" {
		return nil, fmt.Errorf("account_name is required for AzureProvider")
	}
	accountKey, ok := config["account_key"].(string)
	if !ok || accountKey == "" {
		return nil, fmt.Errorf("account_key is required for AzureProvider")
	}
	containerName, ok := config["container_name"].(string)
	if !ok || containerName == "" {
		return nil, fmt.Errorf("container_name is required for AzureProvider")
	}

	serviceURL, _ := config["service_url"].(string)
	if serviceURL == "" {
		serviceURL = fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	}

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Errorf(context.Background(), "Failed to create Azure shared key credential", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create Azure shared key credential: %w", err)
	}

	// Client for all operations
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		log.Errorf(context.Background(), "Failed to create Azure Blob client", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	// Service client specifically for generating user delegation keys for SAS if needed, or service level properties
	serviceClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		log.Errorf(context.Background(), "Failed to create Azure Blob service client", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create Azure Blob service client: %w", err)
	}

	log.Infof(context.Background(), "AzureProvider initialized", map[string]any{"account": accountName, "container": containerName})
	return &azureProvider{
		client:        client,
		serviceClient: serviceClient,
		containerName: containerName,
		accountName:   accountName,
		logger:        log,
	}, nil
}

// getContainerClient returns a client for the specific container.
func (p *azureProvider) getContainerClient() *container.Client {
	return p.client.ServiceClient().NewContainerClient(p.containerName)
}

// getBlobClient returns a client for a specific blob.
func (p *azureProvider) getBlobClient(key string) *blob.Client {
	return p.getContainerClient().NewBlobClient(key)
}

// Upload uploads a file to Azure Blob Storage.
func (p *azureProvider) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *port.UploadOptions) (*port.FileObject, error) {
	if key == "" {
		return nil, fmt.Errorf("upload key cannot be empty")
	}

	blobClient := p.getBlobClient(key)

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

	// Convert map[string]string to map[string]*string for Azure SDK
	metadata := make(map[string]*string)
	if opts != nil && opts.Metadata != nil {
		for k, v := range opts.Metadata {
			value := v
			metadata[k] = &value
		}
	}

	uploadOpts := &blockblob.UploadStreamOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
		Metadata: metadata,
	}

	p.logger.Infof(ctx, "Attempting to upload file to Azure Blob Storage", map[string]any{"key": key, "container": p.containerName, "contentType": contentType})

	blockBlobClient := p.getContainerClient().NewBlockBlobClient(key)
	_, err := blockBlobClient.UploadStream(ctx, reader, uploadOpts)

	if err != nil {
		p.logger.Errorf(ctx, "Failed to upload file to Azure Blob Storage", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to upload to Azure Blob Storage key %s: %w", key, err)
	}

	properties, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to get properties after Azure upload", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to get properties for Azure key %s: %w", key, err)
	}

	p.logger.Infof(ctx, "File uploaded successfully to Azure Blob Storage", map[string]any{"key": key})
	return &port.FileObject{
		Key:          key,
		URL:          blobClient.URL(), // This is the direct blob URL
		Size:         *properties.ContentLength,
		ContentType:  *properties.ContentType,
		LastModified: *properties.LastModified,
		ETag:         string(*properties.ETag),
		Provider:     p.ProviderType(),
	}, nil
}

// GetURL returns a publicly accessible URL for the given key.
// This URL is accessible if the container/blob has public access enabled.
func (p *azureProvider) GetURL(ctx context.Context, key string) (string, error) {
	blobClient := p.getBlobClient(key)
	// Check if blob exists
	_, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			p.logger.Warnf(ctx, "Azure blob not found, cannot get URL", map[string]any{"key": key})
			return "", fmt.Errorf("azure blob %s not found: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to get Azure blob properties for URL", map[string]any{"key": key, "error": err})
		return "", fmt.Errorf("failed to check Azure blob %s: %w", key, err)
	}
	return blobClient.URL(), nil
}

// GetSignedURL generates a time-limited SAS URL for accessing a private blob.
func (p *azureProvider) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	blobClient := p.getBlobClient(key)

	// User delegation SAS is more secure if Azure AD is set up.
	// For simplicity, using account key SAS here.
	// To use User Delegation SAS, you'd get a user delegation key from serviceClient.GetUserDelegationKey()
	// and then use blob.GetSASURL (with UserDelegationCredential).

	permissions := sas.BlobPermissions{Read: true}
	startTime := time.Now().Add(-10 * time.Minute) // SAS start time, slightly in the past
	expiryTime := time.Now().Add(duration)

	sasOptions := &blob.GetSASURLOptions{
		StartTime: &startTime,
	}
	sasURL, err := blobClient.GetSASURL(permissions, expiryTime, sasOptions)
	if err != nil {
		p.logger.Errorf(ctx, "Failed to generate Azure Blob SAS URL", map[string]any{"key": key, "error": err})
		return "", fmt.Errorf("failed to generate Azure SAS URL for key %s: %w", key, err)
	}

	p.logger.Infof(ctx, "Generated Azure Blob SAS URL", map[string]any{"key": key, "duration": duration})
	return sasURL, nil
}

// Delete removes a file from Azure Blob Storage.
func (p *azureProvider) Delete(ctx context.Context, key string) error {
	blobClient := p.getBlobClient(key)
	_, err := blobClient.Delete(ctx, nil)
	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			p.logger.Warnf(ctx, "Attempted to delete non-existent Azure blob", map[string]any{"key": key})
			return nil // Consistent with other providers: no error if not found
		}
		p.logger.Errorf(ctx, "Failed to delete Azure blob", map[string]any{"key": key, "error": err})
		return fmt.Errorf("failed to delete Azure blob %s: %w", key, err)
	}
	p.logger.Infof(ctx, "Azure blob deleted successfully", map[string]any{"key": key})
	return nil
}

// GetObject retrieves file information (metadata) from Azure Blob Storage.
func (p *azureProvider) GetObject(ctx context.Context, key string) (*port.FileObject, error) {
	blobClient := p.getBlobClient(key)
	properties, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			p.logger.Warnf(ctx, "Azure blob not found for GetObject", map[string]any{"key": key})
			return nil, fmt.Errorf("azure blob %s not found: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to get Azure blob properties", map[string]any{"key": key, "error": err})
		return nil, fmt.Errorf("failed to get Azure blob properties for %s: %w", key, err)
	}

	return &port.FileObject{
		Key:          key,
		URL:          blobClient.URL(),
		Size:         *properties.ContentLength,
		ContentType:  *properties.ContentType,
		LastModified: *properties.LastModified,
		ETag:         string(*properties.ETag),
		Provider:     p.ProviderType(),
	}, nil
}

// Download downloads a file from Azure Blob Storage.
func (p *azureProvider) Download(ctx context.Context, key string) (io.ReadCloser, *port.FileObject, error) {
	blobClient := p.getBlobClient(key)
	downloadResponse, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			p.logger.Warnf(ctx, "Azure blob not found for Download", map[string]any{"key": key})
			return nil, nil, fmt.Errorf("azure blob %s not found for download: %w", key, err)
		}
		p.logger.Errorf(ctx, "Failed to download Azure blob stream", map[string]any{"key": key, "error": err})
		return nil, nil, fmt.Errorf("failed to download Azure blob %s: %w", key, err)
	}

	// Get properties to populate FileObject
	properties, propErr := blobClient.GetProperties(ctx, nil)
	if propErr != nil {
		// Log error but proceed with download if stream was obtained
		p.logger.Errorf(ctx, "Failed to get properties during Azure download, but stream obtained", map[string]any{"key": key, "error": propErr})
		// Try to close the body if properties failed, to avoid resource leak from DownloadStream
		if downloadResponse.Body != nil {
			downloadResponse.Body.Close()
		}
		return nil, nil, fmt.Errorf("failed to get properties for Azure blob %s during download: %w", key, propErr)
	}

	fileObject := &port.FileObject{
		Key:          key,
		URL:          blobClient.URL(),
		Size:         *properties.ContentLength,
		ContentType:  *properties.ContentType,
		LastModified: *properties.LastModified,
		ETag:         string(*properties.ETag),
		Provider:     p.ProviderType(),
	}
	p.logger.Infof(ctx, "Prepared Azure blob for download", map[string]any{"key": key})
	return downloadResponse.Body, fileObject, nil
}

// ProviderType returns the type of the adapters provider.
func (p *azureProvider) ProviderType() port.StorageProviderType {
	return port.ProviderAzure
}

// Helper to read all content into a buffer if needed (example, not directly used by UploadStream with io.Reader)
func readToBuffer(r io.Reader) (*bytes.Buffer, int64, error) {
	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(r)
	if err != nil {
		return nil, 0, err
	}
	return buf, size, nil
}
