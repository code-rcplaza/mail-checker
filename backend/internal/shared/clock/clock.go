// Package clock provides the real time source. Domain and use cases receive it
// injected as func() time.Time so tests can control time.
package clock

import "time"

// Now returns the current time in UTC (dates are stored as unix UTC).
func Now() time.Time { return time.Now().UTC() }
