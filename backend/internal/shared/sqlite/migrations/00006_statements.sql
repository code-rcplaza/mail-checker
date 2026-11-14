-- statements: una cartola por usuario, banco y fecha de facturación. bad_lines es un
-- arreglo JSON de números de línea; nunca guarda texto de la cartola.

-- +goose Up
CREATE TABLE statements (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    public_id    TEXT    NOT NULL UNIQUE,
    user_id      INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    bank         TEXT    NOT NULL,
    filename     TEXT    NOT NULL,
    kind         TEXT    NOT NULL CHECK (kind IN ('pdf', 'xlsx', 'csv', 'text')),
    billing_date INTEGER NOT NULL,
    imported_at  INTEGER NOT NULL,
    bad_lines    TEXT    NOT NULL DEFAULT '[]',
    UNIQUE (user_id, bank, billing_date)
);

-- +goose Down
DROP TABLE statements;
