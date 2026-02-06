package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	AnnouncementsExchange = "announcements"
	LogsExchange          = "logs"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
	mu      sync.Mutex
	url     string
}

func NewRabbitMQPublisher(url string, logger *zap.Logger) (*RabbitMQPublisher, error) {
	p := &RabbitMQPublisher{
		url:    url,
		logger: logger,
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	return p, nil
}

func NewRabbitMQPublisherWithRetry(url string, logger *zap.Logger, maxRetries int, retryDelay time.Duration) (*RabbitMQPublisher, error) {
	var p *RabbitMQPublisher
	var err error

	for i := 0; i < maxRetries; i++ {
		p, err = NewRabbitMQPublisher(url, logger)
		if err == nil {
			return p, nil
		}

		logger.Warn("failed to connect to RabbitMQ, retrying...",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err))

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

func (p *RabbitMQPublisher) connect() error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchanges
	if err := channel.ExchangeDeclare(
		AnnouncementsExchange,
		"fanout",
		false, // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("failed to declare announcements exchange: %w", err)
	}

	if err := channel.ExchangeDeclare(
		LogsExchange,
		"direct",
		false, // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("failed to declare logs exchange: %w", err)
	}

	p.conn = conn
	p.channel = channel

	p.logger.Info("connected to RabbitMQ")
	return nil
}

func (p *RabbitMQPublisher) ensureConnected() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn == nil || p.conn.IsClosed() {
		p.logger.Info("reconnecting to RabbitMQ")
		return p.connect()
	}
	return nil
}

func (p *RabbitMQPublisher) PublishAnnouncement(ctx context.Context, message string) error {
	if err := p.ensureConnected(); err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		AnnouncementsExchange,
		"", // routing key (ignored for fanout)
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (p *RabbitMQPublisher) PublishToOperator(ctx context.Context, rmqQueue string, message string) error {
	if err := p.ensureConnected(); err != nil {
		return err
	}

	// The routing key is the operator's rmq_queue field from the database
	return p.channel.PublishWithContext(
		ctx,
		LogsExchange,
		rmqQueue, // routing key = operator's rmq_queue
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (p *RabbitMQPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}

	p.logger.Info("RabbitMQ connection closed")
	return nil
}
