package domain

import (
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var (
	accountA = int64(1)
	accountB = int64(2)
)

// paidCharge is a Netflix-like charge on the given month and day of 2026.
func paidCharge(externalID string, account *int64, last4 string, month time.Month, day int) shared.Charge {
	source := shared.SourceEmail
	if account == nil {
		source = shared.SourceStatement
	}
	return shared.Charge{
		Bank:       "bci",
		Source:     source,
		ExternalID: externalID,
		AccountID:  account,
		OccurredAt: time.Date(2026, month, day, 12, 0, 0, 0, time.UTC),
		Amount:     1299000,
		Currency:   shared.CurrencyCLP,
		Merchant:   "NETFLIX",
		CardLast4:  last4,
	}
}

func originIDs(origins [][]shared.Charge) [][]string {
	ids := make([][]string, 0, len(origins))
	for _, origin := range origins {
		ids = append(ids, externalIDs(origin))
	}
	return ids
}

func assertOrigins(t *testing.T, charges []shared.Charge, want [][]string) {
	t.Helper()
	if got := originIDs(SplitByOrigin(charges)); !reflect.DeepEqual(got, want) {
		t.Errorf("origins = %v, want %v", got, want)
	}
}

func TestSplitByOrigin_OneOriginPerMailAccount(t *testing.T) {
	charges := []shared.Charge{
		paidCharge("a1", &accountA, "4821", time.August, 1),
		paidCharge("b1", &accountB, "4821", time.August, 3),
		paidCharge("a2", &accountA, "4821", time.September, 1),
		paidCharge("b2", &accountB, "4821", time.September, 3),
	}

	assertOrigins(t, charges, [][]string{{"a1", "a2"}, {"b1", "b2"}})
}

func TestSplitByOrigin_OneOriginPerCard(t *testing.T) {
	charges := []shared.Charge{
		paidCharge("x1", &accountA, "4821", time.August, 1),
		paidCharge("y1", &accountA, "1234", time.August, 3),
		paidCharge("x2", &accountA, "4821", time.September, 1),
		paidCharge("y2", &accountA, "1234", time.September, 3),
	}

	assertOrigins(t, charges, [][]string{{"x1", "x2"}, {"y1", "y2"}})
}

func TestSplitByOrigin_StatementChargeJoinsTheOriginItChainsWith(t *testing.T) {
	charges := []shared.Charge{
		paidCharge("a1", &accountA, "4821", time.August, 1),
		paidCharge("s1", nil, "4821", time.September, 1),
	}

	assertOrigins(t, charges, [][]string{{"a1", "s1"}})
}

func TestSplitByOrigin_StatementChargeThatChainsWithNoOriginStaysApart(t *testing.T) {
	unchained := paidCharge("s1", nil, "4821", time.August, 10)
	unchained.Amount = 299000
	charges := []shared.Charge{paidCharge("a1", &accountA, "4821", time.August, 1), unchained}

	assertOrigins(t, charges, [][]string{{"a1"}, {"s1"}})
}

func TestSplitByOrigin_EmptyCardDoesNotSeparate(t *testing.T) {
	charges := []shared.Charge{
		paidCharge("a1", &accountA, "4821", time.August, 1),
		paidCharge("a2", &accountA, "", time.September, 1),
	}

	assertOrigins(t, charges, [][]string{{"a1", "a2"}})
}

func TestSplitByOrigin_LeavesTheCallerSliceUntouched(t *testing.T) {
	charges := []shared.Charge{
		paidCharge("b1", &accountB, "4821", time.September, 3),
		paidCharge("a1", &accountA, "4821", time.August, 1),
	}
	before := slices.Clone(charges)

	SplitByOrigin(charges)

	if !reflect.DeepEqual(charges, before) {
		t.Errorf("charges changed to %v", externalIDs(charges))
	}
}

func TestChargesMatch_ComparesCardsOnlyWhenBothHaveThem(t *testing.T) {
	tests := []struct {
		name                      string
		statementLast4, mailLast4 string
		want                      bool
	}{
		{"different cards", "4821", "1234", false},
		{"same card", "4821", "4821", true},
		{"statement without card", "", "1234", true},
		{"email without card", "4821", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statement := statementCharge("s1", "STREAMING CO", 599000, 5)
			statement.CardLast4 = tt.statementLast4
			email := emailCharge("e1", "STREAMING CO", 599000, 3)
			email.CardLast4 = tt.mailLast4

			if got := ChargesMatch(statement, email, identifyFake, normalizeFake); got != tt.want {
				t.Errorf("ChargesMatch = %v, want %v", got, tt.want)
			}
		})
	}
}
