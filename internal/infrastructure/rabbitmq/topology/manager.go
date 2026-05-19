package topology

import (
	"log"

	"github.com/conmeo200/Golang-V1/internal/core/constant"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupTopology initializes all exchanges, queues, and bindings for the application.
// This ensures that all components use the same infrastructure configuration.
func SetupTopology(rmq *rabbitmq.RabbitMQ) error {
	conn, err := rmq.GetConnection()
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("🛠️ Initializing RabbitMQ Topology...")

	// 1. Declare Exchanges
	exchanges := []string{
		constant.ExchangePayment,
		constant.ExchangeRetry,
		constant.ExchangeDLX,
	}
	for _, ex := range exchanges {
		if err := ch.ExchangeDeclare(ex, "direct", true, false, false, false, nil); err != nil {
			return err
		}
	}

	// 2. Declare Queues
	
	// Main Payment Queue
	// Note: We don't declare x-dead-letter args here if we want the consumer to control them,
	// BUT for a centralized topology, it's better to declare them here.
	paymentArgs := amqp.Table{
		"x-dead-letter-exchange":    constant.ExchangeRetry,
		"x-dead-letter-routing-key": constant.RoutingPaymentCompletedRetry,
	}
	if _, err := ch.QueueDeclare(constant.QueuePaymentCompleted, true, false, false, false, paymentArgs); err != nil {
		return err
	}

	// Retry Queue
	retryArgs := amqp.Table{
		"x-dead-letter-exchange":    constant.ExchangePayment,
		"x-dead-letter-routing-key": constant.RoutingPaymentCompleted,
		"x-message-ttl":             5000, // 5 seconds
	}
	if _, err := ch.QueueDeclare(constant.QueuePaymentCompletedRetry, true, false, false, false, retryArgs); err != nil {
		return err
	}

	// Dead Letter Queue
	if _, err := ch.QueueDeclare(constant.QueuePaymentCompletedDLQ, true, false, false, false, nil); err != nil {
		return err
	}

	// 3. Bindings
	
	// Payment Completed (Main)
	if err := ch.QueueBind(constant.QueuePaymentCompleted, constant.RoutingPaymentCompleted, constant.ExchangePayment, false, nil); err != nil {
		return err
	}

	// Payment Retry
	if err := ch.QueueBind(constant.QueuePaymentCompletedRetry, constant.RoutingPaymentCompletedRetry, constant.ExchangeRetry, false, nil); err != nil {
		return err
	}

	// Payment DLQ
	if err := ch.QueueBind(constant.QueuePaymentCompletedDLQ, constant.RoutingPaymentCompletedFailed, constant.ExchangeDLX, false, nil); err != nil {
		return err
	}

	log.Println("✅ RabbitMQ Topology initialized successfully")
	return nil
}
