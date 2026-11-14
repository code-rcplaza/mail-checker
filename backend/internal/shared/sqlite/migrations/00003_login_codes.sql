-- login_codes: un solo código vigente por usuario; se guarda su hash, nunca el código.

-- +goose Up
CREATE TABLE login_codes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    code_hash  TEXT    NOT NULL,
    expires_at INTEGER NOT NULL,
    attempts   INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0)
);

-- +goose Down
DROP TABLE login_codes;
