package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func registerUser(users *usersByEmail, form domain.RegistrationForm) (domain.User, error) {
	return registerWith(newDeliveryFixture(), users, form)
}

func registerWith(fixture deliveryFixture, users *usersByEmail, form domain.RegistrationForm) (domain.User, error) {
	useCase := NewRegisterUser(RegisterUserDependencies{Users: users, Delivery: fixture.delivery(), NewPublicID: func() string { return "user-new" }, Now: fixture.clock.Now})
	return useCase.Execute(context.Background(), form)
}

func TestRegisterUser_SendsACodeToTheNewOrTheExistingOwner(t *testing.T) {
	for _, users := range []*usersByEmail{usersWith(), usersWith(ana)} {
		fixture := newDeliveryFixture()

		user, err := registerWith(fixture, users, domain.RegistrationForm{Name: "Ana", Email: "ANA@example.com", AcceptedTerms: true})

		sent := fixture.sender.sent
		if err != nil || len(sent) != 1 || sent[0].email != ana.Email || !fixture.codes.codes[user.ID].Matches(sent[0].code) {
			t.Errorf("sent = %+v, err = %v; want the stored code sent to %s", sent, err, ana.Email)
		}
	}
}

func TestRegisterUser_SendFailureKeepsTheUser(t *testing.T) {
	errProvider := errors.New("provider down")
	fixture := newDeliveryFixture()
	fixture.sender.err = errProvider
	users := usersWith()

	_, err := registerWith(fixture, users, domain.RegistrationForm{Name: "Ana", Email: "ana@example.com", AcceptedTerms: true})

	assertWrapped(t, err, errProvider)
	if _, isThere := users.users["ana@example.com"]; !isThere {
		t.Errorf("users = %v, want the user kept so login can send another code", users.users)
	}
}

func TestRegisterUser_OverTheLimitSendsNothing(t *testing.T) {
	fixture := newDeliveryFixture()
	fixture.sends["ana@example.com"] = []time.Time{testNow.Add(-30 * time.Minute), testNow.Add(-20 * time.Minute), testNow.Add(-10 * time.Minute)}

	_, err := registerWith(fixture, usersWith(), domain.RegistrationForm{Name: "Ana", Email: "ana@example.com", AcceptedTerms: true})

	assertKeyedError(t, err, shared.ErrTooMany, domain.KeyTooManySends)
	if len(fixture.sender.sent) != 0 {
		t.Errorf("sent = %+v, want nothing", fixture.sender.sent)
	}
}

func TestRegisterUser_SavesANewUserWithPublicIDAndTermsDate(t *testing.T) {
	users := usersWith()

	got, err := registerUser(users, domain.RegistrationForm{Name: " Beto ", Email: "Beto@Example.com", AcceptedTerms: true})

	want := domain.User{ID: 1, PublicID: "user-new", Name: "Beto", Email: "beto@example.com", AcceptedTermsAt: testNow}
	if err != nil || got != want {
		t.Fatalf("Execute = %+v, %v; want %+v", got, err, want)
	}
	if users.users["beto@example.com"] != want {
		t.Errorf("saved = %+v, want %+v", users.users["beto@example.com"], want)
	}
}

func TestRegisterUser_InvalidFormDoesNotTouchTheRepository(t *testing.T) {
	users := usersWith()

	_, err := registerUser(users, domain.RegistrationForm{Name: "Beto", Email: "beto@example.com", AcceptedTerms: false})

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyTermsRequired)
	if users.calls != 0 {
		t.Errorf("repository calls = %d, want 0", users.calls)
	}
}

func TestRegisterUser_RegisteredEmailReturnsTheOwnerWithoutSavingAnother(t *testing.T) {
	for _, email := range []string{"ana@example.com", "ANA@Example.com"} {
		t.Run(email, func(t *testing.T) {
			users := usersWith(ana)

			got, err := registerUser(users, domain.RegistrationForm{Name: "Otra Ana", Email: email, AcceptedTerms: true})

			if err != nil || got != ana {
				t.Errorf("Execute = %+v, %v; want the existing user without error", got, err)
			}
			if len(users.users) != 1 || users.users[ana.Email] != ana {
				t.Errorf("users = %+v, want only the existing one, unchanged", users.users)
			}
		})
	}
}
