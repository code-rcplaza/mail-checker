package domain

import (
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// LineDetailCharges is how many recent charges the detail of a line lists.
const LineDetailCharges = 6

// RiseWithImpact is the rise of a line with what it costs.
type RiseWithImpact struct {
	Rise   Rise
	Impact Impact
}

type LineDetail struct {
	Cycle       LineCycle
	LastCharges []shared.Charge // newest first
	Rise        *RiseWithImpact // nil when the line has no rise
}

func BuildLineDetail(line Subscription) LineDetail {
	return LineDetail{
		Cycle:       lineCycleOf(line),
		LastCharges: lastChargesNewestFirst(line.Charges),
		Rise:        riseWithImpactOf(line),
	}
}

// lineCycleOf takes a yearly line's day and month from its last charge, since
// one charge a year has no typical day. A line without charges (a manual one)
// takes the day of its next charge.
func lineCycleOf(line Subscription) LineCycle {
	if len(line.Charges) == 0 {
		return LineCycle{Cadence: line.Cadence, Day: line.NextCharge.Day()}
	}
	if line.Cadence == CadenceYearly {
		last := lastOf(line.Charges).OccurredAt
		return LineCycle{Cadence: line.Cadence, Day: last.Day(), Month: last.Month()}
	}
	return LineCycle{Cadence: line.Cadence, Day: medianDayOfMonth(line.Charges)}
}

// lastChargesNewestFirst expects the charges oldest first, as lines keep them.
func lastChargesNewestFirst(charges []shared.Charge) []shared.Charge {
	newest := slices.Clone(charges[max(0, len(charges)-LineDetailCharges):])
	slices.Reverse(newest)
	return newest
}

func riseWithImpactOf(line Subscription) *RiseWithImpact {
	rise, hasRise := DetectRise(line.Charges)
	if !hasRise {
		return nil
	}
	return &RiseWithImpact{Rise: rise, Impact: RiseImpact(rise, line.Cadence)}
}
