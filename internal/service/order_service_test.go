package service

import (
	"context"
	"testing"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		amount := 100.0
		idempotencyKey := "key123"

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(&model.Order{
			UserID:         userID,
			Amount:         amount,
			IdempotencyKey: idempotencyKey,
		}, nil)

		order, err := service.CreateOrder(ctx, userID, amount, idempotencyKey)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, userID, order.UserID)
		assert.Equal(t, amount, order.Amount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Amount", func(t *testing.T) {
		userID := uuid.New()
		order, err := service.CreateOrder(ctx, userID, -10.0, "")

		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, "amount must be greater than zero", err.Error())
	})
}

func TestGetOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		orderUUID := uuid.New()
		mockRepo.On("GetByUUID", ctx, orderUUID).Return(&model.Order{UUID: orderUUID}, nil)

		order, err := service.GetOrder(ctx, orderUUID)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, orderUUID, order.UUID)
	})
}

func TestUpdateOrderStatus(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		orderUUID := uuid.New()
		order := &model.Order{UUID: orderUUID, Status: "pending"}

		mockRepo.On("GetByUUID", ctx, orderUUID).Return(order, nil)
		mockRepo.On("Update", ctx, order).Return(nil)

		err := service.UpdateOrderStatus(ctx, orderUUID, "completed", "paid")

		assert.NoError(t, err)
		assert.Equal(t, "completed", order.Status)
		assert.Equal(t, "paid", order.PaymentStatus)
	})

	t.Run("Order Not Found", func(t *testing.T) {
		orderUUID := uuid.New()
		mockRepo.On("GetByUUID", ctx, orderUUID).Return(nil, nil)

		err := service.UpdateOrderStatus(ctx, orderUUID, "completed", "paid")

		assert.Error(t, err)
		assert.Equal(t, "order not found", err.Error())
	})
}
