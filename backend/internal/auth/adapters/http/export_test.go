package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/export"
	authsqlite "github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	statementssqlite "github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/sqlite"
	statementsusecases "github.com/code-rcplaza/mail-checker/backend/internal/statements/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/catalog"
	subscriptionssqlite "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite"
	subscriptionsusecases "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

// copiesOver gathers the copy from every feature over the test database.
func copiesOver(db *sql.DB) export.UserCopies {
	now := func() time.Time { return testNow }
	decisions := subscriptionssqlite.NewDecisions(db, now)
	aliases := subscriptionssqlite.NewAliases(db)
	manuals := subscriptionssqlite.NewManuals(db)
	overrides := subscriptionssqlite.NewOverrides(db)
	charges := subscriptionssqlite.NewChargeReader(db)
	return export.NewUserCopies(export.CopySources{
		Me:         usecases.NewGetMe(authsqlite.NewUsers(db, now)),
		Statements: statementsusecases.NewListStatements(statementssqlite.NewStatements(db)),
		Saved: subscriptionsusecases.NewListSavedChoices(subscriptionsusecases.SavedChoicesSources{
			Charges: charges, Decisions: decisions, Aliases: aliases, Manuals: manuals, Overrides: overrides,
		}),
		Summary: subscriptionsusecases.NewGetSummary(subscriptionsusecases.SummarySources{
			Lines:     subscriptionsusecases.ChargeLines{Charges: charges, Catalog: catalog.Catalog{}, Now: now},
			Decisions: decisions, Aliases: aliases, Rates: subscriptionssqlite.NewExchangeRates(db),
			KeptRises: subscriptionssqlite.NewPriceAlertDecisions(db), Manuals: manuals, Overrides: overrides,
			Hidden: subscriptionssqlite.NewHiddenLines(db), Cancellations: subscriptionssqlite.NewCancellations(db),
		}),
	})
}

func TestExport(t *testing.T) {
	cases := []struct {
		query           string
		wantType        string
		wantDisposition string
	}{
		{"", "text/csv; charset=utf-8", `attachment; filename="fokko-copia.csv"`},
		{"?format=csv", "text/csv; charset=utf-8", `attachment; filename="fokko-copia.csv"`},
		{"?format=json", "application/json", `attachment; filename="fokko-copia.json"`},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			app := newAuthApp(t)
			cookie := app.signIn(t)

			got := app.do(http.MethodGet, "/api/export"+tc.query, "", cookie)

			if got.Code != http.StatusOK || got.Header().Get("Content-Type") != tc.wantType || got.Header().Get("Content-Disposition") != tc.wantDisposition {
				t.Errorf("export = %d %q %q, want 200 %q %q", got.Code, got.Header().Get("Content-Type"), got.Header().Get("Content-Disposition"), tc.wantType, tc.wantDisposition)
			}
		})
	}
}

func TestExport_Contents(t *testing.T) {
	app := newAuthApp(t)
	cookie := app.signIn(t)

	csv := app.do(http.MethodGet, "/api/export", "", cookie)
	jsonCopy := app.do(http.MethodGet, "/api/export?format=json", "", cookie)

	if !strings.HasPrefix(csv.Body.String(), "name,alias,amount,currency,cadence,nextCharge\n") {
		t.Errorf("csv = %q, want the subscriptions header", csv.Body.String())
	}
	var document struct {
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(jsonCopy.Body.Bytes(), &document); err != nil || document.Profile.Email != "ana@example.com" {
		t.Errorf("json copy = %s, want the profile of ana", jsonCopy.Body.String())
	}
}

func TestExport_UnknownFormatOrNoSession(t *testing.T) {
	app := newAuthApp(t)
	cookie := app.signIn(t)

	unknown := app.do(http.MethodGet, "/api/export?format=xlsx", "", cookie)
	anonymous := app.do(http.MethodGet, "/api/export", "")

	if unknown.Code != http.StatusBadRequest || errorOf(t, unknown).Error.Fields["format"] != "auth.export_format_unknown" {
		t.Errorf("unknown format = %d %s, want 400 with the format field", unknown.Code, unknown.Body.String())
	}
	if anonymous.Code != http.StatusUnauthorized {
		t.Errorf("export without session = %d, want 401", anonymous.Code)
	}
}
