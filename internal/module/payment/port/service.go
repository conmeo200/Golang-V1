package port

import (
	"context"
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
)

type PaymentService interface {
	ListAllTransactions(ctx context.Context) ([]model.Payment, error)
	GetPaymentByUUID(ctx context.Context, paymentID uuid.UUID) (*model.Payment, error)
	GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error)
	GetPaymentByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error)
	
	CreatePendingPayment(ctx context.Context, payment *model.Payment) error
	CreatePayment(ctx context.Context, payment *model.Payment, eventIDs ...uuid.UUID) error
	HandleWebhookEvent(ctx context.Context, providerID string, eventType string, eventID uuid.UUID, payload map[string]interface{}) error
	UpdatePaymentStatus(ctx context.Context, txUUID uuid.UUID, status string) error
}
