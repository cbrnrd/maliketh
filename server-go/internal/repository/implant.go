package repository

import (
	"context"
	"fmt"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/database/repo"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImplantRepository struct {
	q    *repo.Queries
	pool *pgxpool.Pool
}

func NewImplantRepository(pool *pgxpool.Pool) *ImplantRepository {
	return &ImplantRepository{
		q:    repo.New(pool),
		pool: pool,
	}
}

func (r *ImplantRepository) GetByID(ctx context.Context, implantID string) (*model.Implant, error) {
	imp, err := r.q.GetImplantByID(ctx, implantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get implant: %w", err)
	}
	return mapImplant(imp), nil
}

func (r *ImplantRepository) GetByPublicKey(ctx context.Context, publicKey string) (*model.Implant, error) {
	imp, err := r.q.GetImplantByPublicKey(ctx, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get implant by public key: %w", err)
	}
	return mapImplant(imp), nil
}

// Resolve resolves an implant identifier (ID or alias) to an implant.
func (r *ImplantRepository) Resolve(ctx context.Context, identifier string) (*model.Implant, error) {
	imp, err := r.q.ResolveImplantIdentifier(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve implant: %w", err)
	}
	return mapImplant(imp), nil
}

func (r *ImplantRepository) Create(ctx context.Context, imp *model.Implant) (*model.Implant, error) {
	// Generate a random ID if not set
	if imp.ImplantID == "" {
		imp.ImplantID = crypto.RandomID(8)
	}

	created, err := r.q.CreateImplant(ctx, repo.CreateImplantParams{
		ImplantID: imp.ImplantID,
		Hostname:  &imp.Hostname,
		Ip:        &imp.IP,
		Os:        &imp.OS,
		Arch:      &imp.Arch,
		Username:  &imp.Username,
		ServerSk:  imp.ServerSK,
		ImplantPk: imp.ImplantPK,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create implant: %w", err)
	}
	return mapImplant(created), nil
}

func (r *ImplantRepository) UpdateCheckin(ctx context.Context, implantID, ip string) error {
	return r.q.UpdateImplantCheckin(ctx, repo.UpdateImplantCheckinParams{
		ImplantID: implantID,
		Ip:        &ip,
	})
}

func (r *ImplantRepository) Delete(ctx context.Context, implantID string) error {
	return r.q.DeleteImplant(ctx, implantID)
}

func (r *ImplantRepository) List(ctx context.Context) ([]*model.Implant, error) {
	imps, err := r.q.ListImplants(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list implants: %w", err)
	}

	result := make([]*model.Implant, len(imps))
	for i, imp := range imps {
		result[i] = mapImplant(imp)
	}
	return result, nil
}

func (r *ImplantRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountImplants(ctx)
}

// Alias operations

func (r *ImplantRepository) CreateAlias(ctx context.Context, implantID, alias string) (*model.ImplantAlias, error) {
	created, err := r.q.CreateImplantAlias(ctx, repo.CreateImplantAliasParams{
		ImplantID: implantID,
		Alias:     alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create alias: %w", err)
	}
	return &model.ImplantAlias{
		ID:        int64(created.ID),
		ImplantID: created.ImplantID,
		Alias:     created.Alias,
	}, nil
}

func (r *ImplantRepository) GetAliasByName(ctx context.Context, alias string) (*model.ImplantAlias, error) {
	a, err := r.q.GetAliasByName(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}
	return &model.ImplantAlias{
		ID:        int64(a.ID),
		ImplantID: a.ImplantID,
		Alias:     a.Alias,
	}, nil
}

func (r *ImplantRepository) GetAliasesByImplantID(ctx context.Context, implantID string) ([]*model.ImplantAlias, error) {
	aliases, err := r.q.GetAliasesByImplantID(ctx, implantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get aliases: %w", err)
	}

	result := make([]*model.ImplantAlias, len(aliases))
	for i, a := range aliases {
		result[i] = &model.ImplantAlias{
			ID:        int64(a.ID),
			ImplantID: a.ImplantID,
			Alias:     a.Alias,
		}
	}
	return result, nil
}

func (r *ImplantRepository) DeleteAlias(ctx context.Context, alias string) error {
	return r.q.DeleteAlias(ctx, alias)
}

func mapImplant(imp repo.Implant) *model.Implant {
	result := &model.Implant{
		ID:        int64(imp.ID),
		ImplantID: imp.ImplantID,
		ServerSK:  imp.ServerSk,
		ImplantPK: imp.ImplantPk,
		CreatedAt: imp.CreatedAt.Time,
		LastSeen:  imp.LastSeen.Time,
	}

	if imp.Hostname != nil {
		result.Hostname = *imp.Hostname
	}
	if imp.Ip != nil {
		result.IP = *imp.Ip
	}
	if imp.Os != nil {
		result.OS = *imp.Os
	}
	if imp.Arch != nil {
		result.Arch = *imp.Arch
	}
	if imp.Username != nil {
		result.Username = *imp.Username
	}

	return result
}
