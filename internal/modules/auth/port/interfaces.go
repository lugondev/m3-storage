package port

import (
	"context"
	"time"

	"github.com/lugondev/m3-storage/internal/infra/jwt"
	"github.com/lugondev/m3-storage/internal/modules/auth/domain"

	"github.com/google/uuid"
)

// UserRepository defines the contract for user data persistence
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error

	// UpdateFailedAttempts updates the failed login attempts
	UpdateFailedAttempts(ctx context.Context, id uuid.UUID, attempts int) error

	// LockUser locks a user account until the specified time
	LockUser(ctx context.Context, id uuid.UUID, lockedUntil *time.Time) error
}

// UserProfileRepository defines the contract for user profile data persistence
type UserProfileRepository interface {
	// Create creates a new user profile
	Create(ctx context.Context, profile *domain.UserProfile) error

	// GetByUserID retrieves a user profile by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error)

	// Update updates an existing user profile
	Update(ctx context.Context, profile *domain.UserProfile) error

	// Delete deletes a user profile
	Delete(ctx context.Context, userID uuid.UUID) error
}

// AuthService defines the contract for authentication operations
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)

	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)

	// RefreshToken generates new tokens using refresh token
	RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.LoginResponse, error)

	// ChangePassword changes user's password
	ChangePassword(ctx context.Context, userID uuid.UUID, req *domain.ChangePasswordRequest) error

	// ForgotPassword initiates password reset process
	ForgotPassword(ctx context.Context, req *domain.ForgotPasswordRequest) error

	// ResetPassword resets password using reset token
	ResetPassword(ctx context.Context, req *domain.ResetPasswordRequest) error

	// GetProfile retrieves user profile
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, *domain.UserProfile, error)

	// UpdateProfile updates user profile
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *domain.UpdateProfileRequest) error

	// ValidateToken validates JWT token and returns claims
	ValidateToken(ctx context.Context, tokenString string) (*jwt.JWTClaims, error)
}
