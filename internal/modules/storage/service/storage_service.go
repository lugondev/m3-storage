package service

import (
	"context"
	"sync"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/storage/domain"
	"github.com/lugondev/m3-storage/internal/modules/storage/dto"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
	"github.com/lugondev/m3-storage/internal/shared/errors"
)

// StorageService defines the interface for storage business logic
type StorageService interface {
	CheckHealth(ctx context.Context, req *dto.HealthCheckRequest) (*dto.HealthCheckResponse, error)
	CheckHealthAll(ctx context.Context) (*dto.HealthCheckAllResponse, error)
}

type storageService struct {
	factory port.StorageFactory
	logger  logger.Logger
}

// NewStorageService creates a new instance of StorageService
func NewStorageService(factory port.StorageFactory, logger logger.Logger) StorageService {
	return &storageService{
		factory: factory,
		logger:  logger.WithFields(map[string]any{"component": "StorageService"}),
	}
}

// CheckHealth checks the health of a specific storage provider
func (s *storageService) CheckHealth(ctx context.Context, req *dto.HealthCheckRequest) (*dto.HealthCheckResponse, error) {
	if req.ProviderType == "" {
		return nil, errors.NewBadRequestError("provider_type is required")
	}

	// Convert to domain type
	domainProviderType := domain.StorageProviderType(req.ProviderType)

	// Validate provider type
	if !s.isValidProviderType(domainProviderType) {
		return nil, errors.NewBadRequestError("invalid provider type")
	}

	// Convert back to port type for factory (adapter layer)
	portProviderType := port.StorageProviderType(req.ProviderType)
	provider, err := s.factory.CreateProvider(portProviderType)
	if err != nil {
		s.logger.Errorf(ctx, "Failed to create storage provider", map[string]any{"error": err, "provider_type": domainProviderType})
		return nil, errors.NewBadRequestError("invalid provider type")
	}

	err = provider.CheckHealth(ctx)
	if err != nil {
		s.logger.Errorf(ctx, "Health check failed", map[string]any{"error": err})
		return &dto.HealthCheckResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	return &dto.HealthCheckResponse{
		Status: "healthy",
	}, nil
}

// CheckHealthAll checks the health of all configured storage providers
func (s *storageService) CheckHealthAll(ctx context.Context) (*dto.HealthCheckAllResponse, error) {
	results := make(map[string]dto.HealthCheckResponse)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// List of all provider types to check (using domain types)
	providers := []domain.StorageProviderType{
		domain.ProviderS3,
		domain.ProviderCloudflareR2,
		domain.ProviderLocal,
		domain.ProviderFirebase,
		domain.ProviderAzure,
		domain.ProviderDiscord,
		domain.ProviderScaleway,
		domain.ProviderBackBlaze,
	}

	for _, providerType := range providers {
		wg.Add(1)
		go func(pType domain.StorageProviderType) {
			defer wg.Done()

			req := &dto.HealthCheckRequest{
				ProviderType: string(pType),
			}

			response, err := s.CheckHealth(ctx, req)
			if err != nil {
				mutex.Lock()
				results[string(pType)] = dto.HealthCheckResponse{
					Status:  "error",
					Message: err.Error(),
				}
				mutex.Unlock()
				return
			}

			mutex.Lock()
			results[string(pType)] = *response
			mutex.Unlock()
		}(providerType)
	}

	wg.Wait()

	return &dto.HealthCheckAllResponse{
		Providers: results,
	}, nil
}

// isValidProviderType validates if the provider type is supported
func (s *storageService) isValidProviderType(providerType domain.StorageProviderType) bool {
	validTypes := []domain.StorageProviderType{
		domain.ProviderS3,
		domain.ProviderCloudflareR2,
		domain.ProviderLocal,
		domain.ProviderFirebase,
		domain.ProviderAzure,
		domain.ProviderDiscord,
		domain.ProviderScaleway,
		domain.ProviderBackBlaze,
	}

	for _, validType := range validTypes {
		if providerType == validType {
			return true
		}
	}
	return false
}
