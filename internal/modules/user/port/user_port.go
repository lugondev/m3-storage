package port

import (
	"github.com/lugondev/m3-storage/internal/modules/user/domain"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	RegisterUser(email, password, firstName, lastName string) (*domain.User, error)
	GetUserByID(userID string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GenerateAPIKey(userID string) (string, error)
	GetUserByAPIKey(apiKey string) (*domain.User, error)
	CanUpload(userID string, fileSize int64) (bool, error)
	RecordUpload(userID string, fileSize int64) error
	// Add other methods like UpdateUser, DeleteUser, AuthenticateUser, etc.
}

// UserRepository defines the interface for user data persistence.
// This will be implemented by an adapter in the infra layer.
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(userID string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByAPIKey(apiKey string) (*domain.User, error)
	Update(user *domain.User) error // Needed for updating quotas and API key
	// Add other methods like Delete, etc.
}
