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

// SubscriptionDetector classifies groups of charges as of Now, asking the
// catalog whether a single unknown charge reads like a subscription.
type SubscriptionDetector struct {
	LooksLikeSubscription LooksLikeSubscriptionFunc
	Now                   time.Time
}

// Detect classifies a group of charges. A known service is never dismissed:
// with too little evidence it is possible and the user decides. An unknown
// merchant is at most possible. Only what may still be charging today is
// asked about: a possible outside its activity window is inactive.
func (d SubscriptionDetector) Detect(group ChargeGroup) Detection {
	if group.IsKnownService {
		return d.detectKnownService(group.Charges)
	}
	return d.detectUnknownMerchant(group.Charges)
}

// detectKnownService looks only at recurring charges, so one-off charges of
// the service do not count towards the cadence.
func (d SubscriptionDetector) detectKnownService(charges []shared.Charge) Detection {
	recurring := AnalyzeKnownServiceAmount(charges).RecurringCharges
	if len(recurring) == 0 {
		return Detection{Status: StatusNotSubscription}
	}
	lastCharge := recurring[len(recurring)-1].OccurredAt
	cadence := recurringCadence(recurring)
	if cadence == CadenceNone {
		return d.possibleWhileRecent(lastCharge, CadenceNone)
	}
	return Detection{Status: d.activityStatus(lastCharge, cadence), Cadence: cadence}
}

// recurringCadence needs MinChargesForSubscription: with one charge it is
// unknown whether the service bills monthly or yearly.
func recurringCadence(recurring []shared.Charge) Cadence {
	if len(recurring) < MinChargesForSubscription {
		return CadenceNone
	}
	return DetectCadence(datesOf(recurring))
}

func (d SubscriptionDetector) detectUnknownMerchant(charges []shared.Charge) Detection {
	if len(charges) == 1 {
		return d.detectSingleUnknownCharge(charges[0])
	}
	cadence := DetectCadence(datesOf(charges))
	if cadence == CadenceNone || !AnalyzeUnknownMerchantAmount(charges).IsConsistent {
		return Detection{Status: StatusNotSubscription}
	}
	return d.possibleWhileRecent(charges[len(charges)-1].OccurredAt, cadence)
}

func (d SubscriptionDetector) detectSingleUnknownCharge(c shared.Charge) Detection {
	if !d.LooksLikeSubscription(c.Merchant) {
		return Detection{Status: StatusNotSubscription}
	}
	return d.possibleWhileRecent(c.OccurredAt, CadenceNone)
}

func (d SubscriptionDetector) possibleWhileRecent(lastCharge time.Time, cadence Cadence) Detection {
	if !d.isRecent(lastCharge, cadence) {
		return Detection{Status: StatusInactive, Cadence: cadence}
	}
	return Detection{Status: StatusPossible, Cadence: cadence}
}

func (d SubscriptionDetector) activityStatus(lastCharge time.Time, cadence Cadence) Status {
	if !d.isRecent(lastCharge, cadence) {
		return StatusInactive
	}
	return StatusActive
}

func (d SubscriptionDetector) isRecent(lastCharge time.Time, cadence Cadence) bool {
	return daysBetween(lastCharge, d.Now) <= activeWithinDays(cadence)
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
