package handler

import (
	"github.com/gofiber/fiber/v2" // Assuming you are using Fiber

	"github.com/lugondev/go-log"
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

	userGroup.Post("/register", h.Register)
	userGroup.Get("/:id", h.GetUserByID)
	// Add other routes like /login, /profile, etc.
}

// Placeholder for request body struct
type RegisterUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// Placeholder for response body struct
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, first name, and last name.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User registration details"
// @Success 201 {object} UserResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 409 {object} ErrorResponse "User with this email already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {

	// Placeholder for error response
	type ErrorResponse struct {
		Message string `json:"message"`
	}

	var req RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request payload"})
	}

	// TODO: Add validation for the request payload (e.g., using a validator library)

	user, err := h.userService.RegisterUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		// Basic error handling, can be more sophisticated
		if err.Error() == "user with this email already exists" {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: "Failed to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
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
