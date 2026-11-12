package clock

import (
	"testing"
	"time"
)

func TestNow_ReturnsCurrentTimeInUTC(t *testing.T) {
	before := time.Now()
	got := Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Errorf("Now() = %v, want between %v and %v", got, before, after)
	}
	if got.Location() != time.UTC {
		t.Errorf("Now().Location() = %v, want UTC", got.Location())
	}
}
