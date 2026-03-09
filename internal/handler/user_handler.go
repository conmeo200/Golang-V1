package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/conmeo200/Golang-V1/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Email   string  `json:"email"`
		Balance float64 `json:"balance"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(req)
	err := h.service.CreateUser(req.Email, req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	log.Println("Param IDS", idStr)
	id, _ := strconv.Atoi(idStr)

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {

	log.Println("List User")

	users, err := h.service.ListUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(users)
}

type UpdateUserRequest struct {
	Balance float64 `json:"balance"`
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Update user id=%d balance=%f", id, req.Balance)

	err = h.service.Update(uint(id), req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated"))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User deleted")
}
