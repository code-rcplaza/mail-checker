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
	name, err := ParseName(form.Name)
	if err != nil {
		return Registration{}, err
	}
	email, err := ParseEmail(form.Email)
	if err != nil {
		return Registration{}, err
	}
	if err := RequireTerms(form.AcceptedTerms); err != nil {
		return Registration{}, err
	}
	return Registration{Name: name, Email: email}, nil
}

func ParseName(raw string) (string, error) {
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

func RequireTerms(accepted bool) error {
	if !accepted {
		return shared.Invalid(KeyTermsRequired, nil)
	}
	return nil
}

// Initials are the first letter of the first two words, in capitals, for the
// avatar. They work on characters, so an accent or an emoji is never split.
func Initials(name string) string {
	var initials strings.Builder
	for _, word := range firstWords(name, 2) {
		first, _ := utf8.DecodeRuneInString(word)
		initials.WriteRune(unicode.ToUpper(first))
	}
	return initials.String()
}

func firstWords(text string, count int) []string {
	words := strings.Fields(text)
	return words[:min(count, len(words))]
}
