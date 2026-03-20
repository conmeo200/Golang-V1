package service

import (
	"context"
	"errors"
	"time"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo repository.OrderRepo
}

func NewOrderService(repo repository.OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, amount float64, idempotencyKey string) (*model.Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	order := &model.Order{
		UUID:           uuid.New(),
		UserID:         userID,
		Amount:         amount,
		Status:         "pending",
		PaymentStatus:  "unpaid",
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}

	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	return s.repo.GetByUUID(ctx, orderUUID)
}

func (s *OrderService) ListOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderUUID uuid.UUID, status string, paymentStatus string) error {
	order, err := s.repo.GetByUUID(ctx, orderUUID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	order.Status = status
	order.PaymentStatus = paymentStatus
	order.UpdatedAt = time.Now().Unix()

	return s.repo.Update(ctx, order)
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderUUID uuid.UUID) error {
	return s.repo.Delete(ctx, orderUUID)
}
