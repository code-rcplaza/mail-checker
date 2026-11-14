-- exchange_rates: el tipo de cambio que mostró cada cartola, uno por usuario y fecha.
-- Lo escribe statements y lo lee subscriptions, cada una con su adaptador.

-- +goose Up
CREATE TABLE exchange_rates (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id          INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    pesos_per_dollar INTEGER NOT NULL CHECK (pesos_per_dollar > 0),
    date             INTEGER NOT NULL,
    UNIQUE (user_id, date)
);

-- +goose Down
DROP TABLE exchange_rates;
