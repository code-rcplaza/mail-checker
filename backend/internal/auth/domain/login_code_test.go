package domain

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"
)

var testNow = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

type failingRandom struct{ err error }

func (f failingRandom) Read([]byte) (int, error) { return 0, f.err }

func TestIssueLoginCode_HasSixDigits(t *testing.T) {
	cases := []struct {
		name   string
		random io.Reader
		want   string
	}{
		{"leading zeros are kept", bytes.NewReader(make([]byte, 8)), "000000"},
		{"the highest code", bytes.NewReader([]byte{0x0F, 0x42, 0x3F}), "999999"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IssueLoginCode(tc.random, CodeIssue{UserID: 7, Now: testNow})

			if err != nil || got.Code != tc.want {
				t.Errorf("code = %q, %v; want %q", got.Code, err, tc.want)
			}
		})
	}
}

func TestIssueLoginCode_ExpiresInTenMinutesWithoutAttempts(t *testing.T) {
	got, err := IssueLoginCode(bytes.NewReader(make([]byte, 8)), CodeIssue{UserID: 7, Now: testNow})

	want := LoginCode{UserID: 7, CodeHash: HashSecret("000000"), ExpiresAt: testNow.Add(10 * time.Minute), Attempts: 0}
	if err != nil || got.LoginCode != want {
		t.Errorf("login code = %+v, %v; want %+v", got.LoginCode, err, want)
	}
}

func TestIssueLoginCode_WrapsRandomError(t *testing.T) {
	errEntropy := errors.New("no entropy")

	got, err := IssueLoginCode(failingRandom{errEntropy}, CodeIssue{UserID: 7, Now: testNow})

	if !errors.Is(err, errEntropy) || err.Error() == errEntropy.Error() || got != (IssuedCode{}) {
		t.Errorf("IssueLoginCode = %+v, %v; want no code and the error wrapped with context", got, err)
	}
}

func TestLoginCode_Matches(t *testing.T) {
	code := LoginCode{CodeHash: HashSecret("042917")}

	for entered, want := range map[string]bool{"042917": true, "042918": false, "": false} {
		if got := code.Matches(entered); got != want {
			t.Errorf("Matches(%q) = %v, want %v", entered, got, want)
		}
	}
}

func TestLoginCode_IsExpiredAt(t *testing.T) {
	code := LoginCode{ExpiresAt: testNow.Add(CodeLifetime)}
	cases := []struct {
		name string
		at   time.Time
		want bool
	}{
		{"one second before ten minutes", testNow.Add(CodeLifetime - time.Second), false},
		{"exactly ten minutes is expired", testNow.Add(CodeLifetime), true},
		{"after ten minutes", testNow.Add(CodeLifetime + time.Second), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := code.IsExpiredAt(tc.at); got != tc.want {
				t.Errorf("IsExpiredAt = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAttemptsExhausted(t *testing.T) {
	for attempts, want := range map[int]bool{4: false, 5: true, 6: true} {
		if got := AttemptsExhausted(attempts); got != want {
			t.Errorf("AttemptsExhausted(%d) = %v, want %v", attempts, got, want)
		}
	}
}
