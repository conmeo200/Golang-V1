package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Publish(ctx context.Context, exchange, key string, body []byte) error {
	return r.Channel.PublishWithContext(ctx,
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.Channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}
