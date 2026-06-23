package port

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
)

type TransactionService interface {
	ListAllTransactions(ctx context.Context) ([]model.Transaction, error)
	GetTransactionsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Transaction, error)
	CreateTransaction(ctx context.Context, tx *model.Transaction) error
	UpdateTransactionStatus(ctx context.Context, txUUID uuid.UUID, status string) error
}
