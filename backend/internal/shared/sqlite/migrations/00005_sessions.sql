-- sessions: una por dispositivo; se guarda el hash del token, nunca el token.

-- +goose Up
CREATE TABLE sessions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT    NOT NULL UNIQUE,
    expires_at INTEGER NOT NULL
);

CREATE INDEX sessions_user_id ON sessions (user_id);

-- +goose Down
DROP TABLE sessions;
