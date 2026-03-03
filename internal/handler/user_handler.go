package handler

import (
	"net/http"
	"github.com/conmeo200/Golang-V1/internal/service"
	"log"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health1", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health endpoint called")
		h.Health(w, r)
	})
}


func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"status": "ok123123123"}`))
}