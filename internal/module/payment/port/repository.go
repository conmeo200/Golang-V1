package port

import (
	"context"
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	WithTx(tx *gorm.DB) PaymentRepository
	Create(ctx context.Context, payment *model.Payment) error
	Update(ctx context.Context, payment *model.Payment) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error)
	GetByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error)
	ListAll(ctx context.Context) ([]model.Payment, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	DB() *gorm.DB
}

type OutboxEventRepository interface {
	WithTx(tx *gorm.DB) OutboxEventRepository
	Create(ctx context.Context, event *model.OutboxEvents) error
	Update(ctx context.Context, event *model.OutboxEvents) error
	FetchPending(ctx context.Context, limit int) ([]model.OutboxEvents, error)
	MarkAsPublished(ctx context.Context, id interface{}, sentAt int64) error
}

type InboxEventRepository interface {
	WithTx(tx *gorm.DB) InboxEventRepository
	Create(ctx context.Context, event *model.InboxEvent) error
	Update(ctx context.Context, event *model.InboxEvent) error
	HasBeenProcessed(ctx context.Context, eventID uuid.UUID) (bool, error)
}
