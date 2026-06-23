package transaction

import (
	"context"

	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/conmeo200/Golang-V1/internal/module/transaction/port"
	"github.com/google/uuid"
)

type TransactionService struct {
	repo port.TransactionRepository
}

func NewTransactionService(repo port.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) ListAllTransactions(ctx context.Context) ([]model.Transaction, error) {
	return s.repo.ListAll(ctx)
}

func (s *TransactionService) GetTransactionsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Transaction, error) {
	return s.repo.FindByOrderID(ctx, orderID)
}

func (s *TransactionService) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	return s.repo.Create(ctx, tx)
}

func (s *TransactionService) UpdateTransactionStatus(ctx context.Context, txUUID uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, txUUID, status)
}
