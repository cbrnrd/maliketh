package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/builder"
	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/database"
	"github.com/cbrnrd/maliketh/server-go/internal/handler/admin"
	"github.com/cbrnrd/maliketh/server-go/internal/logging"
	"github.com/cbrnrd/maliketh/server-go/internal/messaging"
	"github.com/cbrnrd/maliketh/server-go/internal/repository"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger, err := logging.New(cfg.LogLevel, cfg.Env)
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("starting maliketh admin server",
		zap.String("env", cfg.Env),
		zap.String("log_level", cfg.LogLevel))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.ConnectWithRetry(ctx, cfg.Postgres.DSN(), logger, 5, 5*time.Second)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	repos := repository.New(db.Pool)

	// Initialize publisher (use noop if RabbitMQ not configured)
	var publisher messaging.Publisher = &messaging.NoopPublisher{}
	// TODO: Initialize RabbitMQ publisher when needed

	// Load server private key for operator authentication
	serverPrivateKey := os.Getenv("SERVER_PRIVATE_KEY")
	if serverPrivateKey == "" {
		logger.Warn("SERVER_PRIVATE_KEY not set, operator authentication will fail")
	}

	// Initialize builder
	builderCfg := &builder.Config{
		CppImplantDir: getEnvOrDefault("CPP_IMPLANT_DIR", "/implant"),
		GoImplantDir:  getEnvOrDefault("GO_IMPLANT_DIR", "/go_implant"),
	}
	implantBuilder := builder.New(builderCfg, logger)
	logger.Info("initialized implant builder",
		zap.String("cpp_dir", builderCfg.CppImplantDir),
		zap.String("go_dir", builderCfg.GoImplantDir))

	srv := server.New(cfg.Admin.Addr(), logger)

	srv.Router().Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Health(r.Context()); err != nil {
			server.Error(w, http.StatusServiceUnavailable, "database unhealthy")
			return
		}
		server.Success(w, "healthy")
	})

	// Register admin routes
	adminHandler := admin.NewHandler(
		repos, cfg, publisher, logger, serverPrivateKey,
		admin.WithBuilder(implantBuilder),
	)
	adminHandler.RegisterRoutes(srv.Router())

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
