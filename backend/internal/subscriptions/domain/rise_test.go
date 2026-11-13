package domain

import (
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func chargeOf(amount shared.Money, currency shared.Currency, month time.Month) shared.Charge {
	return lineCharge("c-"+month.String(), amount, currency, month, 5)
}

func TestDetectRise(t *testing.T) {
	october := chargeOf(0, shared.CurrencyCLP, time.October).OccurredAt
	cases := []struct {
		name     string
		charges  []shared.Charge
		want     Rise
		wantRise bool
	}{
		{
			"a higher charge is a rise with before, after and since",
			[]shared.Charge{chargeOf(1000000, shared.CurrencyCLP, time.September), chargeOf(1005000, shared.CurrencyCLP, time.October)},
			Rise{Before: 1000000, After: 1005000, Currency: shared.CurrencyCLP, Since: october},
			true,
		},
		{
			"equal charges are not a rise",
			[]shared.Charge{chargeOf(1000000, shared.CurrencyCLP, time.September), chargeOf(1000000, shared.CurrencyCLP, time.October)},
			Rise{}, false,
		},
		{
			"a lower charge is not a rise",
			[]shared.Charge{chargeOf(1000000, shared.CurrencyCLP, time.September), chargeOf(900000, shared.CurrencyCLP, time.October)},
			Rise{}, false,
		},
		{
			"only the last two charges are compared",
			[]shared.Charge{chargeOf(900000, shared.CurrencyCLP, time.August), chargeOf(1000000, shared.CurrencyCLP, time.September), chargeOf(1000000, shared.CurrencyCLP, time.October)},
			Rise{}, false,
		},
		{
			"a change of currency is not a rise",
			[]shared.Charge{chargeOf(2000, shared.CurrencyUSD, time.September), chargeOf(1900000, shared.CurrencyCLP, time.October)},
			Rise{}, false,
		},
		{
			"a single charge is not a rise",
			[]shared.Charge{chargeOf(1000000, shared.CurrencyCLP, time.October)},
			Rise{}, false,
		},
		{
			"a jump of 50 % is a plan change",
			[]shared.Charge{chargeOf(2000, shared.CurrencyUSD, time.September), chargeOf(3000, shared.CurrencyUSD, time.October)},
			Rise{Before: 2000, After: 3000, Currency: shared.CurrencyUSD, Since: october, IsPlanChange: true},
			true,
		},
		{
			"exactly 25 % is still a plain rise: only above the threshold is a plan change",
			[]shared.Charge{chargeOf(2000, shared.CurrencyUSD, time.September), chargeOf(2500, shared.CurrencyUSD, time.October)},
			Rise{Before: 2000, After: 2500, Currency: shared.CurrencyUSD, Since: october},
			true,
		},
		{
			"a small rise is not a plan change",
			[]shared.Charge{chargeOf(10000, shared.CurrencyUSD, time.September), chargeOf(10050, shared.CurrencyUSD, time.October)},
			Rise{Before: 10000, After: 10050, Currency: shared.CurrencyUSD, Since: october},
			true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasRise := DetectRise(tc.charges)

			if hasRise != tc.wantRise || got != tc.want {
				t.Errorf("DetectRise = (%+v, %v), want (%+v, %v)", got, hasRise, tc.want, tc.wantRise)
			}
		})
	}
}

func risingCharges(prefix, merchant string, before, after shared.Money, dayOfMonth int) []shared.Charge {
	return slices.Concat(
		scenarioCharges(prefix, merchant, nil, before, shared.CurrencyCLP, dayOfMonth, time.August, time.September),
		scenarioCharges(prefix, merchant, nil, after, shared.CurrencyCLP, dayOfMonth, time.October),
	)
}

func TestBuildSummary_CarriesTheRisesWithTheirLineKey(t *testing.T) {
	charges := slices.Concat(
		risingCharges("netflix", "NETFLIX", 1299000, 1349000, 25),
		risingCharges("apple", "APPLE COM BILL", 469000, 499000, 9),
		scenarioCharges("claude", "ANTHROPIC CLAUDE", nil, 2000, shared.CurrencyUSD, 20, lastMonths...),
	)

	got := summaryOfCharges(charges, nil)

	want := map[GroupKey]shared.Money{"service:Netflix||d25": 1349000, "umbrella:App Store||d09": 499000}
	if len(got.Rises) != len(want) {
		t.Fatalf("rises = %+v, want %d", got.Rises, len(want))
	}
	for _, rise := range got.Rises {
		if after, isExpected := want[rise.Key]; !isExpected || rise.After != after {
			t.Errorf("rise = %+v, want one of %v", rise, want)
		}
	}
}

func TestBuildSummary_WithoutRisesTheListIsEmpty(t *testing.T) {
	charges := scenarioCharges("netflix", "NETFLIX", nil, 1299000, shared.CurrencyCLP, 25, lastMonths...)

	got := summaryOfCharges(charges, nil)

	if len(got.Rises) != 0 {
		t.Errorf("rises = %+v, want none", got.Rises)
	}
}

func TestRiseImpact(t *testing.T) {
	cases := []struct {
		name        string
		rise        Rise
		cadence     Cadence
		wantMonthly shared.Money
		wantYearly  shared.Money
	}{
		{"monthly line: the difference every month", Rise{Before: 1000000, After: 1050000}, CadenceMonthly, 50000, 600000},
		{"yearly line: the difference once a year, a twelfth per month rounded half up", Rise{Before: 12000000, After: 13000000}, CadenceYearly, 83333, 1000000},
		{"yearly line: half a hundredth rounds up", Rise{Before: 0, After: 6}, CadenceYearly, 1, 6},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := RiseImpact(tc.rise, tc.cadence)

			if got.Monthly != tc.wantMonthly || got.Yearly != tc.wantYearly {
				t.Errorf("RiseImpact = %+v, want monthly %d and yearly %d", got, tc.wantMonthly, tc.wantYearly)
			}
		})
	}
}
