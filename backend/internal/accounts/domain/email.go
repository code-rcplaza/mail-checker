package domain

import (
	"strings"
	"unicode"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// ParseEmail normalizes an address so the same mailbox is not added twice with
// different casing. Any domain is accepted: Gmail and Google Workspace both work.
func ParseEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if !isEmailShaped(email) {
		return "", shared.Invalid(KeyEmailInvalid, nil)
	}
	return email, nil
}

// isEmailShaped checks only "something@domain": the IMAP login is the real test.
func isEmailShaped(email string) bool {
	local, domain, found := strings.Cut(email, "@")
	if !found || local == "" || domain == "" || strings.Contains(domain, "@") {
		return false
	}
	return !strings.ContainsFunc(email, unicode.IsSpace)
}
