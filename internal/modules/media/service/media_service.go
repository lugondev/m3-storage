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
	"gorm.io/gorm"

	"github.com/lugondev/m3-storage/internal/modules/media/domain"
	"github.com/lugondev/m3-storage/internal/modules/media/port"
	storagePort "github.com/lugondev/m3-storage/internal/modules/storage/port"
	"github.com/lugondev/m3-storage/internal/shared/utils"
)

type mediaService struct {
	db             *gorm.DB
	logger         logger.Logger
	storageFactory storagePort.StorageFactory
}

// NewMediaService creates a new MediaService.
func NewMediaService(db *gorm.DB, appLogger logger.Logger, storageFactory storagePort.StorageFactory) port.MediaService {
	return &mediaService{
		db:             db,
		logger:         appLogger.WithFields(map[string]any{"component": "MediaService"}),
		storageFactory: storageFactory,
	}
}

// UploadFile implements port.MediaService.
func (s *mediaService) UploadFile(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader, providerName string, mediaTypeHint string) (*domain.Media, error) {
	s.logger.Info(ctx, "Starting file upload process", map[string]any{
		"userID":        userID.String(),
		"fileName":      fileHeader.Filename,
		"providerName":  providerName,
		"mediaTypeHint": mediaTypeHint,
	})

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
		s.logger.Info(ctx, "No provider specified, using default provider", map[string]any{"defaultProvider": string(defaultProviderType)})
		storageProvider, err = s.storageFactory.CreateProvider(defaultProviderType)
	} else {
		s.logger.Info(ctx, "Using specified provider", map[string]any{"providerName": providerName})
		storageProvider, err = s.storageFactory.CreateProvider(storagePort.StorageProviderType(providerName))
	}

	if err != nil {
		s.logger.Error(ctx, "Failed to get adapters provider", map[string]any{"error": err, "providerName": providerName})
		return nil, fmt.Errorf("failed to get adapters provider '%s': %w", providerName, err)
	}
	actualProviderName := string(storageProvider.ProviderType())
	s.logger.Info(ctx, "Using adapters provider", map[string]any{"provider": actualProviderName})

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
				s.logger.Warn(ctx, "Could not determine media type from extension or content type", map[string]any{"fileName": fileHeader.Filename})
			}
		}
	}
	if determinedMediaType == "" {
		determinedMediaType = "other" // Default if still undetermined
	}
	s.logger.Info(ctx, "Determined media type", map[string]any{"mediaType": determinedMediaType})

	// 3. Create adapters path: {userID}/{mediaType}/{date}/{fileName}
	// Sanitize filename to prevent path traversal or invalid characters
	safeFileName := filepath.Base(fileHeader.Filename) // Ensures only the filename part is used

	dateStr := time.Now().Format("20060102") // YYYYMMDD
	storagePathKey := fmt.Sprintf("%s/%s/%s/%s", userID.String(), determinedMediaType, dateStr, safeFileName)
	s.logger.Info(ctx, "Generated adapters path key", map[string]any{"storagePathKey": storagePathKey})

	// 4. Upload file
	file, err := fileHeader.Open()
	if err != nil {
		s.logger.Error(ctx, "Failed to open file header", map[string]any{"error": err})
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
		s.logger.Error(ctx, "Failed to upload file to provider", map[string]any{"error": err, "provider": actualProviderName, "path": storagePathKey})
		return nil, fmt.Errorf("failed to upload file to provider '%s': %w", actualProviderName, err)
	}
	s.logger.Info(ctx, "File uploaded successfully", map[string]any{"fileURL": fileObject.URL, "signedURL": fileObject.SignedURL})

	// 5. Create media metadata
	// Use fileObject.URL or fileObject.SignedURL depending on whether you want public or temporary access
	// For now, let's assume PublicURL should be the direct URL if available, otherwise SignedURL or an internal identifier.
	// This might need adjustment based on how you want to expose URLs.
	publicAccessURL := fileObject.URL
	if publicAccessURL == "" && fileObject.SignedURL != "" {
		// Fallback to signed URL if direct public URL is not available, though this is temporary.
		// Consider if this is the desired behavior.
		// publicAccessURL = fileObject.SignedURL
		s.logger.Warn(ctx, "Public URL not available, consider using signed URL or other mechanism", map[string]any{"key": storagePathKey})
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
		s.logger.Error(ctx, "Failed to save media metadata to database", map[string]any{"error": err})
		// Optional: Attempt to delete the uploaded file from adapters if DB save fails
		// if delErr := provider.Delete(storagePath); delErr != nil {
		// s.logger.Error("Failed to delete uploaded file after DB error", map[string]any{"error": delErr, "path": storagePath})
		// }
		return nil, fmt.Errorf("failed to save media metadata: %w", err)
	}
	s.logger.Info(ctx, "Media metadata saved to database", map[string]any{"mediaID": mediaEntity.ID.String()})

	return mediaEntity, nil
}

