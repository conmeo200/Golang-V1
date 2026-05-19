package constant

const (
	// Exchanges
	ExchangePayment = "payment.exchange"
	ExchangeRetry   = "retry.exchange"
	ExchangeDLX     = "dlx.exchange"

	// Queues
	QueuePaymentCompleted      = "payment_completed_queue"
	QueuePaymentCompletedRetry = "payment_completed_queue.retry.queue"
	QueuePaymentCompletedDLQ   = "payment_completed_queue.dlq"
	QueueDLQMonitor            = "dlq_monitor_queue.dlq"

	// Routing Keys
	RoutingPaymentCompleted       = "PaymentCompleted"
	RoutingPaymentCompletedRetry  = "payment_completed_queue.retry"
	RoutingPaymentCompletedFailed = "payment_completed_queue.failed"
)
