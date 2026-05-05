package api

import (
	"net/http"
)

func RegisterAPIRoutes(mux *http.ServeMux, userHandler *UserHandler, authHandler *AuthHandler, orderHandler *OrderHandler, paymentHandler *PaymentHandler, webhookHandler *WebhookHandler) {
	RegisterUserRoutes(mux, userHandler)
	RegisterAuthRoutes(mux, authHandler)
	RegisterOrderRoutes(mux, orderHandler)
	RegisterPaymentRoutes(mux, paymentHandler, webhookHandler)
}
