package domain

import "time"

// CleanupCutoff says what no longer serves: codes and sessions that expired by
// ExpiredBy, and sends at or before SentBy, which the limit no longer counts.
type CleanupCutoff struct {
	ExpiredBy time.Time
	SentBy    time.Time
}

// CleanupCutoffAt uses the same edges as validating: expired at the instant
// itself, and a send of exactly an hour ago already out of the limit.
func CleanupCutoffAt(now time.Time) CleanupCutoff {
	return CleanupCutoff{ExpiredBy: now, SentBy: now.Add(-SendWindow)}
}
