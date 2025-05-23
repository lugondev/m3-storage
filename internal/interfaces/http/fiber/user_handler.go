package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lugondev/m3-storage/internal/domain/user"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(app *fiber.App) {
	api := app.Group("/api/users")
	api.Get("/:id", h.GetUser)
	api.Post("/auth", h.Authenticate)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	resp, err := h.userService.GetUser(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(resp)
}

// AuthRequest represents the authentication request body
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User   *user.User `json:"user"`
	APIKey string     `json:"api_key,omitempty"`
}

func (h *UserHandler) Authenticate(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request payload")
	}

	// Authenticate user
	authenticatedUser, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	// Return user info (without password)
	response := AuthResponse{
		User:   authenticatedUser,
		APIKey: authenticatedUser.ApiKey,
	}

	return c.JSON(response)
}
