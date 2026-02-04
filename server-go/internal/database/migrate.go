package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

type Migrator struct {
	m      *migrate.Migrate
	logger *zap.Logger
}

func NewMigrator(dsn, migrationsPath string, logger *zap.Logger) (*Migrator, error) {
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{
		m:      m,
		logger: logger,
	}, nil
}

func (m *Migrator) Up() error {
	m.logger.Info("running database migrations")

	if err := m.m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("no migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	m.logger.Info("migrations applied successfully")
	return nil
}

func (m *Migrator) Down() error {
	m.logger.Info("rolling back database migrations")

	if err := m.m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("no migrations to roll back")
			return nil
		}
		return fmt.Errorf("rollback failed: %w", err)
	}

	m.logger.Info("migrations rolled back successfully")
	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.m.Version()
}

func (m *Migrator) Close() error {
	srcErr, dbErr := m.m.Close()
	if srcErr != nil {
		return srcErr
	}
	return dbErr
}
