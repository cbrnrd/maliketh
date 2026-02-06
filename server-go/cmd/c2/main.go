package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/database"
	c2handler "github.com/cbrnrd/maliketh/server-go/internal/handler/c2"
	"github.com/cbrnrd/maliketh/server-go/internal/logging"
	"github.com/cbrnrd/maliketh/server-go/internal/messaging"
	"github.com/cbrnrd/maliketh/server-go/internal/repository"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

func ptr[T any](v T) *T {
	return &v
}

func main() {
	profilePath := flag.String("profile", "", "path to C2 profile YAML")
	disableRMQ := flag.Bool("no-rmq", false, "disable RabbitMQ messaging")
	flag.Parse()

	if *profilePath == "" {
		profilePath = ptr(os.Getenv("C2_PROFILE"))
	}

	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger, err := logging.New(cfg.LogLevel, cfg.Env)
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logger.Sync()

	var profile *config.MaleableProfile
	if *profilePath != "" {
		profile, err = config.LoadProfile(*profilePath)
		if err != nil {
			logger.Fatal("failed to load profile", zap.Error(err))
		}
		logger.Info("loaded C2 profile", zap.String("name", profile.Name))
	} else {
		profile = config.DefaultProfile()
		logger.Info("using default C2 profile")
	}

	logger.Info("starting maliketh C2 server",
		zap.String("env", cfg.Env),
		zap.String("log_level", cfg.LogLevel),
		zap.Any("profile", profile))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to database
	db, err := database.ConnectWithRetry(ctx, cfg.Postgres.DSN(), logger, 5, 5*time.Second)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	repos := repository.New(db.Pool)

	// Initialize messaging
	var publisher messaging.Publisher
	if *disableRMQ {
		logger.Warn("RabbitMQ disabled")
		publisher = &messaging.NoopPublisher{}
	} else {
		publisher, err = messaging.NewRabbitMQPublisherWithRetry(cfg.RabbitMQ.URL(), logger, 5, 5*time.Second)
		if err != nil {
			logger.Warn("failed to connect to RabbitMQ, continuing without messaging", zap.Error(err))
			publisher = &messaging.NoopPublisher{}
		} else {
			defer publisher.Close()
		}
	}

	// Create server
	srv := server.New(cfg.C2.Addr(), logger)

	// Root redirect
	srv.Router().Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, profile.ServerConfig.RedirectURL, http.StatusFound)
	})

	// Health check
	srv.Router().Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Health(r.Context()); err != nil {
			server.Error(w, http.StatusServiceUnavailable, "database unhealthy")
			return
		}
		server.Success(w, "healthy")
	})

	// Register C2 handlers
	c2Handler := c2handler.NewHandler(repos, profile, publisher, logger)
	c2Handler.RegisterRoutes(srv.Router())

	logger.Info("C2 routes registered",
		zap.String("base_path", profile.Routes.BasePath),
		zap.String("register", profile.Routes.Register.Path),
		zap.String("checkin", profile.Routes.Checkin.Path),
		zap.String("results", profile.Routes.TaskResults.Path))

	// Start server
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// Wait for interrupt
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
