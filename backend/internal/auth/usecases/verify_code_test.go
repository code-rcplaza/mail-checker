package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const anasCode = "042917"

// verifyFixture: Ana has the code 042917, issued at testNow.
type verifyFixture struct {
	codes *codesByUser
	clock *testClock
}

func newVerifyFixture() verifyFixture {
	codes := newCodesByUser()
	codes.codes[ana.ID] = domain.LoginCode{UserID: ana.ID, CodeHash: domain.HashSecret(anasCode), ExpiresAt: testNow.Add(domain.CodeLifetime)}
	return verifyFixture{codes: codes, clock: &testClock{now: testNow}}
}

func (f verifyFixture) verifyCode() VerifyCode {
	return NewVerifyCode(VerifyCodeDependencies{Users: usersWith(ana), Codes: f.codes, Now: f.clock.Now})
}

func (f verifyFixture) verify(entry CodeEntry) (domain.User, error) {
	return f.verifyCode().Execute(context.Background(), entry)
}

func TestVerifyCode_RightCodeReturnsTheUserAndIsConsumed(t *testing.T) {
	fixture := newVerifyFixture()

	got, err := fixture.verify(CodeEntry{Email: " ANA@Example.com", Code: anasCode})

	if err != nil || got != ana {
		t.Fatalf("Execute = %+v, %v; want %+v", got, err, ana)
	}
	_, err = fixture.verify(CodeEntry{Email: ana.Email, Code: anasCode})
	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCodeInvalid)
}

func TestVerifyCode_WrongCodeCountsAnAttempt(t *testing.T) {
	fixture := newVerifyFixture()

	_, err := fixture.verify(CodeEntry{Email: ana.Email, Code: "000000"})

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCodeInvalid)
	if got := fixture.codes.codes[ana.ID].Attempts; got != 1 {
		t.Errorf("attempts = %d, want 1", got)
	}
}

func TestVerifyCode_FifthWrongAttemptInvalidatesTheCode(t *testing.T) {
	fixture := newVerifyFixture()
	for range domain.CodeMaxAttempts - 1 {
		_, _ = fixture.verify(CodeEntry{Email: ana.Email, Code: "000000"})
	}

	_, err := fixture.verify(CodeEntry{Email: ana.Email, Code: "000000"})
	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCodeAttemptsExceeded)

	_, err = fixture.verify(CodeEntry{Email: ana.Email, Code: anasCode})
	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCodeInvalid)
}

func TestVerifyCode_ExpiredCode(t *testing.T) {
	cases := []struct {
		name  string
		after time.Duration
		want  shared.MessageKey
	}{
		{"more than ten minutes is expired", domain.CodeLifetime + time.Second, domain.KeyCodeExpired},
		{"exactly ten minutes is already expired", domain.CodeLifetime, domain.KeyCodeExpired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newVerifyFixture()
			fixture.clock.advance(tc.after)

			_, err := fixture.verify(CodeEntry{Email: ana.Email, Code: anasCode})

			assertKeyedError(t, err, shared.ErrInvalid, tc.want)
		})
	}
}

func TestVerifyCode_EmailWithoutCodeDoesNotRevealIfItExists(t *testing.T) {
	for _, email := range []string{"nadie@example.com", ana.Email} {
		t.Run(email, func(t *testing.T) {
			fixture := newVerifyFixture()
			delete(fixture.codes.codes, ana.ID)

			_, err := fixture.verify(CodeEntry{Email: email, Code: anasCode})

			assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCodeInvalid)
		})
	}
}
