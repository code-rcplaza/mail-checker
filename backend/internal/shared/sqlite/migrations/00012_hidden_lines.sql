-- hidden_lines: las líneas que el usuario dejó de seguir.

-- +goose Up
CREATE TABLE hidden_lines (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    line_key  TEXT    NOT NULL,
    hidden_at INTEGER NOT NULL,
    UNIQUE (user_id, line_key)
);

-- +goose Down
DROP TABLE hidden_lines;
