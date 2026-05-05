package app

import (
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/stripe"
	"github.com/conmeo200/Golang-V1/internal/module/auth"
	"github.com/conmeo200/Golang-V1/internal/module/order"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/conmeo200/Golang-V1/internal/module/user"
	"github.com/conmeo200/Golang-V1/internal/transport/http/api"
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

	// 2. Services
	userSvc := user.NewUserService(userRepo)
	authSvc := auth.NewAuthService(authRepo, userRepo, tokenRepo)
	orderSvc := order.NewOrderService(orderRepo, nil) // Producer can be added if needed
	paymentSvc := payment.NewPaymentService(paymentRepo, outboxRepo, inboxRepo)
	stripeSvc := stripe.NewStripeService()

	// 3. Handlers
	userHandler := api.NewUserHandler(userSvc)
	authHandler := api.NewAuthHandler(authSvc)
	orderHandler := api.NewOrderHandler(orderSvc)
	paymentHandler := api.NewPaymentHandler(paymentSvc, stripeSvc)
	webhookHandler := api.NewWebhookHandler(paymentSvc, stripeSvc)

	// 4. Router
	mux := http.NewServeMux()
	api.RegisterAPIRoutes(mux, userHandler, authHandler, orderHandler, paymentHandler, webhookHandler)

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
