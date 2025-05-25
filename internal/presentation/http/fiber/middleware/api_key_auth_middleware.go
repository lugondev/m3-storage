package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	APIKeyHeader   = "X-API-Key"
	UserContextKey = "user" // Key to store user object in Fiber context
)

// APIKeyAuthMiddleware creates a Fiber middleware for API key authentication.
func APIKeyAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get(APIKeyHeader)

		if apiKey == "" {
			// Allow specific public routes to bypass API key check if needed.
			// Example: if c.Path() == "/public/route" { return c.Next() }
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing API key",
			})
		}

		// Potentially, you might want to remove a "Bearer " prefix if you use it
		apiKey = strings.TrimSpace(apiKey)
		if !strings.HasPrefix(apiKey, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid API key format",
			})
		}
		// Store user information in context for subsequent handlers
		// c.Locals(UserContextKey, user)

		return c.Next()
	}
}
