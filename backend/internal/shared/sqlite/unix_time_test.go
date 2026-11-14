package sqlite

import (
	"testing"
	"time"
)

func TestTimeFromUnix_IsUTC(t *testing.T) {
	got := TimeFromUnix(1763035215)

	if !got.Equal(time.Date(2025, 11, 13, 12, 0, 15, 0, time.UTC)) || got.Location() != time.UTC {
		t.Errorf("TimeFromUnix = %v, want 2025-11-13 12:00:15 UTC", got)
	}
}
