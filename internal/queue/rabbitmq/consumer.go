package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
