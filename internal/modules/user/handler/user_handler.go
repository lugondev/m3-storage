package handler

import (
	"github.com/gofiber/fiber/v2" // Assuming you are using Fiber

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/user/port"
	// Import a DTO package if you create one for request/response structs
)

// UserHandler handles HTTP requests for user-related operations.
type UserHandler struct {
	userService port.UserService
	logger      logger.Logger
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService port.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger: logger.WithFields(map[string]any{
			"module": "user",
		}),
	}
}

// RegisterRoutes registers the user routes to the Fiber app.
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	userGroup := app.Group("/users") // Or /api/v1/users

	userGroup.Get("/:id", h.GetUserByID)
	// Add other routes like /login, /profile, etc.
}

// Placeholder for response body struct
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Placeholder for error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve user details by their unique ID.
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "Successfully retrieved user"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {

	userID := c.Params("id")
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}
