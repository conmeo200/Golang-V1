package router

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/handler"
	"github.com/conmeo200/Golang-V1/internal/help"
)

func RegisterUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

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

	})

	mux.HandleFunc("/find-first-email", help.Method(http.MethodPost, h.FindFirstByEmail))
}

// func RegisterUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {

// 	mux.HandleFunc("/users", help.Method(http.MethodGet, h.ListUser))
// 	mux.HandleFunc("/users", help.Method(http.MethodPost, h.CreateUser))
// 	mux.HandleFunc("/users", help.Method(http.MethodPut, h.UpdateUser))
// 	mux.HandleFunc("/users", help.Method(http.MethodDelete, h.DeleteUser))
// 	mux.HandleFunc("/find-first-email", help.Method(http.MethodPost, h.FindFirstByEmail))
// }