package domain

import (
	"errors"
	"strings"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func assertInvalidWithKey(t *testing.T, err error, want shared.MessageKey) {
	t.Helper()
	var keyed shared.KeyedError
	if !errors.Is(err, shared.ErrInvalid) || !errors.As(err, &keyed) || keyed.Key != want {
		t.Errorf("error = %v, want invalid with key %q", err, want)
	}
}

func TestParseRegistration_AcceptsAValidUser(t *testing.T) {
	got, err := ParseRegistration(RegistrationForm{Name: "  Ana Pérez ", Email: "  Ana@Example.COM ", AcceptedTerms: true})

	want := Registration{Name: "Ana Pérez", Email: "ana@example.com"}
	if err != nil || got != want {
		t.Errorf("ParseRegistration = %+v, %v; want %+v, nil", got, err, want)
	}
}

func TestParseRegistration_RejectsInvalidNameOrEmail(t *testing.T) {
	cases := []struct {
		name string
		form RegistrationForm
		want shared.MessageKey
	}{
		{"empty name", RegistrationForm{Name: "", Email: "ana@example.com", AcceptedTerms: true}, KeyNameRequired},
		{"blank name", RegistrationForm{Name: "   ", Email: "ana@example.com", AcceptedTerms: true}, KeyNameRequired},
		{"name longer than 60", RegistrationForm{Name: strings.Repeat("ñ", NameMaxLength+1), Email: "ana@example.com", AcceptedTerms: true}, KeyNameTooLong},
		{"malformed email", RegistrationForm{Name: "Ana", Email: "ana@", AcceptedTerms: true}, KeyEmailInvalid},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseRegistration(tc.form)

			if got != (Registration{}) {
				t.Errorf("registration = %+v, want zero", got)
			}
			assertInvalidWithKey(t, err, tc.want)
		})
	}
}

func TestParseRegistration_NameOfExactly60IsAccepted(t *testing.T) {
	name := strings.Repeat("ñ", NameMaxLength)

	got, err := ParseRegistration(RegistrationForm{Name: name, Email: "ana@example.com", AcceptedTerms: true})

	if err != nil || got.Name != name {
		t.Errorf("ParseRegistration = %+v, %v; want the 60-character name", got, err)
	}
}

func TestParseRegistration_RequiresTheTerms(t *testing.T) {
	_, err := ParseRegistration(RegistrationForm{Name: "Ana", Email: "ana@example.com", AcceptedTerms: false})

	assertInvalidWithKey(t, err, KeyTermsRequired)
}

func TestParseEmail(t *testing.T) {
	cases := []struct {
		raw     string
		want    string
		wantErr bool
	}{
		{"  User@Example.com ", "user@example.com", false},
		{"user@company.cl", "user@company.cl", false},
		{"", "", true},
		{"user", "", true},
		{"@example.com", "", true},
		{"a@b@example.com", "", true},
		{"us er@example.com", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := ParseEmail(tc.raw)

			if got != tc.want {
				t.Errorf("ParseEmail(%q) = %q, want %q", tc.raw, got, tc.want)
			}
			if tc.wantErr {
				assertInvalidWithKey(t, err, KeyEmailInvalid)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("ParseEmail(%q) error = %v, want nil", tc.raw, err)
			}
		})
	}
}
