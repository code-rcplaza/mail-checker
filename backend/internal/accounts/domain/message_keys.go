package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	KeyEmailInvalid       shared.MessageKey = "accounts.email_invalid"
	KeyAppPasswordInvalid shared.MessageKey = "accounts.app_password_invalid"
	KeyEmailTaken         shared.MessageKey = "accounts.email_taken"
	KeyNotFound           shared.MessageKey = "accounts.not_found"
	KeyPasswordUnreadable shared.MessageKey = "accounts.password_unreadable"
)
