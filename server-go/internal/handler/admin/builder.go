package admin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/builder"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

// handleBuildImplant builds an implant with the provided options.
func (h *Handler) handleBuildImplant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		// Check if builder is available
		if h.builder == nil {
			server.InternalError(w, "Builder not configured")
			return
		}

		// Parse build options from request
		var optsMap map[string]any
		if err := json.NewDecoder(r.Body).Decode(&optsMap); err != nil {
			server.BadRequest(w, "Invalid request body")
			return
		}

		// Convert to BuildOptions
		opts, err := builder.FromMap(optsMap)
		if err != nil {
			server.BadRequest(w, err.Error())
			return
		}

		h.logger.Info("building implant",
			zap.String("operator", operator.Username),
			zap.String("implant_type", string(opts.ImplantType)),
			zap.String("target", string(opts.TargetOS)+"/"+string(opts.TargetArch)))

		// Build the implant
		result, err := h.builder.Build(r.Context(), operator.Username, opts)
		if err != nil {
			h.logger.Error("failed to build implant",
				zap.String("operator", operator.Username),
				zap.Error(err))
			server.InternalError(w, "Error building implant")
			return
		}

		h.logger.Info("implant build complete",
			zap.String("operator", operator.Username),
			zap.Int("artifact_size", len(result.Artifact)),
			zap.Duration("build_time", result.BuildTime))

		// Return base64-encoded implant
		server.SuccessData(w, map[string]any{
			"status":     true,
			"implant":    base64.StdEncoding.EncodeToString(result.Artifact),
			"filename":   result.Filename,
			"build_time": result.BuildTime.String(),
		})
	}
}

// handleBuilderInfo returns information about available build targets.
func (h *Handler) handleBuilderInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		if h.builder == nil {
			server.SuccessData(w, map[string]any{
				"status":    true,
				"available": false,
				"targets":   map[string][]string{},
			})
			return
		}

		available := h.builder.Available()
		targets := make(map[string][]string)
		for implantType, archs := range available {
			targets[string(implantType)] = archs
		}

		server.SuccessData(w, map[string]any{
			"status":    true,
			"available": len(targets) > 0,
			"targets":   targets,
			"defaults": map[string]any{
				"implant_type":          string(builder.DefaultBuildOptions.ImplantType),
				"target_os":             string(builder.DefaultBuildOptions.TargetOS),
				"target_arch":           string(builder.DefaultBuildOptions.TargetArch),
				"initial_sleep_seconds": builder.DefaultBuildOptions.InitialSleepSeconds,
				"schtask_persist":       builder.DefaultBuildOptions.SchtaskPersist,
				"use_antidebug":         builder.DefaultBuildOptions.UseAntidebug,
				"use_antivm":            builder.DefaultBuildOptions.UseAntiVM,
				"scheduled_task_name":   builder.DefaultBuildOptions.ScheduledTaskName,
				"register_max_retries":  builder.DefaultBuildOptions.RegisterMaxRetries,
			},
		})
	}
}
