package domain

import (
	"errors"
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

func TestParseEmail_TrimsAndLowercases(t *testing.T) {
	got, err := ParseEmail("  User@Example.com ")

	if err != nil || got != "user@example.com" {
		t.Errorf("ParseEmail = %q, %v; want user@example.com, nil", got, err)
	}
}

func TestParseEmail_AcceptsAnyDomain(t *testing.T) {
	got, err := ParseEmail("user@company.cl")

	if err != nil || got != "user@company.cl" {
		t.Errorf("ParseEmail = %q, %v; want user@company.cl, nil", got, err)
	}
}

func TestParseEmail_RejectsMalformed(t *testing.T) {
	for _, raw := range []string{"", "user", "user@", "@example.com", "a@b@example.com", "us er@example.com"} {
		t.Run(raw, func(t *testing.T) {
			got, err := ParseEmail(raw)

			if got != "" {
				t.Errorf("ParseEmail(%q) = %q, want empty", raw, got)
			}
			assertInvalidWithKey(t, err, "accounts.email_invalid")
		})
	}
}
