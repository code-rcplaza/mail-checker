-- subscription_decisions: lo que el usuario decidió sobre una línea posible.
-- La cadencia va solo con una confirmada.

-- +goose Up
CREATE TABLE subscription_decisions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_key  TEXT    NOT NULL,
    decision   TEXT    NOT NULL CHECK (decision IN ('confirmed', 'dismissed')),
    cadence    TEXT    CHECK (cadence IN ('monthly', 'yearly')),
    decided_at INTEGER NOT NULL,
    UNIQUE (user_id, group_key),
    CHECK ((decision = 'confirmed') = (cadence IS NOT NULL))
);

-- +goose Down
DROP TABLE subscription_decisions;
