package domain

import (
	"errors"
	"testing"
	"time"
)

var rateDate = time.Date(2026, time.October, 20, 0, 0, 0, 0, time.UTC)

func TestDeriveExchangeRate(t *testing.T) {
	cases := []struct {
		name     string
		pesos    Money
		dollars  Money
		want     Money
		wantRate bool
	}{
		{"CLP 190.000 for US$ 200,00 is 950 pesos per dollar", 19000000, 20000, 95000, true},
		{"half a hundredth rounds up", 1005, 200, 503, true},
		{"under half a hundredth rounds down", 1003, 300, 334, true},
		{"zero dollars give no rate", 19000000, 0, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasRate, err := DeriveExchangeRate(tc.pesos, tc.dollars, rateDate)

			if err != nil {
				t.Fatalf("DeriveExchangeRate error = %v, want nil", err)
			}
			if hasRate != tc.wantRate {
				t.Fatalf("has rate = %v, want %v", hasRate, tc.wantRate)
			}
			if !tc.wantRate {
				return
			}
			want := ExchangeRate{PesosPerDollar: tc.want, Date: rateDate}
			if got != want {
				t.Errorf("DeriveExchangeRate = %+v, want %+v", got, want)
			}
		})
	}
}

func TestDeriveExchangeRate_RejectsNegativeAmounts(t *testing.T) {
	cases := []struct {
		name    string
		pesos   Money
		dollars Money
	}{
		{"negative pesos", -19000000, 20000},
		{"negative dollars", 19000000, -20000},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, hasRate, err := DeriveExchangeRate(tc.pesos, tc.dollars, rateDate)

			if !errors.Is(err, ErrInvalid) || hasRate {
				t.Errorf("DeriveExchangeRate = (has rate %v, error %v), want no rate and ErrInvalid", hasRate, err)
			}
		})
	}
}
