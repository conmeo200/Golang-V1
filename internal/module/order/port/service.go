package port

import (
	"context"
	"github.com/conmeo200/Golang-V1/internal/domain/dto"
	"github.com/conmeo200/Golang-V1/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	WithTx(tx *gorm.DB) OrderService
	CreateOrder(ctx context.Context, userID uuid.UUID, amount float64, idempotencyKey string) (*model.Order, error)
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	ListOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	ListAllOrders(ctx context.Context) ([]model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderUUID uuid.UUID, status string, paymentStatus string) error
	DeleteOrder(ctx context.Context, orderUUID uuid.UUID) error
	ProcessOrder(event dto.OrderMessage) error
	DB() *gorm.DB
}
