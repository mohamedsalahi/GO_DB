-- name: CreateTask :one
INSERT INTO tasks (
    user_id, title, description, status, due_date
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, title, description, status, due_date, created_at, updated_at;

-- name: GetTaskByID :one
SELECT id, user_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE id = $1 LIMIT 1;

-- name: ListTasksByUserID :many
SELECT id, user_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateTask :one
UPDATE tasks
SET 
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    status = COALESCE($4, status),
    due_date = COALESCE($5, due_date),
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, title, description, status, due_date, created_at, updated_at;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;

-- name: ListAllTasks :many
SELECT id, user_id, title, description, status, due_date, created_at, updated_at
FROM tasks
ORDER BY created_at DESC;
