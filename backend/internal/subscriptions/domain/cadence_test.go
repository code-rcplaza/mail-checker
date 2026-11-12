package domain

import (
	"slices"
	"testing"
	"time"
)

// datesAt builds dates the given number of days after a fixed start.
func datesAt(days ...int) []time.Time {
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	dates := make([]time.Time, 0, len(days))
	for _, d := range days {
		dates = append(dates, start.AddDate(0, 0, d))
	}
	return dates
}

func TestDetectCadence_FindsCadence(t *testing.T) {
	cases := []struct {
		name  string
		dates []time.Time
		want  Cadence
	}{
		{"three dates 30 days apart are monthly", datesAt(0, 30, 60), CadenceMonthly},
		{"30, 61 and 30 days (one missed month) are monthly", datesAt(0, 30, 91, 121), CadenceMonthly},
		{"two dates 365 days apart are yearly", datesAt(0, 365), CadenceYearly},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DetectCadence(tc.dates); got != tc.want {
				t.Errorf("DetectCadence = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDetectCadence_FindsNoCadence(t *testing.T) {
	cases := []struct {
		name string
		// nearMiss differs only in the breaking interval and does have a
		// cadence, so a case cannot pass just because nothing is detected.
		nearMiss []time.Time
		dates    []time.Time
	}{
		{"an interval of 120 days (more than 3 months)", datesAt(0, 30, 141), datesAt(0, 30, 150)},
		{"a single date", datesAt(0, 30), datesAt(0)},
		{"two dates 10 days apart", datesAt(0, 24), datesAt(0, 10)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if DetectCadence(tc.nearMiss) == CadenceNone {
				t.Fatalf("near miss has no cadence; the case would prove nothing")
			}

			if got := DetectCadence(tc.dates); got != CadenceNone {
				t.Errorf("DetectCadence = %q, want none", got)
			}
		})
	}
}

func TestDetectCadence_UnorderedDatesGiveSameResultAndInputIsUntouched(t *testing.T) {
	dates := datesAt(60, 0, 30)
	original := slices.Clone(dates)

	got := DetectCadence(dates)

	if got != CadenceMonthly {
		t.Errorf("DetectCadence = %q, want %q", got, CadenceMonthly)
	}
	if !slices.Equal(dates, original) {
		t.Errorf("input was modified: %v, want %v", dates, original)
	}
}
