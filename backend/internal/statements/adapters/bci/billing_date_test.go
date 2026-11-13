package bci

import (
	"strings"
	"testing"
	"time"
)

func TestBillingDate(t *testing.T) {
	parser := newParser(t)
	midnight := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, parser.location)
	}
	october := readFixture(t, "statement_october.txt")
	cases := []struct {
		name     string
		text     string
		want     time.Time
		wantDate bool
	}{
		{"the statement date at midnight in Santiago", october, midnight(2026, time.October, 21), true},
		{"another closing day", readFixture(t, "statement_september.txt"), midnight(2026, time.September, 17), true},
		{"no statement date label", withoutLinesContaining(october, "FECHA ESTADO DE CUENTA"), time.Time{}, false},
		{"an unreadable date", strings.ReplaceAll(october, "FECHA ESTADO DE CUENTA\t21/10/2026", "FECHA ESTADO DE CUENTA\t21/13/2026"), time.Time{}, false},
		{"a label without a value", "FECHA ESTADO DE CUENTA", time.Time{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasDate := parser.BillingDate(tc.text)

			if hasDate != tc.wantDate || !got.Equal(tc.want) {
				t.Errorf("BillingDate = (%v, %v), want (%v, %v)", got, hasDate, tc.want, tc.wantDate)
			}
		})
	}
}
