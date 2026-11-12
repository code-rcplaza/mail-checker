package usecases

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
)

// SyncAll reads the purchase emails of every account, for every registered
// bank, and stores them as charges.
type SyncAll struct {
	accounts ports.AccountDirectory
	mail     ports.MailReader
	parsers  []ports.BankEmailParser
	charges  ports.ChargeWriter
	now      func() time.Time
}

// NewSyncAll copies the bank registry so it stays read-only after main builds it.
func NewSyncAll(
	accounts ports.AccountDirectory,
	mail ports.MailReader,
	parsers []ports.BankEmailParser,
	charges ports.ChargeWriter,
	now func() time.Time,
) SyncAll {
	return SyncAll{accounts: accounts, mail: mail, parsers: slices.Clone(parsers), charges: charges, now: now}
}

func (u SyncAll) Execute(ctx context.Context, months int) ([]domain.SyncResult, error) {
	credentials, err := u.readCredentials(ctx)
	if err != nil {
		return nil, err
	}
	since := domain.SyncSince(u.now(), months)
	results := make([]domain.SyncResult, 0, len(credentials))
	for _, credential := range credentials {
		result, err := u.syncAccount(ctx, credential, since)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (u SyncAll) readCredentials(ctx context.Context) ([]domain.Credential, error) {
	credentials, err := u.accounts.Credentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("read account credentials: %w", err)
	}
	return credentials, nil
}

func (u SyncAll) syncAccount(ctx context.Context, credential domain.Credential, since time.Time) (domain.SyncResult, error) {
	result := domain.SyncResult{AccountPublicID: credential.AccountPublicID, Email: credential.Email}
	for _, parser := range u.parsers {
		found, inserted, err := u.syncBank(ctx, credential, parser, since)
		if err != nil {
			return domain.SyncResult{}, err
		}
		result.Found += found
		result.New += inserted
	}
	return result, nil
}

// syncBank returns how many emails it read and how many became new charges.
func (u SyncAll) syncBank(ctx context.Context, credential domain.Credential, parser ports.BankEmailParser, since time.Time) (found, inserted int, err error) {
	emails, err := u.readBankEmails(ctx, credential, parser, since)
	if err != nil {
		return 0, 0, err
	}
	for _, email := range emails {
		isNew, err := u.saveRecognized(ctx, credential, parser, email)
		if err != nil {
			return 0, 0, err
		}
		if isNew {
			inserted++
		}
	}
	return len(emails), inserted, nil
}

func (u SyncAll) readBankEmails(ctx context.Context, credential domain.Credential, parser ports.BankEmailParser, since time.Time) ([]domain.RawEmail, error) {
	query := parser.Query()
	query.Since = since
	emails, err := u.mail.Read(ctx, credential, query)
	if err != nil {
		return nil, fmt.Errorf("read %s emails of account %s: %w", parser.Bank(), credential.AccountPublicID, err)
	}
	return emails, nil
}

// saveRecognized stores the email as a charge when its bank's parser reads it.
// Error messages name the account by public id, never by address.
func (u SyncAll) saveRecognized(ctx context.Context, credential domain.Credential, parser ports.BankEmailParser, email domain.RawEmail) (bool, error) {
	parsed, isRecognized := parser.Parse(email)
	if !isRecognized {
		return false, nil
	}
	charge, err := domain.ChargeFromEmail(parsed, email, credential.AccountID, parser.Bank())
	if err != nil {
		return false, fmt.Errorf("convert %s email of account %s: %w", parser.Bank(), credential.AccountPublicID, err)
	}
	isNew, err := u.charges.Insert(ctx, charge)
	if err != nil {
		return false, fmt.Errorf("insert charge of account %s: %w", credential.AccountPublicID, err)
	}
	return isNew, nil
}
