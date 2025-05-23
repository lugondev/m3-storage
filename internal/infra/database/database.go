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

	log.Println("Database connection established successfully.")

	return db, nil
}
