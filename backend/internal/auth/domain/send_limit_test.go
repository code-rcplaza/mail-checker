package domain

import (
	"testing"
	"time"
)

func TestCanSendCode(t *testing.T) {
	minutesAgo := func(minutes int) time.Time { return testNow.Add(-time.Duration(minutes) * time.Minute) }
	cases := []struct {
		name  string
		sends []time.Time
		want  bool
	}{
		{"two sends in the last hour allow a third", []time.Time{minutesAgo(50), minutesAgo(5)}, true},
		{"three sends in the last hour block a fourth", []time.Time{minutesAgo(50), minutesAgo(20), minutesAgo(5)}, false},
		{"a send older than an hour does not count", []time.Time{minutesAgo(61), minutesAgo(20), minutesAgo(5)}, true},
		{"a send of exactly an hour ago does not count", []time.Time{minutesAgo(60), minutesAgo(20), minutesAgo(5)}, true},
		{"no sends allow one", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CanSendCode(tc.sends, testNow); got != tc.want {
				t.Errorf("CanSendCode = %v, want %v", got, tc.want)
			}
		})
	}
}
