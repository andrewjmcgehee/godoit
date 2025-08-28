-- +goose Up
-- Database schema for godoit todo application
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    priority TEXT CHECK(priority IN ('P0', 'P1', 'P2')) DEFAULT 'P2',
    completed BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster queries on completed status
CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);

-- Index for priority filtering
CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);

-- +goose Down
DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_completed;
DROP TABLE IF EXISTS todos;
