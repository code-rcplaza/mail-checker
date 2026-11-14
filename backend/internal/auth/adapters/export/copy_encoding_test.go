package export

import (
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
	"time"

	authdomain "github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	statementsdomain "github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	subscriptionsdomain "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	subscriptionsusecases "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

func internalID(id int64) *int64 { return &id }

// sampleCopy uses internal ids that must not leave: 7701 (user), 8802 (statement), 9903 (charge).
func sampleCopy() UserCopy {
	billed := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
	return UserCopy{
		Profile:    authdomain.User{ID: 7701, PublicID: "user-pub", Name: "Ana", Email: "ana@example.com", AcceptedTermsAt: testNow},
		Statements: []statementsdomain.Statement{{ID: 8802, PublicID: "st-pub", UserID: 7701, Bank: "bci", Filename: "c.pdf", BillingDate: billed, ImportedAt: testNow}},
		Saved: subscriptionsusecases.SavedChoices{
			Charges: []shared.Charge{{ID: 9903, UserID: 7701, StatementID: internalID(8802), ExternalID: "st-pub:3", Bank: "bci", Source: shared.SourceStatement,
				OccurredAt: billed, Amount: 799000, Currency: shared.CurrencyCLP, Merchant: "NETFLIX", Card: "credito", CardLast4: "1234"}},
			Decisions: map[subscriptionsdomain.GroupKey]subscriptionsdomain.Decision{"k2": {Kind: subscriptionsdomain.DecisionConfirmed, Cadence: subscriptionsdomain.CadenceMonthly}, "k1": {Kind: subscriptionsdomain.DecisionDismissed}},
			Overrides: map[subscriptionsdomain.GroupKey]subscriptionsdomain.LineOverride{"k1": {BillingDay: 9}},
		},
	}
}

func TestEncodeJSON_HasPublicIDsOnlyAndUnixDates(t *testing.T) {
	body, err := EncodeJSON(sampleCopy())
	if err != nil {
		t.Fatalf("EncodeJSON error = %v", err)
	}
	compact := strings.Join(strings.Fields(string(body)), "")

	for _, internal := range []string{"7701", "8802", "9903", `"userId"`, "st-pub:3", "externalId"} {
		if strings.Contains(compact, internal) {
			t.Errorf("copy contains %s:\n%s", internal, body)
		}
	}
	for _, want := range []string{
		`"profile":{"id":"user-pub","name":"Ana","email":"ana@example.com","acceptedTermsAt":` + fmt.Sprint(testNow.Unix()) + `}`,
		`"billingDate":` + fmt.Sprint(time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC).Unix()),
		`"statementId":"st-pub"`,
		`"amount":{"amount":799000,"currency":"CLP"}`,
		`"decisions":[{"lineKey":"k1","decision":"dismissed"},{"lineKey":"k2","decision":"confirmed","cadence":"monthly"}]`,
		`"overrides":[{"lineKey":"k1","amount":null,"billingDay":9}]`,
	} {
		if !strings.Contains(compact, want) {
			t.Errorf("copy lacks %s:\n%s", want, body)
		}
	}
}

func TestEncodeCSV_OneRowPerSubscription(t *testing.T) {
	next := time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)
	subscriptions := []subscriptionsdomain.Subscription{
		{Parent: "service:Netflix", Alias: "Padres", Amount: 799050, Currency: shared.CurrencyCLP, Cadence: subscriptionsdomain.CadenceMonthly, NextCharge: next},
		{Parent: "merchant:SHOP, \"INC\"\nTWO:USD", Amount: 1299, Currency: shared.CurrencyUSD},
	}

	body, err := EncodeCSV(subscriptions)
	if err != nil {
		t.Fatalf("EncodeCSV error = %v", err)
	}

	rows, err := csv.NewReader(strings.NewReader(string(body))).ReadAll()
	want := [][]string{
		{"name", "alias", "amount", "currency", "cadence", "nextCharge"},
		{"Netflix", "Padres", "7990.50", "CLP", "monthly", "2026-11-20"},
		{"SHOP, \"INC\"\nTWO", "", "12.99", "USD", "", ""},
	}
	if err != nil || fmt.Sprint(rows) != fmt.Sprint(want) {
		t.Errorf("rows = %q, %v; want %q (quotes, commas and newlines escaped)", rows, err, want)
	}
}

func TestEncodeCSV_NeutralizesFormulas(t *testing.T) {
	for _, start := range []string{"=", "+", "-", "@"} {
		t.Run(start, func(t *testing.T) {
			formula := start + "HYPERLINK(\"x\")"
			body, err := EncodeCSV([]subscriptionsdomain.Subscription{{Parent: "manual:m1", Alias: formula, Amount: 100, Currency: shared.CurrencyCLP}})
			if err != nil {
				t.Fatalf("EncodeCSV error = %v", err)
			}

			rows, err := csv.NewReader(strings.NewReader(string(body))).ReadAll()
			if err != nil || len(rows) != 2 || rows[1][1] != "'"+formula {
				t.Errorf("rows = %q, %v; want the alias prefixed with a quote", rows, err)
			}
		})
	}
}
