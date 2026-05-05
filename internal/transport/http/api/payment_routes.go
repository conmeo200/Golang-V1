package api

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/transport/http/common/middleware"
)

func RegisterPaymentRoutes(mux *http.ServeMux, handler *PaymentHandler, webhookHandler *WebhookHandler) {
	mux.Handle("/api/v1/payments/process", middleware.RequireAPIKey(help.Method(http.MethodPost, handler.ProcessPayment)))
	mux.Handle("/api/v1/payments/refund", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, handler.RefundPayment))))
	mux.Handle("/api/v1/webhooks/stripe", help.Method(http.MethodPost, webhookHandler.HandleStripeWebhook))
}
