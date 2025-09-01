package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lugondev/m3-storage/internal/infra/jwt"
	"github.com/lugondev/m3-storage/internal/modules/auth/domain"
	"github.com/lugondev/m3-storage/internal/modules/auth/port"
	"github.com/lugondev/m3-storage/internal/shared/errors"

	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// MaxFailedAttempts before locking account
	MaxFailedAttempts = 5
	// AccountLockDuration for locked accounts
	AccountLockDuration = 30 * time.Minute
	// AccessTokenDuration for access tokens
	AccessTokenDuration = 15 * time.Minute
	// RefreshTokenDuration for refresh tokens
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
)

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	userRepo        port.UserRepository
	userProfileRepo port.UserProfileRepository
	jwtService      *jwt.JWTService
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo port.UserRepository,
	userProfileRepo port.UserProfileRepository,
	jwtService *jwt.JWTService,
) port.AuthService {
	return &AuthServiceImpl{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
		jwtService:      jwtService,
	}
}

// Register creates a new user account
func (s *AuthServiceImpl) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewConflictError("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to hash password")
	}

	// Create user entity
	user := &domain.User{
		ID:            uuid.New(),
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Status:        domain.UserStatusActive,
		EmailVerified: false, // In production, require email verification
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.WrapError(err, 500, "failed to create user")
	}

	// Create basic user profile
	profile := &domain.UserProfile{
		UserID:    user.ID,
		Language:  "en", // Default language
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userProfileRepo.Create(ctx, profile); err != nil {
		// Log error but don't fail registration
		// In production, consider using transactions or saga pattern
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthServiceImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Check if account can login
	if !user.CanLogin() {
		if user.IsLocked() {
			return nil, errors.NewUnauthorizedError("account is temporarily locked")
		}
		return nil, errors.NewUnauthorizedError("account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Increment failed attempts
		user.IncrementFailedAttempts()

		// Lock account if max attempts reached
		if user.FailedAttempts >= MaxFailedAttempts {
			user.LockAccount(AccountLockDuration)
		}

		// Update failed attempts in database
		s.userRepo.UpdateFailedAttempts(ctx, user.ID, user.FailedAttempts)
		if user.LockedUntil != nil {
			s.userRepo.LockUser(ctx, user.ID, user.LockedUntil)
		}

		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Reset failed attempts on successful login
	if user.FailedAttempts > 0 {
		user.ResetFailedAttempts()
		s.userRepo.UpdateFailedAttempts(ctx, user.ID, 0)
		s.userRepo.LockUser(ctx, user.ID, nil)
	}

	// Update last login
	user.UpdateLastLogin()
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to generate tokens")
	}

	// Remove sensitive information
	user.PasswordHash = ""

	return &domain.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenDuration.Seconds()),
		User:         user,
	}, nil
}

// RefreshToken generates new tokens using refresh token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid refresh token")
	}

	// Check if it's a refresh token
	isRefreshToken := false
	for _, aud := range claims.Audience {
		if aud == "refresh" {
			isRefreshToken = true
			break
		}
	}
	if !isRefreshToken {
		return nil, errors.NewUnauthorizedError("invalid token type")
	}

	// Get user
	userID := uuid.MustParse(claims.Subject)
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("user not found")
	}

	if !user.CanLogin() {
		return nil, errors.NewUnauthorizedError("account not active")
	}

	// Generate new tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to generate tokens")
	}

	// Remove sensitive information
	user.PasswordHash = ""

	return &domain.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenDuration.Seconds()),
		User:         user,
	}, nil
}

// ChangePassword changes user's password
func (s *AuthServiceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, req *domain.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.NewUnauthorizedError("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalServerError("failed to hash password")
	}

	// Update user password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// ForgotPassword initiates password reset process
func (s *AuthServiceImpl) ForgotPassword(ctx context.Context, req *domain.ForgotPasswordRequest) error {
	// In production, you would:
	// 1. Generate a secure reset token
	// 2. Store it in database with expiration
	// 3. Send email with reset link

	// For now, just check if user exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// TODO: Implement email service integration
	return nil
}

// ResetPassword resets password using reset token
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, req *domain.ResetPasswordRequest) error {
	// TODO: Implement password reset with token validation
	return errors.NewNotImplementedError("password reset not implemented")
}

// GetProfile retrieves user profile
func (s *AuthServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, *domain.UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, errors.NewNotFoundError("user not found")
	}

	profile, err := s.userProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Profile might not exist, return user only
		user.PasswordHash = ""
		return user, nil, nil
	}

	user.PasswordHash = ""
	return user, profile, nil
}

// UpdateProfile updates user profile
func (s *AuthServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req *domain.UpdateProfileRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}

	// Update user basic info
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.WrapError(err, 500, "failed to update user")
	}

	// Get or create profile
	profile, err := s.userProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create new profile
		profile = &domain.UserProfile{
			UserID:    userID,
			CreatedAt: time.Now(),
		}
	}

	// Update profile info
	if req.PhoneNumber != "" {
		profile.PhoneNumber = req.PhoneNumber
	}
	if req.Timezone != "" {
		profile.Timezone = req.Timezone
	}
	if req.Language != "" {
		profile.Language = req.Language
	}
	profile.UpdatedAt = time.Now()

	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = time.Now()
		return s.userProfileRepo.Create(ctx, profile)
	} else {
		return s.userProfileRepo.Update(ctx, profile)
	}
}

// ValidateToken validates JWT token and returns claims
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*jwt.JWTClaims, error) {
	return s.jwtService.ValidateToken(ctx, tokenString)
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// generateTokens creates access and refresh tokens for a user
func (s *AuthServiceImpl) generateTokens(user *domain.User) (*TokenPair, error) {
	now := time.Now()

	// Create access token claims
	accessClaims := &jwt.JWTClaims{
		Email: user.Email,
		RegisteredClaims: jwtLib.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    "m3-storage",
			Audience:  []string{"access"},
			ExpiresAt: jwtLib.NewNumericDate(now.Add(AccessTokenDuration)),
			NotBefore: jwtLib.NewNumericDate(now),
			IssuedAt:  jwtLib.NewNumericDate(now),
			ID:        s.jwtService.GenerateJTI().String(),
		},
	}

	// Create refresh token claims
	refreshClaims := &jwt.JWTClaims{
		Email: user.Email,
		RegisteredClaims: jwtLib.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    "m3-storage",
			Audience:  []string{"refresh"},
			ExpiresAt: jwtLib.NewNumericDate(now.Add(RefreshTokenDuration)),
			NotBefore: jwtLib.NewNumericDate(now),
			IssuedAt:  jwtLib.NewNumericDate(now),
			ID:        s.jwtService.GenerateJTI().String(),
		},
	}

	// Generate and sign tokens
	accessToken, err := s.jwtService.GenerateToken(context.Background(), accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateToken(context.Background(), refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
