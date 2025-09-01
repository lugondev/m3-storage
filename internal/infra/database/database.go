package database

import (
	"fmt"
	"log"
	"time"

	"github.com/lugondev/m3-storage/internal/infra/config"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewDatabaseConnection creates a new GORM database instance based on the provided configuration, instrumented with OpenTelemetry.
func NewDatabaseConnection(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SslMode,
	)

	// Configure GORM logger
	var logLevel gormLogger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = gormLogger.Silent
	case "error":
		logLevel = gormLogger.Error
	case "warn":
		logLevel = gormLogger.Warn
	case "info":
		logLevel = gormLogger.Info
	default: // Treat debug and any other value as Info for GORM
		logLevel = gormLogger.Info
	}

	newLogger := gormLogger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Instrument GORM with OpenTelemetry
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("failed to apply otelgorm plugin: %w", err)
	}
	log.Println("GORM OpenTelemetry instrumentation applied.")

	// Configure Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate database schemas
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	log.Println("Database connection established successfully.")

	return db, nil
}

// autoMigrate performs automatic database schema migration
func autoMigrate(db *gorm.DB) error {
	// First, handle existing users table that might have null emails
	if err := handleExistingUsersTable(db); err != nil {
		return fmt.Errorf("failed to handle existing users table: %w", err)
	}

	return db.AutoMigrate(
		&User{},
		&UserProfile{},
		&AuditLog{},
	)
}

// handleExistingUsersTable handles migration of existing users table
func handleExistingUsersTable(db *gorm.DB) error {
	// Check if users table exists
	if !db.Migrator().HasTable(&User{}) {
		// Table doesn't exist, safe to create new one
		return nil
	}

	log.Println("Found existing users table, preparing migration...")

	// Add email column if missing, then handle null values
	if !db.Migrator().HasColumn(&User{}, "email") {
		log.Println("Adding missing email column")
		sql := "ALTER TABLE users ADD COLUMN email VARCHAR(255) NULL"
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to add email column: %w", err)
		}
	}

	// Update any null/empty emails with placeholder values BEFORE applying NOT NULL constraint
	var count int64
	result := db.Model(&User{}).Where("email IS NULL OR email = ''").Count(&count)
	if result.Error != nil {
		return fmt.Errorf("failed to count null emails: %w", result.Error)
	}

	if count > 0 {
		log.Printf("Found %d users with null/empty emails, updating...", count)
		// Update null/empty emails with placeholder values
		result := db.Exec("UPDATE users SET email = 'user_' || id::text || '@placeholder.com' WHERE email IS NULL OR email = ''")
		if result.Error != nil {
			return fmt.Errorf("failed to update null emails: %w", result.Error)
		}
		log.Printf("Updated %d users with placeholder emails", count)
	}

	// Handle missing columns by adding them with default values
	columnsToAdd := map[string]string{
		"password_hash":   "'needs_reset'",
		"first_name":      "'Unknown'",
		"last_name":       "'User'",
		"status":          "'active'",
		"email_verified":  "false",
		"failed_attempts": "0",
	}

	for column, defaultValue := range columnsToAdd {
		if !db.Migrator().HasColumn(&User{}, column) {
			log.Printf("Adding missing column: %s", column)
			// Use raw SQL to add column with default value
			var sql string
			switch column {
			case "password_hash", "first_name", "last_name", "status":
				sql = fmt.Sprintf("ALTER TABLE users ADD COLUMN IF NOT EXISTS %s VARCHAR(255) DEFAULT %s", column, defaultValue)
			case "email_verified":
				sql = fmt.Sprintf("ALTER TABLE users ADD COLUMN IF NOT EXISTS %s BOOLEAN DEFAULT %s", column, defaultValue)
			case "failed_attempts":
				sql = fmt.Sprintf("ALTER TABLE users ADD COLUMN IF NOT EXISTS %s INTEGER DEFAULT %s", column, defaultValue)
			}

			if err := db.Exec(sql).Error; err != nil {
				return fmt.Errorf("failed to add column %s: %w", column, err)
			}
		}
	}

	// Add timestamp columns if missing
	timestampColumns := []string{"last_login_at", "locked_until"}
	for _, column := range timestampColumns {
		if !db.Migrator().HasColumn(&User{}, column) {
			log.Printf("Adding timestamp column: %s", column)
			sql := fmt.Sprintf("ALTER TABLE users ADD COLUMN IF NOT EXISTS %s TIMESTAMP NULL", column)
			if err := db.Exec(sql).Error; err != nil {
				return fmt.Errorf("failed to add timestamp column %s: %w", column, err)
			}
		}
	}

	log.Println("Users table migration preparation completed")
	return nil
}
