package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lugondev/m3-storage/internal/dependencies"
	"github.com/lugondev/m3-storage/internal/interfaces/http/fiber/middleware"
	"github.com/lugondev/m3-storage/internal/router"

	// Infrastructure Providers & Core Infra
	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/cache"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/infra/database"
	"github.com/lugondev/m3-storage/internal/infra/tracer"

	// External Libs
	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	// Docs
	_ "github.com/lugondev/m3-storage/docs"
)

// @title AuthSys API
// @version 1.0
// @description This is the core API for AuthSys platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// --- Load Configuration ---
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// --- Initialize i18n Bundle ---
	i18nBundle, err := initI18n()
	if err != nil {
		fmt.Printf("Failed to initialize i18n bundle: %v\n", err)
		os.Exit(1)
	}

	// --- Initialize Logger & OpenTelemetry ---
	log, otelErr := logger.NewLogger(&logger.Option{})
	if otelErr != nil {
		fmt.Printf("Failed to create OpenTelemetry logger: %v\n", otelErr)
		os.Exit(1)
	}
	var cleanupOtel func(context.Context) error = func(context.Context) error { return nil } // Default to no-op cleanup

	// Check if Signoz is configured (CollectorURL is present and not the default)
	if cfg.Signoz.CollectorURL != "" {
		// Initialize OpenTelemetry (Tracer & Logger Provider)
		otelInitCfg := tracer.OtelInitConfig{
			ServiceName:  "github.com/lugondev/m3-storage",
			CollectorURL: cfg.Signoz.CollectorURL,
			Insecure:     cfg.Signoz.Insecure,
			Headers:      cfg.Signoz.Headers,
		}
		cleanupOtel = tracer.InitOtel(otelInitCfg)
	}

	// Defer OTel cleanup (will be no-op if not configured)
	defer func() {
		if err := cleanupOtel(context.Background()); err != nil {
			// Use fmt for critical cleanup errors as logger might be compromised
			fmt.Printf("Error during OpenTelemetry cleanup: %v\n", err)
		}
	}()

	// Defer logger sync (works for both OtelLogger and ZapLogger implementations)
	defer func() {
		if err := log.Sync(); err != nil {
			// Use fmt as logger might be compromised if Sync fails
			fmt.Printf("Failed to sync logger: %v\n", err)
		}
	}()

	// --- Initialize Database ---
	// Pass the initialized logger (either Otel or Zap)
	db, sqlDB, err := database.InitializeDatabase(cfg, log)
	database.ExitOnError(log, "Database initialization failed", err)
	defer database.CloseSqlDB(sqlDB) // Close the underlying sql.DB

	// --- Initialize Redis ---
	redisClient, err := cache.InitializeRedisClient(cfg, log)
	cache.ExitOnError(log, "Redis initialization failed", err)
	defer cache.CloseRedisClient(redisClient, log) // Close the Redis client wrapper

	// --- Build Infrastructure Struct ---
	infra := &dependencies.Infrastructure{
		Config:      &cfg,
		Logger:      log,
		DB:          db,
		RedisClient: redisClient,
	}

	// --- Build Application Dependencies ---
	appDeps, err := dependencies.BuildDependencies(infra)
	if err != nil {
		log.Fatalf(context.Background(), "Failed to build application dependencies: %v", err)
	}

	// --- Initialize Fiber App ---
	app := fiber.New(fiber.Config{
		AppName:           fmt.Sprintf("%s API", cfg.App.Name),
		ErrorHandler:      middleware.ErrorHandler(log),
		StreamRequestBody: true, // Enable streaming for large request bodies
	})

	// --- Setup Middleware ---
	middleware.SetupMiddleware(app, cfg, i18nBundle, log)

	// --- Register API Routes ---
	router.RegisterRoutes(app, &router.RouterConfig{
		AuthMw:         appDeps.AuthMiddleware,
		MediaHandler:   appDeps.MediaHandler,
		UserHandler:    appDeps.UserHandler,
		StorageHandler: appDeps.StorageHandler,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		// Basic DB Ping check (more detailed checks are in HealthService)
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":    "unhealthy",
				"message":   "Database connection failed",
				"error":     err.Error(),
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		// Basic Redis Ping check
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		// Use the underlying client from the wrapper for Ping
		if err := redisClient.Client().Ping(ctx).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":    "unhealthy",
				"message":   "Redis connection failed",
				"error":     err.Error(),
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		// All basic systems operational
		return c.JSON(fiber.Map{
			"status":      "healthy",
			"database":    "ok",
			"redis":       "ok",
			"environment": cfg.App.Env,
			"timestamp":   time.Now(),
		})
	})

	// --- Graceful Shutdown Setup ---
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// --- Start Server ---
	log.Info(context.Background(), "Server starting", map[string]any{
		"port": cfg.App.Port,
	})

	go func() {
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			log.Error(context.Background(), "Server error", map[string]any{
				"error": err,
			})
		}
	}()

	// --- Wait for Interrupt Signal ---
	<-shutdownChan

	// --- Graceful Shutdown ---
	log.Info(context.Background(), "Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Error(context.Background(), "Failed to shutdown server gracefully", map[string]any{
			"error": err,
		})
	}

	log.Info(context.Background(), "Server gracefully stopped")
}

func initI18n() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English) // Default language
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load translation files from the 'locales' directory
	// You might want to make the locale directory configurable
	// and handle errors more gracefully in a production environment.
	_, err := bundle.LoadMessageFile("locales/en.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to load en.toml: %w", err)
	}
	_, err = bundle.LoadMessageFile("locales/vi.toml")
	if err != nil {
		// Log a warning or handle missing locale files as appropriate
		fmt.Printf("Warning: Failed to load vi.toml: %v\n", err)
	}

	return bundle, nil
}
