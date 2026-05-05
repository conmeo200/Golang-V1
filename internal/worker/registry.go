package worker

import (
	"context"
	"log"
	"sync"
)

type Registry struct {
	workers map[string]Worker
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewRegistry() *Registry {
	ctx, cancel := context.WithCancel(context.Background())
	return &Registry{
		workers: make(map[string]Worker),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Register adds a new worker to the registry
func (r *Registry) Register(w Worker) {
	r.workers[w.Name()] = w
}

// StartEnabledWorkers starts only the workers that are enabled in the configuration
func (r *Registry) StartEnabledWorkers(enabledList []string) {
	for _, name := range enabledList {
		if w, exists := r.workers[name]; exists {
			r.wg.Add(1)
			go func(worker Worker) {
				defer r.wg.Done()
				log.Printf("🚀 Starting worker: [%s]", worker.Name())
				if err := worker.Start(r.ctx); err != nil {
					log.Printf("❌ Worker [%s] failed: %v", worker.Name(), err)
				}
			}(w)
		} else {
			log.Printf("⚠️ Warning: Enabled worker [%s] not found in registry", name)
		}
	}
}

// Stop gracefully shuts down all active workers
func (r *Registry) Stop() {
	log.Println("🛑 Sending stop signal to all workers...")
	r.cancel() // Cancel the context passed to all workers
	
	// Also call Stop() on each worker if they implement custom teardown
	for _, w := range r.workers {
		if err := w.Stop(); err != nil {
			log.Printf("⚠️ Worker [%s] error during stop: %v", w.Name(), err)
		}
	}

	r.wg.Wait()
	log.Println("✅ All workers stopped gracefully.")
}
