package app

import (
	"gorm.io/gorm"

	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/conmeo200/Golang-V1/internal/handler"
)

type App struct {
	UserHandler  *handler.UserHandler
	AuthHandler  *handler.AuthHandler
	OrderHandler *handler.OrderHandler
	WebHandler   *handler.WebHandler
}

func NewApp(db *gorm.DB) *App {

	// repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo, tokenRepo)
	orderService := service.NewOrderService(orderRepo)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)
	webHandler  := handler.NewWebHandler()

	return &App{
		UserHandler:  userHandler,
		AuthHandler:  authHandler,
		OrderHandler: orderHandler,
		WebHandler:   webHandler,
	}
}