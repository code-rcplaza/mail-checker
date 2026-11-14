-- charges: los cargos de todas las fuentes (compartida). account_id llega con la
-- tabla accounts, en el post-MVP.

-- +goose Up
CREATE TABLE charges (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    bank         TEXT    NOT NULL,
    source       TEXT    NOT NULL CHECK (source IN ('statement', 'email')),
    external_id  TEXT    NOT NULL,
    statement_id INTEGER REFERENCES statements (id) ON DELETE CASCADE,
    occurred_at  INTEGER NOT NULL,
    amount       INTEGER NOT NULL CHECK (amount >= 0),
    currency     TEXT    NOT NULL CHECK (currency IN ('CLP', 'USD')),
    merchant     TEXT    NOT NULL,
    card         TEXT    NOT NULL DEFAULT '',
    card_last4   TEXT    NOT NULL DEFAULT ''
        CHECK (card_last4 = '' OR (length(card_last4) = 4 AND card_last4 NOT GLOB '*[^0-9]*')),
    UNIQUE (user_id, source, external_id)
);

CREATE INDEX charges_user_occurred_at ON charges (user_id, occurred_at);
CREATE INDEX charges_statement_id ON charges (statement_id);

-- +goose Down
DROP TABLE charges;
