package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
	mediadomain "github.com/lugondev/m3-storage/internal/modules/media/domain"

	"gorm.io/gorm"
)

// ProviderResult holds the results of initializing the database connection.
type ProviderResult struct {
	DB    *gorm.DB
	SqlDB *sql.DB
	Error error
}

// InitializeDatabase sets up the database connection based on the configuration.
// It returns the GORM DB instance, the underlying sql.DB instance, and any error.
func InitializeDatabase(cfg config.Config, log logger.Logger) (*gorm.DB, *sql.DB, error) { // Changed log type
	ctx := context.Background()              // Context for initialization
	db, err := NewDatabaseConnection(cfg.DB) // Use the existing connection function
	if err != nil {
		log.Errorf(ctx, "Failed to connect to database: %v", err)
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the models
	if err := db.AutoMigrate(&mediadomain.Media{}); err != nil {
		log.Errorf(ctx, "Failed to auto-migrate Media model: %v", err)
		return nil, nil, fmt.Errorf("failed to auto-migrate Media model: %w", err)
	}

	// The User, UserProfile, and AuditLog models are auto-migrated in database.go autoMigrate function

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf(ctx, "Failed to get underlying database instance: %v", err)
		// Attempt to close the gorm connection if getting sql.DB fails
		if closeErr := Close(db); closeErr != nil {
			log.Error(ctx, "Failed to close GORM DB after failing to get sql.DB", map[string]any{
				"closeError": closeErr,
			})
		}
		return nil, nil, fmt.Errorf("failed to get underlying database instance: %w", err)
	}

	log.Info(ctx, "Database connection established successfully")
	return db, sqlDB, nil
}

// Close safely closes the database connection.
// This function might already exist in database.go or needs to be added.
// For now, defining it here for completeness of the provider concept.
func Close(db *gorm.DB) error {
	if db == nil {
		return nil // Nothing to close
	}
	sqlDB, err := db.DB()
	if err != nil {
		// Log the error but proceed, as we still want to attempt closing
		fmt.Fprintf(os.Stderr, "Failed to get underlying sql.DB for closing: %v\n", err)
		return fmt.Errorf("failed to get underlying sql.DB for closing: %w", err)
	}
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil // Should not happen if db was not nil, but handle defensively
}

// CloseSqlDB safely closes the sql.DB connection.
func CloseSqlDB(sqlDB *sql.DB) error {
	if sqlDB == nil {
		return nil
	}
	return sqlDB.Close()
}

// ExitOnError logs a fatal error and exits if err is not nil.
// Helper function to reduce repetition in main or setup functions.
func ExitOnError(log logger.Logger, msg string, err error) { // Changed log type
	if err != nil {
		log.Errorf(context.Background(), "%s: %v", msg, err)
		// Consider if Fatalf is more appropriate if the logger interface supports it and exit is desired
		// log.Fatalf(context.Background(), "%s: %v", msg, err) // Alternative if Fatalf exists and is desired
		os.Exit(1) // Keep os.Exit for now
	}
}
