package sqlite

import (
	"context"
	"database/sql"
	"testing"
)

// insertUser adds the user every access row belongs to and returns its id.
func insertUser(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	result, err := db.ExecContext(context.Background(), `INSERT INTO users (public_id, name, email, accepted_terms_at, created_at) VALUES ('u1', 'Ana', 'ana@example.com', 1, 1)`)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("user id: %v", err)
	}
	return id
}

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}

func mustExec(t *testing.T, db *sql.DB, statement string, args ...any) {
	t.Helper()
	if _, err := db.ExecContext(context.Background(), statement, args...); err != nil {
		t.Fatalf("exec %q: %v", statement, err)
	}
}

func TestLoginCodesTable_KeepsOneCodePerUser(t *testing.T) {
	db := openMigratedDB(t)
	userID := insertUser(t, db)
	insert := `INSERT INTO login_codes (user_id, code_hash, expires_at, attempts) VALUES (?, ?, 1, 0)`
	mustExec(t, db, insert, userID, "hash-1")

	if _, err := db.Exec(insert, userID, "hash-2"); err == nil {
		t.Errorf("second code for the same user error = nil, want a UNIQUE violation")
	}
}

func TestAccessTables_RowsGoWithTheirUser(t *testing.T) {
	db := openMigratedDB(t)
	userID := insertUser(t, db)
	mustExec(t, db, `INSERT INTO login_codes (user_id, code_hash, expires_at, attempts) VALUES (?, 'hash', 1, 0)`, userID)
	mustExec(t, db, `INSERT INTO sessions (user_id, token_hash, expires_at) VALUES (?, 'token-hash', 1)`, userID)

	mustExec(t, db, `DELETE FROM users WHERE id = ?`, userID)

	for _, table := range []string{"login_codes", "sessions"} {
		if got := countRows(t, db, table); got != 0 {
			t.Errorf("%s rows after deleting the user = %d, want 0", table, got)
		}
	}
}

func TestSessionsTable_TokenHashIsUnique(t *testing.T) {
	db := openMigratedDB(t)
	userID := insertUser(t, db)
	insert := `INSERT INTO sessions (user_id, token_hash, expires_at) VALUES (?, 'same-hash', 1)`
	mustExec(t, db, insert, userID)

	if _, err := db.Exec(insert, userID); err == nil {
		t.Errorf("second session with the same hash error = nil, want a UNIQUE violation")
	}
}

func TestAccessTables_RequireTheirColumns(t *testing.T) {
	db := openMigratedDB(t)
	userID := insertUser(t, db)
	mustExec(t, db, `INSERT INTO code_sends (email, sent_at) VALUES ('ana@example.com', 1)`)
	cases := []struct {
		column string
		insert string
		args   []any
	}{
		{"code_sends.email", `INSERT INTO code_sends (sent_at) VALUES (1)`, nil},
		{"code_sends.sent_at", `INSERT INTO code_sends (email) VALUES ('ana@example.com')`, nil},
		{"login_codes.code_hash", `INSERT INTO login_codes (user_id, expires_at, attempts) VALUES (?, 1, 0)`, []any{userID}},
		{"login_codes.expires_at", `INSERT INTO login_codes (user_id, code_hash, attempts) VALUES (?, 'h', 0)`, []any{userID}},
		{"login_codes.attempts not negative", `INSERT INTO login_codes (user_id, code_hash, expires_at, attempts) VALUES (?, 'h', 1, -1)`, []any{userID}},
		{"sessions.user_id", `INSERT INTO sessions (token_hash, expires_at) VALUES ('h', 1)`, nil},
		{"sessions.expires_at", `INSERT INTO sessions (user_id, token_hash) VALUES (?, 'h')`, []any{userID}},
	}
	for _, tc := range cases {
		t.Run(tc.column, func(t *testing.T) {
			if _, err := db.Exec(tc.insert, tc.args...); err == nil {
				t.Errorf("insert breaking %s error = nil, want a constraint violation", tc.column)
			}
		})
	}
}
