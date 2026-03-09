package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Stratify-Systems/ThreatSIM/internal/core"
	goredis "github.com/redis/go-redis/v9"
)

// Stream implements core.EventStream using Redis Streams.
//
// Redis Streams is a log-based data structure perfect for event streaming:
//   - Append-only (events are never lost)
//   - Consumer groups (multiple readers)
//   - Lightweight (single Redis instance)
//
// For production at scale, swap this for the Kafka implementation.
type Stream struct {
	client *goredis.Client
	group  string // consumer group name
}

// NewStream creates a new Redis Streams event stream.
func NewStream(addr string, password string, db int) (*Stream, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Stream{
		client: client,
		group:  "threatsim",
	}, nil
}

// Publish sends an event to a Redis stream topic.
func (s *Stream) Publish(ctx context.Context, topic string, event core.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return s.client.XAdd(ctx, &goredis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Err()
}

// Subscribe listens for events on a Redis stream topic.
// It creates a consumer group if one doesn't exist, then continuously
// reads new events and passes them to the handler.
func (s *Stream) Subscribe(ctx context.Context, topic string, handler core.EventHandler) error {
	// Create consumer group (ignore error if it already exists)
	s.client.XGroupCreateMkStream(ctx, topic, s.group, "0").Err()

	consumerName := fmt.Sprintf("consumer-%d", time.Now().UnixNano())

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read new messages from the stream
		results, err := s.client.XReadGroup(ctx, &goredis.XReadGroupArgs{
			Group:    s.group,
			Consumer: consumerName,
			Streams:  []string{topic, ">"},
			Count:    10,
			Block:    time.Second,
		}).Result()

		if err != nil {
			if err == goredis.Nil {
				continue // No new messages, try again
			}
			// If context was cancelled, return cleanly
			if ctx.Err() != nil {
				return ctx.Err()
			}
			continue
		}

		for _, stream := range results {
			for _, msg := range stream.Messages {
				data, ok := msg.Values["data"].(string)
				if !ok {
					continue
				}

				var event core.Event
				if err := json.Unmarshal([]byte(data), &event); err != nil {
					continue // Skip malformed events
				}

				if err := handler(ctx, event); err != nil {
					// Log but don't stop on handler errors
					fmt.Printf("event handler error: %v\n", err)
				}

				// Acknowledge the message
				s.client.XAck(ctx, topic, s.group, msg.ID)
			}
		}
	}
}

// Close shuts down the Redis connection.
func (s *Stream) Close() error {
	return s.client.Close()
}
