package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo     repository.OrderRepo
	producer *rabbitmq.Producer
}

func NewOrderService(repo repository.OrderRepo, producer *rabbitmq.Producer) *OrderService {
	return &OrderService{
		repo:     repo,
		producer: producer,
	}
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
		IdempotencyKey: uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}

	_, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// Publish Order Created Event
	if err := s.PublishOrderCreated(ctx, order); err != nil {
		log.Printf("failed to publish order created event: %v", err)
	}

	return order, nil
}

func (s *OrderService) PublishOrderCreated(ctx context.Context, order *model.Order) error {
	if s.producer == nil {
		return nil
	}

	msg := dto.RabbitMQResponse{
		Event: "order.created",
		Data: dto.OrderMessage{
			OrderUUID: order.UUID.String(),
			UserID:    order.UserID.String(),
			Amount:    order.Amount,
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.producer.PublishOrderCreated(ctx, body)
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

func (s *OrderService) ProcessOrder(event dto.OrderMessage) error {
	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		return errors.New("invalid order uuid: " + event.OrderUUID)
	}

	order, err := s.GetOrder(context.Background(), orderUUID)
	if err != nil {
		return err
	}

	// ❗ idempotency check
	if order.Status == "completed" {
		return nil
	}

	// fake payment processing
	time.Sleep(2 * time.Second)

	// update order
	order.Status 		= "completed"
	order.PaymentStatus = "paid"

	return s.UpdateOrderStatus(context.Background(), order.UUID, order.Status, order.PaymentStatus)
}