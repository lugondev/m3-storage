package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
	"github.com/lugondev/m3-storage/internal/shared/errors"
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
// @Failure default {object} errors.Error
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
		return errors.NewBadRequestError("invalid provider type")
	}

	err = provider.CheckHealth(c.Context())
	if err != nil {
		h.logger.Errorf(c.Context(), "Health check failed", map[string]any{"error": err})
		return errors.NewInternalServerError("health check failed")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}

// CheckHealthAll godoc
// @Summary Check all storage providers health
// @Description Check if all configured storage providers are healthy and accessible
// @Tags storage
// @Accept json
// @Produce json
// @Success 200 {object} map[string]map[string]string
// @Failure default {object} errors.Error
// @Router /storage/health/all [get]
func (h *StorageHandler) CheckHealthAll(c *fiber.Ctx) error {
	results := make(map[string]map[string]string)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// List of all provider types to check
	providers := []port.StorageProviderType{
		port.ProviderS3,
		port.ProviderCloudflareR2,
		port.ProviderLocal,
		port.ProviderFirebase,
		port.ProviderAzure,
		port.ProviderDiscord,
		port.ProviderScaleway,
		port.ProviderBackBlaze,
	}

	for _, providerType := range providers {
		wg.Add(1)
		go func(pType port.StorageProviderType) {
			defer wg.Done()

			status := make(map[string]string)
			provider, err := h.factory.CreateProvider(pType)

			if err != nil {
				status["status"] = "error"
				status["message"] = "Failed to create provider: " + err.Error()
				mutex.Lock()
				results[string(pType)] = status
				mutex.Unlock()
				return
			}

			err = provider.CheckHealth(c.Context())
			if err != nil {
				status["status"] = "error"
				status["message"] = err.Error()
			} else {
				status["status"] = "healthy"
			}

			mutex.Lock()
			results[string(pType)] = status
			mutex.Unlock()
		}(providerType)
	}

	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"providers": results,
	})
}
