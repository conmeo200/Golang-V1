package router

import (
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/handler"
)

func New(userHandler *handler.UserHandler) *http.ServeMux {

	mux := http.NewServeMux()

	// collection routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			log.Println("List Users")
			userHandler.ListUser(w, r)

		case http.MethodPost:
			log.Println("Create User")
			userHandler.CreateUser(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// single resource
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			log.Println("Get User By ID")
			userHandler.GetUser(w, r)

		case http.MethodPut:
			log.Println("Route Update User")
			userHandler.UpdateUser(w, r)

		case http.MethodDelete:
			log.Println("Delete User")
			userHandler.DeleteUser(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}