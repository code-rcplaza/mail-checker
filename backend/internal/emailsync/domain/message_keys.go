package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	KeyMailAuthFailed   shared.MessageKey = "emailsync.mail_auth_failed"
	KeyMailReadFailed   shared.MessageKey = "emailsync.mail_read_failed"
	KeySyncInProgress   shared.MessageKey = "emailsync.sync_in_progress"
	KeyMonthsOutOfRange shared.MessageKey = "emailsync.months_out_of_range"
	KeyLimitOutOfRange  shared.MessageKey = "emailsync.limit_out_of_range"
)
