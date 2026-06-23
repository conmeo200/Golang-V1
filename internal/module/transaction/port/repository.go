package port

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	ListAll(ctx context.Context) ([]model.Transaction, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Transaction, error)
	Create(ctx context.Context, tx *model.Transaction) error
	UpdateStatus(ctx context.Context, txUUID uuid.UUID, status string) error
}
