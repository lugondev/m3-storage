package middleware

import (
	"strings"

	"github.com/lugondev/m3-storage/internal/infra/jwt"
	"github.com/lugondev/m3-storage/internal/shared/constants"
	"github.com/lugondev/m3-storage/internal/shared/errors"

	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

// NewAuthMiddleware creates a new instance of AuthMiddleware.
func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware ensures the request has a valid JWT access token
// and stores the validated claims in the context.
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := m.extractToken(c)
		if tokenString == "" {
			return errors.NewUnauthorizedError("missing or malformed authorization header")
		}

		claims, err := m.jwtService.ValidateToken(c.Context(), tokenString)
		if err != nil {
			// Differentiate between expired and other invalid token errors
			if strings.Contains(err.Error(), "token expired") {
				return errors.NewUnauthorizedError("access token expired")
			}
			return errors.NewUnauthorizedError("invalid access token: " + err.Error())
		}

		// Validate that it's an access token
		isAccessToken := slices.Contains(claims.Audience, "access")
		if !isAccessToken {
			return errors.NewUnauthorizedError("invalid token type: expected access token")
		}

		// Store validated claims in context for later use
		c.Locals(constants.UserClaimsKey, claims)
		c.Locals(constants.UserIDKey, claims.Subject)

		return c.Next()
	}
}

// extractToken gets the JWT token from the Authorization header
func (m *AuthMiddleware) extractToken(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") { // Case-insensitive check for "Bearer"
		return ""
	}

	return parts[1]
}

// VerifyCSRF middleware validates CSRF token
// (Implementation placeholder)
func (m *AuthMiddleware) VerifyCSRF() fiber.Handler {
	// TODO: Implement CSRF protection if needed (typically for web forms, less common for APIs)
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// Helper function to get UserClaims from context
func GetUserClaims(c *fiber.Ctx) (*jwt.JWTClaims, error) {
	claims, ok := c.Locals(constants.UserClaimsKey).(*jwt.JWTClaims)
	if !ok || claims == nil {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

// Helper function to get UserID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals(constants.UserIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		// Try getting from claims as a fallback
		claims, err := GetUserClaims(c)
		if err != nil {
			return uuid.Nil, errors.ErrUnauthorized
		}
		userID = uuid.MustParse(claims.Subject)
		if userID == uuid.Nil {
			return uuid.Nil, errors.ErrUnauthorized
		}
		return userID, nil
	}
	return userID, nil
}
