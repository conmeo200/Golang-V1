package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	rmq     *RabbitMQ
	channel *amqp.Channel
	mu      sync.RWMutex
}

func NewProducer(rmq *RabbitMQ) *Producer {
	p := &Producer{
		rmq: rmq,
	}
	// Automatically recreate channel when the connection reconnects
	rmq.RegisterHook(p.onReconnect)
	return p
}

func (p *Producer) onReconnect(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("❌ Failed to open channel for producer: %v", err)
		return
	}
	p.channel = ch
	log.Println("✅ Producer channel established")
}

func (p *Producer) getChannel() (*amqp.Channel, error) {
	p.mu.RLock()
	ch := p.channel
	p.mu.RUnlock()

	// If channel hasn't been created yet or is closed, recreate it transparently
	if ch == nil || ch.IsClosed() {
		conn, err := p.rmq.GetConnection()
		if err != nil {
			return nil, err
		}

		p.mu.Lock()
		defer p.mu.Unlock()

		// Double-check locking
		if p.channel != nil && !p.channel.IsClosed() {
			return p.channel, nil
		}

		ch, err = conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to open channel: %w", err)
		}
		p.channel = ch
	}
	return p.channel, nil
}

// Publish sends a message to RabbitMQ with DeliveryMode Persistent (Best Practice)
func (p *Producer) Publish(ctx context.Context, exchange, key string, body []byte) error {
	ch, err := p.getChannel()
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx,
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// PublishWithRetry retries publishing if it fails due to network or channel issues
func (p *Producer) PublishWithRetry(ctx context.Context, exchange, key string, body []byte) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := p.Publish(ctx, exchange, key, body)
		if err == nil {
			return nil
		}
		log.Printf("⚠️ Publish failed, retrying %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to publish message after %d retries", maxRetries)
}

func (p *Producer) PublishOrderCreated(ctx context.Context, body []byte) error {
	return p.PublishWithRetry(ctx, "order.exchange", "order.created", body)
}

func (p *Producer) DeclareQueue(name string, args amqp.Table) (amqp.Queue, error) {
	ch, err := p.getChannel()
	if err != nil {
		return amqp.Queue{}, err
	}
	return ch.QueueDeclare(
		name,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		args,
	)
}

func (p *Producer) DeclareExchange(name, kind string) error {
	ch, err := p.getChannel()
	if err != nil {
		return err
	}
	return ch.ExchangeDeclare(
		name,
		kind,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)
}
