package consumers

import (
	"encoding/json"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// extractEventID tries to get ID from headers, then payload, then new UUID
func extractEventID(msg amqp.Delivery) uuid.UUID {
	// 1. Try header
	if id, ok := msg.Headers["event_id"].(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			return parsed
		}
	}

	// 2. Try payload
	var payload struct {
		EventID string `json:"event_id"`
	}
	if err := json.Unmarshal(msg.Body, &payload); err == nil {
		if parsed, err := uuid.Parse(payload.EventID); err == nil {
			return parsed
		}
	}

	// 3. Fallback
	return uuid.New()
}

func extractDeathReason(headers amqp.Table) string {
	deaths, ok := headers["x-death"].([]interface{})
	if !ok || len(deaths) == 0 {
		return "unknown"
	}

	// lấy record cuối (lần chết gần nhất)
	last := deaths[len(deaths)-1]

	deathMap, ok := last.(amqp.Table)
	if !ok {
		return "unknown"
	}

	if reason, ok := deathMap["reason"].(string); ok {
		return reason
	}

	return "unknown"
}
