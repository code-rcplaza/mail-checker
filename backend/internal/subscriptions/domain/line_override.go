package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	KeyOverrideAmountNotPositive shared.MessageKey = "subscriptions.override_amount_not_positive"
	KeyOverrideDayOutOfRange     shared.MessageKey = "subscriptions.override_day_out_of_range"
)

// LineOverride replaces what is shown for a line, never its real charges. A
// zero amount or day means that part is not replaced.
type LineOverride struct {
	Amount     shared.Money
	BillingDay int
}

// OverriddenLine is the override the user set on one line.
type OverriddenLine struct {
	Key      GroupKey
	Override LineOverride
}

// NewLineOverride takes what the user sent; nil means "not given", and
// neither given means removing the override, like an empty alias.
func NewLineOverride(amount *shared.Money, billingDay *int) (LineOverride, error) {
	if amount != nil && *amount <= 0 {
		return LineOverride{}, shared.Invalid(KeyOverrideAmountNotPositive, nil)
	}
	if billingDay != nil && (*billingDay < firstBillingDay || *billingDay > lastBillingDay) {
		return LineOverride{}, shared.Invalid(KeyOverrideDayOutOfRange, nil)
	}
	var override LineOverride
	if amount != nil {
		override.Amount = *amount
	}
	if billingDay != nil {
		override.BillingDay = *billingDay
	}
	return override, nil
}

func (o LineOverride) IsEmpty() bool {
	return o == LineOverride{}
}

// ApplyLineOverrides shows each overridden line with its amount and billing
// day, and recounts what depends on them; the given summary is left untouched.
func ApplyLineOverrides(summary Summary, overrides map[GroupKey]LineOverride) Summary {
	summary = changeEveryLine(summary, func(line Subscription) Subscription {
		return withOverride(line, overrides[line.Key])
	})
	return recountSpend(summary)
}

func withOverride(line Subscription, override LineOverride) Subscription {
	if override.Amount > 0 {
		line.Amount = override.Amount
	}
	if override.BillingDay > 0 {
		line.NextCharge = nextChargeOnDay(line, override.BillingDay)
	}
	return line
}

// nextChargeOnDay keeps a line without a next charge (no cadence, inactive or
// no charges) as it is.
func nextChargeOnDay(line Subscription, billingDay int) time.Time {
	if line.NextCharge.IsZero() || len(line.Charges) == 0 {
		return line.NextCharge
	}
	next, hasNext := NextCharge(lastOf(line.Charges).OccurredAt, LineCycle{Cadence: line.Cadence, Day: billingDay})
	if !hasNext {
		return line.NextCharge
	}
	return next
}
