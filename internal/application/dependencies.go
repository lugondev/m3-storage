package application

import (
	"context"
	"fmt"
	"os"

	"github.com/lugondev/m3-storage/internal/presentation/http/fiber/middleware"

	"gorm.io/gorm"

	"github.com/lugondev/m3-storage/internal/infra/cache"
	"github.com/lugondev/m3-storage/internal/infra/config" // For NewFirebaseStorageService
	infraJWT "github.com/lugondev/m3-storage/internal/infra/jwt"

	logger "github.com/lugondev/go-log"
	sen "github.com/lugondev/send-sen"
	senConfig "github.com/lugondev/send-sen/config"

	// App Ports & Services (Health, Storage)
	appPort "github.com/lugondev/m3-storage/internal/modules/app/port"

	// Media Module
	mediaHandler "github.com/lugondev/m3-storage/internal/modules/media/handler"
	mediaPort "github.com/lugondev/m3-storage/internal/modules/media/port"
	mediaService "github.com/lugondev/m3-storage/internal/modules/media/service"

	// Storage Module - DDD compliant
	storageFactory "github.com/lugondev/m3-storage/internal/modules/storage/factory"
	storageHandler "github.com/lugondev/m3-storage/internal/modules/storage/handler"
	storageService "github.com/lugondev/m3-storage/internal/modules/storage/service"
)

// Infrastructure holds the initialized infrastructure components.
type Infrastructure struct {
	Config      *config.Config
	Logger      logger.Logger
	DB          *gorm.DB
	RedisClient *cache.RedisClient
}

// Application holds the initialized application components (services, handlers, etc.).
type Application struct {
	// Services
	CacheSvc   appPort.CacheService
	StorageSvc storageService.StorageService // DDD-compliant storage service
	JWTSvc     *infraJWT.JWTService
	NotifySvc  sen.NotifyService
	MediaSvc   mediaPort.MediaService

	// Handlers
	MediaHandler   *mediaHandler.MediaHandler
	StorageHandler *storageHandler.StorageHandler

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware
}

// BuildDependencies initializes and wires up all application dependencies.
func BuildDependencies(infra *Infrastructure) (*Application, error) {
	app := &Application{}
	if infra.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	cfg := *infra.Config
	log := infra.Logger
	// db := infra.DB
	redisClient := infra.RedisClient // This is the wrapper *cache.RedisClient
	ctx := context.Background()      // Use background context for initialization logs

	// --- Initialize Repositories ---
	// userRepo := userService.NewUserRepository(infra.DB, log)

	// --- Initialize Shared Infrastructure Services ---
	app.CacheSvc = cache.NewRedisCacheService(redisClient) // Pass the wrapper

	// --- Initialize JWT Service ---
	jwtSvc, err := infraJWT.NewJWTService(cfg.App.Secret)
	if err != nil {
		log.Errorf(ctx, "Failed to initialize JWT service: %v", err)
		return nil, fmt.Errorf("failed to initialize JWT service: %w", err)
	}
	app.JWTSvc = jwtSvc
	log.Info(ctx, "JWT Service initialized successfully")

	// Initialize Notify Service
	app.NotifySvc, err = sen.NewNotifyService(senConfig.Config{
		Adapter:  cfg.Adapter,
		Telegram: cfg.Telegram,
	}, log)
	if err != nil {
		log.Warnf(ctx, "Failed to initialize notification service (continuing without it?): %v", err)
	} else {
		log.Info(ctx, "Notification service initialized successfully")
	}

	// --- Initialize Module Services ---
	log.Info(ctx, "Module services initialized")

	// --- Initialize Middleware ---
	app.AuthMiddleware = middleware.NewAuthMiddleware(app.JWTSvc)
	log.Info(ctx, "Custom middleware initialized")

	// --- Initialize Storage Module (DDD-compliant) ---
	// Initialize Storage Factory
	sFactory := storageFactory.NewStorageFactory(infra.Config, log)

	// Initialize Storage Service (Application Layer)
	app.StorageSvc = storageService.NewStorageService(sFactory, log)
	log.Info(ctx, "Storage service initialized")

	// Initialize Storage Handler (Presentation Layer)
	app.StorageHandler = storageHandler.NewStorageHandler(app.StorageSvc, log)
	log.Info(ctx, "Storage handler initialized")

	// --- Initialize Media Module ---
	app.MediaSvc = mediaService.NewMediaService(infra.DB, log, sFactory)
	app.MediaHandler = mediaHandler.NewMediaHandler(log, app.MediaSvc, infra.Config)
	log.Info(ctx, "Media module initialized")

	log.Info(ctx, "Handlers initialized")

	return app, nil
}

func ExitOnError(log logger.Logger, msg string, err error) {
	if err != nil {
		log.Error(context.Background(), msg, map[string]any{
			"error": err,
		})
		os.Exit(1)
	}
}
