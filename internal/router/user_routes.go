package router

import (
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/help"
	"github.com/conmeo200/Golang-V1/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {

	mux.Handle("/users", middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

	})))

	mux.Handle("/find-first-email", middleware.JWTMiddleware(help.Method(http.MethodPost, h.FindFirstByEmail)))
}