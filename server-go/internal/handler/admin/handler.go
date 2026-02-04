package admin

import (
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/builder"
	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/messaging"
	"github.com/cbrnrd/maliketh/server-go/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	repos     *repository.Repositories
	cfg       *config.Config
	publisher messaging.Publisher
	logger    *zap.Logger
	startTime time.Time
	builder   *builder.Builder

	// Server keys for operator authentication
	serverPrivateKey string
}

type HandlerOption func(*Handler)

// WithBuilder sets the implant builder for the handler.
func WithBuilder(b *builder.Builder) HandlerOption {
	return func(h *Handler) {
		h.builder = b
	}
}

func NewHandler(
	repos *repository.Repositories,
	cfg *config.Config,
	publisher messaging.Publisher,
	logger *zap.Logger,
	serverPrivateKey string,
	opts ...HandlerOption,
) *Handler {
	h := &Handler{
		repos:            repos,
		cfg:              cfg,
		publisher:        publisher,
		logger:           logger,
		startTime:        time.Now(),
		serverPrivateKey: serverPrivateKey,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/op", func(r chi.Router) {
		// Authentication endpoints (no auth required)
		r.Get("/auth/token/request", h.handleRequestToken())
		r.Get("/auth/token/status", h.handleTokenStatus())

		// Authenticated endpoints
		r.Group(func(r chi.Router) {
			r.Use(h.operatorAuthMiddleware)

			// Auth management
			r.Get("/auth/token/revoke", h.handleRevokeToken())

			// Stats
			r.Get("/stats", h.handleStats())

			// Tasks
			r.Get("/tasks/list", h.handleListTasks())
			r.Post("/tasks/add", h.handleAddTask())
			r.Get("/tasks/results/{task_id}", h.handleGetTaskResult())
			r.Delete("/tasks/delete/{task_id}", h.handleDeleteTask())

			// Implants
			r.Get("/implant/list", h.handleListImplants())
			r.Get("/implant/config/{implant_id}", h.handleGetImplantConfig())
			r.Post("/implant/config/{implant_id}", h.handleUpdateImplantConfig())
			r.Delete("/implant/kill/{implant_id}", h.handleKillImplant())

			// Aliases
			r.Post("/implant/{implant_id}/alias/create", h.handleCreateAlias())
			r.Delete("/implant/{implant_id}/alias/delete/{alias}", h.handleDeleteAlias())
			r.Get("/implant/{implant_id}/alias/list", h.handleListAliases())
			r.Get("/implant/alias/resolve/{alias}", h.handleResolveAlias())

			// Builder
			r.Post("/implant/build", h.handleBuildImplant())
			r.Get("/implant/build/info", h.handleBuilderInfo())

			// Admin
			r.Post("/admin/revoke_access", h.handleAdminRevokeOperator())
		})
	})
}
