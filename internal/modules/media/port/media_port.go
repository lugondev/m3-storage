package port

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/lugondev/m3-storage/internal/modules/media/domain"
	"github.com/lugondev/m3-storage/internal/shared/utils"
)

// MediaService defines the interface for media services.
type MediaService interface {
	UploadFile(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader, providerName string, mediaTypeHint string) (*domain.Media, error)
	ListMedia(ctx context.Context, userID uuid.UUID, query *utils.PaginationQuery) (*utils.Pagination, []*domain.Media, error)
	GetMedia(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) (*domain.Media, error)
	GetPublicMedia(ctx context.Context, mediaID uuid.UUID) (*domain.Media, error)
	DeleteMedia(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) error
}
