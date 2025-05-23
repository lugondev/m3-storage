package firebase

import (
	"context"
	"fmt"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// InitializeFirebase initializes the Firebase Admin SDK using the application configuration.
// It returns the Firebase App instance, Auth client, and any error encountered.
func InitializeFirebase(cfg config.FireStoreConfig, log logger.Logger) (*firebase.App, *auth.Client, error) {
	ctx := context.Background()

	// Check if ServiceAccountFile path is provided in config
	if cfg.CredentialsFile == "" {
		log.Warn(ctx, "Firebase ServiceAccountFile path not provided in config. Skipping Firebase Admin SDK initialization.")
		return nil, nil, nil // Return nil for clients, no error if skipping is acceptable
	}

	// Prepare Firebase config and options from the app config
	fbConfig := &firebase.Config{
		ProjectID: cfg.ProjectID,
	}
	opt := option.WithCredentialsFile(cfg.CredentialsFile)

	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, fbConfig, opt)
	if err != nil {
		log.Errorf(ctx, "Error initializing firebase app: %v", err)
		return nil, nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	// Initialize Auth Client
	client, err := app.Auth(ctx)
	if err != nil {
		log.Errorf(ctx, "Error getting Auth client: %v", err)
		// Clean up the initialized app if auth client fails
		// Note: Firebase Go SDK doesn't have an explicit app.Close() or similar.
		// Rely on process termination for cleanup in this case.
		return nil, nil, fmt.Errorf("error getting Auth client: %w", err)
	}

	log.Info(ctx, "Firebase Admin SDK initialized successfully")
	return app, client, nil
}
