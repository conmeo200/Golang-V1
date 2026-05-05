package dead_letter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DeadLetterService struct {
	repo     model.DeadLetterRepo
	producer *rabbitmq.Producer
}

func NewDeadLetterService(repo model.DeadLetterRepo, producer *rabbitmq.Producer) *DeadLetterService {
	return &DeadLetterService{
		repo:     repo,
		producer: producer,
	}
}

func (s *DeadLetterService) ListPending(ctx context.Context) ([]model.DeadLetterEvent, error) {
	return s.repo.ListPending()
}

// ReplayMessage sends the message back to its original exchange and routing key
func (s *DeadLetterService) ReplayMessage(ctx context.Context, id uint) error {
	event, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if event.Status != "pending" {
		return fmt.Errorf("message %d is already %s", id, event.Status)
	}

	var headers amqp.Table
	if len(event.Headers) > 0 {
		if err := json.Unmarshal(event.Headers, &headers); err != nil {
			log.Printf("⚠️ Failed to unmarshal headers for replay: %v", err)
		}
	}

	// Publish back to original exchange
	err = s.producer.Publish(ctx, event.ExchangeName, event.RoutingKey, event.Payload)
	if err != nil {
		return fmt.Errorf("failed to replay message: %w", err)
	}

	return s.repo.MarkAsReplayed(id)
}

func (s *DeadLetterService) ResolveMessage(ctx context.Context, id uint) error {
	return s.repo.MarkAsResolved(id)
}
