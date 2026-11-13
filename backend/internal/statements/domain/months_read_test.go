package domain

import (
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func santiago(t *testing.T) *time.Location {
	t.Helper()
	location, err := time.LoadLocation("America/Santiago")
	if err != nil {
		t.Fatalf("load America/Santiago: %v", err)
	}
	return location
}

func chargeAt(occurredAt time.Time) shared.Charge {
	return shared.Charge{OccurredAt: occurredAt}
}

func TestMonthsRead(t *testing.T) {
	location := santiago(t)
	noon := func(month time.Month, day int) time.Time {
		return time.Date(2026, month, day, 12, 0, 0, 0, location)
	}
	cases := []struct {
		name    string
		charges []shared.Charge
		want    []Month
	}{
		{
			"two months come out in order",
			[]shared.Charge{chargeAt(noon(time.September, 3)), chargeAt(noon(time.August, 20))},
			[]Month{{2026, time.August}, {2026, time.September}},
		},
		{
			"a month with several charges appears once",
			[]shared.Charge{chargeAt(noon(time.August, 1)), chargeAt(noon(time.August, 15)), chargeAt(noon(time.August, 31))},
			[]Month{{2026, time.August}},
		},
		{
			"same month of different years are different months",
			[]shared.Charge{chargeAt(noon(time.January, 5)), chargeAt(time.Date(2025, time.December, 28, 12, 0, 0, 0, location))},
			[]Month{{2025, time.December}, {2026, time.January}},
		},
		{"no charges give an empty list", nil, []Month{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MonthsRead(tc.charges, location)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("MonthsRead = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMonthsRead_CountsTheMonthInSantiagoTime(t *testing.T) {
	location := santiago(t)
	lastNightOfAugust := time.Date(2026, time.August, 31, 23, 30, 0, 0, location).UTC()
	if lastNightOfAugust.Month() != time.September {
		t.Fatalf("fixture: %v is not already September in UTC", lastNightOfAugust)
	}

	got := MonthsRead([]shared.Charge{chargeAt(lastNightOfAugust)}, location)

	want := []Month{{2026, time.August}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MonthsRead = %v, want %v", got, want)
	}
}
