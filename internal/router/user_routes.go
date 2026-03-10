package router

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/handler"
)

func registerUserRoutes(mux *http.ServeMux, userHandler *handler.UserHandler) {

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			userHandler.ListUser(w, r)

		case http.MethodPost:
			userHandler.CreateUser(w, r)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			userHandler.GetUser(w, r)

		case http.MethodPut:
			userHandler.UpdateUser(w, r)

		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		}
	})
}