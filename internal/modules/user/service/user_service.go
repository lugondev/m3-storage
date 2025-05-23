package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors" // For basic error handling
	"time"

	"golang.org/x/crypto/bcrypt" // For password hashing

	"github.com/lugondev/m3-storage/internal/modules/user/domain"
	"github.com/lugondev/m3-storage/internal/modules/user/port"
)

type userService struct {
	userRepo port.UserRepository
	// Add other dependencies like a logger, etc.
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// RegisterUser handles the business logic for registering a new user.
func (s *userService) RegisterUser(email, password, firstName, lastName string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate API Key
	apiKeyBytes := make([]byte, 16) // 32 characters hex string
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return nil, errors.New("failed to generate API key")
	}
	apiKey := hex.EncodeToString(apiKeyBytes)

	newUser := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,   // Default to active
		APIKey:    apiKey, // Store the plain API key for now. Consider hashing if needed.
		// Default quotas will be set by GORM or DB default values
	}

	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Clear password before returning
	newUser.Password = ""
	// newUser.APIKey = "" // Decide if API key should be returned upon registration
	return newUser, nil
}

// GetUserByID retrieves a user by their ID.
func (s *userService) GetUserByID(userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.Password = "" // Clear password
	return user, nil
}

// GetUserByEmail retrieves a user by their email.
func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.Password = "" // Clear password
	return user, nil
}

// GenerateAPIKey generates a new API key for the user and saves it.
func (s *userService) GenerateAPIKey(userID string) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	apiKeyBytes := make([]byte, 16) // 32 characters hex string
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", errors.New("failed to generate API key")
	}
	newAPIKey := hex.EncodeToString(apiKeyBytes)
	user.APIKey = newAPIKey // Consider hashing if you store hashed API keys

	if err := s.userRepo.Update(user); err != nil {
		return "", errors.New("failed to update user with new API key")
	}
	return newAPIKey, nil
}

// GetUserByAPIKey retrieves a user by their API key.
// Note: If API keys are hashed in DB, this logic needs to change.
// This implementation assumes plain text API keys or that FindByAPIKey handles hashing.
func (s *userService) GetUserByAPIKey(apiKey string) (*domain.User, error) {
	user, err := s.userRepo.FindByAPIKey(apiKey)
	if err != nil {
		return nil, errors.New("invalid API key or user not found")
	}
	if user == nil {
		return nil, errors.New("invalid API key or user not found")
	}
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}
	user.Password = "" // Clear password
	return user, nil
}

// isSameDay checks if two times are on the same calendar day in UTC.
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.UTC().Date()
	y2, m2, d2 := t2.UTC().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// CanUpload checks if the user can upload a file of a given size.
func (s *userService) CanUpload(userID string, fileSize int64) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, errors.New("user not found")
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	// Reset daily upload count if it's a new day
	now := time.Now()
	if !isSameDay(user.LastUploadDate, now) {
		user.UploadedFilesToday = 0
		// user.LastUploadDate will be updated in RecordUpload
	}

	// Check adapters quota
	if user.UsedStorageBytes+fileSize > user.MaxStorageBytes {
		return false, errors.New("adapters quota exceeded")
	}

	// Check daily file count
	if user.UploadedFilesToday >= user.MaxFilesPerDay {
		return false, errors.New("daily file upload limit reached")
	}

	// Kích thước tối đa của file đã được kiểm tra ở media_validator.go
	// Nếu muốn có giới hạn kích thước file theo từng user, thêm trường MaxFileSizePerUpload vào user_entity.go
	// và kiểm tra ở đây:
	// if fileSize > user.MaxFileSizePerUpload {
	// 	return false, errors.New("file size exceeds user's allowed maximum")
	// }

	return true, nil
}

// RecordUpload updates the user's quotas after a successful upload.
func (s *userService) RecordUpload(userID string, fileSize int64) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user == nil {
		return errors.New("user not found")
	}

	now := time.Now()
	// Reset daily upload count if it's a new day (double check, might have been done in CanUpload)
	if !isSameDay(user.LastUploadDate, now) {
		user.UploadedFilesToday = 0
	}

	user.UsedStorageBytes += fileSize
	user.UploadedFilesToday++
	user.LastUploadDate = now

	if err := s.userRepo.Update(user); err != nil {
		// Consider how to handle this error, e.g., rollback upload if possible, or log inconsistency.
		return errors.New("failed to update user upload records")
	}
	return nil
}
