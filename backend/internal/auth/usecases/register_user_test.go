package usecases

import (
	"context"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func registerUser(users *usersByEmail, form domain.RegistrationForm) (domain.User, error) {
	clock := &testClock{now: testNow}
	useCase := NewRegisterUser(RegisterUserDependencies{Users: users, NewPublicID: func() string { return "user-new" }, Now: clock.Now})
	return useCase.Execute(context.Background(), form)
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
