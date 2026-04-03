package repository

import (
	"context"
	"errors"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepo interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	ListAll(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (r *OrderRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).First(&order, "uuid = ?", uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) ListAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).Preload("User").Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, "uuid = ?", uuid).Error
}
