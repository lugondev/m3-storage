package auth

import (
	"github.com/lugondev/m3-storage/internal/infra/jwt"
	"github.com/lugondev/m3-storage/internal/modules/auth/handler"
	"github.com/lugondev/m3-storage/internal/modules/auth/port"
	"github.com/lugondev/m3-storage/internal/modules/auth/service"
	"github.com/lugondev/m3-storage/internal/shared/validator"

	"gorm.io/gorm"
)

// Dependencies holds all authentication module dependencies
type Dependencies struct {
	UserRepo        port.UserRepository
	UserProfileRepo port.UserProfileRepository
	AuthService     port.AuthService
	AuthHandler     *handler.AuthHandler
}

// NewDependencies creates and wires all authentication dependencies
func NewDependencies(db *gorm.DB, jwtService *jwt.JWTService, validator validator.Validator) *Dependencies {
	// Repositories
	userRepo := service.NewUserRepository(db)
	userProfileRepo := service.NewUserProfileRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, userProfileRepo, jwtService)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, validator)

	return &Dependencies{
		UserRepo:        userRepo,
		UserProfileRepo: userProfileRepo,
		AuthService:     authService,
		AuthHandler:     authHandler,
	}
}
