package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	KeyNameRequired         shared.MessageKey = "auth.name_required"
	KeyNameTooLong          shared.MessageKey = "auth.name_too_long"
	KeyEmailInvalid         shared.MessageKey = "auth.email_invalid"
	KeyTermsRequired        shared.MessageKey = "auth.terms_required"
	KeyCodeInvalid          shared.MessageKey = "auth.code_invalid"
	KeyCodeExpired          shared.MessageKey = "auth.code_expired"
	KeyCodeAttemptsExceeded shared.MessageKey = "auth.code_attempts_exceeded"
	KeyTooManySends         shared.MessageKey = "auth.too_many_sends"
	KeyCodeFormatInvalid    shared.MessageKey = "auth.code_format_invalid"
	KeyExportFormatUnknown  shared.MessageKey = "auth.export_format_unknown"
)
