package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	Pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(ctx context.Context, dsn string, logger *zap.Logger) (*DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("connected to database")

	return &DB{
		Pool:   pool,
		logger: logger,
	}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		db.logger.Info("database connection closed")
	}
}

func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func ConnectWithRetry(ctx context.Context, dsn string, logger *zap.Logger, maxRetries int, retryDelay time.Duration) (*DB, error) {
	var db *DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = New(ctx, dsn, logger)
		if err == nil {
			return db, nil
		}

		logger.Warn("failed to connect to database, retrying...",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err))

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelay):
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}
