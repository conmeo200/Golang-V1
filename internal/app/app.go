package app

import (
	"gorm.io/gorm"

	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
	api_handler "github.com/conmeo200/Golang-V1/internal/handler/api"
	"github.com/conmeo200/Golang-V1/internal/handler/web/client"
	"github.com/conmeo200/Golang-V1/internal/handler/web/dashboard"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	//"github.com/conmeo200/Golang-V1/internal/worker"
)

type App struct {
	UserHandler      *api_handler.UserHandler
	AuthHandler      *api_handler.AuthHandler
	OrderHandler     *api_handler.OrderHandler
	PaymentHandler   *api_handler.PaymentHandler
	ClientHandler    *client.ClientHandler
	DashboardHandler *dashboard.DashboardHandler
	OrderService     *service.OrderService
}

func NewApp(db *gorm.DB, rabbitMQ *rabbitmq.RabbitMQ) *App {

	// repositories
	userRepo     := repository.NewUserRepository(db)
	authRepo     := repository.NewAuthRepository(db)
	tokenRepo    := repository.NewTokenRepository(db)
	orderRepo    := repository.NewOrderRepository(db)
	roleRepo     := repository.NewRoleRepository(db)
	taxRepo          := repository.NewTaxRepository(db)
	transactionRepo  := repository.NewTransactionRepository(db)
	paymentRepo      := repository.NewPaymentRepository(db)
	//paymentEventRepo := repository.NewPaymentEventRepository(db)
	outboxRepo       := repository.NewOutboxEventRepository(db)
	inboxRepo        := repository.NewInboxEventRepository(db)
	//webhookLogRepo   := repository.NewWebhookLogRepository(db)

	// services
	producer     := rabbitmq.NewProducer(rabbitMQ)
	userService  := service.NewUserService(userRepo)
	authService  := service.NewAuthService(authRepo, userRepo, tokenRepo)
	orderService := service.NewOrderService(orderRepo, producer)
	roleService  := service.NewRoleService(roleRepo)
	logService   := service.NewFileLogService("log")
	taxService         := service.NewTaxService(taxRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	paymentService     := service.NewPaymentService(paymentRepo, outboxRepo, inboxRepo)

	// handlers
	userHandler     := api_handler.NewUserHandler(userService)
	authHandler     := api_handler.NewAuthHandler(authService)
	orderHandler    := api_handler.NewOrderHandler(orderService)
	paymentHandler  := api_handler.NewPaymentHandler()
	clientHandler   := client.NewClientHandler(authService)
	dashboardHandler := dashboard.NewDashboardHandler(authService, roleService, logService, taxService, orderService, transactionService, paymentService)

	// worker manager
	// workerManager := worker.NewManager(rabbitMQ, orderService, outboxRepo)
	// workerManager.Start()

	return &App{
		UserHandler:      userHandler,
		AuthHandler:      authHandler,
		OrderHandler:     orderHandler,
		PaymentHandler:   paymentHandler,
		ClientHandler:    clientHandler,
		DashboardHandler: dashboardHandler,
		OrderService:     orderService,
	}
}