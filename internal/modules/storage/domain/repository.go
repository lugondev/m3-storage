package domain

import "context"

// StorageRepository defines the domain interface for storage data operations
// This follows DDD principles by defining the repository interface in the domain layer
type StorageRepository interface {
	// GetProviderConfig retrieves configuration for a specific provider
	GetProviderConfig(ctx context.Context, providerType StorageProviderType) (map[string]interface{}, error)

	// SaveHealthStatus saves the health status of a provider
	SaveHealthStatus(ctx context.Context, status *HealthStatus) error

	// GetHealthStatus retrieves the last known health status of a provider
	GetHealthStatus(ctx context.Context, providerType StorageProviderType) (*HealthStatus, error)

	// GetAllHealthStatuses retrieves health statuses for all providers
	GetAllHealthStatuses(ctx context.Context) ([]*HealthStatus, error)
}
