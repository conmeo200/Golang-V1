package api

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler/api"
	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux, userHandler *api.UserHandler) {

	// User routes
	mux.Handle("/api/v1/users", middleware.RequireAPIKey(middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			userHandler.ListUser(w, r)

		case http.MethodPost:
			userHandler.CreateUser(w, r)

		case http.MethodPut:
			userHandler.UpdateUser(w, r)

		case http.MethodDelete:
			userHandler.DeleteUser(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}))))

	mux.Handle("/api/v1/users/", middleware.RequireAPIKey(middleware.JWTMiddleware(help.Method(http.MethodGet, userHandler.GetUser))))
}