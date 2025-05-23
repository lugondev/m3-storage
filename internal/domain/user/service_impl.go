package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) Authenticate(email string, password string) (*User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *userService) ValidateAPIKey(apiKey string) (*User, error) {
	return s.repo.GetByAPIKey(apiKey)
}

func (s *userService) CreateUser(email string, password string, name string, role string) (*User, error) {
	// Check if user already exists
	if existing, _ := s.repo.GetByEmail(email); existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate API key
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	now := time.Now().Unix()
	user := &User{
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		ApiKey:    apiKey,
		Role:      role,
		Quota:     1000000000, // 1GB default quota
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) UpdateUser(id uint, name string, quota int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	user.Name = name
	user.Quota = quota
	user.UpdatedAt = time.Now().Unix()

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) ListUsers() ([]*User, error) {
	return s.repo.List()
}

func (s *userService) UpdateQuota(id uint, newQuota int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	user.Quota = newQuota
	user.UpdatedAt = time.Now().Unix()

	return s.repo.Update(user)
}

func (s *userService) DeductQuota(id uint, amount int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if user.Quota < amount {
		return fmt.Errorf("insufficient quota")
	}

	user.Quota -= amount
	user.UpdatedAt = time.Now().Unix()

	return s.repo.Update(user)
}

// Helper function to generate random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
