package dependencies

import (
	"context"
	"fmt"
	"os"

	"github.com/lugondev/m3-storage/internal/interfaces/http/fiber/middleware"

	// External Libs
	// Alias for firebase app

	"gorm.io/gorm"

	"github.com/lugondev/m3-storage/internal/infra/cache"
	"github.com/lugondev/m3-storage/internal/infra/config" // For NewFirebaseStorageService
	infraJWT "github.com/lugondev/m3-storage/internal/infra/jwt"
	notifyPort "github.com/lugondev/m3-storage/internal/modules/notify/port"
	notifyService "github.com/lugondev/m3-storage/internal/modules/notify/service"

	logger "github.com/lugondev/go-log"

	// App Ports & Services (Health, Storage)
	appPort "github.com/lugondev/m3-storage/internal/modules/app/port"

	// Media Module
	mediaHandler "github.com/lugondev/m3-storage/internal/modules/media/handler"
	mediaPort "github.com/lugondev/m3-storage/internal/modules/media/port"
	mediaService "github.com/lugondev/m3-storage/internal/modules/media/service"
	storageFactory "github.com/lugondev/m3-storage/internal/modules/storage/factory"

	// User Module
	userHandler "github.com/lugondev/m3-storage/internal/modules/user/handler"
	userPort "github.com/lugondev/m3-storage/internal/modules/user/port"
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
	StorageSvc appPort.StorageService
	JWTSvc     *infraJWT.JWTService
	NotifySvc  notifyPort.NotifyService
	MediaSvc   mediaPort.MediaService
	UserSvc    userPort.UserService

	// Handlers
	MediaHandler mediaPort.MediaHandler
	UserHandler  userHandler.UserHandler

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
	app.NotifySvc, err = notifyService.NewNotifyService(infra.Config, log)
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

	// --- Initialize Handlers ---
	// app.UserHandler = *userHandler.NewUserHandler(app.JWTSvc, log) // Assuming UserHandler is initialized similarly

	// Initialize Media Module
	// Assuming StorageFactory is available or initialized elsewhere (e.g., in Infrastructure or passed to BuildDependencies)
	// For now, let's assume a storageFactoryInstance is available.
	// If it's part of `infra`, it should be `infra.StorageFactory`.
	// This needs to be resolved based on where StorageFactory is actually initialized.
	// Let's assume it's created here for now or passed in.
	// If StorageFactory is part of `infra`, then `infra.StorageFactory` should be used.
	// We need a concrete instance of StorageFactory.
	// For the purpose of this example, let's assume NewStorageFactory() can be called.
	// Initialize Storage Factory with config and logger
	sFactory := storageFactory.NewStorageFactory(infra.Config, log)

	app.MediaSvc = mediaService.NewMediaService(infra.DB, log, sFactory) // Pass infra.DB
	app.MediaHandler = mediaHandler.NewMediaHandler(log, app.MediaSvc)
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
