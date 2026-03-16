package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {

	// User routes
	mux.Handle("/api/v1/users", middleware.RequireAPIKey(middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			h.ListUser(w, r)

		case http.MethodPost:
			h.CreateUser(w, r)

		case http.MethodPut:
			h.UpdateUser(w, r)

		case http.MethodDelete:
			h.DeleteUser(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

	}))))

	mux.Handle("/api/v1/find-first-email", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodPost, h.FindFirstByEmail))))
}