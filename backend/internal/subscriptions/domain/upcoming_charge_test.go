package domain

import (
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// upcomingOn is a subscription whose next charge is the given number of days after detectionNow.
func upcomingOn(key GroupKey, alias string, daysAhead int) Subscription {
	day := time.Date(detectionNow.Year(), detectionNow.Month(), detectionNow.Day(), 0, 0, 0, 0, time.UTC)
	return Subscription{Key: key, Parent: key, Alias: alias, Status: StatusActive, Cadence: CadenceMonthly, NextCharge: day.AddDate(0, 0, daysAhead)}
}

func TestUpcomingCharges_FromTheNearestToTheFarthest(t *testing.T) {
	cases := []struct {
		name  string
		lines []Subscription
		want  []GroupKey
	}{
		{"different dates go from the nearest", []Subscription{upcomingOn("c", "", 20), upcomingOn("a", "", 3), upcomingOn("b", "", 9)}, []GroupKey{"a", "b", "c"}},
		{"the same date is broken by name", []Subscription{upcomingOn("x", "Spotify", 5), upcomingOn("y", "Netflix", 5)}, []GroupKey{"y", "x"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := slices.Clone(tc.lines)

			got := keysOf(UpcomingCharges(tc.lines, detectionNow))

			if !slices.Equal(got, tc.want) {
				t.Errorf("UpcomingCharges = %v, want %v", got, tc.want)
			}
			if !slices.EqualFunc(tc.lines, original, func(a, b Subscription) bool { return a.Key == b.Key }) {
				t.Errorf("the given lines changed")
			}
		})
	}
}

func TestUpcomingCharges_LeavesOutWhatIsPastTheHorizon(t *testing.T) {
	withoutNextCharge := Subscription{Key: "none", Parent: "none", Status: StatusActive}
	lines := []Subscription{upcomingOn("today", "", 0), upcomingOn("last", "", UpcomingHorizonDays), upcomingOn("past", "", UpcomingHorizonDays+1), withoutNextCharge}

	got := keysOf(UpcomingCharges(lines, detectionNow))

	if want := []GroupKey{"today", "last"}; !slices.Equal(got, want) {
		t.Errorf("UpcomingCharges = %v, want %v (horizon of %d days)", got, want, UpcomingHorizonDays)
	}
}

func TestUpcomingCharges_IncludesTheManualSubscriptions(t *testing.T) {
	manual, err := NewManualSubscription("0193-manual", ManualSubscriptionInput{Name: "Gimnasio", Amount: 2500000, Currency: string(shared.CurrencyCLP), BillingDay: 20})
	if err != nil {
		t.Fatalf("NewManualSubscription error = %v", err)
	}
	summary := AddManualSubscriptions(Summary{Subscriptions: []Subscription{upcomingOn("service:Netflix||d25", "", 13)}, AsOf: detectionNow}, []ManualSubscription{manual})

	got := keysOf(UpcomingCharges(summary.Subscriptions, detectionNow))

	if want := []GroupKey{"manual:0193-manual", "service:Netflix||d25"}; !slices.Equal(got, want) {
		t.Errorf("UpcomingCharges = %v, want the manual on the 20th before Netflix on the 25th", got)
	}
}
