package ports

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/google/uuid"
)

type PaymentService interface {
	ListAllTransactions(ctx context.Context) ([]model.Payment, error)
	GetPaymentByUUID(ctx context.Context, paymentID uuid.UUID) (*model.Payment, error)
	HandleWebhookEvent(ctx context.Context, providerID string, eventType string, eventID uuid.UUID, payload map[string]interface{}) error
	UpdatePaymentStatus(ctx context.Context, txUUID uuid.UUID, status string) error
}
