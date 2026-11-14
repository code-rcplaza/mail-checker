package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
)

// fillEveryTable gives the user one row in each table that holds their data.
func fillEveryTable(t *testing.T, db *sql.DB, email string) int64 {
	t.Helper()
	userID := sqlitetest.Users{DB: db}.Insert(t, email)
	statements := []string{
		`INSERT INTO login_codes (user_id, code_hash, expires_at) VALUES (?1, 'h', 1)`,
		`INSERT INTO sessions (user_id, token_hash, expires_at) VALUES (?1, 'token-' || ?1, 1)`,
		`INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at) VALUES ('st-' || ?1, ?1, 'bci', 'c.pdf', 'pdf', 1, 1)`,
		`INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant) VALUES (?1, 'bci', 'statement', 'x', (SELECT id FROM statements WHERE user_id = ?1), 1, 1, 'CLP', 'M')`,
		`INSERT INTO charges (user_id, bank, source, external_id, occurred_at, amount, currency, merchant) VALUES (?1, 'bci', 'email', 'y', 1, 1, 'CLP', 'M')`,
		`INSERT INTO subscription_decisions (user_id, group_key, decision, decided_at) VALUES (?1, 'k', 'dismissed', 1)`,
		`INSERT INTO line_aliases (user_id, line_key, alias) VALUES (?1, 'k', 'a')`,
		`INSERT INTO line_overrides (user_id, line_key, amount) VALUES (?1, 'k', 1)`,
		`INSERT INTO manual_subscriptions (user_id, public_id, name, amount, currency, billing_day) VALUES (?1, 'm-' || ?1, 'n', 1, 'CLP', 1)`,
		`INSERT INTO hidden_lines (user_id, line_key, hidden_at) VALUES (?1, 'k', 1)`,
		`INSERT INTO cancellations (user_id, line_key, marked_at) VALUES (?1, 'k', 1)`,
		`INSERT INTO price_alert_decisions (user_id, line_key, amount) VALUES (?1, 'k', 1)`,
		`INSERT INTO exchange_rates (user_id, pesos_per_dollar, date) VALUES (?1, 1, 1)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement, userID); err != nil {
			t.Fatalf("fill %q: %v", statement, err)
		}
	}
	if _, err := db.Exec(`INSERT INTO code_sends (email, sent_at) VALUES (?, 1)`, email); err != nil {
		t.Fatalf("fill code_sends: %v", err)
	}
	return userID
}

// tablesWithUserID lists every table with a user_id column, so a new table is checked too.
func tablesWithUserID(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT m.name FROM sqlite_master m JOIN pragma_table_info(m.name) c ON c.name = 'user_id' WHERE m.type = 'table'`)
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table: %v", err)
		}
		tables = append(tables, name)
	}
	return tables
}

// rowsOf counts what is left of the user in every table, users and code_sends included.
func rowsOf(t *testing.T, db *sql.DB, owner userOwner) map[string]int {
	t.Helper()
	counts := make(map[string]int)
	for _, table := range tablesWithUserID(t, db) {
		counts[table] = countRows(t, db, rowCount{query: `SELECT count(*) FROM ` + table + ` WHERE user_id = ?`, arg: owner.id})
	}
	counts["users"] = countRows(t, db, rowCount{query: `SELECT count(*) FROM users WHERE id = ?`, arg: owner.id})
	counts["code_sends"] = countRows(t, db, rowCount{query: `SELECT count(*) FROM code_sends WHERE email = ?`, arg: owner.email})
	return counts
}

type userOwner struct {
	id    int64
	email string
}

type rowCount struct {
	query string
	arg   any
}

func countRows(t *testing.T, db *sql.DB, count rowCount) int {
	t.Helper()
	var n int
	if err := db.QueryRow(count.query, count.arg).Scan(&n); err != nil {
		t.Fatalf("%s: %v", count.query, err)
	}
	return n
}

func TestUsers_DeleteLeavesNothingOfTheUser(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)
	ana := userOwner{id: fillEveryTable(t, db, "ana@example.com"), email: "ana@example.com"}
	beto := userOwner{id: fillEveryTable(t, db, "beto@example.com"), email: "beto@example.com"}
	if tables := tablesWithUserID(t, db); len(tables) != 12 {
		t.Fatalf("tables with user_id = %v, want the 12 filled ones", tables)
	}

	if err := NewUsers(db, nil).Delete(context.Background(), ana.id); err != nil {
		t.Fatalf("Delete error = %v, want nil", err)
	}

	for table, n := range rowsOf(t, db, ana) {
		if n != 0 {
			t.Errorf("%s keeps %d rows of the deleted user", table, n)
		}
	}
	for table, n := range rowsOf(t, db, beto) {
		if n == 0 {
			t.Errorf("%s lost the rows of the other user", table)
		}
	}
}

func TestUsers_DeleteUnknownIsNotFound(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)
	beto := userOwner{id: fillEveryTable(t, db, "beto@example.com"), email: "beto@example.com"}

	err := NewUsers(db, nil).Delete(context.Background(), beto.id+1)

	if !errors.Is(err, shared.ErrNotFound) {
		t.Fatalf("Delete error = %v, want ErrNotFound", err)
	}
	if rowsOf(t, db, beto)["code_sends"] != 1 {
		t.Errorf("deleting an unknown user touched the code sends of another")
	}
}

func TestUsers_DeleteIsAllOrNothing(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)
	ana := userOwner{id: fillEveryTable(t, db, "ana@example.com"), email: "ana@example.com"}
	// A trigger makes the user's own delete fail after the code sends are gone.
	if _, err := db.Exec(`CREATE TRIGGER block_user_delete BEFORE DELETE ON users BEGIN SELECT RAISE(ABORT, 'blocked'); END`); err != nil {
		t.Fatalf("create trigger: %v", err)
	}

	if err := NewUsers(db, nil).Delete(context.Background(), ana.id); err == nil {
		t.Fatalf("Delete error = nil, want the blocked delete")
	}

	if rowsOf(t, db, ana)["code_sends"] != 1 || rowsOf(t, db, ana)["users"] != 1 {
		t.Errorf("rows = %v, want everything kept after the failed delete", rowsOf(t, db, ana))
	}
}
