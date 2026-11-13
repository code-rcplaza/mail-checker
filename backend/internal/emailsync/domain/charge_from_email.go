package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// ChargeFromEmail turns what the bank parser read into a charge: the
// Message-ID identifies it and the email date stands in when the notice has none.
func ChargeFromEmail(parsed shared.ParsedCharge, email RawEmail, userID, accountID int64, bank shared.Bank) (shared.Charge, error) {
	return parsed.ToCharge(shared.ChargeOrigin{
		UserID:       userID,
		Bank:         bank,
		Source:       shared.SourceEmail,
		ExternalID:   email.MessageID,
		AccountID:    &accountID,
		FallbackDate: email.ReceivedAt,
	})
}
