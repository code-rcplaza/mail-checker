package sqlite

import (
	"strings"
	"testing"
)

// schemaCase is one insert into a subscriptions table and whether the base must reject it.
type schemaCase struct {
	name       string
	insert     string
	args       []any
	wantReject bool
}

func runSchemaCases(t *testing.T, cases []schemaCase) {
	t.Helper()
	db := openMigratedDB(t)
	userID := insertUser(t, db)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]any{userID}, tc.args...)
			_, err := db.Exec(tc.insert, args...)

			if (err != nil) != tc.wantReject {
				t.Errorf("insert error = %v, want rejected %v", err, tc.wantReject)
			}
		})
	}
}

func TestSubscriptionDecisionsTable(t *testing.T) {
	insert := `INSERT INTO subscription_decisions (user_id, group_key, decision, cadence, decided_at) VALUES (?, ?, ?, ?, 1)`
	runSchemaCases(t, []schemaCase{
		{"a confirmed decision with cadence", insert, []any{"k1", "confirmed", "monthly"}, false},
		{"a dismissed decision without cadence", insert, []any{"k2", "dismissed", nil}, false},
		{"the same key again", insert, []any{"k1", "dismissed", nil}, true},
		{"an unknown decision", insert, []any{"k3", "later", nil}, true},
		{"confirmed without cadence", insert, []any{"k4", "confirmed", nil}, true},
		{"dismissed with cadence", insert, []any{"k5", "dismissed", "yearly"}, true},
		{"an unknown cadence", insert, []any{"k6", "confirmed", "weekly"}, true},
	})
}

func TestLineAliasesTable(t *testing.T) {
	insert := `INSERT INTO line_aliases (user_id, line_key, alias) VALUES (?, ?, ?)`
	runSchemaCases(t, []schemaCase{
		{"an alias of 40 characters", insert, []any{"k1", strings.Repeat("ñ", 40)}, false},
		{"an alias of 41 characters", insert, []any{"k2", strings.Repeat("ñ", 41)}, true},
		{"the same key again", insert, []any{"k1", "Padres"}, true},
	})
}

func TestLineOverridesTable(t *testing.T) {
	insert := `INSERT INTO line_overrides (user_id, line_key, amount, billing_day) VALUES (?, ?, ?, ?)`
	runSchemaCases(t, []schemaCase{
		{"only an amount", insert, []any{"k1", 599000, nil}, false},
		{"only a day", insert, []any{"k2", nil, 31}, false},
		{"day 0", insert, []any{"k3", nil, 0}, true},
		{"day 32", insert, []any{"k4", nil, 32}, true},
		{"amount 0", insert, []any{"k5", 0, nil}, true},
		{"the same key again", insert, []any{"k1", 1, nil}, true},
	})
}

func TestManualSubscriptionsTable(t *testing.T) {
	insert := `INSERT INTO manual_subscriptions (user_id, public_id, name, amount, currency, billing_day) VALUES (?, ?, 'Gimnasio', ?, ?, ?)`
	runSchemaCases(t, []schemaCase{
		{"a valid manual", insert, []any{"m1", 2500000, "CLP", 5}, false},
		{"the same public id", insert, []any{"m1", 2500000, "CLP", 5}, true},
		{"amount 0", insert, []any{"m2", 0, "CLP", 5}, true},
		{"an unknown currency", insert, []any{"m3", 100, "EUR", 5}, true},
		{"day 32", insert, []any{"m4", 100, "USD", 32}, true},
	})
}

func TestHiddenLinesAndCancellationsTables(t *testing.T) {
	hide := `INSERT INTO hidden_lines (user_id, line_key, hidden_at) VALUES (?, ?, 1)`
	mark := `INSERT INTO cancellations (user_id, line_key, marked_at) VALUES (?, ?, 1)`
	runSchemaCases(t, []schemaCase{
		{"a hidden line", hide, []any{"k1"}, false},
		{"the same hidden line again", hide, []any{"k1"}, true},
		{"a marked line", mark, []any{"k1"}, false},
		{"the same marked line again", mark, []any{"k1"}, true},
	})
}

func TestPriceAlertDecisionsTable(t *testing.T) {
	insert := `INSERT INTO price_alert_decisions (user_id, line_key, amount) VALUES (?, ?, ?)`
	runSchemaCases(t, []schemaCase{
		{"a kept amount", insert, []any{"k1", 1050000}, false},
		{"another amount of the same line", insert, []any{"k1", 1100000}, false},
		{"the same amount again", insert, []any{"k1", 1050000}, true},
		{"amount 0", insert, []any{"k2", 0}, true},
	})
}

func TestExchangeRatesTable(t *testing.T) {
	insert := `INSERT INTO exchange_rates (user_id, pesos_per_dollar, date) VALUES (?, ?, ?)`
	runSchemaCases(t, []schemaCase{
		{"a rate", insert, []any{95012, 100}, false},
		{"another date", insert, []any{95500, 200}, false},
		{"the same date again", insert, []any{96000, 100}, true},
		{"rate 0", insert, []any{0, 300}, true},
		{"a negative rate", insert, []any{-1, 400}, true},
	})
}

func TestSubscriptionsTables_GoWithTheirUser(t *testing.T) {
	db := openMigratedDB(t)
	userID := insertUser(t, db)
	mustExec(t, db, `INSERT INTO subscription_decisions (user_id, group_key, decision, decided_at) VALUES (?, 'k', 'dismissed', 1)`, userID)
	mustExec(t, db, `INSERT INTO line_aliases (user_id, line_key, alias) VALUES (?, 'k', 'Padres')`, userID)
	mustExec(t, db, `INSERT INTO line_overrides (user_id, line_key, amount) VALUES (?, 'k', 1)`, userID)
	mustExec(t, db, `INSERT INTO manual_subscriptions (user_id, public_id, name, amount, currency, billing_day) VALUES (?, 'm', 'G', 1, 'CLP', 1)`, userID)
	mustExec(t, db, `INSERT INTO hidden_lines (user_id, line_key, hidden_at) VALUES (?, 'k', 1)`, userID)
	mustExec(t, db, `INSERT INTO cancellations (user_id, line_key, marked_at) VALUES (?, 'k', 1)`, userID)
	mustExec(t, db, `INSERT INTO price_alert_decisions (user_id, line_key, amount) VALUES (?, 'k', 1)`, userID)
	mustExec(t, db, `INSERT INTO exchange_rates (user_id, pesos_per_dollar, date) VALUES (?, 95012, 1)`, userID)

	mustExec(t, db, `DELETE FROM users WHERE id = ?`, userID)

	for _, table := range []string{"subscription_decisions", "line_aliases", "line_overrides", "manual_subscriptions", "hidden_lines", "cancellations", "price_alert_decisions", "exchange_rates"} {
		if got := countRows(t, db, table); got != 0 {
			t.Errorf("%s rows after deleting the user = %d, want 0", table, got)
		}
	}
}
