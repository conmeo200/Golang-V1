package api

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler/api"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterOrderRoutes(mux *http.ServeMux, orderHandler *api.OrderHandler) {
	mux.Handle("/api/v1/orders", middleware.RequireAPIKey(middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderHandler.ListOrders(w, r)
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodPut:
			orderHandler.UpdateOrder(w, r)
		case http.MethodDelete:
			orderHandler.DeleteOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/v1/orders/get", middleware.RequireAPIKey(middleware.JWTMiddleware(http.HandlerFunc(orderHandler.GetOrder))))
}
