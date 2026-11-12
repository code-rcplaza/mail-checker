package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// SyncSince is how far back the sync searches: months before now.
func SyncSince(now time.Time, months int) time.Time {
	return now.AddDate(0, -months, 0)
}

// The period the sync can search, in months back from today; the same range
// as the summary.
const (
	MinSyncMonths = 1
	MaxSyncMonths = 24
)

func ValidateSyncMonths(months int) error {
	if months < MinSyncMonths || months > MaxSyncMonths {
		return shared.Invalid(KeyMonthsOutOfRange, nil)
	}
	return nil
}
