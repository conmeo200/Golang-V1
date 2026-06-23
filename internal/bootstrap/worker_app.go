package bootstrap

import (
	"log"

	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq"
	"github.com/conmeo200/Golang-V1/internal/infrastructure/rabbitmq/topology"
	"github.com/conmeo200/Golang-V1/internal/worker"
	"github.com/conmeo200/Golang-V1/internal/worker/consumers"
	"github.com/conmeo200/Golang-V1/internal/worker/jobs"
	paymentOut "github.com/conmeo200/Golang-V1/internal/module/payment/adapter/out"
)

type WorkerApp struct {
	Container *Container
	Registry  *worker.Registry
}

func NewWorkerApp(container *Container) *WorkerApp {
	// 1.5 Setup RabbitMQ Topology
	if err := topology.SetupTopology(container.RabbitMQ); err != nil {
		log.Printf("⚠️ Failed to setup RabbitMQ topology: %v", err)
	}

	// For the outbox worker and consumer, we need the specific implementation or the port.
	// We'll instantiate outbox/inbox repos directly or type assert from the container if needed.
	outboxRepo := paymentOut.NewOutboxEventRepository(container.DB)
	inboxRepo := paymentOut.NewInboxEventRepository(container.DB)

	// 3. Initialize Registry
	registry := worker.NewRegistry()

	// 4. Register Workers
	registry.Register(jobs.NewOutboxWorker(outboxRepo, rabbitmq.NewProducer(container.RabbitMQ)))
	registry.Register(jobs.NewReconciliationWorker(container.PaymentRepo, container.PaymentService))

	// Register Consumers
	registry.Register(consumers.NewPaymentConsumer(container.OrderService, rabbitmq.NewConsumer(container.RabbitMQ), inboxRepo))

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

	//a.Registry.StartEnabledWorkers(enabled)
}

func (a *WorkerApp) Stop() {
	a.Registry.Stop()
}
