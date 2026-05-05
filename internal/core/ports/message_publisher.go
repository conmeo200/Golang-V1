package ports

import "context"

type MessagePublisher interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}
