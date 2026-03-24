package dto

import (
	"github.com/conmeo200/Golang-V1/internal/model"
)

type OrderResponse struct {
	UUID           string  `json:"uuid"`
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	PaymentStatus  string  `json:"payment_status"`
	IdempotencyKey string  `json:"idempotency_key"`
	ProcessedAt    int64   `json:"processed_at"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
}

type CreateOrderRequest struct {
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	IdempotencyKey string  `json:"idempotency_key" validate:"required"`
}

type UpdateOrderRequest struct {
	Status        string `json:"status" validate:"required,oneof=pending completed cancelled"`
	PaymentStatus string `json:"payment_status" validate:"required,oneof=unpaid paid failed"`
}

func ToOrderResponse(order *model.Order) OrderResponse {
	return OrderResponse{
		UUID:           order.UUID.String(),
		UserID:         order.UserID.String(),
		Amount:         order.Amount,
		Status:         order.Status,
		PaymentStatus:  order.PaymentStatus,
		IdempotencyKey: order.IdempotencyKey,
		ProcessedAt:    order.ProcessedAt,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}
}

func ToOrderResponsesArray(orders []model.Order) []OrderResponse {
	var responses []OrderResponse
	for _, order := range orders {
		responses = append(responses, ToOrderResponse(&order))
	}
	return responses
}
