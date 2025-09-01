package router

import (
	authHandler "github.com/lugondev/m3-storage/internal/modules/auth/handler"
	mediaHandler "github.com/lugondev/m3-storage/internal/modules/media/handler"
	storageHandler "github.com/lugondev/m3-storage/internal/modules/storage/handler"
	"github.com/lugondev/m3-storage/internal/presentation/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// RouterConfig holds all the dependencies needed for route registration
type RouterConfig struct {
	AuthMw         *middleware.AuthMiddleware
	AuthHandler    *authHandler.AuthHandler
	MediaHandler   *mediaHandler.MediaHandler
	StorageHandler *storageHandler.StorageHandler
}

// RegisterRoutes centralizes all API route registrations following DDD principles.
// It organizes routes by domain modules and maintains clear separation of concerns.
func RegisterRoutes(app *fiber.App, config *RouterConfig) {
	// Infrastructure routes (non-domain specific)
	registerInfrastructureRoutes(app)

	// API versioning - follows DDD by keeping domain routes versioned
	v1 := app.Group("/api/v1")

	// Register domain-specific route groups
	registerAuthRoutes(v1, config.AuthHandler)
	registerMediaRoutes(v1, config.AuthMw, config.MediaHandler)
	registerStorageRoutes(v1, config.StorageHandler)
}

// registerInfrastructureRoutes handles non-domain specific routes
func registerInfrastructureRoutes(app *fiber.App) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)
}

// registerAuthRoutes handles all authentication domain routes
func registerAuthRoutes(api fiber.Router, handler *authHandler.AuthHandler) {
	authRoutes := api.Group("/auth")

	// Public authentication routes (no auth required)
	authRoutes.Post("/register", handler.Register)
	authRoutes.Post("/login", handler.Login)
	authRoutes.Post("/refresh", handler.RefreshToken)
	authRoutes.Post("/forgot-password", handler.ForgotPassword)

	// Protected authentication routes (auth required)
	// These will be handled by middleware in the handler level for now
	authRoutes.Get("/profile", handler.GetProfile)
	authRoutes.Put("/profile", handler.UpdateProfile)
	authRoutes.Post("/change-password", handler.ChangePassword)
	authRoutes.Post("/logout", handler.Logout)
}

// registerMediaRoutes handles all media domain routes
// This follows DDD by grouping routes by domain context
func registerMediaRoutes(api fiber.Router, authMw *middleware.AuthMiddleware, handler *mediaHandler.MediaHandler) {
	mediaRoutes := api.Group("/media")
	// Media upload operations - core domain functionality
	mediaRoutes.Post("/upload", authMw.RequireAuth(), handler.UploadFile)

	// TODO: Add other media operations following RESTful patterns
	mediaRoutes.Get("/", authMw.RequireAuth(), handler.ListMedia)
	mediaRoutes.Get("/:id", authMw.RequireAuth(), handler.GetMedia)
	mediaRoutes.Get("/:id/file", authMw.RequireAuth(), handler.ServeLocalFile)
	mediaRoutes.Delete("/:id", authMw.RequireAuth(), handler.DeleteMedia)

	// Public routes - no authentication required
	mediaRoutes.Get("/public/:id/file", handler.ServePublicLocalFile)
}

// registerStorageRoutes handles storage-related routes
func registerStorageRoutes(api fiber.Router, handler *storageHandler.StorageHandler) {
	storageRoutes := api.Group("/storage")
	storageRoutes.Get("/providers", handler.ListProviders)
	storageRoutes.Get("/health", handler.CheckHealth)
	storageRoutes.Get("/health/all", handler.CheckHealthAll)
}
