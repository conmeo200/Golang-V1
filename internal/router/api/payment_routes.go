package api

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler/api"
)

func RegisterPaymentRoutes(mux *http.ServeMux, handler *api.PaymentHandler) {
	mux.Handle("/api/v1/payments/process", http.HandlerFunc(handler.ProcessPayment))
}
