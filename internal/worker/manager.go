package worker

import (
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq/consumers"
	"github.com/conmeo200/Golang-V1/internal/service"
)

type Manager struct {
	rabbitMQ      *rabbitmq.RabbitMQ
	orderService service.OrderServiceInterface
}

func NewManager(rabbitMQ *rabbitmq.RabbitMQ, orderService service.OrderServiceInterface) *Manager {
	return &Manager{
		rabbitMQ:      rabbitMQ,
		orderService: orderService,
	}
}

func (m *Manager) Start() {
	// 1. Start Order Consumer
	// Note: We use the alias 'consumer' because the package name in that folder is 'consumer'
	orderConsumer := consumer.NewOrderConsumer(m.rabbitMQ, m.orderService)
	go orderConsumer.StartOrder()

	// 2. Add more workers here easily in the future...
}
