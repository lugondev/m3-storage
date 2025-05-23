package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	logger "github.com/lugondev/go-log" // Import custom logger
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/lugondev/m3-storage/internal/modules/media/domain"
	"github.com/lugondev/m3-storage/internal/modules/media/port"
	storagePort "github.com/lugondev/m3-storage/internal/modules/storage/port"
)

type mediaService struct {
	db             *gorm.DB
	logger         logger.Logger // Keep as *zap.Logger for internal use
	storageFactory storagePort.StorageFactory
	// mediaValidator *MediaValidator // Assuming MediaValidator is in the same package or imported
}

// NewMediaService creates a new MediaService.
func NewMediaService(db *gorm.DB, appLogger logger.Logger, storageFactory storagePort.StorageFactory) port.MediaService {
	return &mediaService{
		db:             db,
		logger:         appLogger,
		storageFactory: storageFactory,
	}
}

// UploadFile implements port.MediaService.
func (s *mediaService) UploadFile(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader, providerName string, mediaTypeHint string) (*domain.Media, error) {
	s.logger.Info(ctx, "Starting file upload process",
		zap.String("userID", userID.String()),
		zap.String("fileName", fileHeader.Filename),
		zap.String("providerName", providerName),
		zap.String("mediaTypeHint", mediaTypeHint),
	)

	// 1. Get adapters provider
	// Assuming StorageFactory has a GetProvider method that takes providerName string and returns StorageProvider
	// If providerName is empty, the factory should return the default provider.
	// The actual method name on StorageFactory might be different, e.g., CreateProvider based on a type.
	// For now, let's assume a method GetProvider(name string) (StorageProvider, error) exists or can be added.
	// If StorageFactory only has CreateProvider(type, config), we'd need to map providerName to a type and get config.
	// This part might need adjustment based on the actual StorageFactory implementation.
	var storageProvider storagePort.StorageProvider
	var err error

	// Get storage provider - if no provider specified, use default (local)
	if providerName == "" {
		defaultProviderType := storagePort.ProviderLocal // Default to local storage
		s.logger.Info(ctx, "No provider specified, using default provider", zap.String("defaultProvider", string(defaultProviderType)))
		storageProvider, err = s.storageFactory.CreateProvider(defaultProviderType)
	} else {
		s.logger.Info(ctx, "Using specified provider", zap.String("providerName", providerName))
		storageProvider, err = s.storageFactory.CreateProvider(storagePort.StorageProviderType(providerName))
	}

	if err != nil {
		s.logger.Error(ctx, "Failed to get adapters provider", zap.Error(err), zap.String("providerName", providerName))
		return nil, fmt.Errorf("failed to get adapters provider '%s': %w", providerName, err)
	}
	actualProviderName := string(storageProvider.ProviderType())
	s.logger.Info(ctx, "Using adapters provider", zap.String("provider", actualProviderName))

	// 2. Determine media type
	determinedMediaType := mediaTypeHint
	if determinedMediaType == "" {
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType != "" {
			determinedMediaType = strings.Split(contentType, "/")[0] // "image/png" -> "image"
		} else {
			// Fallback: try to guess from extension
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg":
				determinedMediaType = "image"
			case ".mp4", ".mov", ".avi", ".mkv", ".webm":
				determinedMediaType = "video"
			case ".mp3", ".wav", ".ogg":
				determinedMediaType = "audio"
			case ".pdf", ".doc", ".docx", ".txt", ".csv", ".xls", ".xlsx", ".ppt", ".pptx":
				determinedMediaType = "document"
			default:
				determinedMediaType = "other"
				s.logger.Warn(ctx, "Could not determine media type from extension or content type", zap.String("fileName", fileHeader.Filename))
			}
		}
	}
	if determinedMediaType == "" {
		determinedMediaType = "other" // Default if still undetermined
	}
	s.logger.Info(ctx, "Determined media type", zap.String("mediaType", determinedMediaType))

	// 3. Create adapters path: {userID}/{mediaType}/{date}/{fileName}
	// Sanitize filename to prevent path traversal or invalid characters
	safeFileName := filepath.Base(fileHeader.Filename) // Ensures only the filename part is used

	dateStr := time.Now().Format("20060102") // YYYYMMDD
	storagePathKey := fmt.Sprintf("%s/%s/%s/%s", userID.String(), determinedMediaType, dateStr, safeFileName)
	s.logger.Info(ctx, "Generated adapters path key", zap.String("storagePathKey", storagePathKey))

	// 4. Upload file
	file, err := fileHeader.Open()
	if err != nil {
		s.logger.Error(ctx, "Failed to open file header", zap.Error(err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	uploadOpts := &storagePort.UploadOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
		// Metadata:    nil, // Add custom metadata if needed
		// ACL:         "",  // Set ACL if needed, e.g., "public-read"
	}

	fileObject, err := storageProvider.Upload(ctx, storagePathKey, file, fileHeader.Size, uploadOpts)
	if err != nil {
		s.logger.Error(ctx, "Failed to upload file to provider", zap.Error(err), zap.String("provider", actualProviderName), zap.String("path", storagePathKey))
		return nil, fmt.Errorf("failed to upload file to provider '%s': %w", actualProviderName, err)
	}
	s.logger.Info(ctx, "File uploaded successfully", zap.String("fileURL", fileObject.URL), zap.String("signedURL", fileObject.SignedURL))

	// 5. Create media metadata
	// Use fileObject.URL or fileObject.SignedURL depending on whether you want public or temporary access
	// For now, let's assume PublicURL should be the direct URL if available, otherwise SignedURL or an internal identifier.
	// This might need adjustment based on how you want to expose URLs.
	publicAccessURL := fileObject.URL
	if publicAccessURL == "" && fileObject.SignedURL != "" {
		// Fallback to signed URL if direct public URL is not available, though this is temporary.
		// Consider if this is the desired behavior.
		// publicAccessURL = fileObject.SignedURL
		s.logger.Warn(ctx, "Public URL not available, consider using signed URL or other mechanism", zap.String("key", storagePathKey))
		// For now, if no direct public URL, we might store an empty string or an internal reference.
		// Or, the application logic might always generate a signed URL on demand when access is needed.
		// Let's assume for now that if fileObject.URL is empty, we store it as such.
	}

	mediaEntity := domain.NewMedia(
		userID,
		safeFileName,
		storagePathKey, // Store the relative path (key)
		fileHeader.Size,
		determinedMediaType,
		actualProviderName,
		publicAccessURL, // This could be fileObject.URL or a generated signed URL
	)

	// 6. Save metadata to database
	if err := s.db.Create(mediaEntity).Error; err != nil {
		s.logger.Error(ctx, "Failed to save media metadata to database", zap.Error(err))
		// Optional: Attempt to delete the uploaded file from adapters if DB save fails
		// if delErr := provider.Delete(storagePath); delErr != nil {
		// 	s.logger.Error("Failed to delete uploaded file after DB error", zap.Error(delErr), zap.String("path", storagePath))
		// }
		return nil, fmt.Errorf("failed to save media metadata: %w", err)
	}
	s.logger.Info(ctx, "Media metadata saved to database", zap.String("mediaID", mediaEntity.ID.String()))

	return mediaEntity, nil
}
