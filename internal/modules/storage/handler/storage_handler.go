package handler

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

type StorageHandler struct {
	factory port.StorageFactory
	logger  logger.Logger
}

func NewStorageHandler(factory port.StorageFactory, logger logger.Logger) *StorageHandler {
	return &StorageHandler{
		factory: factory,
		logger:  logger.WithFields(map[string]any{"component": "StorageHandler"}),
	}
}

// CheckHealth godoc
// @Summary Check storage provider health
// @Description Check if the storage provider is healthy and accessible
// @Tags storage
// @Accept json
// @Produce json
// @Param provider_type query port.StorageProviderType true "Storage provider type"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /storage/health [get]
func (h *StorageHandler) CheckHealth(c *fiber.Ctx) error {
	rawProviderType := c.Query("provider_type")
	if rawProviderType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "provider_type is required",
		})
	}

	providerType := port.StorageProviderType(rawProviderType)
	provider, err := h.factory.CreateProvider(providerType)
	if err != nil {
		h.logger.Errorf(c.Context(), "Failed to create storage provider", map[string]any{"error": err, "provider_type": providerType})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid provider type: " + rawProviderType,
		})
	}

	err = provider.CheckHealth(c.Context())
	if err != nil {
		h.logger.Errorf(c.Context(), "Health check failed", map[string]any{"error": err})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}
