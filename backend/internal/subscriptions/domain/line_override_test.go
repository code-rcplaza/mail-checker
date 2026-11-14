package domain

import (
	"errors"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func moneyPointer(m shared.Money) *shared.Money { return &m }

func intPointer(n int) *int { return &n }

func TestNewLineOverride(t *testing.T) {
	cases := []struct {
		name      string
		amount    *shared.Money
		day       *int
		want      LineOverride
		wantEmpty bool
	}{
		{"an amount above zero", moneyPointer(1000000), nil, LineOverride{Amount: 1000000}, false},
		{"a day between 1 and 31", nil, intPointer(31), LineOverride{BillingDay: 31}, false},
		{"both", moneyPointer(1000000), intPointer(1), LineOverride{Amount: 1000000, BillingDay: 1}, false},
		{"neither means no override", nil, nil, LineOverride{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewLineOverride(tc.amount, tc.day)

			if err != nil || got != tc.want || got.IsEmpty() != tc.wantEmpty {
				t.Errorf("NewLineOverride = (%+v, empty %v, %v), want (%+v, empty %v, nil)", got, got.IsEmpty(), err, tc.want, tc.wantEmpty)
			}
		})
	}
}

func TestNewLineOverride_RejectsInvalidValues(t *testing.T) {
	cases := []struct {
		name   string
		amount *shared.Money
		day    *int
		want   shared.MessageKey
	}{
		{"zero amount", moneyPointer(0), nil, KeyOverrideAmountNotPositive},
		{"negative amount", moneyPointer(-1), nil, KeyOverrideAmountNotPositive},
		{"day 0", nil, intPointer(0), KeyOverrideDayOutOfRange},
		{"day 32", nil, intPointer(32), KeyOverrideDayOutOfRange},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewLineOverride(tc.amount, tc.day)

			var keyed shared.KeyedError
			if !errors.Is(err, shared.ErrInvalid) || !errors.As(err, &keyed) || keyed.Key != tc.want {
				t.Errorf("error = %v, want ErrInvalid with %q", err, tc.want)
			}
		})
	}
}

func TestApplyLineOverrides(t *testing.T) {
	charges := scenarioCharges("netflix", "NETFLIX", nil, 1299000, shared.CurrencyCLP, 25, lastMonths...)
	original := slices.Clone(charges)
	summary := summaryOfCharges(charges, nil)
	key := summary.Subscriptions[0].Key
	overrides := map[GroupKey]LineOverride{key: {Amount: 1599000, BillingDay: 5}}

	got := ApplyLineOverrides(summary, overrides)

	line := got.Subscriptions[0]
	if line.Amount != 1599000 || got.MonthlySpend[shared.CurrencyCLP] != 1599000 {
		t.Errorf("amount = %d, spend = %v; want the override 1599000 in both", line.Amount, got.MonthlySpend)
	}
	if want := time.Date(2026, time.December, 5, 0, 0, 0, 0, time.UTC); !line.NextCharge.Equal(want) {
		t.Errorf("next charge = %v, want %v (the override day)", line.NextCharge, want)
	}
	if !slices.Equal(line.Charges, original) {
		t.Errorf("charges changed: %+v", line.Charges)
	}
	if summary.Subscriptions[0].Amount != 1299000 {
		t.Errorf("the given summary changed")
	}
}
