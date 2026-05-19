package app

import (
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/stripe"
	"github.com/conmeo200/Golang-V1/internal/module/auth"
	"github.com/conmeo200/Golang-V1/internal/module/news"
	"github.com/conmeo200/Golang-V1/internal/module/order"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/conmeo200/Golang-V1/internal/module/role"
	"github.com/conmeo200/Golang-V1/internal/module/tax"
	"github.com/conmeo200/Golang-V1/internal/module/transaction"
	"github.com/conmeo200/Golang-V1/internal/module/user"
	"github.com/conmeo200/Golang-V1/internal/transport/http/api"
	"github.com/conmeo200/Golang-V1/internal/transport/http/middleware"
	"github.com/conmeo200/Golang-V1/internal/transport/http/web"
)

type APIApp struct {
	Container *bootstrap.Container
	Server    *http.Server
}

func NewAPIApp(container *bootstrap.Container) *APIApp {
	  // 1. Repos
	
	userRepo := persistence.NewUserRepository(container.DB)
	authRepo := persistence.NewAuthRepository(container.DB)
	tokenRepo := persistence.NewTokenRepository(container.DB)
	orderRepo := persistence.NewOrderRepository(container.DB)
	paymentRepo := persistence.NewPaymentRepository(container.DB)
	outboxRepo := persistence.NewOutboxEventRepository(container.DB)
	inboxRepo := persistence.NewInboxEventRepository(container.DB)
	newsRepo := news.NewNewsRepository(container.DB)
	roleRepo := persistence.NewRoleRepository(container.DB)
	taxRepo := persistence.NewTaxRepository(container.DB)
	transactionRepo := persistence.NewTransactionRepository(container.DB)

	// 2. Services
	userSvc := user.NewUserService(userRepo)
	authSvc := auth.NewAuthService(authRepo, userRepo, tokenRepo)
	orderSvc := order.NewOrderService(orderRepo, nil)
	paymentSvc := payment.NewPaymentService(paymentRepo, outboxRepo, inboxRepo)
	stripeSvc := stripe.NewStripeService()
	newsSvc := news.NewNewsService(newsRepo)
	roleSvc := role.NewRoleService(roleRepo)
	
	// Seed roles and permissions
	if err := roleSvc.SeedDefaultPermissions(); err != nil {
		log.Printf("⚠️ Warning: Failed to seed default permissions: %v", err)
	}

	taxSvc := tax.NewTaxService(taxRepo)
	transactionSvc := transaction.NewTransactionService(transactionRepo)

	// 3. Handlers
	userHandler := api.NewUserHandler(userSvc)
	authHandler := api.NewAuthHandler(authSvc)
	orderHandler := api.NewOrderHandler(orderSvc)
	paymentHandler := api.NewPaymentHandler(paymentSvc, stripeSvc)
	webhookHandler := api.NewWebhookHandler(paymentSvc, stripeSvc)

	// Web Handlers
	clientHandler := web.NewClientHandler(authSvc)
	dashboardHandler := web.NewDashboardHandler(
		authSvc,
		roleSvc,
		taxSvc,
		orderSvc,
		transactionSvc,
		paymentSvc,
		userSvc,
		newsSvc,
	)

	// 4. Router
	mux := http.NewServeMux()
	
	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(authSvc, roleSvc)

	api.RegisterAPIRoutes(mux, userHandler, authHandler, orderHandler, paymentHandler, webhookHandler)
	web.RegisterWebRoutes(mux, clientHandler, dashboardHandler, paymentHandler, authMiddleware)
	log.Println("✅ API and Web routes registered successfully")

	server := &http.Server{
		Addr:    ":" + container.Config.Port,
		Handler: mux,
	}

	return &APIApp{
		Container: container,
		Server:    server,
	}
}

func (a *APIApp) Run() {
	log.Printf("🚀 API Server starting on %s", a.Server.Addr)
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ API Server failed: %v", err)
	}
}

func (a *APIApp) Stop() {
	log.Println("🛑 Stopping API Server...")
	// Add graceful shutdown logic here
}
