package domain

import "time"

// CancellationMark is a line the user marked "to cancel", and since when.
type CancellationMark struct {
	Key      GroupKey
	MarkedAt time.Time
}

// ApplyCancellations shows on each marked line the date it was marked "to
// cancel"; the given summary is left untouched. A marked line still adds to
// the spend: the app cannot know whether it was cancelled until it stops
// being charged.
func ApplyCancellations(summary Summary, marks map[GroupKey]time.Time) Summary {
	return changeEveryLine(summary, func(line Subscription) Subscription {
		line.MarkedForCancellationAt = marks[line.Key]
		return line
	})
}
