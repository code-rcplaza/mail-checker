package domain

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// NameMaxLength counts characters, not bytes: "ñ" is one.
const NameMaxLength = 60

type User struct {
	ID              int64
	PublicID        string // UUID v7, the one exposed
	Name            string
	Email           string // normalized by ParseEmail
	AcceptedTermsAt time.Time
}

// RegistrationForm is what the person typed to sign up.
type RegistrationForm struct {
	Name          string
	Email         string
	AcceptedTerms bool
}

// Registration is a form that passed validation, with name and email normalized.
type Registration struct {
	Name  string
	Email string
}

func ParseRegistration(form RegistrationForm) (Registration, error) {
	name, err := parseName(form.Name)
	if err != nil {
		return Registration{}, err
	}
	email, err := ParseEmail(form.Email)
	if err != nil {
		return Registration{}, err
	}
	if !form.AcceptedTerms {
		return Registration{}, shared.Invalid(KeyTermsRequired, nil)
	}
	return Registration{Name: name, Email: email}, nil
}

func parseName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", shared.Invalid(KeyNameRequired, nil)
	}
	if utf8.RuneCountInString(name) > NameMaxLength {
		return "", shared.Invalid(KeyNameTooLong, map[string]any{"max": NameMaxLength})
	}
	return name, nil
}

// ParseEmail normalizes an address so the same person is one user whatever the
// casing. It checks only "something@domain": the code that arrives is the real test.
func ParseEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if !isEmailShaped(email) {
		return "", shared.Invalid(KeyEmailInvalid, nil)
	}
	return email, nil
}

func isEmailShaped(email string) bool {
	local, domain, found := strings.Cut(email, "@")
	if !found || local == "" || domain == "" || strings.Contains(domain, "@") {
		return false
	}
	return !strings.ContainsFunc(email, unicode.IsSpace)
}
