package messaging

import "context"

// Publisher defines the interface for publishing messages to operators.
type Publisher interface {
	// PublishAnnouncement sends a message to all operators.
	PublishAnnouncement(ctx context.Context, message string) error

	// PublishToOperator sends a message to a specific operator's RMQ queue.
	PublishToOperator(ctx context.Context, rmqQueue string, message string) error

	// Close closes the publisher connection.
	Close() error
}

// NoopPublisher is a no-op implementation of Publisher for when messaging is disabled.
type NoopPublisher struct{}

func (n *NoopPublisher) PublishAnnouncement(ctx context.Context, message string) error {
	return nil
}

func (n *NoopPublisher) PublishToOperator(ctx context.Context, rmqQueue string, message string) error {
	return nil
}

func (n *NoopPublisher) Close() error {
	return nil
}
