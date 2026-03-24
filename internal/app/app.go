package app

import (
	"gorm.io/gorm"

	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
	api_handler "github.com/conmeo200/Golang-V1/internal/handler/api"
	web_handler "github.com/conmeo200/Golang-V1/internal/handler/web"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
)

type App struct {
	UserHandler  *api_handler.UserHandler
	AuthHandler  *api_handler.AuthHandler
	OrderHandler *api_handler.OrderHandler
	WebHandler   *web_handler.WebHandler
	OrderService *service.OrderService
}

func NewApp(db *gorm.DB, rabbitMQ *rabbitmq.RabbitMQ) *App {

	// repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// services
	producer := rabbitmq.NewProducer(rabbitMQ)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo, tokenRepo)
	orderService := service.NewOrderService(orderRepo, producer)

	// handlers
	userHandler := api_handler.NewUserHandler(userService)
	authHandler := api_handler.NewAuthHandler(authService)
	orderHandler := api_handler.NewOrderHandler(orderService)
	webHandler  := web_handler.NewWebHandler()

	return &App{
		UserHandler:  userHandler,
		AuthHandler:  authHandler,
		OrderHandler: orderHandler,
		WebHandler:   webHandler,
		OrderService: orderService,
	}
}