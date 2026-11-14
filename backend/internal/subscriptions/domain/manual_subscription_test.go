package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestNewManualSubscription_AcceptsValidData(t *testing.T) {
	got, err := NewManualSubscription("manual-1", ManualSubscriptionInput{Name: "  Gimnasio  ", Amount: 500000, Currency: "CLP", BillingDay: 15})

	want := ManualSubscription{PublicID: "manual-1", Name: "Gimnasio", Amount: 500000, Currency: shared.CurrencyCLP, BillingDay: 15}
	if err != nil || got != want {
		t.Errorf("NewManualSubscription = (%+v, %v), want (%+v, nil)", got, err, want)
	}
}

func TestNewManualSubscription_RejectsInvalidData(t *testing.T) {
	cases := []struct {
		name       string
		manualName string
		amount     shared.Money
		currency   string
		day        int
		want       shared.MessageKey
	}{
		{"empty name", "   ", 500000, "CLP", 15, KeyManualNameMissing},
		{"name too long", strings.Repeat("ñ", MaxLineAliasRunes+1), 500000, "CLP", 15, KeyManualNameTooLong},
		{"zero amount", "Gimnasio", 0, "CLP", 15, KeyManualAmountNotPositive},
		{"negative amount", "Gimnasio", -1, "CLP", 15, KeyManualAmountNotPositive},
		{"day 0", "Gimnasio", 500000, "CLP", 0, KeyManualDayOutOfRange},
		{"day 32", "Gimnasio", 500000, "CLP", 32, KeyManualDayOutOfRange},
		{"another currency", "Gimnasio", 500000, "EUR", 15, shared.KeyCurrencyUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewManualSubscription("manual-1", ManualSubscriptionInput{Name: tc.manualName, Amount: tc.amount, Currency: tc.currency, BillingDay: tc.day})

			var keyed shared.KeyedError
			if !errors.Is(err, shared.ErrInvalid) || !errors.As(err, &keyed) || keyed.Key != tc.want {
				t.Errorf("error = %v, want ErrInvalid with %q", err, tc.want)
			}
		})
	}
}

func TestNextManualCharge(t *testing.T) {
	today := func(month time.Month, dayOfMonth int) time.Time {
		return time.Date(2026, month, dayOfMonth, 15, 30, 0, 0, time.UTC)
	}
	cases := []struct {
		name       string
		billingDay int
		now        time.Time
		want       time.Time
	}{
		{"a day still to come this month", 20, today(time.November, 10), day(2026, time.November, 20)},
		{"a day already past goes to next month", 5, today(time.November, 10), day(2026, time.December, 5)},
		{"today is the next charge", 10, today(time.November, 10), day(2026, time.November, 10)},
		{"day 31 in a month of 30 days", 31, today(time.November, 10), day(2026, time.November, 30)},
		{"past day in december goes to january", 5, today(time.December, 10), day(2027, time.January, 5)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NextManualCharge(tc.billingDay, tc.now); !got.Equal(tc.want) {
				t.Errorf("NextManualCharge = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAddManualSubscriptions(t *testing.T) {
	now := time.Date(2026, time.November, 10, 12, 0, 0, 0, time.UTC)
	gym := ManualSubscription{PublicID: "manual-1", Name: "Gimnasio", Amount: 500000, Currency: shared.CurrencyCLP, BillingDay: 20}
	vpn := ManualSubscription{PublicID: "manual-2", Name: "VPN", Amount: 500, Currency: shared.CurrencyUSD, BillingDay: 3}
	detected := Subscription{Key: "service:Netflix||d25", Status: StatusActive, Cadence: CadenceMonthly, Amount: 1299000, Currency: shared.CurrencyCLP}
	summary := Summary{Subscriptions: []Subscription{detected}, MonthlySpend: MonthlySpend([]Subscription{detected}), AsOf: now}

	got := AddManualSubscriptions(summary, []ManualSubscription{gym, vpn})

	if len(got.Subscriptions) != 3 {
		t.Fatalf("subscriptions = %d, want the detected one and two manual", len(got.Subscriptions))
	}
	manual := got.Subscriptions[1]
	if manual.Key != "manual:manual-1" || manual.Alias != "Gimnasio" || manual.Status != StatusActive || manual.Cadence != CadenceMonthly ||
		!manual.NextCharge.Equal(day(2026, time.November, 20)) || len(manual.Charges) != 0 {
		t.Errorf("manual line = %+v, want an active monthly line named Gimnasio, next charge on the 20th and no charges", manual)
	}
	wantSpend := map[shared.Currency]shared.Money{shared.CurrencyCLP: 1799000, shared.CurrencyUSD: 500}
	if got.MonthlySpend[shared.CurrencyCLP] != wantSpend[shared.CurrencyCLP] || got.MonthlySpend[shared.CurrencyUSD] != wantSpend[shared.CurrencyUSD] {
		t.Errorf("monthly spend = %v, want %v", got.MonthlySpend, wantSpend)
	}
	if len(summary.Subscriptions) != 1 {
		t.Errorf("the given summary changed")
	}
}
