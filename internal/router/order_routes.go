package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterOrderRoutes(mux *http.ServeMux, h *handler.OrderHandler) {
	mux.Handle("/api/v1/orders", middleware.RequireAPIKey(middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListOrders(w, r)
		case http.MethodPost:
			h.CreateOrder(w, r)
		case http.MethodPut:
			h.UpdateOrder(w, r)
		case http.MethodDelete:
			h.DeleteOrder(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/v1/orders/get", middleware.RequireAPIKey(middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.GetOrder(w, r)
	}))))
}
