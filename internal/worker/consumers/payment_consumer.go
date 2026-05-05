package consumers

import (
	"context"
	"encoding/json"
	//"fmt"
	"log"
	"time"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/prometheus"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/module/order"
	"github.com/google/uuid"
	prometheus_client "github.com/prometheus/client_golang/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
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
		QueueName:    "payment_completed_queue",
		ExchangeName: "payment.exchange",
		RoutingKey:   "PaymentCompleted",
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
	timer := prometheus_client.NewTimer(prometheus.MessageProcessingDuration.WithLabelValues("payment_completed_queue"))
	defer timer.ObserveDuration()

	log.Printf("📥 [PaymentConsumer] Received message: %s", string(msg.Body))

	// Generate deterministic eventID from payload to ensure idempotency
	eventID := uuid.NewSHA1(uuid.NameSpaceOID, msg.Body)

	var event struct {
		Payload struct {
			OrderID string `json:"order_id"`
			Status  string `json:"status"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ [PaymentConsumer] Failed to unmarshal message: %v", err)
		return nil // Don't retry if JSON is invalid
	}

	orderUUID, err := uuid.Parse(event.Payload.OrderID)
	if err != nil {
		log.Printf("❌ [PaymentConsumer] Invalid Order UUID: %v", err)
		return nil
	}

	// 1. Inbox Idempotency Check
	ctx := context.Background()
	exists, err := c.inboxRepo.HasBeenProcessed(ctx, eventID)
	if err != nil {
		prometheus.MessagesConsumedTotal.WithLabelValues("payment_completed_queue", "failed").Inc()
		log.Printf("❌ [PaymentConsumer] Failed to check inbox for event %s: %v", eventID, err)
		return err // Retry
	}
	if exists {
		prometheus.MessagesConsumedTotal.WithLabelValues("payment_completed_queue", "success").Inc()
		log.Printf("⚠️ [PaymentConsumer] Event %s already processed, skipping", eventID)
		return nil // Ack
	}

	if event.Payload.Status != "SUCCESS" {
		log.Printf("⚠️ [PaymentConsumer] Payment status is not SUCCESS: %s", event.Payload.Status)
		// Even if not success, we might want to record it in inbox to avoid reprocessing
		c.inboxRepo.Create(ctx, &model.InboxEvent{
			EventID:     eventID,
			EventType:   "PaymentCompleted",
			Payload:     msg.Body,
			Status:      "PROCESSED",
			ProcessedAt: time.Now().Unix(),
		})
		prometheus.MessagesConsumedTotal.WithLabelValues("payment_completed_queue", "success").Inc()
		return nil
	}

	// 2. Update Order status
	err = c.orderService.UpdateOrderStatus(ctx, orderUUID, "completed", "paid")
	if err != nil {
		log.Printf("❌ [PaymentConsumer] Failed to update order status: %v", err)
		return err // Trigger retry
	}

	// 3. Save to Inbox
	err = c.inboxRepo.Create(ctx, &model.InboxEvent{
		EventID:     eventID,
		EventType:   "PaymentCompleted",
		Payload:     msg.Body,
		Status:      "PROCESSED",
		ProcessedAt: time.Now().Unix(),
	})
	if err != nil {
		prometheus.MessagesConsumedTotal.WithLabelValues("payment_completed_queue", "failed").Inc()
		log.Printf("❌ [PaymentConsumer] Failed to save to inbox: %v", err)
		return err // Trigger retry so we don't lose the idempotency record
	}

	prometheus.MessagesConsumedTotal.WithLabelValues("payment_completed_queue", "success").Inc()
	log.Printf("✅ [PaymentConsumer] Order %s successfully marked as COMPLETED and PAID", orderUUID)
	return nil
}
