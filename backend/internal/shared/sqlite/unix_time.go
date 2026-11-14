package sqlite

import "time"

// TimeFromUnix reads a stored date: unix seconds are always UTC.
func TimeFromUnix(seconds int64) time.Time { return time.Unix(seconds, 0).UTC() }
