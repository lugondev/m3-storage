package port

import (
	"context"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/lugondev/m3-storage/internal/modules/media/domain"
)

// MediaHandler defines the interface for media handlers.
type MediaHandler interface {
	UploadFile(c *fiber.Ctx) error
}

// MediaService defines the interface for media services.
type MediaService interface {
	UploadFile(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader, providerName string, mediaTypeHint string) (*domain.Media, error)
}
