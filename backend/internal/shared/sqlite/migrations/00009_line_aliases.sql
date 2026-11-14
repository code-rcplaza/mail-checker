-- line_aliases: el nombre que el usuario le puso a una línea (hasta 40 caracteres).

-- +goose Up
CREATE TABLE line_aliases (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id  INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    line_key TEXT    NOT NULL,
    alias    TEXT    NOT NULL CHECK (length(alias) <= 40),
    UNIQUE (user_id, line_key)
);

-- +goose Down
DROP TABLE line_aliases;
