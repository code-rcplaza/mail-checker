-- name: CreateUser :one
INSERT INTO users (public_id, name, email, accepted_terms_at, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: FindUserByID :one
SELECT *
FROM users
WHERE id = ?;

-- name: SaveLoginCode :exec
INSERT INTO login_codes (user_id, code_hash, expires_at, attempts)
VALUES (?, ?, ?, ?)
ON CONFLICT (user_id) DO UPDATE SET
    code_hash  = excluded.code_hash,
    expires_at = excluded.expires_at,
    attempts   = excluded.attempts;

-- name: FindLoginCode :one
SELECT * FROM login_codes WHERE user_id = ?;

-- name: AddFailedAttempt :one
UPDATE login_codes SET attempts = attempts + 1 WHERE user_id = ? RETURNING attempts;

-- name: DeleteLoginCode :exec
DELETE FROM login_codes WHERE user_id = ?;

-- name: RecordCodeSend :exec
INSERT INTO code_sends (email, sent_at) VALUES (?, ?);

-- name: ListCodeSendsSince :many
SELECT sent_at FROM code_sends
WHERE email = ? AND sent_at > ?
ORDER BY sent_at;

-- name: CreateSession :exec
INSERT INTO sessions (user_id, token_hash, expires_at) VALUES (?, ?, ?);

-- name: FindSessionByTokenHash :one
SELECT * FROM sessions WHERE token_hash = ?;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token_hash = ?;

-- name: DeleteExpiredLoginCodes :exec
DELETE FROM login_codes WHERE expires_at <= ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at <= ?;

-- name: DeleteOldCodeSends :exec
DELETE FROM code_sends WHERE sent_at <= ?;
