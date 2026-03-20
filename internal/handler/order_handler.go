package handler

import (
	"encoding/json"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/google/uuid"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("CreateOrder decode error: %v", err)
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid User ID", "INVALID_USER_ID"))
		return
	}

	order, err := h.service.CreateOrder(r.Context(), userUUID, req.Amount, req.IdempotencyKey)
	if err != nil {
		logger.ErrorLogger.Printf("CreateOrder error: %v", err)
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusCreated, dto.ToOrderResponse(order), "Order created successfully")
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "UUID is required", "UUID_REQUIRED"))
		return
	}

	orderUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid UUID", "INVALID_UUID"))
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderUUID)
	if err != nil {
		logger.ErrorLogger.Printf("GetOrder error: %v", err)
		dto.RespondWithError(w, err)
		return
	}

	if order == nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusNotFound, "Order not found", "ORDER_NOT_FOUND"))
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, dto.ToOrderResponse(order), "Order found")
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "User ID is required", "USER_ID_REQUIRED"))
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid User ID", "INVALID_USER_ID"))
		return
	}

	orders, err := h.service.ListOrdersByUserID(r.Context(), userUUID)
	if err != nil {
		logger.ErrorLogger.Printf("ListOrders error: %v", err)
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, dto.ToOrderResponsesArray(orders), "Order list")
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "UUID is required", "UUID_REQUIRED"))
		return
	}

	orderUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid UUID", "INVALID_UUID"))
		return
	}

	var req dto.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondWithError(w, dto.ErrInvalidRequest)
		return
	}

	err = h.service.UpdateOrderStatus(r.Context(), orderUUID, req.Status, req.PaymentStatus)
	if err != nil {
		logger.ErrorLogger.Printf("UpdateOrder error: %v", err)
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, nil, "Order updated successfully")
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "UUID is required", "UUID_REQUIRED"))
		return
	}

	orderUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		dto.RespondWithError(w, dto.NewAppError(http.StatusBadRequest, "Invalid UUID", "INVALID_UUID"))
		return
	}

	err = h.service.DeleteOrder(r.Context(), orderUUID)
	if err != nil {
		logger.ErrorLogger.Printf("DeleteOrder error: %v", err)
		dto.RespondWithError(w, err)
		return
	}

	dto.RespondWithSuccess(w, http.StatusOK, nil, "Order deleted successfully")
}
