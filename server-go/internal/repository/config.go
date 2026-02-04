package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/database/repo"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigRepository struct {
	q    *repo.Queries
	pool *pgxpool.Pool
}

func NewConfigRepository(pool *pgxpool.Pool) *ConfigRepository {
	return &ConfigRepository{
		q:    repo.New(pool),
		pool: pool,
	}
}

func (r *ConfigRepository) GetByImplantID(ctx context.Context, implantID string) (*model.ImplantConfig, error) {
	cfg, err := r.q.GetImplantConfig(ctx, implantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get implant config: %w", err)
	}
	return mapConfig(cfg), nil
}

func (r *ConfigRepository) Create(ctx context.Context, cfg *model.ImplantConfig) (*model.ImplantConfig, error) {
	var killDate pgtype.Timestamptz
	if cfg.KillDate != nil {
		killDate = pgtype.Timestamptz{Time: *cfg.KillDate, Valid: true}
	}

	created, err := r.q.CreateImplantConfig(ctx, repo.CreateImplantConfigParams{
		ImplantID:             cfg.ImplantID,
		Cookie:                cfg.Cookie,
		KillDate:              killDate,
		UserAgent:             &cfg.UserAgent,
		AutoSelfDestruct:      cfg.AutoSelfDestruct,
		SleepTime:             int32(cfg.SleepTime),
		Jitter:                float32(cfg.Jitter),
		MaxRetries:            int32(cfg.MaxRetries),
		RetryWait:             int32(cfg.RetryWait),
		RetryJitter:           float32(cfg.RetryJitter),
		EncKey:                cfg.EncKey,
		TailoringHashFunction: &cfg.TailoringHashFunc,
		TailoringHashRounds:   ptrInt32(int32(cfg.TailoringHashRounds)),
		TailoringHashes:       cfg.TailoringHashes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create implant config: %w", err)
	}
	return mapConfig(created), nil
}

// CreateFromProfile creates a config from a malleable profile.
func (r *ConfigRepository) CreateFromProfile(ctx context.Context, implantID, encKey string, profile *config.MaleableProfile) (*model.ImplantConfig, error) {
	cfg := &model.ImplantConfig{
		ImplantID:           implantID,
		Cookie:              profile.Globals.ImplantIDCookie,
		UserAgent:           profile.ImplantConfig.UserAgent,
		AutoSelfDestruct:    profile.ImplantConfig.AutoSelfDestruct,
		SleepTime:           profile.ImplantConfig.Sleep,
		Jitter:              profile.ImplantConfig.Jitter,
		MaxRetries:          profile.ImplantConfig.MaxRetries,
		RetryWait:           profile.ImplantConfig.RetryWait,
		RetryJitter:         profile.ImplantConfig.RetryJitter,
		EncKey:              encKey,
		TailoringHashFunc:   profile.ImplantConfig.TailoringHashFunc,
		TailoringHashRounds: profile.ImplantConfig.TailoringHashRounds,
		TailoringHashes:     profile.ImplantConfig.TailoringHashes,
	}

	if profile.ImplantConfig.KillDate != "" {
		if t, err := time.Parse(time.RFC3339, profile.ImplantConfig.KillDate); err == nil {
			cfg.KillDate = &t
		}
	}

	return r.Create(ctx, cfg)
}

// CreateMinimal creates a minimal config with default values.
func (r *ConfigRepository) CreateMinimal(ctx context.Context, implantID, cookie, encKey string) (*model.ImplantConfig, error) {
	cfg := &model.ImplantConfig{
		ImplantID:           implantID,
		Cookie:              cookie,
		SleepTime:           60,
		Jitter:              0.1,
		MaxRetries:          -1,
		EncKey:              encKey,
		TailoringHashFunc:   "sha256",
		TailoringHashRounds: 100,
		TailoringHashes:     []string{},
	}
	return r.Create(ctx, cfg)
}

func (r *ConfigRepository) Update(ctx context.Context, implantID string, updates map[string]any) error {
	params := repo.UpdateImplantConfigParams{
		ImplantID: implantID,
	}

	// Map the updates to the params
	if v, ok := updates["cookie"].(string); ok {
		params.Cookie = &v
	}
	if v, ok := updates["user_agent"].(string); ok {
		params.UserAgent = &v
	}
	if v, ok := updates["auto_self_destruct"].(bool); ok {
		params.AutoSelfDestruct = &v
	}
	if v, ok := updates["sleep_time"].(int); ok {
		val := int32(v)
		params.SleepTime = &val
	}
	if v, ok := updates["jitter"].(float64); ok {
		val := float32(v)
		params.Jitter = &val
	}
	if v, ok := updates["max_retries"].(int); ok {
		val := int32(v)
		params.MaxRetries = &val
	}
	if v, ok := updates["retry_wait"].(int); ok {
		val := int32(v)
		params.RetryWait = &val
	}
	if v, ok := updates["retry_jitter"].(float64); ok {
		val := float32(v)
		params.RetryJitter = &val
	}
	if v, ok := updates["tailoring_hash_function"].(string); ok {
		params.TailoringHashFunction = &v
	}
	if v, ok := updates["tailoring_hash_rounds"].(int); ok {
		val := int32(v)
		params.TailoringHashRounds = &val
	}
	if v, ok := updates["tailoring_hashes"].([]string); ok {
		params.TailoringHashes = v
	}

	return r.q.UpdateImplantConfig(ctx, params)
}

func (r *ConfigRepository) Delete(ctx context.Context, implantID string) error {
	return r.q.DeleteImplantConfig(ctx, implantID)
}

func mapConfig(cfg repo.ImplantConfig) *model.ImplantConfig {
	result := &model.ImplantConfig{
		ID:               int64(cfg.ID),
		ImplantID:        cfg.ImplantID,
		Cookie:           cfg.Cookie,
		AutoSelfDestruct: cfg.AutoSelfDestruct,
		SleepTime:        int(cfg.SleepTime),
		Jitter:           float64(cfg.Jitter),
		MaxRetries:       int(cfg.MaxRetries),
		RetryWait:        int(cfg.RetryWait),
		RetryJitter:      float64(cfg.RetryJitter),
		EncKey:           cfg.EncKey,
	}

	if cfg.KillDate.Valid {
		result.KillDate = &cfg.KillDate.Time
	}
	if cfg.UserAgent != nil {
		result.UserAgent = *cfg.UserAgent
	}
	if cfg.TailoringHashFunction != nil {
		result.TailoringHashFunc = *cfg.TailoringHashFunction
	}
	if cfg.TailoringHashRounds != nil {
		result.TailoringHashRounds = int(*cfg.TailoringHashRounds)
	}
	if cfg.TailoringHashes != nil {
		result.TailoringHashes = cfg.TailoringHashes
	} else {
		result.TailoringHashes = []string{}
	}

	return result
}

func ptrInt32(v int32) *int32 {
	return &v
}
