package domain

import (
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func day(year int, month time.Month, dayOfMonth int) time.Time {
	return time.Date(year, month, dayOfMonth, 0, 0, 0, 0, time.UTC)
}

func TestNextCharge(t *testing.T) {
	cases := []struct {
		name       string
		lastCharge time.Time
		typicalDay int
		cadence    Cadence
		want       time.Time
	}{
		{"monthly: a month after the last charge", day(2026, time.September, 1), 1, CadenceMonthly, day(2026, time.October, 1)},
		{"monthly: a late charge keeps the typical day", day(2026, time.September, 3), 1, CadenceMonthly, day(2026, time.October, 1)},
		{"monthly: an early charge does not repeat the same month", day(2026, time.September, 30), 1, CadenceMonthly, day(2026, time.November, 1)},
		{"monthly: december goes to january", day(2026, time.December, 15), 15, CadenceMonthly, day(2027, time.January, 15)},
		{"yearly: a year after the last charge", day(2026, time.March, 5), 5, CadenceYearly, day(2027, time.March, 5)},
		{"yearly: 29 february lands on the 28th of a common year", day(2024, time.February, 29), 29, CadenceYearly, day(2025, time.February, 28)},
		{"short month: day 31 lands on the 30th", day(2026, time.March, 31), 31, CadenceMonthly, day(2026, time.April, 30)},
		{"short month: day 30 lands on the last day of february", day(2026, time.January, 30), 30, CadenceMonthly, day(2026, time.February, 28)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasNext := NextCharge(tc.lastCharge, tc.typicalDay, tc.cadence)

			if !hasNext || !got.Equal(tc.want) {
				t.Errorf("NextCharge = (%v, %v), want (%v, true)", got, hasNext, tc.want)
			}
		})
	}
}

func TestNextCharge_WithoutCadenceThereIsNone(t *testing.T) {
	got, hasNext := NextCharge(day(2026, time.September, 1), 1, CadenceNone)

	if hasNext || !got.IsZero() {
		t.Errorf("NextCharge = (%v, %v), want none", got, hasNext)
	}
}

func TestBuildSummary_ActiveSubscriptionCarriesItsNextCharge(t *testing.T) {
	charges := scenarioCharges("netflix", "NETFLIX", nil, 1299000, shared.CurrencyCLP, 25, lastMonths...)

	got := summaryOfCharges(charges, nil)

	if len(got.Subscriptions) != 1 {
		t.Fatalf("subscriptions = %d, want 1", len(got.Subscriptions))
	}
	want := time.Date(2026, time.November, 25, 0, 0, 0, 0, time.UTC)
	if next := got.Subscriptions[0].NextCharge; !next.Equal(want) {
		t.Errorf("next charge = %v, want %v", next, want)
	}
}

func TestBuildSummary_UndecidedPossibleWithoutCadenceHasNoNextCharge(t *testing.T) {
	charges := scenarioCharges("claude", "ANTHROPIC CLAUDE", nil, 2000, shared.CurrencyUSD, 20, time.October)

	got := summaryOfCharges(charges, nil)

	if len(got.Possible) != 1 {
		t.Fatalf("possible = %d, want 1", len(got.Possible))
	}
	if next := got.Possible[0].NextCharge; !next.IsZero() {
		t.Errorf("next charge = %v, want none", next)
	}
}
