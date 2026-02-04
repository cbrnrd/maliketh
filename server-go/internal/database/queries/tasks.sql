-- name: GetTaskByID :one
SELECT * FROM tasks WHERE task_id = $1;

-- name: CreateTask :one
INSERT INTO tasks (
    task_id, operator_name, implant_id, opcode, args, status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOldestPendingTaskForImplant :one
SELECT * FROM tasks
WHERE implant_id = $1 AND status = 'CREATED'
ORDER BY created_at ASC
LIMIT 1;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = $2, read_at = $3
WHERE task_id = $1;

-- name: UpdateTaskResult :exec
UPDATE tasks
SET status = $2, output = $3, executed_at = NOW()
WHERE task_id = $1;

-- name: ListTasksByOperator :many
SELECT * FROM tasks
WHERE operator_name = $1
ORDER BY created_at DESC;

-- name: ListTasksByImplant :many
SELECT * FROM tasks
WHERE implant_id = $1
ORDER BY created_at DESC;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE task_id = $1;

-- name: CountActiveTasks :one
SELECT COUNT(*) FROM tasks
WHERE operator_name = $1 AND status IN ('CREATED', 'TASKED');

-- name: CountTotalTasks :one
SELECT COUNT(*) FROM tasks WHERE operator_name = $1;
