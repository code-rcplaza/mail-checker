package usecases

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// SyncAll reads the purchase emails of every account, for every registered
// bank, and stores them as charges.
type SyncAll struct {
	accounts ports.AccountDirectory
	mail     ports.MailReader
	parsers  []ports.BankEmailParser
	charges  ports.ChargeWriter
	unparsed ports.UnparsedRepository
	now      func() time.Time
}

// NewSyncAll copies the bank registry so it stays read-only after main builds it.
func NewSyncAll(
	accounts ports.AccountDirectory,
	mail ports.MailReader,
	parsers []ports.BankEmailParser,
	charges ports.ChargeWriter,
	unparsed ports.UnparsedRepository,
	now func() time.Time,
) SyncAll {
	return SyncAll{accounts: accounts, mail: mail, parsers: slices.Clone(parsers), charges: charges, unparsed: unparsed, now: now}
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
		outcomes, err := u.syncBank(ctx, credential, parser, since)
		if err != nil {
			return domain.SyncResult{}, err
		}
		result = countOutcomes(result, outcomes)
	}
	return result, nil
}

// emailOutcome is what happened to one email of the sync.
type emailOutcome int

const (
	outcomeNew emailOutcome = iota
	outcomeAlreadySaved
	outcomeUnparsed
)

func (u SyncAll) syncBank(ctx context.Context, credential domain.Credential, parser ports.BankEmailParser, since time.Time) ([]emailOutcome, error) {
	emails, err := u.readBankEmails(ctx, credential, parser, since)
	if err != nil {
		return nil, err
	}
	outcomes := make([]emailOutcome, 0, len(emails))
	for _, email := range emails {
		outcome, err := u.saveEmail(ctx, credential, parser, email)
		if err != nil {
			return nil, err
		}
		outcomes = append(outcomes, outcome)
	}
	return outcomes, nil
}

func countOutcomes(result domain.SyncResult, outcomes []emailOutcome) domain.SyncResult {
	for _, outcome := range outcomes {
		result.Found++
		switch outcome {
		case outcomeNew:
			result.New++
		case outcomeUnparsed:
			result.Unparsed++
		}
	}
	return result
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

// saveEmail stores the email as a charge, or as unparsed when its bank's
// parser cannot read it or what it read is not a valid charge. Error
// messages name the account by public id, never by address.
func (u SyncAll) saveEmail(ctx context.Context, credential domain.Credential, parser ports.BankEmailParser, email domain.RawEmail) (emailOutcome, error) {
	parsed, isRecognized := parser.Parse(email)
	if !isRecognized {
		return u.saveUnparsed(ctx, credential, parser, email)
	}
	charge, err := domain.ChargeFromEmail(parsed, email, credential.AccountID, parser.Bank())
	if errors.Is(err, shared.ErrInvalid) {
		return u.saveUnparsed(ctx, credential, parser, email)
	}
	if err != nil {
		return 0, fmt.Errorf("convert %s email of account %s: %w", parser.Bank(), credential.AccountPublicID, err)
	}
	return u.saveCharge(ctx, credential, email, charge)
}

// saveCharge forgets the email as unparsed once it becomes a charge, since a
// newer parser may now read what an older one could not.
func (u SyncAll) saveCharge(ctx context.Context, credential domain.Credential, email domain.RawEmail, charge shared.Charge) (emailOutcome, error) {
	isNew, err := u.charges.Insert(ctx, charge)
	if err != nil {
		return 0, fmt.Errorf("insert charge of account %s: %w", credential.AccountPublicID, err)
	}
	if !isNew {
		return outcomeAlreadySaved, nil
	}
	if err := u.unparsed.Forget(ctx, email.MessageID); err != nil {
		return 0, fmt.Errorf("forget unparsed email of account %s: %w", credential.AccountPublicID, err)
	}
	return outcomeNew, nil
}

func (u SyncAll) saveUnparsed(ctx context.Context, credential domain.Credential, parser ports.BankEmailParser, email domain.RawEmail) (emailOutcome, error) {
	unparsed := domain.UnparsedFromEmail(email, credential.AccountID, parser.Bank())
	if err := u.unparsed.Save(ctx, unparsed); err != nil {
		return 0, fmt.Errorf("save unparsed email of account %s: %w", credential.AccountPublicID, err)
	}
	return outcomeUnparsed, nil
}
