package domain

import (
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// lineWithCharges bills amount on day of each given month of 2026, oldest first.
func lineWithCharges(cadence Cadence, amounts []shared.Money, day int, months ...time.Month) Subscription {
	charges := make([]shared.Charge, 0, len(months))
	for i, month := range months {
		charges = append(charges, lineCharge(month.String(), amounts[i], shared.CurrencyCLP, month, day))
	}
	return Subscription{Key: "service:Netflix||d05", Status: StatusActive, Cadence: cadence, Charges: charges}
}

func sameAmount(amount shared.Money, count int) []shared.Money {
	return slices.Repeat([]shared.Money{amount}, count)
}

func TestBuildLineDetail_Cycle(t *testing.T) {
	cases := []struct {
		name string
		line Subscription
		want LineCycle
	}{
		{"a monthly line says monthly and its typical day", lineWithCharges(CadenceMonthly, sameAmount(1299000, 3), 5, time.August, time.September, time.October), LineCycle{Cadence: CadenceMonthly, Day: 5}},
		{"a yearly line says yearly with its day and month", lineWithCharges(CadenceYearly, sameAmount(9900000, 1), 14, time.March), LineCycle{Cadence: CadenceYearly, Day: 14, Month: time.March}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildLineDetail(tc.line).Cycle

			if got != tc.want {
				t.Errorf("cycle = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestBuildLineDetail_LastCharges(t *testing.T) {
	allMonths := []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October}
	cases := []struct {
		name string
		line Subscription
		want []string
	}{
		{"ten charges list the six newest, newest first", lineWithCharges(CadenceMonthly, sameAmount(1299000, 10), 5, allMonths...), []string{"October", "September", "August", "July", "June", "May"}},
		{"two charges list both", lineWithCharges(CadenceMonthly, sameAmount(1299000, 2), 5, time.September, time.October), []string{"October", "September"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := slices.Clone(tc.line.Charges)

			got := BuildLineDetail(tc.line).LastCharges

			ids := make([]string, 0, len(got))
			for _, c := range got {
				ids = append(ids, c.ExternalID)
			}
			if !slices.Equal(ids, tc.want) {
				t.Errorf("last charges = %v, want %v", ids, tc.want)
			}
			if !slices.Equal(tc.line.Charges, original) {
				t.Errorf("the charges of the line changed")
			}
		})
	}
}

func TestBuildLineDetail_Rise(t *testing.T) {
	steady := lineWithCharges(CadenceMonthly, sameAmount(1299000, 3), 5, time.August, time.September, time.October)
	rising := lineWithCharges(CadenceMonthly, []shared.Money{1299000, 1299000, 1349000}, 5, time.August, time.September, time.October)

	if got := BuildLineDetail(steady).Rise; got != nil {
		t.Errorf("rise of a steady line = %+v, want none", got)
	}
	got := BuildLineDetail(rising).Rise
	if got == nil {
		t.Fatalf("rise of a rising line = nil, want one")
	}
	since := time.Date(2026, time.October, 5, 12, 0, 0, 0, time.UTC)
	if got.Rise.Before != 1299000 || got.Rise.After != 1349000 || !got.Rise.Since.Equal(since) {
		t.Errorf("rise = %+v, want 1299000 → 1349000 since %v", got.Rise, since)
	}
	if want := (Impact{Monthly: 50000, Yearly: 600000}); got.Impact != want {
		t.Errorf("impact = %+v, want %+v", got.Impact, want)
	}
}
