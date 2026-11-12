package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestAccountsKeysExist(t *testing.T) {
	cases := []struct {
		key  shared.MessageKey
		want shared.MessageKey
	}{
		{KeyEmailInvalid, "accounts.email_invalid"},
		{KeyAppPasswordInvalid, "accounts.app_password_invalid"},
		{KeyEmailTaken, "accounts.email_taken"},
		{KeyNotFound, "accounts.not_found"},
		{KeyPasswordUnreadable, "accounts.password_unreadable"},
	}
	for _, tc := range cases {
		if tc.key != tc.want {
			t.Errorf("key = %q, want %q", tc.key, tc.want)
		}
	}
}
