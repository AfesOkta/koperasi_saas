package event

import (
	"context"
	"log"
	"sync"
)

// MemoryEventBus implements an in-memory event bus for local development.
type MemoryEventBus struct {
	subscribers []func(Event) error
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewMemoryEventBus creates a new in-memory event bus.
func NewMemoryEventBus() *MemoryEventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &MemoryEventBus{
		subscribers: make([]func(Event) error, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Publish sends an event to all registered subscribers.
func (b *MemoryEventBus) Publish(ctx context.Context, evt Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	log.Printf("📤 [Memory] Event published: %s (aggregate_id=%d)", evt.Type, evt.AggregateID)

	for _, handler := range b.subscribers {
		// Run in a goroutine to simulate async behavior like a real message queue
		go func(h func(Event) error, e Event) {
			if err := h(e); err != nil {
				log.Printf("❌ [Memory] Error handling event %s: %v", e.Type, err)
			}
		}(handler, evt)
	}

	return nil
}

// Consume registers a handler for events.
func (b *MemoryEventBus) Consume(ctx context.Context, handler func(Event) error) {
	b.mu.Lock()
	b.subscribers = append(b.subscribers, handler)
	b.mu.Unlock()

	// Keep alive until context is cancelled
	<-ctx.Done()
}

// Close cancels the context and prevents further event dispatching.
func (b *MemoryEventBus) Close() error {
	b.cancel()
	return nil
}
