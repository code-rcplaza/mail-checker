package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	phoneToken  = "phone-token"
	laptopToken = "laptop-token"
)

// sessionsOfAna: Ana is signed in on her phone and her laptop until 30 days from testNow.
func sessionsOfAna() sessionsByHash {
	expiresAt := testNow.Add(domain.SessionLifetime)
	return sessionsByHash{
		domain.HashSecret(phoneToken):  {UserID: ana.ID, TokenHash: domain.HashSecret(phoneToken), ExpiresAt: expiresAt},
		domain.HashSecret(laptopToken): {UserID: ana.ID, TokenHash: domain.HashSecret(laptopToken), ExpiresAt: expiresAt},
	}
}

func validate(sessions sessionsByHash, clock *testClock, token string) (int64, error) {
	return NewValidateSession(sessions, clock.Now).Execute(context.Background(), token)
}

func TestValidateSession(t *testing.T) {
	cases := []struct {
		name    string
		after   time.Duration
		token   string
		want    int64
		wantErr bool
	}{
		{"a valid token returns its user", 0, phoneToken, ana.ID, false},
		{"an expired token is unauthorized", domain.SessionLifetime, phoneToken, 0, true},
		{"an unknown token is unauthorized", 0, "stolen-guess", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clock := &testClock{now: testNow}
			clock.advance(tc.after)

			got, err := validate(sessionsOfAna(), clock, tc.token)

			if got != tc.want {
				t.Errorf("userID = %d, want %d", got, tc.want)
			}
			if tc.wantErr {
				assertKeyedError(t, err, shared.ErrUnauthorized, shared.KeyUnauthorized)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Execute error = %v, want nil", err)
			}
		})
	}
}

func TestCloseSession(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{"an open session", phoneToken},
		{"an unknown session", "never-opened"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sessions := sessionsOfAna()
			clock := &testClock{now: testNow}

			err := NewCloseSession(sessions).Execute(context.Background(), tc.token)

			if err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}
			if _, err := validate(sessions, clock, tc.token); err == nil {
				t.Errorf("token %q still validates after closing", tc.token)
			}
			if got, err := validate(sessions, clock, laptopToken); err != nil || got != ana.ID {
				t.Errorf("laptop session = %d, %v; want it still valid", got, err)
			}
		})
	}
}

func TestCloseSession_ClosingTwiceDoesNotFail(t *testing.T) {
	sessions := sessionsOfAna()
	closeSession := NewCloseSession(sessions)

	for range 2 {
		if err := closeSession.Execute(context.Background(), phoneToken); err != nil {
			t.Fatalf("Execute error = %v, want nil", err)
		}
	}
}
