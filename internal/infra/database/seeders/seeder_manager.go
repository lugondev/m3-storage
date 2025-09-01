package seeders

import (
	"log"

	"gorm.io/gorm"
)

// SeederManager manages all seeders
type SeederManager struct {
	db         *gorm.DB
	userSeeder *UserSeeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(db *gorm.DB) *SeederManager {
	return &SeederManager{
		db:         db,
		userSeeder: NewUserSeeder(db),
	}
}

// SeedAll runs all seeders
func (sm *SeederManager) SeedAll() error {
	log.Println("Starting database seeding...")

	// Seed users
	if err := sm.userSeeder.Seed(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

// SeedTestData seeds test data for development
func (sm *SeederManager) SeedTestData() error {
	log.Println("Starting test data seeding...")

	// Seed test users
	if err := sm.userSeeder.SeedTestData(); err != nil {
		return err
	}

	log.Println("Test data seeding completed successfully!")
	return nil
}

// SeedProduction seeds minimal production data
func (sm *SeederManager) SeedProduction() error {
	log.Println("Starting production data seeding...")

	// Only seed essential users for production
	if err := sm.userSeeder.Seed(); err != nil {
		return err
	}

	log.Println("Production data seeding completed successfully!")
	return nil
}
