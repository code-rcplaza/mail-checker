package sqlite

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.ChargeReader = ChargeReader{}

var since = time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)

// insertStatementCharge stores a charge the way the statements repo does.
func insertStatementCharge(t *testing.T, db *sql.DB, charge shared.Charge) int64 {
	t.Helper()
	result, err := db.Exec(`INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at) VALUES (?, ?, 'bci', 'c.pdf', 'pdf', ?, 1)`,
		charge.ExternalID, charge.UserID, charge.OccurredAt.Unix())
	if err != nil {
		t.Fatalf("insert statement: %v", err)
	}
	statementID, _ := result.LastInsertId()
	_, err = db.Exec(`INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant, card, card_last4) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		charge.UserID, charge.Bank, charge.Source, charge.ExternalID, statementID, charge.OccurredAt.Unix(), charge.Amount, charge.Currency, charge.Merchant, charge.Card, charge.CardLast4)
	if err != nil {
		t.Fatalf("insert charge: %v", err)
	}
	return statementID
}

func TestChargeReader_ListSince(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)
	users := sqlitetest.Users{DB: db}
	anaID := users.Insert(t, "ana@example.com")
	betoID := users.Insert(t, "beto@example.com")
	charge := func(userID int64, externalID string, occurredAt time.Time) shared.Charge {
		return shared.Charge{UserID: userID, Bank: "bci", Source: shared.SourceStatement, ExternalID: externalID, OccurredAt: occurredAt,
			Amount: 1290000, Currency: shared.CurrencyUSD, Merchant: "OPENAI *CHATGPT", Card: "credito", CardLast4: "4321"}
	}
	onSince := charge(anaID, "ana:1", since)
	after := charge(anaID, "ana:2", since.AddDate(0, 1, 0))
	want := []shared.Charge{onSince, after}
	for _, stored := range []shared.Charge{charge(anaID, "ana:0", since.Add(-time.Second)), onSince, after, charge(betoID, "beto:1", since)} {
		id := insertStatementCharge(t, db, stored)
		for i := range want {
			if want[i].ExternalID == stored.ExternalID {
				want[i].StatementID = &id
			}
		}
	}

	got, err := NewChargeReader(db).ListSince(context.Background(), anaID, since)

	for i := range got {
		got[i].ID = 0 // assigned by the database
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ListSince = %+v, %v; want only Ana's charges from %v: %+v", got, err, since, want)
	}
}
