package worker

import "context"

// Worker defines the standard interface for all background jobs and consumers
type Worker interface {
	Name() string
	Start(ctx context.Context) error
	Stop() error
}
