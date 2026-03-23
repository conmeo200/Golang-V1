package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	rmq     *RabbitMQ
	channel *amqp.Channel
	mu      sync.RWMutex
}

func NewConsumer(rmq *RabbitMQ) *Consumer {
	return &Consumer{
		rmq: rmq,
	}
}

// ConsumerConfig configuration for the abstract consumer engine
type ConsumerConfig struct {
	QueueName         string
	RoutingKey        string
	ExchangeName      string
	PrefetchCount     int
	ConsumerName      string

	// Retry constraints
	RetryTTL          int // in milliseconds, e.g., 5000
	MaxRetries        int // e.g., 3

	// DLX and Retry Topology strings (Optional, defaults apply if empty)
	DLXExchange       string // default: "dlx.exchange"
	DLXQueue          string // default: queueName + ".dlq"
	DLXRoutingKey     string // default: queueName + ".failed"
	
	RetryExchange     string // default: "retry.exchange"
	RetryQueue        string // default: queueName + ".retry.queue"
	RetryRoutingKey   string // default: queueName + ".retry"

	// Handler is the business logic function. 
	// Returning an error triggers the Retry/DLQ mechanism.
	Handler           func(msg amqp.Delivery) error 
}

func (cfg *ConsumerConfig) setDefaults() {
	if cfg.DLXExchange == "" {
		cfg.DLXExchange = "dlx.exchange"
	}
	if cfg.DLXQueue == "" {
		cfg.DLXQueue = cfg.QueueName + ".dlq"
	}
	if cfg.DLXRoutingKey == "" {
		cfg.DLXRoutingKey = cfg.QueueName + ".failed"
	}
	if cfg.RetryExchange == "" {
		cfg.RetryExchange = "retry.exchange"
	}
	if cfg.RetryQueue == "" {
		cfg.RetryQueue = cfg.QueueName + ".retry.queue"
	}
	if cfg.RetryRoutingKey == "" {
		cfg.RetryRoutingKey = cfg.QueueName + ".retry"
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.PrefetchCount == 0 {
		cfg.PrefetchCount = 1
	}
}

// Start wraps the consumer in a self-healing background loop handling all topology
func (c *Consumer) Start(cfg ConsumerConfig) {
	cfg.setDefaults()

	go func() {
		for {
			err := c.consumeLoop(cfg)
			if err != nil {
				log.Printf("⚠️ Consumer [%s] stopped: %v. Restarting in 5s...", cfg.ConsumerName, err)
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func (c *Consumer) consumeLoop(cfg ConsumerConfig) error {
	// 1. Setup DLX
	if err := c.SetupDLX(cfg.DLXExchange, cfg.DLXQueue, cfg.DLXRoutingKey); err != nil {
		return err
	}

	// 2. Setup Delayed Retry
	if err := c.SetupDelayedRetry(cfg.RetryExchange, cfg.RetryQueue, cfg.RetryRoutingKey, cfg.ExchangeName, cfg.RoutingKey, cfg.RetryTTL); err != nil {
		return err
	}

	// 3. Declare Main Queue (Failures go to retry.exchange)
	args := amqp.Table{
		"x-dead-letter-exchange":    cfg.RetryExchange,
		"x-dead-letter-routing-key": cfg.RetryRoutingKey,
	}
	if _, err := c.DeclareQueue(cfg.QueueName, args); err != nil {
		return err
	}

	// 4. Bind Main Queue
	if err := c.QueueBind(cfg.QueueName, cfg.RoutingKey, cfg.ExchangeName); err != nil {
		return err
	}

	// 5. QoS
	if err := c.Qos(cfg.PrefetchCount); err != nil {
		return err
	}

	// 6. Start Message Channel
	msgs, err := c.Consume(cfg.QueueName, cfg.ConsumerName, false)
	if err != nil {
		return err
	}

	log.Printf("🚀 Consumer [%s] started successfully on queue [%s]...", cfg.ConsumerName, cfg.QueueName)

	for msg := range msgs {
		err := cfg.Handler(msg)
		if err != nil {
			retryCount := c.GetRetryCount(msg)

			if retryCount >= cfg.MaxRetries {
				log.Printf("❌ Consumer [%s]: Max retry reached, sending to final DLQ", cfg.ConsumerName)
				
				ch, err := c.GetUnderlyingChannel()
				if err == nil {
					ch.Publish(
						cfg.DLXExchange,
						cfg.DLXRoutingKey,
						false, false,
						amqp.Publishing{
							Headers:      msg.Headers,
							Body:         msg.Body,
							DeliveryMode: amqp.Persistent,
						},
					)
					msg.Ack(false)
				} else {
					msg.Nack(false, true)
				}
				continue
			}

			log.Printf("🔁 Consumer [%s]: Processing failed, triggering delayed retry via Nack: %d", cfg.ConsumerName, retryCount+1)
			msg.Nack(false, false) // goes to retry.exchange safely
			continue
		}

		msg.Ack(false)
	}

	return fmt.Errorf("underlying rabbitmq delivery channel closed")
}

func (c *Consumer) getChannel() (*amqp.Channel, error) {
	c.mu.RLock()
	ch := c.channel
	c.mu.RUnlock()

	if ch == nil || ch.IsClosed() {
		conn, err := c.rmq.GetConnection()
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		defer c.mu.Unlock()

		if c.channel != nil && !c.channel.IsClosed() {
			return c.channel, nil
		}

		ch, err = conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to open channel: %w", err)
		}
		c.channel = ch
	}
	return c.channel, nil
}

func (c *Consumer) Consume(queueName, consumerName string, autoAck bool) (<-chan amqp.Delivery, error) {
	ch, err := c.getChannel()
	if err != nil {
		return nil, err
	}
	return ch.Consume(queueName, consumerName, autoAck, false, false, false, nil)
}

func (c *Consumer) Qos(prefetchCount int) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}
	return ch.Qos(prefetchCount, 0, false)
}

func (c *Consumer) GetRetryCount(msg amqp.Delivery) int {
	if deaths, ok := msg.Headers["x-death"].([]interface{}); ok {
		if len(deaths) > 0 {
			if deathInfo, ok := deaths[0].(amqp.Table); ok {
				if count, ok := deathInfo["count"].(int64); ok {
					return int(count)
				}
			}
		}
	}
	return 0
}

func (c *Consumer) SetupDLX(dlxName, dlqName, routingKey string) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return ch.QueueBind(dlqName, routingKey, dlxName, false, nil)
}

func (c *Consumer) SetupDelayedRetry(retryExchange, retryQueue, retryRoutingKey, mainExchange, mainRoutingKey string, ttlMs int) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(retryExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    mainExchange,
		"x-dead-letter-routing-key": mainRoutingKey,
		"x-message-ttl":             ttlMs,
	}
	_, err = ch.QueueDeclare(retryQueue, true, false, false, false, args)
	if err != nil {
		return err
	}

	return ch.QueueBind(retryQueue, retryRoutingKey, retryExchange, false, nil)
}

func (c *Consumer) DeclareQueue(name string, args amqp.Table) (amqp.Queue, error) {
	ch, err := c.getChannel()
	if err != nil {
		return amqp.Queue{}, err
	}
	return ch.QueueDeclare(name, true, false, false, false, args)
}

func (c *Consumer) QueueBind(name, key, exchange string) error {
	ch, err := c.getChannel()
	if err != nil {
		return err
	}
	return ch.QueueBind(name, key, exchange, false, nil)
}

func (c *Consumer) GetUnderlyingChannel() (*amqp.Channel, error) {
	return c.getChannel()
}