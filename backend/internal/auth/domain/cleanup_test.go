package domain

import (
	"testing"
	"time"
)

func TestCleanupCutoffAt(t *testing.T) {
	got := CleanupCutoffAt(testNow)

	want := CleanupCutoff{ExpiredBy: testNow, SentBy: testNow.Add(-time.Hour)}
	if got != want {
		t.Errorf("CleanupCutoffAt = %+v, want %+v", got, want)
	}
}
