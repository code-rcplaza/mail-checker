-- manual_subscriptions: las que el usuario agregó a mano.

-- +goose Up
CREATE TABLE manual_subscriptions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    public_id   TEXT    NOT NULL UNIQUE,
    name        TEXT    NOT NULL,
    amount      INTEGER NOT NULL CHECK (amount > 0),
    currency    TEXT    NOT NULL CHECK (currency IN ('CLP', 'USD')),
    billing_day INTEGER NOT NULL CHECK (billing_day BETWEEN 1 AND 31)
);

-- +goose Down
DROP TABLE manual_subscriptions;
