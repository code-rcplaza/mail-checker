package domain

import (
	"testing"
	"time"
)

func TestApplyCancellations(t *testing.T) {
	const marked GroupKey = "service:Netflix||d25"
	markedAt := detectionNow.AddDate(0, 0, -2)
	summary := summaryOfCharges(netflixAndAppStoreRising(25), nil)

	got := ApplyCancellations(summary, map[GroupKey]time.Time{marked: markedAt})

	for _, line := range got.Subscriptions {
		want := time.Time{}
		if line.Key == marked {
			want = markedAt
		}
		if !line.MarkedForCancellationAt.Equal(want) {
			t.Errorf("%s: marked at = %v, want %v", line.Key, line.MarkedForCancellationAt, want)
		}
	}
	if len(got.Subscriptions) != 2 {
		t.Errorf("subscriptions = %v, want both lines", keysOf(got.Subscriptions))
	}
	for _, line := range summary.Subscriptions {
		if !line.MarkedForCancellationAt.IsZero() {
			t.Errorf("the given summary changed")
		}
	}
}
