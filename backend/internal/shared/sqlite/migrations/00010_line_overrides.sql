-- line_overrides: el monto o el día de cobro que el usuario corrigió; NULL = no lo corrigió.

-- +goose Up
CREATE TABLE line_overrides (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    line_key    TEXT    NOT NULL,
    amount      INTEGER CHECK (amount > 0),
    billing_day INTEGER CHECK (billing_day BETWEEN 1 AND 31),
    UNIQUE (user_id, line_key)
);

-- +goose Down
DROP TABLE line_overrides;
