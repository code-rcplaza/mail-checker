package domain

import (
	"strings"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// appPasswordLength is the length of a Google app password without its spaces.
const appPasswordLength = 16

// ParseAppPassword accepts the password as Google shows it, in groups of four.
func ParseAppPassword(raw string) (string, error) {
	password := strings.Join(strings.Fields(raw), "")
	if !isAppPasswordShaped(password) {
		return "", shared.Invalid(KeyAppPasswordInvalid, nil)
	}
	return password, nil
}

func isAppPasswordShaped(password string) bool {
	if len(password) != appPasswordLength {
		return false
	}
	return !strings.ContainsFunc(password, func(r rune) bool { return r < 'a' || r > 'z' })
}
