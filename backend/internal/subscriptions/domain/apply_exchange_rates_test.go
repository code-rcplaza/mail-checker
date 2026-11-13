package domain

import (
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestApplyExchangeRates_EachLineCarriesItsAmountInPesos(t *testing.T) {
	october := []shared.ExchangeRate{rateOn(time.October, 95000)}
	cases := []struct {
		name     string
		amount   shared.Money
		currency shared.Currency
		rates    []shared.ExchangeRate
		want     *PesosAmount
	}{
		{"dollars with a rate show their approximate equivalent", 2000, shared.CurrencyUSD, october, &PesosAmount{Amount: 1900000, IsApproximate: true}},
		{"dollars without a rate show only the original", 2000, shared.CurrencyUSD, nil, nil},
		{"pesos are their own exact equivalent", 599000, shared.CurrencyCLP, nil, &PesosAmount{Amount: 599000}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			line := Subscription{Key: "line", Amount: tc.amount, Currency: tc.currency}
			summary := Summary{Subscriptions: []Subscription{line}, Possible: []Subscription{line}, Inactive: []Subscription{line}}

			got := ApplyExchangeRates(summary, tc.rates)

			for _, section := range [][]Subscription{got.Subscriptions, got.Possible, got.Inactive} {
				s := section[0]
				if s.Amount != tc.amount || s.Currency != tc.currency {
					t.Errorf("original = %d %s, want %d %s", s.Amount, s.Currency, tc.amount, tc.currency)
				}
				if !samePesosAmount(s.InPesos, tc.want) {
					t.Errorf("in pesos = %+v, want %+v", s.InPesos, tc.want)
				}
			}
		})
	}
}

func samePesosAmount(a, b *PesosAmount) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func TestApplyExchangeRates_LeavesTheGivenSummaryUntouched(t *testing.T) {
	lines := []Subscription{{Key: "line", Amount: 2000, Currency: shared.CurrencyUSD}}
	summary := Summary{Subscriptions: lines}

	ApplyExchangeRates(summary, []shared.ExchangeRate{rateOn(time.October, 95000)})

	if lines[0].InPesos != nil {
		t.Errorf("the caller's line changed to %+v", lines[0].InPesos)
	}
}
