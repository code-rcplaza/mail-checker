package domain

import (
	"maps"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func activeSubscription(amount shared.Money, currency shared.Currency, cadence Cadence) Subscription {
	return Subscription{Key: "test", Status: StatusActive, Cadence: cadence, Amount: amount, Currency: currency}
}

func possibleSubscription(amount shared.Money, currency shared.Currency, decision *Decision) Subscription {
	return Subscription{Key: "test", Status: StatusPossible, Amount: amount, Currency: currency, Decision: decision}
}

func TestMonthlySpend(t *testing.T) {
	cases := []struct {
		name          string
		subscriptions []Subscription
		want          map[shared.Currency]shared.Money
	}{
		{
			"active monthly subscriptions add up per currency",
			[]Subscription{
				activeSubscription(659000, shared.CurrencyCLP, CadenceMonthly),
				activeSubscription(999000, shared.CurrencyCLP, CadenceMonthly),
				activeSubscription(2000, shared.CurrencyUSD, CadenceMonthly),
			},
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 1658000, shared.CurrencyUSD: 2000},
		},
		{
			"active yearly subscription adds a twelfth",
			[]Subscription{activeSubscription(12000000, shared.CurrencyCLP, CadenceYearly)},
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 1000000},
		},
		{
			"yearly amount not divisible by 12 rounds half up to the hundredth",
			[]Subscription{activeSubscription(9990, shared.CurrencyUSD, CadenceYearly)}, // 832.5
			map[shared.Currency]shared.Money{shared.CurrencyUSD: 833},
		},
		{
			"possible confirmed as monthly adds its amount",
			[]Subscription{possibleSubscription(2000, shared.CurrencyUSD, &Decision{Kind: DecisionConfirmed, Cadence: CadenceMonthly})},
			map[shared.Currency]shared.Money{shared.CurrencyUSD: 2000},
		},
		{
			"possible confirmed as yearly adds a twelfth",
			[]Subscription{possibleSubscription(12000, shared.CurrencyUSD, &Decision{Kind: DecisionConfirmed, Cadence: CadenceYearly})},
			map[shared.Currency]shared.Money{shared.CurrencyUSD: 1000},
		},
		{
			// The active one keeps the expected total non-empty, so the case
			// cannot pass with an empty result.
			"inactive, undecided and dismissed do not add",
			[]Subscription{
				activeSubscription(100000, shared.CurrencyCLP, CadenceMonthly),
				{Key: "test", Status: StatusInactive, Cadence: CadenceMonthly, Amount: 500000, Currency: shared.CurrencyCLP},
				possibleSubscription(700000, shared.CurrencyCLP, nil),
				possibleSubscription(900000, shared.CurrencyCLP, &Decision{Kind: DecisionDismissed}),
			},
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 100000},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MonthlySpend(tc.subscriptions)

			if !maps.Equal(got, tc.want) {
				t.Errorf("MonthlySpend = %v, want %v", got, tc.want)
			}
		})
	}
}
