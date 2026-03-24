package consumer

import (
	"encoding/json"
	"log"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/logger"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderConsumer struct {
	consumer     *rabbitmq.Consumer
	orderService service.OrderServiceInterface
}

func NewOrderConsumer(rabbitMQ *rabbitmq.RabbitMQ, orderService service.OrderServiceInterface) *OrderConsumer {
	return &OrderConsumer{
		consumer:     rabbitmq.NewConsumer(rabbitMQ),
		orderService: orderService,
	}
}

// StartOrder initializes the order consumer via the abstract RabbitMQ engine
func (c *OrderConsumer) StartOrder() {
	c.consumer.Start(rabbitmq.ConsumerConfig{
		QueueName:     "order.queue",
		ExchangeName:  "order.exchange",
		RoutingKey:    "order.created",
		ConsumerName:  "order-consumer",
		PrefetchCount: 1,
		RetryTTL:      5000,
		MaxRetries:    3,

		// The abstract Engine handles Nacks and auto DLX/Retry queue publications if Handler returns an error
		Handler: c.handleOrderMessage,
	})
}

// handleOrderMessage contains ONLY the business logic for processing an incoming message
func (c *OrderConsumer) handleOrderMessage(msg amqp.Delivery) error {
	var wrapper struct {
		Data dto.OrderMessage `json:"data"`
	}

	if err := json.Unmarshal(msg.Body, &wrapper); err != nil {
		logger.ErrorLogger.Printf("Unmarshal error in consumer: %v", err)
		// Return error so the Engine sends it to the DLX/Retry flow
		return err 
	}

	dataBody := wrapper.Data
	log.Printf("📥 Processing order UUID: %s", dataBody.OrderUUID)

	err := c.orderService.ProcessOrder(dataBody)
	if err != nil {
		// Just return the processing error, the Engine will magically Nack it and handle the retry
		return err
	}

	log.Printf("✅ Successfully processed order UUID: %s", dataBody.OrderUUID)
	return nil
}