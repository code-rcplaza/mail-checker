// Package ports declares what the email sync use cases need from the outside.
package ports

import (
	"context"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type AccountDirectory interface {
	Credentials(ctx context.Context) ([]domain.Credential, error)
	// RecordSync stores when the account last synced and its error key (empty when fine).
	RecordSync(ctx context.Context, accountID int64, at time.Time, syncErrKey shared.MessageKey) error
}

type MailReader interface {
	Read(ctx context.Context, credential domain.Credential, query domain.MailQuery) ([]domain.RawEmail, error)
}

// BankEmailParser knows how one bank's purchase emails look.
type BankEmailParser interface {
	Bank() shared.Bank
	Query() domain.MailQuery
	Parse(email domain.RawEmail) (shared.ParsedCharge, bool)
}

// ChargeWriter inserts idempotently by source and external id.
type ChargeWriter interface {
	Insert(ctx context.Context, charge shared.Charge) (inserted bool, err error)
}

// UnparsedRepository keeps the emails no parser could read; Forget removes
// one once it is recognized.
type UnparsedRepository interface {
	Save(ctx context.Context, email domain.UnparsedEmail) error
	Forget(ctx context.Context, externalID string) error
}
