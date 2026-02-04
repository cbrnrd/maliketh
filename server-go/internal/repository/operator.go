package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/database/repo"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OperatorRepository struct {
	q    *repo.Queries
	pool *pgxpool.Pool
}

func NewOperatorRepository(pool *pgxpool.Pool) *OperatorRepository {
	return &OperatorRepository{
		q:    repo.New(pool),
		pool: pool,
	}
}

func (r *OperatorRepository) GetByUsername(ctx context.Context, username string) (*model.Operator, error) {
	op, err := r.q.GetOperatorByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator: %w", err)
	}
	return mapOperator(op), nil
}

func (r *OperatorRepository) GetByAuthToken(ctx context.Context, token string) (*model.Operator, error) {
	op, err := r.q.GetOperatorByAuthToken(ctx, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator by token: %w", err)
	}
	return mapOperator(op), nil
}

func (r *OperatorRepository) Create(ctx context.Context, op *model.Operator) (*model.Operator, error) {
	created, err := r.q.CreateOperator(ctx, repo.CreateOperatorParams{
		Username:    op.Username,
		PublicKey:   op.PublicKey,
		VerifyKey:   op.VerifyKey,
		LoginSecret: op.LoginSecret,
		RmqQueue:    &op.RMQQueue,
		Role:        string(op.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create operator: %w", err)
	}
	return mapOperator(created), nil
}

func (r *OperatorRepository) UpdateAuthToken(ctx context.Context, username, token string, expiry time.Time) error {
	return r.q.UpdateOperatorAuthToken(ctx, repo.UpdateOperatorAuthTokenParams{
		Username:        username,
		AuthToken:       &token,
		AuthTokenExpiry: pgtype.Timestamptz{Time: expiry, Valid: true},
	})
}

func (r *OperatorRepository) ClearAuthToken(ctx context.Context, username string) error {
	return r.q.UpdateOperatorAuthToken(ctx, repo.UpdateOperatorAuthTokenParams{
		Username:        username,
		AuthToken:       nil,
		AuthTokenExpiry: pgtype.Timestamptz{Valid: false},
	})
}

func (r *OperatorRepository) UpdateLastLogin(ctx context.Context, username string) error {
	return r.q.UpdateOperatorLastLogin(ctx, username)
}

func (r *OperatorRepository) Revoke(ctx context.Context, username string) error {
	return r.q.RevokeOperator(ctx, username)
}

func (r *OperatorRepository) List(ctx context.Context) ([]*model.Operator, error) {
	ops, err := r.q.ListOperators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list operators: %w", err)
	}

	result := make([]*model.Operator, len(ops))
	for i, op := range ops {
		result[i] = mapOperator(op)
	}
	return result, nil
}

func (r *OperatorRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountOperators(ctx)
}

// GenerateAuthToken creates a new auth token for an operator.
func (r *OperatorRepository) GenerateAuthToken(ctx context.Context, username string, duration time.Duration) (string, error) {
	token, err := crypto.RandomHex(64) // 128 hex characters
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	expiry := time.Now().Add(duration)
	if err := r.UpdateAuthToken(ctx, username, token, expiry); err != nil {
		return "", err
	}

	return token, nil
}

func mapOperator(op repo.Operator) *model.Operator {
	result := &model.Operator{
		ID:          int64(op.ID),
		Username:    op.Username,
		PublicKey:   op.PublicKey,
		VerifyKey:   op.VerifyKey,
		LoginSecret: op.LoginSecret,
		CreatedAt:   op.CreatedAt.Time,
		Role:        model.OperatorRole(op.Role),
		Revoked:     op.Revoked,
	}

	if op.AuthToken != nil {
		result.AuthToken = op.AuthToken
	}
	if op.AuthTokenExpiry.Valid {
		result.AuthTokenExpiry = &op.AuthTokenExpiry.Time
	}
	if op.LastLogin.Valid {
		result.LastLogin = &op.LastLogin.Time
	}
	if op.RmqQueue != nil {
		result.RMQQueue = *op.RmqQueue
	}

	return result
}
