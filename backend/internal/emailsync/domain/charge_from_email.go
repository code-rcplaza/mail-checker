package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// BankEmail is an email read from a user's account while looking for one
// bank's notices.
type BankEmail struct {
	Email     RawEmail
	UserID    int64
	AccountID int64
	Bank      shared.Bank
}

// ChargeFromEmail turns what the bank parser read into a charge: the
// Message-ID identifies it and the email date stands in when the notice has none.
func ChargeFromEmail(parsed shared.ParsedCharge, read BankEmail) (shared.Charge, error) {
	return parsed.ToCharge(shared.ChargeOrigin{
		UserID:       read.UserID,
		Bank:         read.Bank,
		Source:       shared.SourceEmail,
		ExternalID:   read.Email.MessageID,
		AccountID:    &read.AccountID,
		FallbackDate: read.Email.ReceivedAt,
	})
}
