// Package ports declares what the email sync use cases need from the outside.
package ports

import (
	"context"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type AccountDirectory interface {
	Credentials(ctx context.Context) ([]domain.Credential, error)
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

type ChargeWriter interface {
	Insert(ctx context.Context, charge shared.Charge) (inserted bool, err error)
}
