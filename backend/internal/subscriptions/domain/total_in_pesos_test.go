package domain

import (
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestTotalInPesos(t *testing.T) {
	rate := rateOn(time.October, 95000)
	cases := []struct {
		name    string
		spend   map[shared.Currency]shared.Money
		hasRate bool
		want    Total
	}{
		{
			"pesos plus converted dollars, marked approximate",
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 1000000, shared.CurrencyUSD: 2000},
			true,
			Total{Amount: 2900000, IsApproximate: true},
		},
		{
			"only pesos is exact",
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 1000000},
			true,
			Total{Amount: 1000000},
		},
		{
			"dollars without a rate: only the pesos, marked incomplete",
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 1000000, shared.CurrencyUSD: 2000},
			false,
			Total{Amount: 1000000, IsIncomplete: true},
		},
		{
			"only pesos without a rate is complete",
			map[shared.Currency]shared.Money{shared.CurrencyCLP: 1000000},
			false,
			Total{Amount: 1000000},
		},
		{"no spend is zero", nil, false, Total{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := TotalInPesos(tc.spend, rate, tc.hasRate)

			if got != tc.want {
				t.Errorf("TotalInPesos = %+v, want %+v", got, tc.want)
			}
		})
	}
}
