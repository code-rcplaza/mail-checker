package domain

import (
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func rateOn(month time.Month, pesosPerDollar shared.Money) shared.ExchangeRate {
	return shared.ExchangeRate{PesosPerDollar: pesosPerDollar, Date: time.Date(2026, month, 20, 0, 0, 0, 0, time.UTC)}
}

func TestLatestExchangeRate(t *testing.T) {
	cases := []struct {
		name     string
		rates    []shared.ExchangeRate
		want     shared.ExchangeRate
		wantRate bool
	}{
		{
			"the most recent date wins",
			[]shared.ExchangeRate{rateOn(time.September, 94000), rateOn(time.October, 95000), rateOn(time.August, 96000)},
			rateOn(time.October, 95000),
			true,
		},
		{
			"same date: the higher rate wins, whatever the order",
			[]shared.ExchangeRate{rateOn(time.October, 95000), rateOn(time.October, 95500), rateOn(time.October, 94000)},
			rateOn(time.October, 95500),
			true,
		},
		{"an empty list has no rate", nil, shared.ExchangeRate{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := slices.Clone(tc.rates)

			got, hasRate := LatestExchangeRate(tc.rates)

			if hasRate != tc.wantRate || got != tc.want {
				t.Errorf("LatestExchangeRate = (%+v, %v), want (%+v, %v)", got, hasRate, tc.want, tc.wantRate)
			}
			if !reflect.DeepEqual(tc.rates, original) {
				t.Errorf("rates changed to %v, want the original %v", tc.rates, original)
			}
		})
	}
}

func TestConvertToPesos(t *testing.T) {
	rate := rateOn(time.October, 95000)
	cases := []struct {
		name            string
		amount          shared.Money
		currency        shared.Currency
		rate            shared.ExchangeRate
		want            shared.Money
		wantApproximate bool
	}{
		{"US$ 20,00 at 950 is CLP 19.000", 2000, shared.CurrencyUSD, rate, 1900000, true},
		{"half a hundredth rounds up", 1, shared.CurrencyUSD, shared.ExchangeRate{PesosPerDollar: 95050}, 951, true},
		{"under half a hundredth rounds down", 1, shared.CurrencyUSD, shared.ExchangeRate{PesosPerDollar: 95049}, 950, true},
		{"US$ 0 is CLP 0", 0, shared.CurrencyUSD, rate, 0, true},
		{"pesos stay the same and are exact", 599000, shared.CurrencyCLP, rate, 599000, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, isApproximate, err := ConvertToPesos(tc.amount, tc.currency, tc.rate)

			if err != nil {
				t.Fatalf("ConvertToPesos error = %v, want nil", err)
			}
			if got != tc.want || isApproximate != tc.wantApproximate {
				t.Errorf("ConvertToPesos = (%d, approximate %v), want (%d, approximate %v)", got, isApproximate, tc.want, tc.wantApproximate)
			}
		})
	}
}

func TestConvertToPesos_RejectsOtherCurrencies(t *testing.T) {
	_, _, err := ConvertToPesos(2000, shared.Currency("EUR"), rateOn(time.October, 95000))

	if !errors.Is(err, shared.ErrInvalid) {
		t.Errorf("ConvertToPesos error = %v, want ErrInvalid", err)
	}
}
