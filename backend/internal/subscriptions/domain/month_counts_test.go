package domain

import (
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func santiagoTime(t *testing.T, year int, month time.Month, dayOfMonth, hour, minute int) time.Time {
	t.Helper()
	location, err := time.LoadLocation("America/Santiago")
	if err != nil {
		t.Fatalf("load America/Santiago: %v", err)
	}
	return time.Date(year, month, dayOfMonth, hour, minute, 0, 0, location)
}

func subscriptionChargedAt(status Status, chargedAt ...time.Time) Subscription {
	charges := make([]shared.Charge, 0, len(chargedAt))
	for _, at := range chargedAt {
		charges = append(charges, shared.Charge{OccurredAt: at})
	}
	return Subscription{Status: status, Cadence: CadenceMonthly, Charges: charges}
}

func TestCountChargedThisMonth(t *testing.T) {
	now := santiagoTime(t, 2026, time.November, 12, 12, 0)
	thisMonth := santiagoTime(t, 2026, time.November, 5, 10, 0)
	lastMonth := santiagoTime(t, 2026, time.October, 5, 10, 0)
	lastNightOfOctober := santiagoTime(t, 2026, time.October, 31, 23, 30).UTC()
	confirmed := subscriptionChargedAt(StatusPossible, thisMonth)
	confirmed.Decision = &Decision{Kind: DecisionConfirmed, Cadence: CadenceMonthly}
	cases := []struct {
		name          string
		subscriptions []Subscription
		want          int
	}{
		{
			"three charged this month and one not",
			[]Subscription{
				subscriptionChargedAt(StatusActive, lastMonth, thisMonth),
				subscriptionChargedAt(StatusActive, thisMonth),
				confirmed,
				subscriptionChargedAt(StatusActive, lastMonth),
			},
			3,
		},
		{
			"the last night of the previous month in Santiago does not count",
			[]Subscription{subscriptionChargedAt(StatusActive, lastNightOfOctober)},
			0,
		},
		{
			"an undecided possible does not count",
			[]Subscription{subscriptionChargedAt(StatusPossible, thisMonth)},
			0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CountChargedThisMonth(tc.subscriptions, now)

			if got != tc.want {
				t.Errorf("CountChargedThisMonth = %d, want %d", got, tc.want)
			}
		})
	}
}

func subscriptionWithNextCharge(nextCharge time.Time) Subscription {
	return Subscription{Status: StatusActive, Cadence: CadenceMonthly, NextCharge: nextCharge}
}

func TestCountRemainingThisMonth(t *testing.T) {
	now := santiagoTime(t, 2026, time.November, 10, 12, 0)
	nextOn := func(month time.Month, dayOfMonth int) Subscription {
		return subscriptionWithNextCharge(santiagoTime(t, 2026, month, dayOfMonth, 0, 0))
	}
	undecided := nextOn(time.November, 20)
	undecided.Status = StatusPossible
	cases := []struct {
		name          string
		subscriptions []Subscription
		want          int
	}{
		{
			"two still to come this month and one next month",
			[]Subscription{nextOn(time.November, 20), nextOn(time.November, 25), nextOn(time.December, 3)},
			2,
		},
		{"a next charge today is not counted: it may already be charged", []Subscription{nextOn(time.November, 10)}, 0},
		{"a next charge tomorrow counts", []Subscription{nextOn(time.November, 11)}, 1},
		{"a next charge that already passed does not count", []Subscription{nextOn(time.November, 3)}, 0},
		{"already charged this month with the next one in the following month", []Subscription{nextOn(time.December, 5)}, 0},
		{"a line without a next charge does not count", []Subscription{subscriptionWithNextCharge(time.Time{})}, 0},
		{"an undecided possible does not count", []Subscription{undecided}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CountRemainingThisMonth(tc.subscriptions, now)

			if got != tc.want {
				t.Errorf("CountRemainingThisMonth = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestBuildSummary_CountsTheChargesOfTheMonth(t *testing.T) {
	alreadyCharged := []time.Month{time.September, time.October, time.November}
	charges := slices.Concat(
		scenarioCharges("netflix", "NETFLIX", nil, 1299000, shared.CurrencyCLP, 5, alreadyCharged...),
		scenarioCharges("claude", "ANTHROPIC CLAUDE", nil, 2000, shared.CurrencyUSD, 8, alreadyCharged...),
		scenarioCharges("apple", "APPLE COM BILL", nil, 469000, shared.CurrencyCLP, 10, alreadyCharged...),
		scenarioCharges("google-small", "GOOGLE YOUTUBE", nil, 1100000, shared.CurrencyCLP, 20, lastMonths...),
		scenarioCharges("google-large", "GOOGLE YOUTUBE", nil, 2170000, shared.CurrencyCLP, 25, lastMonths...),
	)

	got := summaryOfCharges(charges, nil)

	if got.ChargedThisMonth != 3 || got.RemainingThisMonth != 2 {
		t.Errorf("charged = %d, remaining = %d; want 3 charged and 2 remaining", got.ChargedThisMonth, got.RemainingThisMonth)
	}
}
