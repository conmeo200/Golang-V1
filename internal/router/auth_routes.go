package router

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/handler"
	//"github.com/conmeo200/Golang-V1/internal/help"
)

func RegisterAuthRoutes(mux *http.ServeMux, h *handler.AuthHandler) {

	// mux.HandleFunc("/users", help.Method(http.MethodGet, h.ListUser))
	// mux.HandleFunc("/users", help.Method(http.MethodPost, h.CreateUser))
	// mux.HandleFunc("/users", help.Method(http.MethodPut, h.UpdateUser))
	// mux.HandleFunc("/users", help.Method(http.MethodDelete, h.DeleteUser))
	// mux.HandleFunc("/find-first-email", help.Method(http.MethodPost, h.FindFirstByEmail))
}