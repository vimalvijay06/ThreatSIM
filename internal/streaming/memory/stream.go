package memory

import (
	"context"
	"sync"

	"github.com/Stratify-Systems/ThreatSIM/internal/core"
)

// Stream is an in-memory implementation of core.EventStream.
//
// Use this when:
//   - Running without Redis/Kafka (local development)
//   - Running tests
//   - Quick demos where you don't want external dependencies
//
// Events are delivered directly to subscribers in-process.
type Stream struct {
	mu          sync.RWMutex
	subscribers map[string][]core.EventHandler
	closed      bool
}

// NewStream creates a new in-memory event stream.
func NewStream() *Stream {
	return &Stream{
		subscribers: make(map[string][]core.EventHandler),
	}
}

// Publish sends an event directly to all subscribers on the topic.
func (s *Stream) Publish(ctx context.Context, topic string, event core.Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil
	}

	handlers := s.subscribers[topic]
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			// Log but don't stop — other subscribers should still get the event
			continue
		}
	}

	return nil
}

// Subscribe registers a handler for events on a topic.
// Unlike Redis/Kafka, this blocks until the context is cancelled.
func (s *Stream) Subscribe(ctx context.Context, topic string, handler core.EventHandler) error {
	s.mu.Lock()
	s.subscribers[topic] = append(s.subscribers[topic], handler)
	s.mu.Unlock()

	// Block until context is done (subscriber stays active)
	<-ctx.Done()
	return ctx.Err()
}

// Close marks the stream as closed.
func (s *Stream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}
