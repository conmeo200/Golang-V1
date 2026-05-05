package payment

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service PaymentServiceInterface
}

func (s *PaymentServiceTestSuite) SetupSuite() {
	// 1. Initialize SQLite in-memory
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to open database: %v", err)
	}
	s.db = db

	// 2. Migrate schemas
	err = db.AutoMigrate(
		&model.Payment{},
		&model.InboxEvent{},
		&model.OutboxEvents{},
	)
	if err != nil {
		s.T().Fatalf("failed to migrate database: %v", err)
	}

	// 3. Initialize repositories and service
	paymentRepo := persistence.NewPaymentRepository(db)
	outboxRepo := persistence.NewOutboxEventRepository(db)
	inboxRepo := persistence.NewInboxEventRepository(db)

	s.service = NewPaymentService(paymentRepo, outboxRepo, inboxRepo)
}

func (s *PaymentServiceTestSuite) SetupTest() {
	// Clean tables before each test
	s.db.Exec("DELETE FROM payments")
	s.db.Exec("DELETE FROM inbox_events")
	s.db.Exec("DELETE FROM outbox_events")
}

func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}

func (s *PaymentServiceTestSuite) TestCreatePayment_Success() {
	ctx := context.Background()
	orderID := uuid.New()
	payment := &model.Payment{
		UUID:          uuid.New(),
		OrderID:       orderID,
		Amount:        100.50,
		Currency:      "USD",
		PaymentMethod: "Stripe",
		Status:        "PENDING", // Should be SUCCESS after CreatePayment
		// Descriptio:   "Test Payment",
		CreatedAt:     time.Now().Unix(),
	}

	eventID := uuid.New()

	// Execute
	err := s.service.CreatePayment(ctx, payment, eventID)

	// Asserts
	assert.NoError(s.T(), err)

	// 1. Check Payment table
	var dbPayment model.Payment
	err = s.db.First(&dbPayment, "uuid = ?", payment.UUID).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "SUCCESS", dbPayment.Status)
	assert.Equal(s.T(), 100.50, dbPayment.Amount)

	// 2. Check InboxEvent table
	var dbInbox model.InboxEvent
	err = s.db.First(&dbInbox, "event_id = ?", eventID).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "PaymentCompleted", dbInbox.EventType)
	assert.Equal(s.T(), "PROCESSED", dbInbox.Status)

	// 3. Check OutboxEvent table
	var dbOutbox model.OutboxEvents
	err = s.db.First(&dbOutbox, "event_id = ?", eventID).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "PaymentCompleted", dbOutbox.EventType)
	assert.Equal(s.T(), "PENDING", dbOutbox.Status)

	// 4. Verify Outbox Payload
	var payload map[string]interface{}
	err = json.Unmarshal(dbOutbox.Payload, &payload)
	assert.NoError(s.T(), err)
	
	innerPayload := payload["payload"].(map[string]interface{})
	assert.Equal(s.T(), payment.UUID.String(), innerPayload["payment_uuid"])
	assert.Equal(s.T(), orderID.String(), innerPayload["order_id"])
}

func (s *PaymentServiceTestSuite) TestCreatePayment_Idempotency() {
	ctx := context.Background()
	eventID := uuid.New()
	
	payment1 := &model.Payment{
		UUID:    uuid.New(),
		OrderID: uuid.New(),
		Status:  "PENDING",
	}

	// First call
	err := s.service.CreatePayment(ctx, payment1, eventID)
	assert.NoError(s.T(), err)

	// Second call with SAME eventID but DIFFERENT payment
	payment2 := &model.Payment{
		UUID:    uuid.New(),
		OrderID: uuid.New(),
		Status:  "PENDING",
	}
	err = s.service.CreatePayment(ctx, payment2, eventID)

	// Should fail because eventID in inbox_events is UNIQUE
	assert.Error(s.T(), err)
	
	// Verify second payment was NOT created (Transactional Integrity)
	var count int64
	s.db.Model(&model.Payment{}).Where("uuid = ?", payment2.UUID).Count(&count)
	assert.Equal(s.T(), int64(0), count)
}
