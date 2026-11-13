package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	KeyBankNotDetected     shared.MessageKey = "statements.bank_not_detected"
	KeyBankUnknown         shared.MessageKey = "statements.bank_unknown"
	KeyAlreadyImported     shared.MessageKey = "statements.already_imported"
	KeyNoCharges           shared.MessageKey = "statements.no_charges"
	KeyFileEmpty           shared.MessageKey = "statements.file_empty"
	KeyNotFound            shared.MessageKey = "statements.not_found"
	KeyFileKindUnknown     shared.MessageKey = "statements.file_kind_unknown"
	KeyBillingDateRequired shared.MessageKey = "statements.billing_date_required"
)
