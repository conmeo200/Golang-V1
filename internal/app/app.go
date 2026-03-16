package app

import (
	"gorm.io/gorm"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
)

// App holds the application's dependencies.
type App struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
	WebHandler  *handler.WebHandler
}

// NewApp creates a new App with all dependencies initialized.
func NewApp(db *gorm.DB) *App {

	// Initialize repositories (using interfaces)
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db) // Returns AuthRepository interface
	tokenRepo := repository.NewTokenRepository(db) // Returns TokenRepository interface

	// Initialize services (using interfaces)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, tokenRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	webHandler := handler.NewWebHandler()

	return &App{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		WebHandler:  webHandler,
	}
}
