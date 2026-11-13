package domain

import (
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Rise is a charge higher than the previous one of its line.
type Rise struct {
	Key          GroupKey // the line; set when the summary is built
	Before       shared.Money
	After        shared.Money
	Currency     shared.Currency
	Since        time.Time // the charge with the new amount
	IsPlanChange bool
}

// Impact is what a rise costs per month and per year.
type Impact struct {
	Monthly shared.Money
	Yearly  shared.Money
}

// DetectRise compares the last two charges of a line (ascending by date):
// any increase in the same currency is a rise. A change of currency is a plan
// change, not a rise. A jump above PlanChangeThresholdPercent is also marked
// as a plan change; exactly the threshold is still a plain rise.
func DetectRise(charges []shared.Charge) (Rise, bool) {
	if len(charges) < 2 {
		return Rise{}, false
	}
	previous, last := charges[len(charges)-2], charges[len(charges)-1]
	if previous.Currency != last.Currency || last.Amount <= previous.Amount {
		return Rise{}, false
	}
	return Rise{
		Before:       previous.Amount,
		After:        last.Amount,
		Currency:     last.Currency,
		Since:        last.OccurredAt,
		IsPlanChange: isAmountJump(previous.Amount, last.Amount),
	}, true
}

// RiseImpact: a monthly line pays the difference twelve times a year; a
// yearly one pays it once, and its monthly share is a twelfth.
func RiseImpact(rise Rise, cadence Cadence) Impact {
	difference := rise.After - rise.Before
	if cadence == CadenceYearly {
		return Impact{Monthly: divideRoundingHalfUp(difference, monthsPerYear), Yearly: difference}
	}
	return Impact{Monthly: difference, Yearly: difference * monthsPerYear}
}

// risesOf gives the rise of each subscription that has one, with its line key.
func risesOf(subscriptions []Subscription) []Rise {
	var rises []Rise
	for _, s := range subscriptions {
		rise, hasRise := DetectRise(s.Charges)
		if !hasRise {
			continue
		}
		rise.Key = s.Key
		rises = append(rises, rise)
	}
	return rises
}

// LineRise finds the rise of the line with key, whatever section it is in;
// a one-off group is not a line.
func LineRise(groups []ChargeGroup, key GroupKey) (Rise, bool) {
	for _, group := range groups {
		if group.Key != key || group.IsOneOff {
			continue
		}
		rise, hasRise := DetectRise(group.Charges)
		rise.Key = key
		return rise, hasRise
	}
	return Rise{}, false
}

// HideKeptRises drops the rises the user kept at that same new amount; a
// later rise above it shows again.
func HideKeptRises(summary Summary, kept map[GroupKey]shared.Money) Summary {
	summary.Rises = slices.DeleteFunc(slices.Clone(summary.Rises), func(rise Rise) bool {
		keptAmount, isKept := kept[rise.Key]
		return isKept && keptAmount == rise.After
	})
	return summary
}
