package router

import (
	"github.com/lugondev/m3-storage/internal/interfaces/http/fiber/middleware"
	mediaPort "github.com/lugondev/m3-storage/internal/modules/media/port"
	userHandler "github.com/lugondev/m3-storage/internal/modules/user/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// RouterConfig holds all the dependencies needed for route registration
type RouterConfig struct {
	AuthMw       *middleware.AuthMiddleware
	MediaHandler mediaPort.MediaHandler
	UserHandler  *userHandler.UserHandler // Use concrete type from handler package
}

// RegisterRoutes centralizes all API route registrations following DDD principles.
// It organizes routes by domain modules and maintains clear separation of concerns.
func RegisterRoutes(app *fiber.App, config *RouterConfig) {
	// Infrastructure routes (non-domain specific)
	registerInfrastructureRoutes(app)

	// API versioning - follows DDD by keeping domain routes versioned
	v1 := app.Group("/api/v1")

	// Register domain-specific route groups
	registerMediaRoutes(v1, config.AuthMw, config.MediaHandler)
	registerUserRoutes(v1, config.AuthMw, config.UserHandler)
}

// registerInfrastructureRoutes handles non-domain specific routes
func registerInfrastructureRoutes(app *fiber.App) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)
}

// registerMediaRoutes handles all media domain routes
// This follows DDD by grouping routes by domain context
func registerMediaRoutes(api fiber.Router, authMw *middleware.AuthMiddleware, handler mediaPort.MediaHandler) {
	if handler == nil {
		return // Gracefully handle missing handlers
	}

	mediaRoutes := api.Group("/media")

	// Media upload operations - core domain functionality
	mediaRoutes.Post("/upload", handler.UploadFile)

	// TODO: Add other media operations following RESTful patterns
	// mediaRoutes.Get("/", authMw.RequireAuth(), handler.ListMedia)
	// mediaRoutes.Get("/:id", authMw.RequireAuth(), handler.GetMedia)
	// mediaRoutes.Delete("/:id", authMw.RequireAuth(), handler.DeleteMedia)
}

// registerUserRoutes handles all user domain routes
// This follows DDD by keeping user operations in their own context
func registerUserRoutes(api fiber.Router, authMw *middleware.AuthMiddleware, handler *userHandler.UserHandler) {
	if handler == nil {
		return // Gracefully handle missing handlers
	}

	userRoutes := api.Group("/users")

	// User management operations
	userRoutes.Post("/register", handler.Register)
	userRoutes.Get("/:id", authMw.RequireAuth(), handler.GetUserByID)

	// TODO: Add authentication routes
	// authRoutes := api.Group("/auth")
	// authRoutes.Post("/login", handler.Login)
	// authRoutes.Post("/refresh", handler.RefreshToken)
	// authRoutes.Post("/logout", authMw.RequireAuth(), handler.Logout)
}
