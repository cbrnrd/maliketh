package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/database/repo"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	q    *repo.Queries
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		q:    repo.New(pool),
		pool: pool,
	}
}

func (r *TaskRepository) GetByID(ctx context.Context, taskID string) (*model.Task, error) {
	task, err := r.q.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return mapTask(task), nil
}

func (r *TaskRepository) Create(ctx context.Context, operatorName, implantID string, opcode model.Opcode, args any) (*model.Task, error) {
	taskID := crypto.RandomID(16)

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal args: %w", err)
	}

	created, err := r.q.CreateTask(ctx, repo.CreateTaskParams{
		TaskID:       taskID,
		OperatorName: operatorName,
		ImplantID:    implantID,
		Opcode:       int32(opcode),
		Args:         argsJSON,
		Status:       string(model.TaskStatusCreated),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return mapTask(created), nil
}

func (r *TaskRepository) GetOldestPendingForImplant(ctx context.Context, implantID string) (*model.Task, error) {
	task, err := r.q.GetOldestPendingTaskForImplant(ctx, implantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending task: %w", err)
	}
	return mapTask(task), nil
}

func (r *TaskRepository) MarkAsTasked(ctx context.Context, taskID string) error {
	now := time.Now()
	return r.q.UpdateTaskStatus(ctx, repo.UpdateTaskStatusParams{
		TaskID: taskID,
		Status: string(model.TaskStatusTasked),
		ReadAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
}

func (r *TaskRepository) SetResult(ctx context.Context, taskID string, status model.TaskStatus, output string) error {
	return r.q.UpdateTaskResult(ctx, repo.UpdateTaskResultParams{
		TaskID: taskID,
		Status: string(status),
		Output: &output,
	})
}

func (r *TaskRepository) Delete(ctx context.Context, taskID string) error {
	return r.q.DeleteTask(ctx, taskID)
}

func (r *TaskRepository) ListByOperator(ctx context.Context, operatorName string) ([]*model.Task, error) {
	tasks, err := r.q.ListTasksByOperator(ctx, operatorName)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	result := make([]*model.Task, len(tasks))
	for i, task := range tasks {
		result[i] = mapTask(task)
	}
	return result, nil
}

func (r *TaskRepository) ListByImplant(ctx context.Context, implantID string) ([]*model.Task, error) {
	tasks, err := r.q.ListTasksByImplant(ctx, implantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	result := make([]*model.Task, len(tasks))
	for i, task := range tasks {
		result[i] = mapTask(task)
	}
	return result, nil
}

func (r *TaskRepository) CountActive(ctx context.Context, operatorName string) (int64, error) {
	return r.q.CountActiveTasks(ctx, operatorName)
}

func (r *TaskRepository) CountTotal(ctx context.Context, operatorName string) (int64, error) {
	return r.q.CountTotalTasks(ctx, operatorName)
}

func mapTask(task repo.Task) *model.Task {
	result := &model.Task{
		ID:           int64(task.ID),
		TaskID:       task.TaskID,
		OperatorName: task.OperatorName,
		ImplantID:    task.ImplantID,
		Opcode:       model.Opcode(task.Opcode),
		Args:         task.Args,
		Status:       model.TaskStatus(task.Status),
		CreatedAt:    task.CreatedAt.Time,
	}

	if task.Output != nil {
		result.Output = task.Output
	}
	if task.ReadAt.Valid {
		result.ReadAt = &task.ReadAt.Time
	}
	if task.ExecutedAt.Valid {
		result.ExecutedAt = &task.ExecutedAt.Time
	}

	return result
}
