package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/conmeo200/Golang-V1/internal/module/payment/port"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentProvider defines the contract for all payment gateways
type PaymentProvider interface {
	AuthorizePayment(amount float64, currency string, orderID string) (map[string]interface{}, error)
}

// PaymentFactory delegates payment processing to the correct provider

type EventMessage struct {
	EventID    uuid.UUID              `json:"event_id"`
	EventType  string                 `json:"event_type"`
	OccurredAt int64                  `json:"occurred_at"`
	Payload    map[string]interface{} `json:"payload"`
}

type PaymentService struct {
	repo       port.PaymentRepository
	outboxRepo port.OutboxEventRepository
	inboxRepo  port.InboxEventRepository
}

func NewPaymentService(
	repo port.PaymentRepository,
	outboxRepo port.OutboxEventRepository,
	inboxRepo port.InboxEventRepository,
) *PaymentService {
	return &PaymentService{
		repo      	   : repo,
		outboxRepo	   : outboxRepo,
		inboxRepo 	   : inboxRepo,
	}
}

func (s *PaymentService) ListAllTransactions(ctx context.Context) ([]model.Payment, error) {
	return s.repo.ListAll(ctx)
}

func (s *PaymentService) GetPaymentByUUID(ctx context.Context, paymentID uuid.UUID) (*model.Payment, error) {
	return s.repo.GetByUUID(ctx, paymentID)
}

func (s *PaymentService) GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}

func (s *PaymentService) GetPaymentByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error) {
	return s.repo.GetByProviderPaymentID(ctx, providerID)
}

func (s *PaymentService) CreatePendingPayment(ctx context.Context, payment *model.Payment) error {
	payment.Status 	  = "PENDING"
	payment.CreatedAt = time.Now().Unix()

	return s.repo.Create(ctx, payment)
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment *model.Payment, eventIDs ...uuid.UUID) error {
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

func (s *PaymentService) HandleWebhookEvent(ctx context.Context, providerID string, eventType string, eventID uuid.UUID, payload map[string]interface{}) error {
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

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, txUUID uuid.UUID, status string) error {
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

func (s *PaymentService) DeletePayment(ctx context.Context, uuid uuid.UUID) error {
	return s.repo.Delete(ctx, uuid)
}


