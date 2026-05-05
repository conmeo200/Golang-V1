package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/stripe"
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
	case bootstrap.PaymentMethodStripe:
		return stripe.NewStripeService(), nil
	// TODO: Add PayPal when ready
	// case bootstrap.PaymentMethodPayPal:
	// 	return NewPayPalService(), nil
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", method)
	}
}

func (f *PaymentFactory) GetProviderName(method string) string {
	switch method {
	case bootstrap.PaymentMethodStripe:
		return "Stripe"
	// TODO: Add PayPal when ready
	// case bootstrap.PaymentMethodPayPal:
	// 	return "PayPal"
	default:
		return "Unknown"
	}
}

type PaymentServiceInterface interface {
	ListAllTransactions(ctx context.Context) ([]model.Payment, error)
	GetPaymentByUUID(ctx context.Context, paymentID uuid.UUID) (*model.Payment, error)
	GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error)
	GetPaymentByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error)
	
	CreatePendingPayment(ctx context.Context, payment *model.Payment) error
	CreatePayment(ctx context.Context, payment *model.Payment, eventIDs ...uuid.UUID) error
	HandleWebhookEvent(ctx context.Context, providerID string, eventType string, eventID uuid.UUID, payload map[string]interface{}) error
	UpdatePaymentStatus(ctx context.Context, txUUID uuid.UUID, status string) error
}

type EventMessage struct {
	EventID    uuid.UUID              `json:"event_id"`
	EventType  string                 `json:"event_type"`
	OccurredAt int64                  `json:"occurred_at"`
	Payload    map[string]interface{} `json:"payload"`
}

type paymentService struct {
	repo       persistence.PaymentRepo
	outboxRepo persistence.OutboxEventRepo
	inboxRepo  persistence.InboxEventRepo
}

func NewPaymentService(
	repo persistence.PaymentRepo,
	outboxRepo persistence.OutboxEventRepo,
	inboxRepo persistence.InboxEventRepo,
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

func (s *paymentService) GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}

func (s *paymentService) GetPaymentByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error) {
	return s.repo.GetByProviderPaymentID(ctx, providerID)
}

func (s *paymentService) CreatePendingPayment(ctx context.Context, payment *model.Payment) error {
	payment.Status = "PENDING"
	payment.CreatedAt = time.Now().Unix()

	return s.repo.Create(ctx, payment)
}

func (s *paymentService) CreatePayment(ctx context.Context, payment *model.Payment, eventIDs ...uuid.UUID) error {
	var eventID uuid.UUID
	if len(eventIDs) > 0 {
		eventID = eventIDs[0]
	} else {
		eventID = uuid.New()
	}
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

func (s *paymentService) HandleWebhookEvent(ctx context.Context, providerID string, eventType string, eventID uuid.UUID, payload map[string]interface{}) error {
	now := time.Now().Unix()

	return s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Fetch Payment
		var payment model.Payment
		if err := tx.First(&payment, "provider_payment_id = ?", providerID).Error; err != nil {
			return fmt.Errorf("payment not found for provider ID %s: %w", providerID, err)
		}

		// Calculate new status and outbox event type
		var newStatus string
		var outboxEventType string

		switch eventType {
		case "payment_intent.succeeded":
			newStatus = "SUCCESS"
			outboxEventType = "PaymentCompleted"
		case "payment_intent.payment_failed":
			newStatus = "FAILED"
			outboxEventType = "PaymentFailed"
		case "charge.refunded":
			newStatus = "REFUNDED"
			outboxEventType = "PaymentRefunded"
		default:
			return fmt.Errorf("unhandled webhook event type: %s", eventType)
		}

		// 2. Update Payment Status
		if !model.CanTransitionPayment(payment.Status, newStatus) {
			return fmt.Errorf("invalid payment transition from %s to %s", payment.Status, newStatus)
		}
		payment.Status = newStatus
		if err := s.repo.WithTx(tx).Update(ctx, &payment); err != nil {
			return err
		}

		// 3. Insert Inbox Event (Idempotency)
		inboxEvent := &model.InboxEvent{
			EventID:     eventID,
			EventType:   eventType,
			Payload:     nil, 
			Status:      "PROCESSED",
			CreatedAt:   now,
			ProcessedAt: now,
		}
		
		if err := s.inboxRepo.WithTx(tx).Create(ctx, inboxEvent); err != nil {
			return err
		}

		// 4. Insert Outbox Event
		if payload == nil {
			payload = map[string]interface{}{}
		}
		payload["payment_uuid"] = payment.UUID
		payload["order_id"] = payment.OrderID
		payload["amount"] = payment.Amount
		payload["currency"] = payment.Currency
		payload["status"] = payment.Status

		msg := EventMessage{
			EventID:    eventID,
			EventType:  outboxEventType,
			OccurredAt: now,
			Payload:    payload,
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		outboxEvent := &model.OutboxEvents{
			EventID:     eventID,
			EventType:   outboxEventType,
			Payload:     msgBytes,
			Status:      "PENDING",
			RetryCount:  0,
			CreatedAt:   now,
			NextRetryAt: now, 
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
	if !model.CanTransitionPayment(payment.Status, status) {
		return model.ErrInvalidPaymentStateTransition
	}
	payment.Status = status
	return s.repo.Update(ctx, payment)
}

func (s *paymentService) DeletePayment(ctx context.Context, uuid uuid.UUID) error {
	return s.repo.Delete(ctx, uuid)
}


