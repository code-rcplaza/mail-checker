package sqlite

import (
	"database/sql"
	"testing"
)

const insertStatement = `INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at, bad_lines) VALUES (?, ?, ?, 'cartola.pdf', 'pdf', ?, 1, '[]')`

const insertCharge = `INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant, card, card_last4) VALUES (?, 'bci', 'statement', ?, ?, 1, ?, ?, 'NETFLIX', 'credito', ?)`

// statementWithOwner returns a database with a user and one statement of theirs.
func statementWithOwner(t *testing.T) (db *sql.DB, userID, statementID int64) {
	t.Helper()
	db = openMigratedDB(t)
	userID = insertUser(t, db)
	result, err := db.Exec(insertStatement, "st-1", userID, "bci", 100)
	if err != nil {
		t.Fatalf("insert statement: %v", err)
	}
	statementID, _ = result.LastInsertId()
	return db, userID, statementID
}

func TestStatementsTable_OnePerUserBankAndBillingDate(t *testing.T) {
	db, userID, _ := statementWithOwner(t)
	cases := []struct {
		name     string
		publicID string
		bank     string
		billing  int64
		wantErr  bool
	}{
		{"same user, bank and date is rejected", "st-2", "bci", 100, true},
		{"the same public id is rejected", "st-1", "bci", 200, true},
		{"another date is accepted", "st-3", "bci", 200, false},
		{"another bank is accepted", "st-4", "falabella", 100, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := db.Exec(insertStatement, tc.publicID, userID, tc.bank, tc.billing)

			if (err != nil) != tc.wantErr {
				t.Errorf("insert error = %v, want error %v", err, tc.wantErr)
			}
		})
	}
}

func TestChargesTable_Constraints(t *testing.T) {
	db, userID, statementID := statementWithOwner(t)
	mustExec(t, db, insertCharge, userID, "st-1:1", statementID, 599000, "CLP", "1234")
	cases := []struct {
		name       string
		externalID string
		amount     int64
		currency   string
		last4      string
		wantErr    bool
	}{
		{"a repeated external id is rejected", "st-1:1", 599000, "CLP", "", true},
		{"a negative amount is rejected", "st-1:2", -1, "CLP", "", true},
		{"an unknown currency is rejected", "st-1:3", 100, "EUR", "", true},
		{"three digits are rejected", "st-1:4", 100, "CLP", "123", true},
		{"a letter in the digits is rejected", "st-1:5", 100, "CLP", "12a4", true},
		{"no digits are accepted", "st-1:6", 100, "USD", "", false},
		{"zero is accepted", "st-1:7", 0, "CLP", "9876", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := db.Exec(insertCharge, userID, tc.externalID, statementID, tc.amount, tc.currency, tc.last4)

			if (err != nil) != tc.wantErr {
				t.Errorf("insert error = %v, want error %v", err, tc.wantErr)
			}
		})
	}
}

func TestChargesTable_GoWithTheirStatementAndUser(t *testing.T) {
	cases := []struct {
		name      string
		statement string
		ownerOf   func(userID, statementID int64) int64
	}{
		{"deleting the statement", `DELETE FROM statements WHERE id = ?`, func(_, statementID int64) int64 { return statementID }},
		{"deleting the user", `DELETE FROM users WHERE id = ?`, func(userID, _ int64) int64 { return userID }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, userID, statementID := statementWithOwner(t)
			mustExec(t, db, insertCharge, userID, "st-1:1", statementID, 100, "CLP", "")

			mustExec(t, db, tc.statement, tc.ownerOf(userID, statementID))

			if got := countRows(t, db, "charges"); got != 0 {
				t.Errorf("charges after %s = %d, want 0", tc.name, got)
			}
		})
	}
}
