package domain

import "time"

const (
	SendLimit  = 3
	SendWindow = time.Hour
)

// CanSendCode counts only the sends inside the last hour; one of exactly an hour ago is out.
func CanSendCode(sends []time.Time, now time.Time) bool {
	windowStart := now.Add(-SendWindow)
	recent := 0
	for _, sentAt := range sends {
		if sentAt.After(windowStart) {
			recent++
		}
	}
	return recent < SendLimit
}

// CodeSend is one code request to an email, counted whether or not the email has an account.
type CodeSend struct {
	Email  string
	SentAt time.Time
}
