package handler

import (
	"github.com/lugondev/m3-storage/internal/modules/auth/domain"
	"github.com/lugondev/m3-storage/internal/modules/auth/port"
	"github.com/lugondev/m3-storage/internal/presentation/http/fiber/middleware"
	"github.com/lugondev/m3-storage/internal/shared/errors"
	"github.com/lugondev/m3-storage/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication related HTTP requests
type AuthHandler struct {
	authService port.AuthService
	validator   validator.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService port.AuthService, validator validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration request"
// @Success 201 {object} domain.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.NewBadRequestError("invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return errors.NewValidationError(err.Error())
	}

	user, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	// Remove sensitive data before response
	user.PasswordHash = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
		"message": "User registered successfully",
	})
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.NewBadRequestError("invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return errors.NewValidationError(err.Error())
	}

	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Login successful",
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new tokens using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.NewBadRequestError("invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return errors.NewValidationError(err.Error())
	}

	response, err := h.authService.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Token refreshed successfully",
	})
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Authentication
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	user, profile, err := h.authService.GetProfile(c.Context(), userID)
	if err != nil {
		return err
	}

	response := fiber.Map{
		"user": user,
	}

	if profile != nil {
		response["profile"] = profile
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Profile retrieved successfully",
	})
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update current user profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body domain.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req domain.UpdateProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.NewBadRequestError("invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return errors.NewValidationError(err.Error())
	}

	if err := h.authService.UpdateProfile(c.Context(), userID, &req); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
	})
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change current user password
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body domain.ChangePasswordRequest true "Change password request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req domain.ChangePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.NewBadRequestError("invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return errors.NewValidationError(err.Error())
	}

	if err := h.authService.ChangePassword(c.Context(), userID, &req); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}

// ForgotPassword handles forgot password request
// @Summary Forgot password
// @Description Initiate password reset process
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req domain.ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.NewBadRequestError("invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return errors.NewValidationError(err.Error())
	}

	if err := h.authService.ForgotPassword(c.Context(), &req); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password reset instructions sent to your email",
	})
}

// Logout handles user logout (token invalidation would be handled by client or Redis blacklist)
// @Summary User logout
// @Description Logout user (client should discard tokens)
// @Tags Authentication
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a production system, you might want to:
	// 1. Add token to blacklist in Redis
	// 2. Log the logout event
	// 3. Invalidate all refresh tokens for the user

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}
