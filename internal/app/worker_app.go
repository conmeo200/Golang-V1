package app

import (
	"log"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/persistence"
	"github.com/conmeo200/Golang-V1/internal/module/order"
	"github.com/conmeo200/Golang-V1/internal/module/payment"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/worker"
	"github.com/conmeo200/Golang-V1/internal/worker/consumers"
	"github.com/conmeo200/Golang-V1/internal/worker/jobs"
)

type WorkerApp struct {
	Container *bootstrap.Container
	Registry  *worker.Registry
}

func NewWorkerApp(container *bootstrap.Container) *WorkerApp {
	// 1. Initialize Repositories
	outboxRepo := persistence.NewOutboxEventRepository(container.DB)
	inboxRepo := persistence.NewInboxEventRepository(container.DB)
	paymentRepo := persistence.NewPaymentRepository(container.DB)
	orderRepo := persistence.NewOrderRepository(container.DB)
	deadLetterRepo := persistence.NewDeadLetterRepository(container.DB)

	// 2. Initialize Services
	orderSvc := order.NewOrderService(orderRepo, nil) // passing nil for producer if not needed here
	paymentSvc := payment.NewPaymentService(paymentRepo, outboxRepo, inboxRepo)

	// 3. Initialize Registry
	registry := worker.NewRegistry()

	// 4. Register Workers
	registry.Register(jobs.NewOutboxWorker(outboxRepo, rabbitmq.NewProducer(container.RabbitMQ)))
	registry.Register(jobs.NewReconciliationWorker(paymentRepo, paymentSvc))
	
	// Register Consumers
	registry.Register(consumers.NewPaymentConsumer(orderSvc, rabbitmq.NewConsumer(container.RabbitMQ), inboxRepo))
	registry.Register(consumers.NewDLQMonitor(rabbitmq.NewConsumer(container.RabbitMQ), deadLetterRepo, "payment_completed_queue.dlq"))

	return &WorkerApp{
		Container: container,
		Registry:  registry,
	}
}

func (a *WorkerApp) Run() {
	enabled := a.Container.Config.WorkersEnabled
	if len(enabled) == 0 {
		log.Println("⚠️ No workers enabled in config (WORKERS_ENABLED is empty)")
		return
	}

	a.Registry.StartEnabledWorkers(enabled)
}

func (a *WorkerApp) Stop() {
	a.Registry.Stop()
}
