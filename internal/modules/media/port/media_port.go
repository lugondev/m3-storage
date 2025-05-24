package port

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/lugondev/m3-storage/internal/modules/media/domain"
)

// MediaService defines the interface for media services.
type MediaService interface {
	UploadFile(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader, providerName string, mediaTypeHint string) (*domain.Media, error)
}
