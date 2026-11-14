-- code_sends: cada petición de código por correo, exista o no la cuenta, para el
-- límite de envíos por hora. No referencia a users a propósito.

-- +goose Up
CREATE TABLE code_sends (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    email   TEXT    NOT NULL COLLATE NOCASE,
    sent_at INTEGER NOT NULL
);

CREATE INDEX code_sends_email_sent_at ON code_sends (email, sent_at);

-- +goose Down
DROP TABLE code_sends;
