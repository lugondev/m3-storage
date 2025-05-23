package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/lugondev/m3-storage/internal/modules/media/domain"
)

const (
	// DefaultMaxImageSize defines the default maximum size for image files (e.g., 5MB).
	DefaultMaxImageSize = 5 * 1024 * 1024
	// DefaultMaxVideoSize defines the default maximum size for video files (e.g., 50MB).
	DefaultMaxVideoSize = 50 * 1024 * 1024
	// DefaultMaxAudioSize defines the default maximum size for audio files (e.g., 10MB).
	DefaultMaxAudioSize = 10 * 1024 * 1024
	// DefaultMaxDocumentSize defines the default maximum size for document files (e.g., 10MB).
	DefaultMaxDocumentSize = 10 * 1024 * 1024
)

// MediaValidator provides methods to validate media files.
type MediaValidator struct {
	// You can add configurable max sizes here if needed, e.g.:
	// MaxImageSize int64
	// MaxVideoSize int64
}

// NewMediaValidator creates a new MediaValidator.
func NewMediaValidator() *MediaValidator {
	return &MediaValidator{}
}

// ValidateFile checks if the uploaded file is valid based on its extension and size.
func (v *MediaValidator) ValidateFile(fileHeader *multipart.FileHeader) (domain.MediaType, error) {
	if fileHeader == nil {
		return "", errors.New("file header is nil")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	mediaType := domain.GetMediaTypeFromExtension(ext)

	if mediaType == "" {
		return "", errors.New("unsupported file type: " + ext)
	}

	// Validate size
	fileSize := fileHeader.Size
	var maxSize int64

	switch {
	case strings.HasPrefix(string(mediaType), "image/"):
		maxSize = DefaultMaxImageSize // Or v.MaxImageSize if configurable
	case strings.HasPrefix(string(mediaType), "video/"):
		maxSize = DefaultMaxVideoSize // Or v.MaxVideoSize
	case strings.HasPrefix(string(mediaType), "audio/"):
		maxSize = DefaultMaxAudioSize // Or v.MaxAudioSize
	case mediaType == domain.MediaTypePDF,
		mediaType == domain.MediaTypeDOC,
		mediaType == domain.MediaTypeDOCX,
		mediaType == domain.MediaTypeTXT,
		mediaType == domain.MediaTypeMD:
		maxSize = DefaultMaxDocumentSize // Or v.MaxDocumentSize
	default:
		return "", errors.New("cannot determine max size for media type: " + string(mediaType))
	}

	if fileSize == 0 {
		return "", errors.New("file is empty")
	}

	if fileSize > maxSize {
		return "", errors.New("file size exceeds the limit of " + formatBytes(maxSize))
	}

	return mediaType, nil
}

// formatBytes is a helper function to format byte size into a human-readable string.
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// You might want to add more specific validation logic,
// for example, using libraries to check actual file content signatures (magic numbers)
// instead of relying solely on extensions for better security and accuracy.
