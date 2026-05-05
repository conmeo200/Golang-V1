package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DLQMonitor struct {
	rabbitConsumer *rabbitmq.Consumer
	deadLetterRepo model.DeadLetterRepo
	targetQueue    string
}

func NewDLQMonitor(rabbitConsumer *rabbitmq.Consumer, deadLetterRepo model.DeadLetterRepo, targetQueue string) *DLQMonitor {
	return &DLQMonitor{
		rabbitConsumer: rabbitConsumer,
		deadLetterRepo: deadLetterRepo,
		targetQueue:    targetQueue,
	}
}

func (m *DLQMonitor) Name() string {
	return "dlq_monitor_" + m.targetQueue
}

func (m *DLQMonitor) Start(ctx context.Context) error {
	config := rabbitmq.ConsumerConfig{
		QueueName:    m.targetQueue,
		ConsumerName: "dlq_monitor_" + m.targetQueue,
		Handler:      m.handleMessage,
		// DLQ Monitor doesn't need its own DLQ or Retries
		MaxRetries:   0, 
	}

	m.rabbitConsumer.Start(ctx, config)
	return nil
}

func (m *DLQMonitor) Stop() error {
	return nil
}

func (m *DLQMonitor) handleMessage(msg amqp.Delivery) error {
	log.Printf("🚨 [DLQMonitor] Received DEAD LETTER from queue: %s", m.targetQueue)
	
	headersJSON, _ := json.Marshal(msg.Headers)
	
	// Extract reason from x-death header if available
	reason := "unknown"
	if deaths, ok := msg.Headers["x-death"].([]interface{}); ok && len(deaths) > 0 {
		if death, ok := deaths[0].(amqp.Table); ok {
			if r, ok := death["reason"].(string); ok {
				reason = r
			}
		}
	}

	// Try to get a meaningful EventID from body or headers
	eventID := uuid.New() // default to new
	
	event := &model.DeadLetterEvent{
		EventID:      eventID,
		QueueName:    m.targetQueue,
		ExchangeName: msg.Exchange,
		RoutingKey:   msg.RoutingKey,
		Payload:      msg.Body,
		Headers:      headersJSON,
		Reason:       reason,
		Status:       "pending",
	}

	if err := m.deadLetterRepo.Create(event); err != nil {
		log.Printf("❌ [DLQMonitor] Failed to save dead letter to DB: %v", err)
		return err
	}

	log.Printf("✅ [DLQMonitor] Dead letter saved to DB with ID: %d", event.ID)
	return nil
}
