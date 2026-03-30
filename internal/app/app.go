package app

import (
	"gorm.io/gorm"

	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
	api_handler "github.com/conmeo200/Golang-V1/internal/handler/api"
	"github.com/conmeo200/Golang-V1/internal/handler/web/client"
	"github.com/conmeo200/Golang-V1/internal/handler/web/dashboard"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
)

type App struct {
	UserHandler      *api_handler.UserHandler
	AuthHandler      *api_handler.AuthHandler
	OrderHandler     *api_handler.OrderHandler
	ClientHandler    *client.ClientHandler
	DashboardHandler *dashboard.DashboardHandler
	OrderService     *service.OrderService
}

func NewApp(db *gorm.DB, rabbitMQ *rabbitmq.RabbitMQ) *App {

	// repositories
	userRepo  := repository.NewUserRepository(db)
	authRepo  := repository.NewAuthRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	roleRepo  := repository.NewRoleRepository(db)

	// services
	producer 	 := rabbitmq.NewProducer(rabbitMQ)
	userService  := service.NewUserService(userRepo)
	authService  := service.NewAuthService(authRepo, userRepo, tokenRepo)
	orderService := service.NewOrderService(orderRepo, producer)
	roleService  := service.NewRoleService(roleRepo)

	// handlers
	userHandler  := api_handler.NewUserHandler(userService)
	authHandler  := api_handler.NewAuthHandler(authService)
	orderHandler := api_handler.NewOrderHandler(orderService)
	clientHandler := client.NewClientHandler(authService)
	dashboardHandler := dashboard.NewDashboardHandler(authService, roleService)

	return &App{
		UserHandler:      userHandler,
		AuthHandler:      authHandler,
		OrderHandler:     orderHandler,
		ClientHandler:    clientHandler,
		DashboardHandler: dashboardHandler,
		OrderService:     orderService,
	}
}