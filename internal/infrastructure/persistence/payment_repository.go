package persistence

import (
	"context"
	"errors"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PaymentRepo interface {
	WithTx(tx *gorm.DB) PaymentRepo
	Create(ctx context.Context, payment *model.Payment) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error)
	GetByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error)
	ListAll(ctx context.Context) ([]model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	DB() *gorm.DB
}

type PaymentRepository struct {
	db *gorm.DB
}


func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) DB() *gorm.DB {
	return r.db
}

func (r *PaymentRepository) WithTx(tx *gorm.DB) PaymentRepo {
	return &PaymentRepository{db: tx}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment)  error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).First(&payment, "uuid = ?", uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &payment, err
}

func (r *PaymentRepository) GetByProviderPaymentID(ctx context.Context, providerID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).First(&payment, "provider_payment_id = ?", providerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &payment, err
}

func (r *PaymentRepository) ListAll(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at desc").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *PaymentRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Payment{}, "uuid = ?", uuid).Error
}



