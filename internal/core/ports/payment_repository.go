package ports

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	WithTx(tx *gorm.DB) PaymentRepository
	Create(ctx context.Context, payment *model.Payment) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error)
	GetByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error)
	ListAll(ctx context.Context) ([]model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	DB() *gorm.DB
}
