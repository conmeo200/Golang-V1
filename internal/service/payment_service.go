package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conmeo200/Golang-V1/internal/config"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentProvider defines the contract for all payment gateways
type PaymentProvider interface {
	AuthorizePayment(amount float64, currency string, orderID string) (map[string]interface{}, error)
}

// PaymentFactory delegates payment processing to the correct provider
type PaymentFactory struct{}

func NewPaymentFactory() *PaymentFactory {
	return &PaymentFactory{}
}

func (f *PaymentFactory) GetProvider(method string) (PaymentProvider, error) {
	switch method {
	case config.PaymentMethodStripe:
		return NewStripeService(), nil
	// TODO: Add PayPal when ready
	// case config.PaymentMethodPayPal:
	// 	return NewPayPalService(), nil
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", method)
	}
}

func (f *PaymentFactory) GetProviderName(method string) string {
	switch method {
	case config.PaymentMethodStripe:
		return "Stripe"
	// TODO: Add PayPal when ready
	// case config.PaymentMethodPayPal:
	// 	return "PayPal"
	default:
		return "Unknown"
	}
}

type PaymentServiceInterface interface {
	ListAllTransactions(ctx context.Context) ([]model.Payment, error)
	GetPaymentByUUID(ctx context.Context, paymentID uuid.UUID) (*model.Payment, error)
	CreatePayment(ctx context.Context, tx *model.Payment) error
	UpdatePaymentStatus(ctx context.Context, txUUID uuid.UUID, status string) error
}

type EventMessage struct {
	EventID    uuid.UUID              `json:"event_id"`
	EventType  string                 `json:"event_type"`
	OccurredAt int64                  `json:"occurred_at"`
	Payload    map[string]interface{} `json:"payload"`
}

type paymentService struct {
	repo       *repository.PaymentRepository
	outboxRepo *repository.OutboxEventRepository
	inboxRepo  *repository.InboxEventRepository
}

func NewPaymentService(
	repo *repository.PaymentRepository,
	outboxRepo *repository.OutboxEventRepository,
	inboxRepo *repository.InboxEventRepository,
) *paymentService {
	return &paymentService{
		repo:       repo,
		outboxRepo: outboxRepo,
		inboxRepo:  inboxRepo,
	}
}

func (s *paymentService) ListAllTransactions(ctx context.Context) ([]model.Payment, error) {
	return s.repo.ListAll(ctx)
}

func (s *paymentService) GetPaymentByUUID(ctx context.Context, paymentID uuid.UUID) (*model.Payment, error) {
	return s.repo.GetByUUID(ctx, paymentID)
}

func (s *paymentService) CreatePayment(ctx context.Context, payment *model.Payment) error {
	eventID := uuid.New()
	now := time.Now().Unix()

	// Start transaction
	return s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update Business Data
		payment.Status = "SUCCESS"
		if err := s.repo.WithTx(tx).Create(ctx, payment); err != nil {
			return err
		}

		// 2. Insert Inbox Event (Idempotency)
		inboxEvent := &model.InboxEvent{
			EventID:     eventID,
			EventType:   "PaymentCompleted",
			Payload:     nil, // Can store request payload if needed
			Status:      "PROCESSED",
			CreatedAt:   now,
			ProcessedAt: now,
		}
		if err := s.inboxRepo.WithTx(tx).Create(ctx, inboxEvent); err != nil {
			return err
		}

		// 3. Insert Outbox Event (Standard Message Format)
		msg := EventMessage{
			EventID:    eventID,
			EventType:  "PaymentCompleted",
			OccurredAt: now,
			Payload: map[string]interface{}{
				"payment_uuid": payment.UUID,
				"order_id":     payment.OrderID,
				"amount":       payment.Amount,
				"currency":     payment.Currency,
				"status":       payment.Status,
			},
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		outboxEvent := &model.OutboxEvents{
			EventID:     eventID,
			EventType:   "PaymentCompleted",
			Payload:     payload,
			Status:      "PENDING",
			RetryCount:  0,
			CreatedAt:   now,
			NextRetryAt: now, // Initial attempt immediate
		}

		if err := s.outboxRepo.WithTx(tx).Create(ctx, outboxEvent); err != nil {
			return err
		}

		return nil
	})
}

func (s *paymentService) UpdatePaymentStatus(ctx context.Context, txUUID uuid.UUID, status string) error {
	payment, err := s.repo.GetByUUID(ctx, txUUID)
	if err != nil {
		return err
	}
	if payment == nil {
		return fmt.Errorf("payment not found")
	}
	payment.Status = status
	return s.repo.Update(ctx, payment)
}

func (s *paymentService) DeletePayment(ctx context.Context, uuid uuid.UUID) error {
	return s.repo.Delete(ctx, uuid)
}


