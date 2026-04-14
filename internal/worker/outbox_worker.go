package worker

import (
	"context"
	"log"
	//"math"
	"time"

	//"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/repository"
)

type OutboxWorker struct {
	outboxRepo *repository.OutboxEventRepository
	producer   *rabbitmq.Producer
	maxRetries int
	limit      int
}

func NewOutboxWorker(outboxRepo *repository.OutboxEventRepository, producer *rabbitmq.Producer) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo: outboxRepo,
		producer:   producer,
		maxRetries: 5,
		limit:      100,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Println("🚀 Outbox Worker started")

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Println("🛑 Outbox Worker stopping...")
	// 		return
	// 	case <-ticker.C:
	// 		w.processOutbox(ctx)
	// 	}
	// }
}

func (w *OutboxWorker) processOutbox(ctx context.Context) {
	// events, err := w.outboxRepo.FetchPending(ctx, w.limit)
	// if err != nil {
	// 	log.Printf("❌ Failed to fetch pending events: %v", err)
	// 	return
	// }

	// if len(events) == 0 {
	// 	return
	// }

	// for _, event := range events {
	// 	w.processEvent(ctx, &event)
	// }
}

// func (w *OutboxWorker) processEvent(ctx context.Context, event *model.OutboxEvent) {
// 	// Publish to RabbitMQ
// 	// We use the EventType as the routing key, and a default exchange for now
// 	err := w.producer.Publish(ctx, "payment.exchange", event.EventType, event.Payload)
	
// 	now := time.Now().Unix()

// 	if err != nil {
// 		log.Printf("⚠️ Failed to publish event %s: %v", event.EventID, err)
// 		w.handleFailure(ctx, event, now)
// 		return
// 	}

// 	// Success
// 	if err := w.outboxRepo.MarkAsPublished(ctx, event.ID, now); err != nil {
// 		log.Printf("❌ Failed to mark event as published: %v", err)
// 	} else {
// 		log.Printf("✅ Event %s published successfully", event.EventID)
// 	}
// }

// func (w *OutboxWorker) handleFailure(ctx context.Context, event *model.OutboxEvent, now int64) {
// 	event.RetryCount++
	
// 	if event.RetryCount >= w.maxRetries {
// 		event.Status = "FAILED"
// 		log.Printf("❌ Event %s reached max retries and marked as FAILED", event.EventID)
// 	} else {
// 		// Exponential backoff: 2^retry_count * 30 seconds
// 		backoff := int64(math.Pow(2, float64(event.RetryCount))) * 30
// 		event.NextRetryAt = now + backoff
// 		log.Printf("🔄 Event %s will be retried in %d seconds (retry %d/%d)", event.EventID, backoff, event.RetryCount, w.maxRetries)
// 	}

// 	if err := w.outboxRepo.Update(ctx, event); err != nil {
// 		log.Printf("❌ Failed to update event status: %v", err)
// 	}
// }
