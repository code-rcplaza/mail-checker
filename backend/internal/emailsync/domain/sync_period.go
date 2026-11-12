package domain

import "time"

// SyncSince is how far back the sync searches: months before now.
func SyncSince(now time.Time, months int) time.Time {
	return now.AddDate(0, -months, 0)
}
