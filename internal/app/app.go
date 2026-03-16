package app

import (
	"gorm.io/gorm"

	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/conmeo200/Golang-V1/internal/handler"
)

type App struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
	WebHandler  *handler.WebHandler
}

func NewApp(db *gorm.DB) *App {

	// repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo, tokenRepo)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	webHandler  := handler.NewWebHandler()

	return &App{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		WebHandler:  webHandler,
	}
}