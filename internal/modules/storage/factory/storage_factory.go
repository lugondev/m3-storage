package factory

import (
	"errors"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/adapters/azure"
	"github.com/lugondev/m3-storage/internal/adapters/discord"
	"github.com/lugondev/m3-storage/internal/adapters/firebase"
	"github.com/lugondev/m3-storage/internal/adapters/local"
	"github.com/lugondev/m3-storage/internal/adapters/s3"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

type storageFactory struct {
	config *config.Config
	logger logger.Logger
}

// NewStorageFactory creates a new instance of StorageFactory.
func NewStorageFactory(cfg *config.Config, log logger.Logger) port.StorageFactory {
	return &storageFactory{
		config: cfg,
		logger: log,
	}
}

// CreateProvider creates a specific storage provider based on the type and config.
// If config parameter is nil, it will use the factory's internal config for the provider.
func (f *storageFactory) CreateProvider(providerType port.StorageProviderType) (port.StorageProvider, error) {
	switch providerType {
	case port.ProviderLocal:
		return local.NewLocalStorageProvider(f.config.LocalStorage)
	case port.ProviderS3:
		return s3.NewS3Provider(f.config.S3, f.logger)
	case port.ProviderCloudflareR2:
		// Cloudflare R2 is S3 compatible, so it will use the S3 provider with R2 config
		return s3.NewS3Provider(f.config.Cloudflare.ToS3Config(), f.logger)
	case port.ProviderFirebase:
		return firebase.NewFirebaseProvider(f.config.FireStore, f.logger)
	case port.ProviderAzure:
		return azure.NewAzureProvider(&f.config.Azure, f.logger)
	case port.ProviderDiscord:
		return discord.NewDiscordProvider(f.config.Discord, f.logger)
	default:
		return nil, errors.New("unsupported storage provider type for default config: " + string(providerType))
	}
}

// GetDefaultProvider returns the default storage provider based on configuration
func (f *storageFactory) GetDefaultProvider() (port.StorageProvider, error) {
	// Default to local storage if no specific provider is configured
	return f.CreateProvider(port.ProviderLocal)
}
