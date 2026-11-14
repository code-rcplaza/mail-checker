-- name: CreateStatement :one
INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at, bad_lines)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: CreateCharge :exec
INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant, card, card_last4)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListStatements :many
SELECT s.id, s.public_id, s.user_id, s.bank, s.filename, s.kind, s.billing_date, s.imported_at, s.bad_lines,
       (SELECT COUNT(*) FROM charges c WHERE c.statement_id = s.id) AS charge_count
FROM statements s
WHERE s.user_id = ?
ORDER BY s.billing_date DESC;

-- name: DeleteStatement :execrows
DELETE FROM statements WHERE user_id = ? AND public_id = ?;

-- name: RecordExchangeRate :exec
INSERT INTO exchange_rates (user_id, pesos_per_dollar, date) VALUES (?, ?, ?)
ON CONFLICT (user_id, date) DO NOTHING;
