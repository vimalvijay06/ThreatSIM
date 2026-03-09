package core

import "context"

// EventStream is the interface for the event streaming layer.
//
// This abstraction allows us to swap implementations:
//   - Redis Streams (default, lightweight, great for local dev)
//   - Kafka (production-scale, distributed)
//
// The pipeline works like this:
//
//	Attack Plugin → Publish(event) → [Stream] → Subscribe(handler) → Detection Engine
type EventStream interface {
	// Publish sends an event to a topic (e.g., "attack-events")
	Publish(ctx context.Context, topic string, event Event) error

	// Subscribe listens for events on a topic and calls the handler for each one
	Subscribe(ctx context.Context, topic string, handler EventHandler) error

	// Close cleanly shuts down the stream connection
	Close() error
}

// EventHandler processes a single event received from the stream
type EventHandler func(ctx context.Context, event Event) error

// Default topic names
const (
	TopicAttackEvents = "attack-events"
	TopicAlerts       = "alerts"
)
