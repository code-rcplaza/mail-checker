package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// The number of unparsed emails the user can ask for at once.
const (
	MinUnparsedLimit = 1
	MaxUnparsedLimit = 100
)

func ValidateUnparsedLimit(limit int) error {
	if limit < MinUnparsedLimit || limit > MaxUnparsedLimit {
		return shared.Invalid(KeyLimitOutOfRange, nil)
	}
	return nil
}

// PreviewMaxRunes caps the preview kept from an unparsed email's body.
const PreviewMaxRunes = 200

type UnparsedEmail struct {
	ExternalID string
	AccountID  int64
	Bank       shared.Bank
	Subject    string
	ReceivedAt time.Time
	Preview    string
}

// UnparsedFromEmail keeps what is needed to recognize an email the parser
// could not read; the body is cut down to a short preview.
func UnparsedFromEmail(read BankEmail) UnparsedEmail {
	return UnparsedEmail{
		ExternalID: read.Email.MessageID,
		AccountID:  read.AccountID,
		Bank:       read.Bank,
		Subject:    read.Email.Subject,
		ReceivedAt: read.Email.ReceivedAt,
		Preview:    previewOf(read.Email.Body),
	}
}

// previewOf cuts by characters, not bytes, so accents and ñ stay whole.
func previewOf(body string) string {
	characters := []rune(body)
	if len(characters) <= PreviewMaxRunes {
		return body
	}
	return string(characters[:PreviewMaxRunes])
}
