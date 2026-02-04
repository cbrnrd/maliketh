package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/database"
	"github.com/cbrnrd/maliketh/server-go/internal/logging"
	"go.uber.org/zap"
)

func main() {
	up := flag.Bool("up", false, "run migrations up")
	down := flag.Bool("down", false, "run migrations down")
	version := flag.Bool("version", false, "show current migration version")
	migrationsPath := flag.String("path", "./migrations", "path to migrations directory")
	flag.Parse()

	if !*up && !*down && !*version {
		fmt.Println("Usage: migrate -up|-down|-version [-path ./migrations]")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.New(cfg.LogLevel, cfg.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	migrator, err := database.NewMigrator(cfg.Postgres.DSN(), *migrationsPath, logger)
	if err != nil {
		logger.Fatal("failed to create migrator", zap.Error(err))
	}
	defer migrator.Close()

	switch {
	case *up:
		if err := migrator.Up(); err != nil {
			logger.Fatal("migration up failed", zap.Error(err))
		}
	case *down:
		if err := migrator.Down(); err != nil {
			logger.Fatal("migration down failed", zap.Error(err))
		}
	case *version:
		v, dirty, err := migrator.Version()
		if err != nil {
			logger.Fatal("failed to get version", zap.Error(err))
		}
		logger.Info("current migration version", zap.Uint("version", v), zap.Bool("dirty", dirty))
	}
}
