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
	ListProviders(ctx context.Context) (*dto.ListProvidersResponse, error)
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
		domain.ProviderMinIO,
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

// ListProviders returns a list of all available storage providers
func (s *storageService) ListProviders(ctx context.Context) (*dto.ListProvidersResponse, error) {
	providers := []dto.ProviderInfo{
		{
			Type:        string(domain.ProviderS3),
			Name:        "Amazon S3",
			Description: "Amazon Simple Storage Service",
		},
		{
			Type:        string(domain.ProviderCloudflareR2),
			Name:        "Cloudflare R2",
			Description: "Cloudflare R2 Object Storage",
		},
		{
			Type:        string(domain.ProviderLocal),
			Name:        "Local Storage",
			Description: "Local file system storage",
		},
		{
			Type:        string(domain.ProviderFirebase),
			Name:        "Firebase Storage",
			Description: "Google Firebase Cloud Storage",
		},
		{
			Type:        string(domain.ProviderAzure),
			Name:        "Azure Blob Storage",
			Description: "Microsoft Azure Blob Storage",
		},
		{
			Type:        string(domain.ProviderDiscord),
			Name:        "Discord CDN",
			Description: "Discord Content Delivery Network",
		},
		{
			Type:        string(domain.ProviderScaleway),
			Name:        "Scaleway Object Storage",
			Description: "Scaleway Object Storage Service",
		},
		{
			Type:        string(domain.ProviderBackBlaze),
			Name:        "Backblaze B2",
			Description: "Backblaze B2 Cloud Storage",
		},
		{
			Type:        string(domain.ProviderMinIO),
			Name:        "MinIO Object Storage",
			Description: "MinIO High Performance Object Storage",
		},
	}

	return &dto.ListProvidersResponse{
		Providers: providers,
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
		domain.ProviderMinIO,
	}

	for _, validType := range validTypes {
		if providerType == validType {
			return true
		}
	}
	return false
}
