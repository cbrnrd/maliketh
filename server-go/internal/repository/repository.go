package repository

import "github.com/jackc/pgx/v5/pgxpool"

// Repositories holds all repository instances.
type Repositories struct {
	Operator *OperatorRepository
	Implant  *ImplantRepository
	Config   *ConfigRepository
	Task     *TaskRepository
}

// New creates all repositories from a connection pool.
func New(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Operator: NewOperatorRepository(pool),
		Implant:  NewImplantRepository(pool),
		Config:   NewConfigRepository(pool),
		Task:     NewTaskRepository(pool),
	}
}
