package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// The period the user can ask about, in months back from today.
const (
	MinPeriodMonths = 1
	MaxPeriodMonths = 24
)

const KeyMonthsOutOfRange shared.MessageKey = "subscriptions.months_out_of_range"

func ValidatePeriodMonths(months int) error {
	if months < MinPeriodMonths || months > MaxPeriodMonths {
		return shared.Invalid(KeyMonthsOutOfRange, nil)
	}
	return nil
}

func PeriodStart(now time.Time, months int) time.Time {
	return now.AddDate(0, -months, 0)
}
