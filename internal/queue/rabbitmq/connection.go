package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/conmeo200/Golang-V1/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	addr            string
	conn            *amqp.Connection
	mu              sync.RWMutex
	done            chan bool
	notifyConnClose chan *amqp.Error
	isConnected     atomic.Bool
	
	// Channels for listeners to know when reconnection happens
	reconnectHooks []func(*amqp.Connection)
}

func NewRabbitMQ(cfg *config.Config) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)

	rmq := &RabbitMQ{
		addr: url,
		done: make(chan bool),
	}

	// Initial connection must succeed or fail fast
	err := rmq.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to initially connect to RabbitMQ: %w", err)
	}

	go rmq.handleReconnect()

	return rmq, nil
}

func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp.Dial(r.addr)
	if err != nil {
		return err
	}

	r.conn = conn
	r.notifyConnClose = make(chan *amqp.Error, 1)
	r.conn.NotifyClose(r.notifyConnClose)
	
	r.isConnected.Store(true)
	log.Println("✅ Connected to RabbitMQ successfully")

	// Trigger hooks on successful connection
	for _, hook := range r.reconnectHooks {
		go hook(conn)
	}

	return nil
}

func (r *RabbitMQ) handleReconnect() {
	for {
		select {
		case <-r.done:
			return
		case err := <-r.notifyConnClose:
			log.Printf("⚠️ RabbitMQ connection closed. Reconnecting... (Error: %v)", err)
			r.isConnected.Store(false)
			
			r.reconnectWithBackoff()
		}
	}
}

func (r *RabbitMQ) reconnectWithBackoff() {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-r.done:
			return
		default:
			err := r.connect()
			if err == nil {
				return
			}
			log.Printf("❌ RabbitMQ reconnect failed: %v. Retrying in %v", err, backoff)
			time.Sleep(backoff)
			
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

// RegisterHook allows components (like Consumers) to register a callback 
// when the connection is established/re-established.
func (r *RabbitMQ) RegisterHook(hook func(*amqp.Connection)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.reconnectHooks = append(r.reconnectHooks, hook)
	
	// If already connected, call it immediately
	if r.IsConnected() {
		go hook(r.conn)
	}
}

func (r *RabbitMQ) IsConnected() bool {
	return r.isConnected.Load()
}

func (r *RabbitMQ) GetConnection() (*amqp.Connection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.IsConnected() || r.conn == nil {
		return nil, fmt.Errorf("rabbitmq is not connected")
	}
	return r.conn, nil
}

func (r *RabbitMQ) Close() {
	close(r.done)
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.conn != nil {
		r.conn.Close()
	}
	r.isConnected.Store(false)
}
