package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lugondev/m3-storage/internal/modules/user/port"
)

const (
	APIKeyHeader   = "X-API-Key"
	UserContextKey = "user" // Key to store user object in Fiber context
)

// APIKeyAuthMiddleware creates a Fiber middleware for API key authentication.
func APIKeyAuthMiddleware(userService port.UserService) fiber.Handler {
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

		user, err := userService.GetUserByAPIKey(apiKey)
		if err != nil {
			// Log the error for internal review if needed
			// log.Printf("API key auth failed for key %s: %v", apiKey, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid API key or user not found",
			})
		}

		if user == nil { // Should be caught by GetUserByAPIKey, but as a safeguard
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid API key",
			})
		}

		if !user.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "User account is inactive",
			})
		}

		// Store user information in context for subsequent handlers
		c.Locals(UserContextKey, user)

		return c.Next()
	}
}