// ListMedia returns paginated media files for a given user
func (s *mediaService) ListMedia(ctx context.Context, userID uuid.UUID, query *utils.PaginationQuery) (*utils.Pagination, []*domain.Media, error) {
	s.logger.Info(ctx, "Listing media files for user", map[string]any{
		"userID":   userID.String(),
		"page":     query.Page,
		"pageSize": query.PageSize,
	})

	// Validate and set default pagination values
	query.ValidateAndSetDefaults()

	var totalItems int64
	if err := s.db.Model(&domain.Media{}).Where("user_id = ?", userID).Count(&totalItems).Error; err != nil {
		s.logger.Error(ctx, "Failed to count total media files", map[string]any{"error": err})
		return nil, nil, fmt.Errorf("failed to count media files: %w", err)
	}

	var mediaFiles []*domain.Media
	if err := s.db.Where("user_id = ?", userID).
		Limit(query.GetLimit()).
		Offset(query.GetOffset()).
		Find(&mediaFiles).Error; err != nil {
		s.logger.Error(ctx, "Failed to list media files", map[string]any{"error": err})
		return nil, nil, fmt.Errorf("failed to list media files: %w", err)
	}

	pagination := utils.NewPagination(*query, totalItems)

	return &pagination, mediaFiles, nil
}

// GetMedia returns a specific media file by ID for a given user
func (s *mediaService) GetMedia(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) (*domain.Media, error) {
	s.logger.Info(ctx, "Getting media file", map[string]any{
		"userID":  userID.String(),
		"mediaID": mediaID.String(),
	})

	var media domain.Media
	if err := s.db.Where("id = ? AND user_id = ?", mediaID, userID).First(&media).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn(ctx, "Media file not found", map[string]any{
				"mediaID": mediaID.String(),
				"userID":  userID.String(),
			})
			return nil, fmt.Errorf("media file not found")
		}
		s.logger.Error(ctx, "Failed to get media file", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to get media file: %w", err)
	}

	return &media, nil
}

// DeleteMedia deletes a specific media file by ID for a given user
func (s *mediaService) DeleteMedia(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) error {
	s.logger.Info(ctx, "Deleting media file", map[string]any{
		"userID":  userID.String(),
		"mediaID": mediaID.String(),
	})

	// First get the media to ensure it exists and belongs to the user
	media, err := s.GetMedia(ctx, userID, mediaID)
	if err != nil {
		return err // Already logged in GetMedia
	}

	// Get the storage provider
	storageProvider, err := s.storageFactory.CreateProvider(storagePort.StorageProviderType(media.Provider))
	if err != nil {
		s.logger.Error(ctx, "Failed to get storage provider", map[string]any{
			"error":    err,
			"provider": media.Provider,
		})
		return fmt.Errorf("failed to get storage provider: %w", err)
	}

	// Delete from storage
	if err := storageProvider.Delete(ctx, media.FilePath); err != nil {
		s.logger.Error(ctx, "Failed to delete file from storage", map[string]any{
			"error":    err,
			"filePath": media.FilePath,
		})
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	if err := s.db.Delete(media).Error; err != nil {
		s.logger.Error(ctx, "Failed to delete media from database", map[string]any{"error": err})
		return fmt.Errorf("failed to delete media from database: %w", err)
	}

	s.logger.Info(ctx, "Media file deleted successfully", map[string]any{"mediaID": mediaID.String()})
	return nil
}
