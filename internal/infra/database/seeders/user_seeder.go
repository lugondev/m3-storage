package seeders

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lugondev/m3-storage/internal/infra/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder handles seeding user data
type UserSeeder struct {
	db *gorm.DB
}

// NewUserSeeder creates a new user seeder
func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// Seed seeds user data into the database
func (s *UserSeeder) Seed() error {
	// Check if users already exist to avoid duplicates
	var count int64
	if err := s.db.Model(&database.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Only seed if no users exist
	if count > 0 {
		log.Println("Users already exist, skipping user seeding...")
		return nil
	}

	log.Println("Seeding users...")

	// Hash password for seeded users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []database.User{
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:          "admin@example.com",
			PasswordHash:   string(hashedPassword),
			FirstName:      "Admin",
			LastName:       "User",
			Status:         "active",
			EmailVerified:  true,
			FailedAttempts: 0,
		},
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:          "user@example.com",
			PasswordHash:   string(hashedPassword),
			FirstName:      "Regular",
			LastName:       "User",
			Status:         "active",
			EmailVerified:  true,
			FailedAttempts: 0,
		},
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:          "test@example.com",
			PasswordHash:   string(hashedPassword),
			FirstName:      "Test",
			LastName:       "User",
			Status:         "active",
			EmailVerified:  false,
			FailedAttempts: 0,
		},
	}

	// Create users in batch
	if err := s.db.Create(&users).Error; err != nil {
		return err
	}

	// Create user profiles for seeded users
	profiles := []database.UserProfile{
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserID:      users[0].ID,
			PhoneNumber: "+1234567890",
			Timezone:    "UTC",
			Language:    "en",
		},
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserID:      users[1].ID,
			PhoneNumber: "+0987654321",
			Timezone:    "UTC",
			Language:    "en",
		},
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserID:      users[2].ID,
			PhoneNumber: "",
			Timezone:    "UTC",
			Language:    "en",
		},
	}

	if err := s.db.Create(&profiles).Error; err != nil {
		return err
	}

	log.Printf("Successfully seeded %d users with profiles", len(users))
	return nil
}

// SeedTestData seeds additional test data for development
func (s *UserSeeder) SeedTestData() error {
	log.Println("Seeding test users...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	testUsers := []database.User{
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:          "locked@example.com",
			PasswordHash:   string(hashedPassword),
			FirstName:      "Locked",
			LastName:       "User",
			Status:         "locked",
			EmailVerified:  true,
			FailedAttempts: 5,
			LockedUntil:    &[]time.Time{time.Now().Add(time.Hour)}[0],
		},
		{
			Base: database.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:          "inactive@example.com",
			PasswordHash:   string(hashedPassword),
			FirstName:      "Inactive",
			LastName:       "User",
			Status:         "inactive",
			EmailVerified:  false,
			FailedAttempts: 0,
		},
	}

	if err := s.db.Create(&testUsers).Error; err != nil {
		return err
	}

	log.Printf("Successfully seeded %d test users", len(testUsers))
	return nil
}
