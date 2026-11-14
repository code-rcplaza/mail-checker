-- cancellations: las líneas marcadas "por cancelar", con la fecha de la primera marca.

-- +goose Up
CREATE TABLE cancellations (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    line_key  TEXT    NOT NULL,
    marked_at INTEGER NOT NULL,
    UNIQUE (user_id, line_key)
);

-- +goose Down
DROP TABLE cancellations;
