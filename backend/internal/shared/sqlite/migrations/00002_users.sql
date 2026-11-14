-- users: dueños de todos los datos. El correo se guarda normalizado y además es
-- único sin distinguir mayúsculas, por si algo escribe sin normalizar.

-- +goose Up
CREATE TABLE users (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    public_id         TEXT    NOT NULL UNIQUE,
    name              TEXT    NOT NULL,
    email             TEXT    NOT NULL UNIQUE COLLATE NOCASE,
    accepted_terms_at INTEGER NOT NULL,
    created_at        INTEGER NOT NULL
);

-- +goose Down
DROP TABLE users;
