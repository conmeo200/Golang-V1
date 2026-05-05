package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MessagesConsumedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_messages_consumed_total",
			Help: "The total number of consumed messages",
		},
		[]string{"status", "queue"},
	)

	MessageProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_message_processing_duration_seconds",
			Help:    "Time spent processing a message",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"queue"},
	)

	OutboxEventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "outbox_events_published_total",
			Help: "The total number of published outbox events",
		},
		[]string{"status"},
	)
)
