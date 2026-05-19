package consumers

import (
	"context"
	"encoding/json"
	"errors"

	//"fmt"
	"log"
	"time"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/prometheus"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/core/constant"
	"github.com/conmeo200/Golang-V1/internal/module/order"
	"github.com/google/uuid"
	prometheus_client "github.com/prometheus/client_golang/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type PaymentConsumer struct {
	orderService   order.OrderServiceInterface
	rabbitConsumer *rabbitmq.Consumer
	inboxRepo      persistence.InboxEventRepo
}

func NewPaymentConsumer(orderService order.OrderServiceInterface, rabbitConsumer *rabbitmq.Consumer, inboxRepo persistence.InboxEventRepo) *PaymentConsumer {
	return &PaymentConsumer{
		orderService:   orderService,
		rabbitConsumer: rabbitConsumer,
		inboxRepo:      inboxRepo,
	}
}

func (c *PaymentConsumer) Name() string {
	return "payment_consumer"
}

func (c *PaymentConsumer) Start(ctx context.Context) error {
	config := rabbitmq.ConsumerConfig{
		QueueName:    constant.QueuePaymentCompleted,
		ExchangeName: constant.ExchangePayment,
		RoutingKey:   constant.RoutingPaymentCompleted,
		ConsumerName: "payment_worker",
		RetryTTL:     5000,
		MaxRetries:   3,
		PrefetchCount: 5,
		Handler:      c.handleMessage,
	}

	c.rabbitConsumer.Start(ctx, config)
	return nil
}

func (c *PaymentConsumer) Stop() error {
	return nil
}

func (c *PaymentConsumer) handleMessage(msg amqp.Delivery) error {
	timer := prometheus_client.NewTimer(
		prometheus.MessageProcessingDuration.WithLabelValues(constant.QueuePaymentCompleted),
	)
	defer timer.ObserveDuration()

	// Extract eventID (header → fallback)
	eventID := extractEventID(msg)

	log.Printf("[PaymentConsumer] Processing event_id=%s queue=%s", eventID, c.Name())

	var event struct {
		Payload struct {
			OrderID string  `json:"order_id"`
			Status  string  `json:"status"`
			Amount  float64 `json:"amount"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("[PaymentConsumer] Invalid JSON for event %s: %v", eventID, err)
		return err // Retry + DLQ
	}

	// Validation
	if event.Payload.OrderID == "" {
		log.Printf("[PaymentConsumer] Missing order_id for event %s", eventID)
		return errors.New("missing order_id")
	}

	orderUUID, err := uuid.Parse(event.Payload.OrderID)
	if err != nil {
		log.Printf("[PaymentConsumer] Invalid order UUID %s: %v", event.Payload.OrderID, err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Idempotency check before transaction (Optimization)
	exists, err := c.inboxRepo.HasBeenProcessed(ctx, eventID)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("[PaymentConsumer] Event %s already processed, skipping", eventID)
		return nil
	}

	// Start Transactional Processing
	return c.orderService.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Double check idempotency within transaction (Critical for high concurrency)
		// Note: We use a raw query or repo method that supports locking if needed, 
		// but simple insert constraint is usually enough.
		
		// 2. Business logic (Only if status is SUCCESS)
		if event.Payload.Status == "SUCCESS" {
			if err := c.orderService.WithTx(tx).UpdateOrderStatus(ctx, orderUUID, "completed", "paid"); err != nil {
				return err
			}
			log.Printf("[PaymentConsumer] Order %s marked as PAID", orderUUID)
		} else {
			log.Printf("[PaymentConsumer] Payment for order %s status: %s (Skipping business logic)", orderUUID, event.Payload.Status)
		}

		// 3. Record in Inbox (Same transaction)
		inboxStatus := "PROCESSED"
		if event.Payload.Status != "SUCCESS" {
			inboxStatus = "IGNORED"
		}

		if err := c.inboxRepo.WithTx(tx).Create(ctx, &model.InboxEvent{
			EventID:     eventID,
			EventType:   "PaymentCompleted",
			Payload:     msg.Body,
			Status:      inboxStatus,
			ProcessedAt: time.Now().Unix(),
		}); err != nil {
			return err
		}

		return nil
	})
}

