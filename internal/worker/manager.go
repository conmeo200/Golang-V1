package worker

import (
	//"context"

	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/queue/rabbitmq/consumers"
	//"github.com/conmeo200/Golang-V1/internal/repository"
	"github.com/conmeo200/Golang-V1/internal/service"
)

type Manager struct {
	rabbitMQ     *rabbitmq.RabbitMQ
	orderService service.OrderServiceInterface
	//outboxRepo   *repository.OutboxEventRepository
}

func NewManager(
	rabbitMQ *rabbitmq.RabbitMQ,
	orderService service.OrderServiceInterface,
	//outboxRepo *repository.OutboxEventRepository,
) *Manager {
	return &Manager{
		rabbitMQ:     rabbitMQ,
		orderService: orderService,
		//outboxRepo:   outboxRepo,
	}
}

func (m *Manager) Start() {
	//ctx := context.Background()

	// 1. Start Order Consumer
	orderConsumer := consumer.NewOrderConsumer(m.rabbitMQ, m.orderService)
	go orderConsumer.StartOrder()

	// 2. Start Outbox Worker
	// producer := rabbitmq.NewProducer(m.rabbitMQ)
	// outboxWorker := NewOutboxWorker(m.outboxRepo, producer)
	// go outboxWorker.Start(ctx)

	// 3. Add more workers here easily in the future...
}
