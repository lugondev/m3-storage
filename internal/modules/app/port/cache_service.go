package port

import (
	"context"
	"time"
)

// CacheService defines the interface for caching operations
type CacheService interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (any, error)

	// Set stores a value in cache with expiration
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Clear removes all values from cache
	Clear(ctx context.Context) error
}
