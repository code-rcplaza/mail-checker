-- price_alert_decisions: historial de los montos nuevos que el usuario mantuvo en una
-- línea; cada monto una vez. El resumen muestra igual el último monto cobrado.

-- +goose Up
CREATE TABLE price_alert_decisions (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id  INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    line_key TEXT    NOT NULL,
    amount   INTEGER NOT NULL CHECK (amount > 0),
    UNIQUE (user_id, line_key, amount)
);

-- +goose Down
DROP TABLE price_alert_decisions;
