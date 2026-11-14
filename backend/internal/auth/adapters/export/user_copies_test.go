package export

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	authsqlite "github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite"
	authusecases "github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	statementssqlite "github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/sqlite"
	statementsusecases "github.com/code-rcplaza/mail-checker/backend/internal/statements/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/catalog"
	subscriptionssqlite "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite"
	subscriptionsusecases "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

var testNow = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

// copiesOver gathers copies from a real database with every feature's repos.
func copiesOver(db *sql.DB) UserCopies {
	now := func() time.Time { return testNow }
	lines := subscriptionsusecases.ChargeLines{Charges: subscriptionssqlite.NewChargeReader(db), Catalog: catalog.Catalog{}, Now: now}
	decisions := subscriptionssqlite.NewDecisions(db, now)
	aliases := subscriptionssqlite.NewAliases(db)
	manuals := subscriptionssqlite.NewManuals(db)
	overrides := subscriptionssqlite.NewOverrides(db)
	summary := subscriptionsusecases.NewGetSummary(subscriptionsusecases.SummarySources{
		Lines: lines, Decisions: decisions, Aliases: aliases, Rates: subscriptionssqlite.NewExchangeRates(db),
		KeptRises: subscriptionssqlite.NewPriceAlertDecisions(db), Manuals: manuals, Overrides: overrides,
		Hidden: subscriptionssqlite.NewHiddenLines(db), Cancellations: subscriptionssqlite.NewCancellations(db),
	})
	return NewUserCopies(CopySources{
		Me:         authusecases.NewGetMe(authsqlite.NewUsers(db, now)),
		Statements: statementsusecases.NewListStatements(statementssqlite.NewStatements(db)),
		Saved: subscriptionsusecases.NewListSavedChoices(subscriptionsusecases.SavedChoicesSources{
			Charges: subscriptionssqlite.NewChargeReader(db), Decisions: decisions, Aliases: aliases, Manuals: manuals, Overrides: overrides,
		}),
		Summary: summary,
	})
}

// seedUser gives the user a statement with two monthly Netflix charges and
// one of each choice, every id and name tagged with their name.
func seedUser(t *testing.T, db *sql.DB, name string) int64 {
	t.Helper()
	userID := sqlitetest.Users{DB: db}.Insert(t, name+"@example.com")
	statements := []string{
		`INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at) VALUES ('st-' || ?2, ?1, 'bci', ?2 || '.pdf', 'pdf', ?3, ?3)`,
		`INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant) VALUES (?1, 'bci', 'statement', 'a', (SELECT id FROM statements WHERE user_id = ?1), ?3 - 2592000, 799000, 'CLP', 'NETFLIX')`,
		`INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant) VALUES (?1, 'bci', 'statement', 'b', (SELECT id FROM statements WHERE user_id = ?1), ?3, 799000, 'CLP', 'NETFLIX')`,
		`INSERT INTO subscription_decisions (user_id, group_key, decision, decided_at) VALUES (?1, 'k-' || ?2, 'dismissed', 1)`,
		`INSERT INTO line_aliases (user_id, line_key, alias) VALUES (?1, 'k-' || ?2, ?2)`,
		`INSERT INTO line_overrides (user_id, line_key, billing_day) VALUES (?1, 'k-' || ?2, 9)`,
		`INSERT INTO manual_subscriptions (user_id, public_id, name, amount, currency, billing_day) VALUES (?1, 'm-' || ?2, ?2, 100000, 'CLP', 5)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement, userID, name, testNow.AddDate(0, 0, -5).Unix()); err != nil {
			t.Fatalf("seed %q: %v", statement, err)
		}
	}
	return userID
}

func TestGather_BringsOnlyTheUsersData(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)
	anaID := seedUser(t, db, "ana")
	seedUser(t, db, "beto")

	got, err := copiesOver(db).Gather(context.Background(), anaID)
	if err != nil {
		t.Fatalf("Gather error = %v, want nil", err)
	}

	body, err := EncodeJSON(got)
	if err != nil {
		t.Fatalf("EncodeJSON error = %v", err)
	}
	if strings.Contains(string(body), "beto") {
		t.Errorf("copy of ana carries data of beto:\n%s", body)
	}
	if got.Profile.Email != "ana@example.com" || len(got.Statements) != 1 || len(got.Saved.Charges) != 2 ||
		len(got.Saved.Decisions) != 1 || len(got.Saved.Aliases) != 1 || len(got.Saved.Manuals) != 1 || len(got.Saved.Overrides) != 1 {
		t.Errorf("copy = %+v, want one of each of ana's records and her two charges", got)
	}
	if len(got.Subscriptions) != 2 {
		t.Errorf("subscriptions = %+v, want Netflix and the manual", got.Subscriptions)
	}
}

func TestGather_UserWithoutDataHasEmptyLists(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)
	userID := sqlitetest.Users{DB: db}.Insert(t, "ana@example.com")
	seedUser(t, db, "beto")

	got, err := copiesOver(db).Gather(context.Background(), userID)
	if err != nil {
		t.Fatalf("Gather error = %v, want nil", err)
	}
	body, err := EncodeJSON(got)
	if err != nil {
		t.Fatalf("EncodeJSON error = %v", err)
	}

	var document map[string]json.RawMessage
	if err := json.Unmarshal(body, &document); err != nil {
		t.Fatalf("decode copy: %v", err)
	}
	for _, list := range []string{"statements", "charges", "decisions", "aliases", "manuals", "overrides"} {
		if string(document[list]) != "[]" {
			t.Errorf("%s = %s, want []", list, document[list])
		}
	}
	if csv, err := EncodeCSV(got.Subscriptions); err != nil || string(csv) != "name,alias,amount,currency,cadence,nextCharge\n" {
		t.Errorf("csv = %q, %v; want only the header", csv, err)
	}
}
