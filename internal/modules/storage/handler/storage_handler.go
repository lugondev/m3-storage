package handler

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/storage/dto"
	"github.com/lugondev/m3-storage/internal/modules/storage/service"
)

type StorageHandler struct {
	storageService service.StorageService
	logger         logger.Logger
}

func NewStorageHandler(storageService service.StorageService, logger logger.Logger) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
		logger:         logger.WithFields(map[string]any{"component": "StorageHandler"}),
	}
}

// CheckHealth godoc
// @Summary Check storage provider health
// @Description Check if the storage provider is healthy and accessible
// @Tags storage
// @Accept json
// @Produce json
// @Param provider_type query string true "Storage provider type"
// @Success 200 {object} dto.HealthCheckResponse
// @Failure default {object} errors.Error
// @Router /storage/health [get]
func (h *StorageHandler) CheckHealth(c *fiber.Ctx) error {
	req := &dto.HealthCheckRequest{
		ProviderType: c.Query("provider_type"),
	}

	response, err := h.storageService.CheckHealth(c.Context(), req)
	if err != nil {
		h.logger.Errorf(c.Context(), "Health check failed", map[string]any{"error": err})
		return err
	}

	if response.Status == "error" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// CheckHealthAll godoc
// @Summary Check all storage providers health
// @Description Check if all configured storage providers are healthy and accessible
// @Tags storage
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthCheckAllResponse
// @Failure default {object} errors.Error
// @Router /storage/health/all [get]
func (h *StorageHandler) CheckHealthAll(c *fiber.Ctx) error {
	response, err := h.storageService.CheckHealthAll(c.Context())
	if err != nil {
		h.logger.Errorf(c.Context(), "Health check all failed", map[string]any{"error": err})
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
