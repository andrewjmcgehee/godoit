CREATE TABLE todos (
    id INTEGER PRIMARY KEY NOT NULL,
    content TEXT NOT NULL,
    priority TEXT CHECK (priority IN ('P0', 'P1', 'P2')) DEFAULT 'P2' NOT NULL,
    completed BOOLEAN DEFAULT FALSE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_todos_completed ON todos (completed);
CREATE INDEX idx_todos_priority ON todos (priority);
