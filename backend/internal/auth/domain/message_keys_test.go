package domain

import (
	"regexp"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var featureReasonFormat = regexp.MustCompile(`^auth\.[a-z]+(_[a-z]+)*$`)

func TestAuthKeys_FollowFeatureReasonFormat(t *testing.T) {
	cases := []struct {
		key  shared.MessageKey
		want shared.MessageKey
	}{
		{KeyNameRequired, "auth.name_required"},
		{KeyNameTooLong, "auth.name_too_long"},
		{KeyEmailInvalid, "auth.email_invalid"},
		{KeyTermsRequired, "auth.terms_required"},
		{KeyCodeInvalid, "auth.code_invalid"},
		{KeyCodeExpired, "auth.code_expired"},
		{KeyCodeAttemptsExceeded, "auth.code_attempts_exceeded"},
		{KeyTooManySends, "auth.too_many_sends"},
	}
	for _, tc := range cases {
		t.Run(string(tc.want), func(t *testing.T) {
			if tc.key != tc.want || !featureReasonFormat.MatchString(string(tc.key)) {
				t.Errorf("key = %q, want %q in auth.<reason> format", tc.key, tc.want)
			}
		})
	}
}
