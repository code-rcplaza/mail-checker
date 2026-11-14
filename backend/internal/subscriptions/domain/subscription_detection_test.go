package domain

import (
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var detectionNow = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)

// looksLikeSubscriptionFake stands in for the catalog's descriptor check.
func detectorFake() SubscriptionDetector {
	return SubscriptionDetector{LooksLikeSubscription: looksLikeSubscriptionFake, Now: detectionNow}
}

func looksLikeSubscriptionFake(merchant string) bool {
	return strings.Contains(merchant, "SUBSCRIPTION") || strings.Contains(merchant, "MEMBERSHIP")
}

// groupOf builds a group whose charges happened the given number of days before detectionNow.
func groupOf(isKnownService bool, merchant string, amounts []shared.Money, daysAgo ...int) ChargeGroup {
	charges := make([]shared.Charge, 0, len(daysAgo))
	for i, ago := range daysAgo {
		c := charge("c", merchant, amounts[i], shared.CurrencyCLP, 1)
		c.OccurredAt = detectionNow.AddDate(0, 0, -ago)
		charges = append(charges, c)
	}
	return ChargeGroup{Key: "test", IsKnownService: isKnownService, Charges: charges}
}

func sameAmounts(amount shared.Money, count int) []shared.Money {
	amounts := make([]shared.Money, count)
	for i := range amounts {
		amounts[i] = amount
	}
	return amounts
}

func TestDetectSubscription(t *testing.T) {
	cases := []struct {
		name  string
		group ChargeGroup
		want  Detection
	}{
		{
			"known monthly service charged 10 days ago is an active monthly subscription",
			groupOf(true, "STREAMING CO", sameAmounts(599000, 3), 70, 40, 10),
			Detection{Status: StatusActive, Cadence: CadenceMonthly},
		},
		{
			"monthly service last charged 50 days ago is an inactive subscription",
			groupOf(true, "STREAMING CO", sameAmounts(599000, 3), 110, 80, 50),
			Detection{Status: StatusInactive, Cadence: CadenceMonthly},
		},
		{
			"yearly service charged 300 days ago is an active yearly subscription",
			groupOf(true, "STREAMING CO", sameAmounts(12000000, 2), 665, 300),
			Detection{Status: StatusActive, Cadence: CadenceYearly},
		},
		{
			"known service with a single charge is possible",
			groupOf(true, "STREAMING CO", sameAmounts(599000, 1), 10),
			Detection{Status: StatusPossible, Cadence: CadenceNone},
		},
		{
			"unknown merchant with two consistent monthly charges is possible",
			groupOf(false, "MUSIC APP", sameAmounts(399000, 2), 40, 10),
			Detection{Status: StatusPossible, Cadence: CadenceMonthly},
		},
		{
			"unknown merchant with one charge and a subscription descriptor is possible",
			groupOf(false, "MUSIC APP MEMBERSHIP", sameAmounts(399000, 1), 10),
			Detection{Status: StatusPossible, Cadence: CadenceNone},
		},
		{
			"unknown merchant with one charge and no subscription descriptor is not a subscription",
			groupOf(false, "CORNER SHOP", sameAmounts(800000, 1), 10),
			Detection{Status: StatusNotSubscription, Cadence: CadenceNone},
		},
		{
			"unknown merchant with inconsistent amounts is not a subscription",
			groupOf(false, "CORNER SHOP", []shared.Money{800000, 1450000, 520000}, 70, 40, 10),
			Detection{Status: StatusNotSubscription, Cadence: CadenceNone},
		},
		{
			"consistent unknown merchant last charged 200 days ago is inactive",
			groupOf(false, "MUSIC APP", sameAmounts(399000, 2), 230, 200),
			Detection{Status: StatusInactive, Cadence: CadenceMonthly},
		},
		{
			"known service with a single charge 300 days ago is still possible",
			groupOf(true, "STREAMING CO", sameAmounts(599000, 1), 300),
			Detection{Status: StatusPossible, Cadence: CadenceNone},
		},
		{
			"known service with a single charge 370 days ago is inactive",
			groupOf(true, "STREAMING CO", sameAmounts(599000, 1), 370),
			Detection{Status: StatusInactive, Cadence: CadenceNone},
		},
		{
			"unknown merchant with a subscription descriptor charged 370 days ago is inactive",
			groupOf(false, "MUSIC APP MEMBERSHIP", sameAmounts(399000, 1), 370),
			Detection{Status: StatusInactive, Cadence: CadenceNone},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := detectorFake().Detect(tc.group)

			if got != tc.want {
				t.Errorf("DetectSubscription = %+v, want %+v", got, tc.want)
			}
		})
	}
}
