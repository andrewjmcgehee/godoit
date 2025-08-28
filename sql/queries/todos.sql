-- name: CreateTodo :one
INSERT INTO todos (content, priority, created_at, updated_at)
VALUES (?, ?, ?, ?)
RETURNING id, content, priority, completed, created_at, updated_at;

-- name: GetActiveTodos :many
SELECT id, content, priority, completed, created_at, updated_at 
FROM todos 
WHERE completed = FALSE 
ORDER BY priority ASC, created_at DESC;

-- name: GetCompletedTodos :many
SELECT id, content, priority, completed, created_at, updated_at 
FROM todos 
WHERE completed = TRUE 
ORDER BY updated_at DESC;

-- name: UpdateTodoContent :exec
UPDATE todos 
SET content = ?, updated_at = ? 
WHERE id = ?;

-- name: UpdateTodoPriority :exec
UPDATE todos 
SET priority = ?, updated_at = ? 
WHERE id = ?;

-- name: ToggleTodoCompleted :exec
UPDATE todos 
SET completed = NOT completed, updated_at = ? 
WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;

-- name: CountActiveTodos :one
SELECT COUNT(*) FROM todos WHERE completed = FALSE;

-- name: CountCompletedTodos :one
SELECT COUNT(*) FROM todos WHERE completed = TRUE;
