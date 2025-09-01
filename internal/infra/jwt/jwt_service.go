package jwt

import (
	"context"
	"fmt"

	"github.com/lugondev/m3-storage/internal/shared/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService is an implementation of the JWTService interface.
type JWTService struct {
	secretKey     []byte
	signingMethod jwt.SigningMethod
}

// NewJWTService creates a new instance of jwtService.
// In a real application, secretKey, issuer, and TTLs should come from configuration.
func NewJWTService(secretKey string) (*JWTService, error) {
	if secretKey == "" {
		return nil, errors.NewValidationError("jwt secret key cannot be empty")
	}

	return &JWTService{
		secretKey:     []byte(secretKey),
		signingMethod: jwt.SigningMethodHS256, // Using HS256, ensure secret is strong
	}, nil
}

// ValidateToken parses and validates a JWT token string.
func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewValidationError(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return s.secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrTokenMalformed {
			return nil, errors.NewValidationError("malformed token")
		} else if err == jwt.ErrTokenExpired {
			// Handle expired token specifically if needed, e.g., for refresh logic
			return nil, errors.NewValidationError("token expired")
		} else if err == jwt.ErrTokenNotValidYet {
			return nil, errors.NewValidationError("token not valid yet")
		} else {
			return nil, errors.NewValidationError("token validation failed")
		}
	}

	if !token.Valid {
		return nil, errors.NewValidationError("invalid token")
	}

	// Check if UserID is a valid UUID (basic sanity check)
	if uuid.MustParse(claims.Subject) == uuid.Nil {
		return nil, errors.NewValidationError("invalid UserID in token claims")
	}

	return claims, nil
}

// GenerateToken creates and signs a JWT token with the given claims
func (s *JWTService) GenerateToken(ctx context.Context, claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(s.signingMethod, claims)
	return token.SignedString(s.secretKey)
}

// GenerateJTI creates a new unique identifier for a JWT token.
func (s *JWTService) GenerateJTI() uuid.UUID {
	return uuid.New()
}
