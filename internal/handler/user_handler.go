package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/logger"
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
		PasswordHash string `json:"password"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("CreateUser decode error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.PasswordHash == "" {
		http.Error(w, "Data Invalid", http.StatusBadRequest)
		return
	}

	logger.AppLogger.Println("CreateUser request received:", req.Email)
	result, err := h.service.CreateUser(r.Context(), req.Email, req.Balance, req.PasswordHash)

	if err != nil {
		logger.ErrorLogger.Printf("CreateUser error from DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.SendSuccess(dto.ToUserResponse(result)))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	logger.AppLogger.Println("Param IDS", idStr)
	//id, _ := strconv.Atoi(idStr)

	user, err := h.service.GetUser(r.Context(), idStr)
	if err != nil {
		logger.ErrorLogger.Printf("GetUser error: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.SendSuccess(dto.ToUserResponse(user)))
}

func (h *UserHandler) FindFirstByEmail(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Email   string  `json:"email"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("FindFirstByEmail decode error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Data Invalid", http.StatusBadRequest)
		return
	}

	user, err := h.service.FindFirstByEmail(r.Context(), req.Email)
	if err != nil {
		logger.ErrorLogger.Printf("FindFirstByEmail error: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.SendSuccess(dto.ToUserResponse(user)))
}


func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.ToUserResponsesArray(users))
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

	logger.AppLogger.Printf("Update user id=%d balance=%f", id, req.Balance)

	err = h.service.UpdateBalance(r.Context(), uint(id), req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := h.service.GetUser(r.Context(), idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.SendSuccess(dto.ToUserResponse(result)))
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

	err = h.service.DeleteUser(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.SendSuccess(nil))
}
