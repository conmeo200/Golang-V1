package repository

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentEventRepository struct {
	db *gorm.DB
}

func NewPaymentEventRepository(db *gorm.DB) *PaymentEventRepository {
	return &PaymentEventRepository{db: db}
}

func (r *PaymentEventRepository) Create(ctx context.Context, event *model.PaymentEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *PaymentEventRepository) ListByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]model.PaymentEvent, error) {
	var events []model.PaymentEvent
	err := r.db.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Order("created_at asc").
		Find(&events).Error
	return events, err
}
