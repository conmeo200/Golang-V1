package worker

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MessagesConsumedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_system_messages_consumed_total",
			Help: "Total number of messages consumed from RabbitMQ",
		},
		[]string{"queue", "status"}, // status: success, failed
	)

	OutboxEventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_system_outbox_events_published_total",
			Help: "Total number of outbox events published to RabbitMQ",
		},
		[]string{"event_type", "status"}, // status: success, failed
	)

	MessageProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_system_message_processing_duration_seconds",
			Help:    "Histogram of message processing duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"queue"},
	)
)
