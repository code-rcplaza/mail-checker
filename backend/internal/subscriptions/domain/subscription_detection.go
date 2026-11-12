package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type Status string

const (
	StatusActive          Status = "active"
	StatusInactive        Status = "inactive"
	StatusPossible        Status = "possible"
	StatusNotSubscription Status = "not_subscription"
)

// Detection rules: a subscription needs at least two recurring charges, and
// it (or a possible one) stays active while its last charge is recent enough
// for its cadence. The yearly window is a year plus margin for a late renewal;
// without a cadence a single charge may still be yearly, so a year is allowed.
const (
	MinChargesForSubscription = 2
	ActiveWithinDaysMonthly   = 45
	ActiveWithinDaysYearly    = 400
	ActiveWithinDaysNoCadence = 365
)

type LooksLikeSubscriptionFunc func(merchant string) bool

type Detection struct {
	Status  Status
	Cadence Cadence
}

// DetectSubscription classifies a group of charges. A known service is never
// dismissed: with too little evidence it is possible and the user decides. An
// unknown merchant is at most possible. Only what may still be charging today
// is asked about: a possible outside its activity window is inactive.
func DetectSubscription(group ChargeGroup, looksLikeSubscription LooksLikeSubscriptionFunc, now time.Time) Detection {
	if group.IsKnownService {
		return detectKnownService(group.Charges, now)
	}
	return detectUnknownMerchant(group.Charges, looksLikeSubscription, now)
}

// detectKnownService looks only at recurring charges, so one-off charges of
// the service do not count towards the cadence.
func detectKnownService(charges []shared.Charge, now time.Time) Detection {
	recurring := AnalyzeKnownServiceAmount(charges).RecurringCharges
	if len(recurring) == 0 {
		return Detection{Status: StatusNotSubscription}
	}
	lastCharge := recurring[len(recurring)-1].OccurredAt
	cadence := recurringCadence(recurring)
	if cadence == CadenceNone {
		return possibleWhileRecent(lastCharge, CadenceNone, now)
	}
	return Detection{Status: activityStatus(lastCharge, cadence, now), Cadence: cadence}
}

// recurringCadence needs MinChargesForSubscription: with one charge it is
// unknown whether the service bills monthly or yearly.
func recurringCadence(recurring []shared.Charge) Cadence {
	if len(recurring) < MinChargesForSubscription {
		return CadenceNone
	}
	return DetectCadence(datesOf(recurring))
}

func detectUnknownMerchant(charges []shared.Charge, looksLikeSubscription LooksLikeSubscriptionFunc, now time.Time) Detection {
	if len(charges) == 1 {
		return detectSingleUnknownCharge(charges[0], looksLikeSubscription, now)
	}
	cadence := DetectCadence(datesOf(charges))
	if cadence == CadenceNone || !AnalyzeUnknownMerchantAmount(charges).IsConsistent {
		return Detection{Status: StatusNotSubscription}
	}
	return possibleWhileRecent(charges[len(charges)-1].OccurredAt, cadence, now)
}

func detectSingleUnknownCharge(c shared.Charge, looksLikeSubscription LooksLikeSubscriptionFunc, now time.Time) Detection {
	if !looksLikeSubscription(c.Merchant) {
		return Detection{Status: StatusNotSubscription}
	}
	return possibleWhileRecent(c.OccurredAt, CadenceNone, now)
}

func possibleWhileRecent(lastCharge time.Time, cadence Cadence, now time.Time) Detection {
	if !isRecent(lastCharge, cadence, now) {
		return Detection{Status: StatusInactive, Cadence: cadence}
	}
	return Detection{Status: StatusPossible, Cadence: cadence}
}

func activityStatus(lastCharge time.Time, cadence Cadence, now time.Time) Status {
	if !isRecent(lastCharge, cadence, now) {
		return StatusInactive
	}
	return StatusActive
}

func isRecent(lastCharge time.Time, cadence Cadence, now time.Time) bool {
	return daysBetween(lastCharge, now) <= activeWithinDays(cadence)
}

func activeWithinDays(cadence Cadence) int {
	switch cadence {
	case CadenceMonthly:
		return ActiveWithinDaysMonthly
	case CadenceYearly:
		return ActiveWithinDaysYearly
	}
	return ActiveWithinDaysNoCadence
}

func datesOf(charges []shared.Charge) []time.Time {
	dates := make([]time.Time, 0, len(charges))
	for _, c := range charges {
		dates = append(dates, c.OccurredAt)
	}
	return dates
}
