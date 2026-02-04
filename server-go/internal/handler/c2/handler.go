package c2

import (
	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/messaging"
	"github.com/cbrnrd/maliketh/server-go/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	repos     *repository.Repositories
	profile   *config.MaleableProfile
	publisher messaging.Publisher
	logger    *zap.Logger
}

func NewHandler(repos *repository.Repositories, profile *config.MaleableProfile, publisher messaging.Publisher, logger *zap.Logger) *Handler {
	return &Handler{
		repos:     repos,
		profile:   profile,
		publisher: publisher,
		logger:    logger,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route(h.profile.Routes.BasePath, func(r chi.Router) {
		// Register endpoint - no auth required
		r.Method(
			h.profile.Routes.Register.Methods[0],
			h.profile.Routes.Register.Path,
			h.handleRegister(),
		)

		// Authenticated endpoints
		r.Group(func(r chi.Router) {
			r.Use(h.implantAuthMiddleware)

			r.Method(
				h.profile.Routes.Checkin.Methods[0],
				h.profile.Routes.Checkin.Path,
				h.handleCheckin(),
			)

			r.Method(
				h.profile.Routes.TaskResults.Methods[0],
				h.profile.Routes.TaskResults.Path,
				h.handleTaskResults(),
			)
		})
	})
}
