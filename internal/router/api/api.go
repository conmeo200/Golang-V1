package api

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/handler/api"
)

func RegisterAPIRoutes(mux *http.ServeMux, userHandler *api.UserHandler, authHandler *api.AuthHandler, orderHandler *api.OrderHandler) {
	RegisterUserRoutes(mux, userHandler)
	RegisterAuthRoutes(mux, authHandler)
	RegisterOrderRoutes(mux, orderHandler)
}
