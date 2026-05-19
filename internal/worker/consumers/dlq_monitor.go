package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/core/constant"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DLQMonitor struct {
	rabbitConsumer *rabbitmq.Consumer
	deadLetterRepo model.DeadLetterRepo
}

func NewDLQMonitor(
	rabbitConsumer *rabbitmq.Consumer,
	deadLetterRepo model.DeadLetterRepo,
) *DLQMonitor {
	return &DLQMonitor{
		rabbitConsumer: rabbitConsumer,
		deadLetterRepo: deadLetterRepo,
	}
}

func (m *DLQMonitor) Name() string {
	return "dlq_monitor_payment"
}

func (m *DLQMonitor) Start(ctx context.Context) error {
	config := rabbitmq.ConsumerConfig{
		QueueName:    constant.QueueDLQMonitor,
		ConsumerName: m.Name(),
		Handler:      m.handleMessage,

		// DLQ monitor: không retry, không DLQ thêm
		DisableRetry: true,
		MaxRetries:   0,
	}

	log.Printf("[DLQMonitor] Starting monitor for queue: %s", constant.QueuePaymentCompletedDLQ)

	m.rabbitConsumer.Start(ctx, config)
	return nil
}

func (m *DLQMonitor) Stop() error {
	return nil
}

func (m *DLQMonitor) handleMessage(msg amqp.Delivery) error {
	eventID := extractEventID(msg)
	reason := extractDeathReason(msg.Headers)

	headersJSON, err := json.Marshal(msg.Headers)
	if err != nil {
		log.Printf("[DLQMonitor] Failed to marshal headers: %v", err)
	}

	log.Printf(
		"[DLQMonitor] queue=%s exchange=%s routing=%s event_id=%s reason=%s",
		constant.QueuePaymentCompletedDLQ,
		msg.Exchange,
		msg.RoutingKey,
		eventID.String(),
		reason,
	)

	// Deduplication (nếu repo hỗ trợ)
	if exists, err := m.deadLetterRepo.Exists(eventID, constant.QueuePaymentCompletedDLQ); err == nil && exists {
		log.Printf("[DLQMonitor] Duplicate detected, skipping event_id=%s", eventID)
		return nil
	}

	event := &model.DeadLetterEvent{
		EventID:      eventID,
		QueueName:    constant.QueuePaymentCompletedDLQ,
		ExchangeName: msg.Exchange,
		RoutingKey:   msg.RoutingKey,
		Payload:      msg.Body,
		Headers:      headersJSON,
		Reason:       reason,
		Status:       "pending",
	}

	if err := m.deadLetterRepo.Create(event); err != nil {
		log.Printf("[DLQMonitor] Failed to save dead letter: %v", err)
		return err
	}

	log.Printf("[DLQMonitor] Saved dead letter ID=%d event_id=%s", event.ID, eventID)

	return nil
}