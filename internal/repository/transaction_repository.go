package repository

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) ListAll(ctx context.Context) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.WithContext(ctx).Order("created_at desc").Preload("Order").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, txUUID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Transaction{}).Where("uuid = ?", txUUID).Update("status", status).Error
}
