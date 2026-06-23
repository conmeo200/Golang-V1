package port

import (
	"context"
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	WithTx(tx *gorm.DB) OrderRepository
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	GetByUUID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	ListAll(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	DB() *gorm.DB
}
