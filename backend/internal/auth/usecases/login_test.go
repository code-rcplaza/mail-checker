package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func login(fixture deliveryFixture, users *usersByEmail, email string) error {
	return NewLogin(users, fixture.delivery()).Execute(context.Background(), email)
}

func TestLogin_SendsACodeToARegisteredEmail(t *testing.T) {
	for _, email := range []string{"ana@example.com", " ANA@Example.com "} {
		t.Run(email, func(t *testing.T) {
			fixture := newDeliveryFixture()

			err := login(fixture, usersWith(ana), email)

			sent := fixture.sender.sent
			if err != nil || len(sent) != 1 || sent[0].email != ana.Email || !fixture.codes.codes[ana.ID].Matches(sent[0].code) {
				t.Errorf("sent = %+v, err = %v; want the stored code sent to %s", sent, err, ana.Email)
			}
		})
	}
}

func TestLogin_InvalidEmailIsRejectedWithoutTouchingTheRepository(t *testing.T) {
	fixture := newDeliveryFixture()
	users := usersWith(ana)

	err := login(fixture, users, "ana@")

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyEmailInvalid)
	if users.calls != 0 || len(fixture.sends) != 0 {
		t.Errorf("repository calls = %d, sends = %v; want none", users.calls, fixture.sends)
	}
}

func TestLogin_UnknownEmailAnswersLikeAKnownOneAndSendsNothing(t *testing.T) {
	fixture := newDeliveryFixture()

	err := login(fixture, usersWith(ana), "nadie@example.com")

	if err != nil || len(fixture.sender.sent) != 0 {
		t.Errorf("err = %v, sent = %+v; want no error and nothing sent", err, fixture.sender.sent)
	}
	if got := len(fixture.sends["nadie@example.com"]); got != 1 {
		t.Errorf("requests counted = %d, want 1 (the limit cannot reveal the email)", got)
	}
}

func TestLogin_OverTheLimitSendsNothing(t *testing.T) {
	for _, email := range []string{ana.Email, "nadie@example.com"} {
		t.Run(email, func(t *testing.T) {
			fixture := newDeliveryFixture()
			fixture.sends[email] = []time.Time{testNow.Add(-30 * time.Minute), testNow.Add(-20 * time.Minute), testNow.Add(-10 * time.Minute)}

			err := login(fixture, usersWith(ana), email)

			assertKeyedError(t, err, shared.ErrTooMany, domain.KeyTooManySends)
			if len(fixture.sender.sent) != 0 {
				t.Errorf("sent = %+v, want nothing", fixture.sender.sent)
			}
		})
	}
}

func TestLogin_WrapsSendError(t *testing.T) {
	errProvider := errors.New("provider down")
	fixture := newDeliveryFixture()
	fixture.sender.err = errProvider

	err := login(fixture, usersWith(ana), ana.Email)

	assertWrapped(t, err, errProvider)
}
