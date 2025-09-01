package service

import (
	"context"
	"time"

	"github.com/lugondev/m3-storage/internal/infra/database"
	"github.com/lugondev/m3-storage/internal/modules/auth/domain"
	"github.com/lugondev/m3-storage/internal/modules/auth/port"
	"github.com/lugondev/m3-storage/internal/shared/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	dbUser := r.domainToDBUser(user)

	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return errors.WrapError(err, 500, "failed to create user")
	}

	// Update domain object with generated fields
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var dbUser database.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user not found")
		}
		return nil, errors.WrapError(err, 500, "failed to get user")
	}

	return r.dbToDomainUser(&dbUser), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser database.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user not found")
		}
		return nil, errors.WrapError(err, 500, "failed to get user")
	}

	return r.dbToDomainUser(&dbUser), nil
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	dbUser := r.domainToDBUser(user)
	dbUser.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(dbUser).Error; err != nil {
		return errors.WrapError(err, 500, "failed to update user")
	}

	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

// Delete soft deletes a user
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&database.User{}, id).Error; err != nil {
		return errors.WrapError(err, 500, "failed to delete user")
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	if err := r.db.WithContext(ctx).Model(&database.User{}).Where("id = ?", id).Update("last_login_at", now).Error; err != nil {
		return errors.WrapError(err, 500, "failed to update last login")
	}

	return nil
}

// UpdateFailedAttempts updates the failed login attempts
func (r *UserRepositoryImpl) UpdateFailedAttempts(ctx context.Context, id uuid.UUID, attempts int) error {
	if err := r.db.WithContext(ctx).Model(&database.User{}).Where("id = ?", id).Update("failed_attempts", attempts).Error; err != nil {
		return errors.WrapError(err, 500, "failed to update failed attempts")
	}

	return nil
}

// LockUser locks a user account until the specified time
func (r *UserRepositoryImpl) LockUser(ctx context.Context, id uuid.UUID, lockedUntil *time.Time) error {
	if err := r.db.WithContext(ctx).Model(&database.User{}).Where("id = ?", id).Update("locked_until", lockedUntil).Error; err != nil {
		return errors.WrapError(err, 500, "failed to lock user")
	}

	return nil
}

// domainToDBUser converts domain user to database user
func (r *UserRepositoryImpl) domainToDBUser(user *domain.User) *database.User {
	return &database.User{
		Base: database.Base{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Status:         string(user.Status),
		EmailVerified:  user.EmailVerified,
		LastLoginAt:    user.LastLoginAt,
		FailedAttempts: user.FailedAttempts,
		LockedUntil:    user.LockedUntil,
	}
}

// dbToDomainUser converts database user to domain user
func (r *UserRepositoryImpl) dbToDomainUser(dbUser *database.User) *domain.User {
	return &domain.User{
		ID:             dbUser.ID,
		Email:          dbUser.Email,
		PasswordHash:   dbUser.PasswordHash,
		FirstName:      dbUser.FirstName,
		LastName:       dbUser.LastName,
		Status:         domain.UserStatus(dbUser.Status),
		EmailVerified:  dbUser.EmailVerified,
		LastLoginAt:    dbUser.LastLoginAt,
		FailedAttempts: dbUser.FailedAttempts,
		LockedUntil:    dbUser.LockedUntil,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}
}

// UserProfileRepositoryImpl implements the UserProfileRepository interface
type UserProfileRepositoryImpl struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new user profile repository
func NewUserProfileRepository(db *gorm.DB) port.UserProfileRepository {
	return &UserProfileRepositoryImpl{db: db}
}

// Create creates a new user profile
func (r *UserProfileRepositoryImpl) Create(ctx context.Context, profile *domain.UserProfile) error {
	dbProfile := r.domainToDBUserProfile(profile)

	if err := r.db.WithContext(ctx).Create(dbProfile).Error; err != nil {
		return errors.WrapError(err, 500, "failed to create user profile")
	}

	// Update domain object with generated fields
	profile.CreatedAt = dbProfile.CreatedAt
	profile.UpdatedAt = dbProfile.UpdatedAt

	return nil
}

// GetByUserID retrieves a user profile by user ID
func (r *UserProfileRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	var dbProfile database.UserProfile

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&dbProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user profile not found")
		}
		return nil, errors.WrapError(err, 500, "failed to get user profile")
	}

	return r.dbToDomainUserProfile(&dbProfile), nil
}

// Update updates an existing user profile
func (r *UserProfileRepositoryImpl) Update(ctx context.Context, profile *domain.UserProfile) error {
	dbProfile := r.domainToDBUserProfile(profile)
	dbProfile.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(dbProfile).Error; err != nil {
		return errors.WrapError(err, 500, "failed to update user profile")
	}

	profile.UpdatedAt = dbProfile.UpdatedAt
	return nil
}

// Delete deletes a user profile
func (r *UserProfileRepositoryImpl) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&database.UserProfile{}).Error; err != nil {
		return errors.WrapError(err, 500, "failed to delete user profile")
	}

	return nil
}

// domainToDBUserProfile converts domain user profile to database user profile
func (r *UserProfileRepositoryImpl) domainToDBUserProfile(profile *domain.UserProfile) *database.UserProfile {
	return &database.UserProfile{
		Base: database.Base{
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		},
		UserID:      profile.UserID,
		Avatar:      profile.Avatar,
		PhoneNumber: profile.PhoneNumber,
		DateOfBirth: profile.DateOfBirth,
		Timezone:    profile.Timezone,
		Language:    profile.Language,
	}
}

// dbToDomainUserProfile converts database user profile to domain user profile
func (r *UserProfileRepositoryImpl) dbToDomainUserProfile(dbProfile *database.UserProfile) *domain.UserProfile {
	return &domain.UserProfile{
		UserID:      dbProfile.UserID,
		Avatar:      dbProfile.Avatar,
		PhoneNumber: dbProfile.PhoneNumber,
		DateOfBirth: dbProfile.DateOfBirth,
		Timezone:    dbProfile.Timezone,
		Language:    dbProfile.Language,
		CreatedAt:   dbProfile.CreatedAt,
		UpdatedAt:   dbProfile.UpdatedAt,
	}
}
