package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestEmailsyncKeysAndAuthError(t *testing.T) {
	cases := []struct {
		key  shared.MessageKey
		want shared.MessageKey
	}{
		{KeyMailAuthFailed, "emailsync.mail_auth_failed"},
		{KeyMailReadFailed, "emailsync.mail_read_failed"},
		{KeySyncInProgress, "emailsync.sync_in_progress"},
		{KeyMonthsOutOfRange, "emailsync.months_out_of_range"},
		{KeyLimitOutOfRange, "emailsync.limit_out_of_range"},
	}
	for _, tc := range cases {
		if tc.key != tc.want {
			t.Errorf("key = %q, want %q", tc.key, tc.want)
		}
	}
	if ErrMailAuthRejected == nil {
		t.Errorf("ErrMailAuthRejected is nil, want a sentinel error")
	}
}
