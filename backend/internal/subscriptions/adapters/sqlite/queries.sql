-- name: ListChargesSince :many
SELECT * FROM charges
WHERE user_id = ? AND occurred_at >= ?
ORDER BY occurred_at, id;

-- name: ListDecisions :many
SELECT group_key, decision, cadence FROM subscription_decisions WHERE user_id = ?;

-- name: SaveDecision :exec
INSERT INTO subscription_decisions (user_id, group_key, decision, cadence, decided_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (user_id, group_key) DO UPDATE SET
    decision   = excluded.decision,
    cadence    = excluded.cadence,
    decided_at = excluded.decided_at;

-- name: DeleteDecision :exec
DELETE FROM subscription_decisions WHERE user_id = ? AND group_key = ?;

-- name: ListAliases :many
SELECT line_key, alias FROM line_aliases WHERE user_id = ?;

-- name: SaveAlias :exec
INSERT INTO line_aliases (user_id, line_key, alias) VALUES (?, ?, ?)
ON CONFLICT (user_id, line_key) DO UPDATE SET alias = excluded.alias;

-- name: DeleteAlias :exec
DELETE FROM line_aliases WHERE user_id = ? AND line_key = ?;

-- name: ListOverrides :many
SELECT line_key, amount, billing_day FROM line_overrides WHERE user_id = ?;

-- name: SaveOverride :exec
INSERT INTO line_overrides (user_id, line_key, amount, billing_day) VALUES (?, ?, ?, ?)
ON CONFLICT (user_id, line_key) DO UPDATE SET
    amount      = excluded.amount,
    billing_day = excluded.billing_day;

-- name: DeleteOverride :exec
DELETE FROM line_overrides WHERE user_id = ? AND line_key = ?;

-- name: ListManuals :many
SELECT public_id, name, amount, currency, billing_day FROM manual_subscriptions
WHERE user_id = ?
ORDER BY id;

-- name: CreateManual :exec
INSERT INTO manual_subscriptions (user_id, public_id, name, amount, currency, billing_day)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateManual :execrows
UPDATE manual_subscriptions SET name = ?, amount = ?, currency = ?, billing_day = ?
WHERE user_id = ? AND public_id = ?;

-- name: DeleteManual :execrows
DELETE FROM manual_subscriptions WHERE user_id = ? AND public_id = ?;

-- name: ListHiddenLines :many
SELECT line_key FROM hidden_lines WHERE user_id = ? ORDER BY id;

-- name: HideLine :exec
INSERT INTO hidden_lines (user_id, line_key, hidden_at) VALUES (?, ?, ?)
ON CONFLICT (user_id, line_key) DO NOTHING;

-- name: ListCancellations :many
SELECT line_key, marked_at FROM cancellations WHERE user_id = ?;

-- name: MarkForCancellation :exec
INSERT INTO cancellations (user_id, line_key, marked_at) VALUES (?, ?, ?)
ON CONFLICT (user_id, line_key) DO NOTHING;

-- name: ListKeptRises :many
SELECT line_key, amount FROM price_alert_decisions WHERE user_id = ? ORDER BY id;

-- name: KeepRise :exec
INSERT INTO price_alert_decisions (user_id, line_key, amount) VALUES (?, ?, ?)
ON CONFLICT (user_id, line_key, amount) DO NOTHING;

-- name: ListExchangeRates :many
SELECT pesos_per_dollar, date FROM exchange_rates WHERE user_id = ? ORDER BY date;
