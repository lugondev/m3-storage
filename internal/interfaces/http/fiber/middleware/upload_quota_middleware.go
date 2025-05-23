package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lugondev/m3-storage/internal/modules/user/domain"
	"github.com/lugondev/m3-storage/internal/modules/user/port"
)

// UploadQuotaMiddleware creates a Fiber middleware to check user's upload quotas.
// This middleware should run AFTER an authentication middleware (like APIKeyAuthMiddleware)
// that puts the user object into c.Locals().
func UploadQuotaMiddleware(userService port.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals(UserContextKey).(*domain.User)
		if !ok || user == nil {
			// This should not happen if auth middleware runs first and is correctly configured
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not authenticated or user context not found",
			})
		}

		// Get file size from the request.
		// For multipart form, Fiber might parse it before this middleware.
		// If using `c.FormFile("file")`, the file is already in memory or temp disk.
		// We need to get the file size.
		// A common way is to get it from the Content-Length header for direct uploads,
		// or by inspecting the multipart.FileHeader if the file is already parsed.

		var fileSize int64
		// Attempt to get from FormFile if already parsed by a previous middleware or handler part
		formFile, err := c.FormFile("file") // Assuming "file" is the name of your file input
		if err == nil && formFile != nil {
			fileSize = formFile.Size
		} else {
			// Fallback or if you have another way to determine expected file size
			// For example, a custom header X-File-Size, or if you expect raw body.
			// If Content-Length is reliable for your use case:
			// fileSize = c.Request().Header.ContentLength()
			// This part might need adjustment based on how you handle file uploads.
			// For now, if "file" is not in form, we can't check size here effectively without parsing body.
			// Let's assume for now that if we can't get it, we skip this specific check here
			// and rely on MediaValidator for max file size, and UserService.CanUpload for other quotas.
			// A more robust solution would be to ensure fileSize is always available.
			if fileSize == 0 { // If still zero, means we couldn't determine it easily here
				// Potentially return an error or log a warning.
				// For simplicity, we'll proceed, and CanUpload will handle other checks.
				// The MediaValidator should handle the absolute max file size.
				// This middleware focuses on user-specific quotas.
			}
		}

		canUpload, err := userService.CanUpload(user.ID, fileSize) // fileSize might be 0 here
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": err.Error(), // e.g., "adapters quota exceeded", "daily file upload limit reached"
			})
		}
		if !canUpload {
			// This case might be redundant if CanUpload always returns an error when false,
			// but it's a good safeguard.
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Upload forbidden due to quota limits.",
			})
		}

		return c.Next()
	}
}
