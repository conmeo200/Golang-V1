package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  service.UserServiceInterface
	validate *validator.Validate
}

func NewUserHandler(s service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service:  s,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Email        string  `json:"email" validate:"required,email"`
		Balance      float64 `json:"balance" validate:"min=0"`
		PasswordHash string  `json:"password" validate:"required,min=6"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("CreateUser decode error: %v", err)
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, err.Error(), "VALIDATION_FAILED"))
		return
	}

	if req.Email == "" || req.PasswordHash == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Data Invalid", "DATA_INVALID"))
		return
	}

	logger.AppLogger.Println("CreateUser request received:", req.Email)
	result, err := h.service.CreateUser(r.Context(), req.Email, req.Balance, req.PasswordHash)

	if err != nil {
		logger.ErrorLogger.Printf("CreateUser error from DB: %v", err)
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusCreated, dto.ToUserResponse(result), "User created successfully")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	logger.AppLogger.Println("Param IDS", idStr)
	//id, _ := strconv.Atoi(idStr)

	user, err := h.service.GetUser(r.Context(), idStr)
	if err != nil {
		logger.ErrorLogger.Printf("GetUser error: %v", err)
		dto.RespondWithError(w, dto.NewAppError(http.StatusNotFound, err.Error(), "USER_NOT_FOUND"))
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, dto.ToUserResponse(user), "User found")
}

func (h *UserHandler) FindFirstByEmail(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Email   string  `json:"email"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("FindFirstByEmail decode error: %v", err)
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		return
	}

	if req.Email == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Email is required", "EMAIL_REQUIRED"))
		return
	}

	user, err := h.service.FindFirstByEmail(r.Context(), req.Email)
	if err != nil {
		logger.ErrorLogger.Printf("FindFirstByEmail error: %v", err)
		dto.RespondWithError(w, dto.NewAppError(http.StatusNotFound, err.Error(), "USER_NOT_FOUND"))
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, dto.ToUserResponse(user), "User found")
}


func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUser(r.Context())
	if err != nil {
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, dto.ToUserResponsesArray(users), "User list")
}

type UpdateUserRequest struct {
	Balance float64 `json:"balance" validate:"min=0"`
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid user ID", "INVALID_ID"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid user ID", "INVALID_ID"))
		return
	}

	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, err.Error(), "VALIDATION_FAILED"))
		return
	}

	logger.AppLogger.Printf("Update user id=%d balance=%f", id, req.Balance)

	err = h.service.UpdateBalance(r.Context(), uint(id), req.Balance)
	if err != nil {
		dto.RespondWithError(w, err)
		return
	}

	result, err := h.service.GetUser(r.Context(), idStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusNotFound, err.Error(), "USER_NOT_FOUND"))
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, dto.ToUserResponse(result), "User updated successfully")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid user ID", "INVALID_ID"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid user ID", "INVALID_ID"))
		return
	}

	err = h.service.DeleteUser(r.Context(), uint(id))
	if err != nil {
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, nil, "User deleted successfully")
}