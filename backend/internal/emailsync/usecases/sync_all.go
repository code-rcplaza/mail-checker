package usecases

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync/atomic"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// SyncAll reads the purchase emails of every account, for every registered
// bank, and stores them as charges.
type SyncAll struct {
	dependencies SyncAllDependencies
	// inProgress is shared by every copy of SyncAll, so only one sync runs at a time.
	inProgress *atomic.Bool
}

// SyncAllDependencies are the ports and the clock a sync needs.
type SyncAllDependencies struct {
	Accounts ports.AccountDirectory
	Mail     ports.MailReader
	Parsers  []ports.BankEmailParser
	Charges  ports.ChargeWriter
	Unparsed ports.UnparsedRepository
	Now      func() time.Time
}

// bankMailbox is one account's mailbox searched for one bank's notices.
type bankMailbox struct {
	credential domain.Credential
	parser     ports.BankEmailParser
}

// emailCharge is an email and the charge read from it.
type emailCharge struct {
	email  domain.RawEmail
	charge shared.Charge
}

// NewSyncAll copies the bank registry so it stays read-only after main builds it.
func NewSyncAll(dependencies SyncAllDependencies) SyncAll {
	dependencies.Parsers = slices.Clone(dependencies.Parsers)
	return SyncAll{dependencies: dependencies, inProgress: &atomic.Bool{}}
}

func (u SyncAll) Execute(ctx context.Context, months int) ([]domain.SyncResult, error) {
	if err := domain.ValidateSyncMonths(months); err != nil {
		return nil, err
	}
	if err := u.startSync(); err != nil {
		return nil, err
	}
	defer u.finishSync()
	return u.syncAccounts(ctx, months)
}

func (u SyncAll) startSync() error {
	if !u.inProgress.CompareAndSwap(false, true) {
		return shared.Conflict(domain.KeySyncInProgress)
	}
	return nil
}

// finishSync runs whether the sync succeeded or failed, so the next one can start.
func (u SyncAll) finishSync() {
	u.inProgress.Store(false)
}

func (u SyncAll) syncAccounts(ctx context.Context, months int) ([]domain.SyncResult, error) {
	credentials, err := u.readCredentials(ctx)
	if err != nil {
		return nil, err
	}
	since := domain.SyncSince(u.dependencies.Now(), months)
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
	credentials, err := u.dependencies.Accounts.Credentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("read account credentials: %w", err)
	}
	return credentials, nil
}

// syncAccount keeps a mail read failure in the account's result so the other
// accounts still sync; a storage failure stops everything.
func (u SyncAll) syncAccount(ctx context.Context, credential domain.Credential, since time.Time) (domain.SyncResult, error) {
	result, err := u.syncBanks(ctx, credential, since)
	if err != nil {
		return domain.SyncResult{}, err
	}
	record := domain.SyncRecord{AccountID: credential.AccountID, At: u.dependencies.Now(), ErrorKey: result.ErrorKey}
	if err := u.recordSync(ctx, credential, record); err != nil {
		return domain.SyncResult{}, err
	}
	return result, nil
}

func (u SyncAll) syncBanks(ctx context.Context, credential domain.Credential, since time.Time) (domain.SyncResult, error) {
	result := domain.SyncResult{AccountPublicID: credential.AccountPublicID, Email: credential.Email}
	for _, parser := range u.dependencies.Parsers {
		mailbox := bankMailbox{credential: credential, parser: parser}
		emails, readErr := u.readBankEmails(ctx, mailbox, since)
		if readErr != nil {
			result.ErrorKey = mailErrorKey(readErr)
			return result, nil
		}
		outcomes, err := u.saveEmails(ctx, mailbox, emails)
		if err != nil {
			return domain.SyncResult{}, err
		}
		result = countOutcomes(result, outcomes)
	}
	return result, nil
}

func mailErrorKey(readErr error) shared.MessageKey {
	if errors.Is(readErr, domain.ErrMailAuthRejected) {
		return domain.KeyMailAuthFailed
	}
	return domain.KeyMailReadFailed
}

func (u SyncAll) recordSync(ctx context.Context, credential domain.Credential, record domain.SyncRecord) error {
	if err := u.dependencies.Accounts.RecordSync(ctx, record); err != nil {
		return fmt.Errorf("record sync of account %s: %w", credential.AccountPublicID, err)
	}
	return nil
}

// emailOutcome is what happened to one email of the sync.
type emailOutcome int

const (
	outcomeNew emailOutcome = iota
	outcomeAlreadySaved
	outcomeUnparsed
)

func (u SyncAll) saveEmails(ctx context.Context, mailbox bankMailbox, emails []domain.RawEmail) ([]emailOutcome, error) {
	outcomes := make([]emailOutcome, 0, len(emails))
	for _, email := range emails {
		outcome, err := u.saveEmail(ctx, mailbox, email)
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

func (u SyncAll) readBankEmails(ctx context.Context, mailbox bankMailbox, since time.Time) ([]domain.RawEmail, error) {
	query := mailbox.parser.Query()
	query.Since = since
	emails, err := u.dependencies.Mail.Read(ctx, mailbox.credential, query)
	if err != nil {
		return nil, fmt.Errorf("read %s emails of account %s: %w", mailbox.parser.Bank(), mailbox.credential.AccountPublicID, err)
	}
	return emails, nil
}

// saveEmail stores the email as a charge, or as unparsed when its bank's
// parser cannot read it or what it read is not a valid charge. Error
// messages name the account by public id, never by address.
func (u SyncAll) saveEmail(ctx context.Context, mailbox bankMailbox, email domain.RawEmail) (emailOutcome, error) {
	read := mailbox.bankEmail(email)
	parsed, isRecognized := mailbox.parser.Parse(email)
	if !isRecognized {
		return u.saveUnparsed(ctx, mailbox, read)
	}
	charge, err := domain.ChargeFromEmail(parsed, read)
	if errors.Is(err, shared.ErrInvalid) {
		return u.saveUnparsed(ctx, mailbox, read)
	}
	if err != nil {
		return 0, fmt.Errorf("convert %s email of account %s: %w", mailbox.parser.Bank(), mailbox.credential.AccountPublicID, err)
	}
	return u.saveCharge(ctx, mailbox, emailCharge{email: email, charge: charge})
}

func (mailbox bankMailbox) bankEmail(email domain.RawEmail) domain.BankEmail {
	return domain.BankEmail{Email: email, UserID: mailbox.credential.UserID, AccountID: mailbox.credential.AccountID, Bank: mailbox.parser.Bank()}
}

// saveCharge forgets the email as unparsed once it becomes a charge, since a
// newer parser may now read what an older one could not.
func (u SyncAll) saveCharge(ctx context.Context, mailbox bankMailbox, read emailCharge) (emailOutcome, error) {
	isNew, err := u.dependencies.Charges.Insert(ctx, read.charge)
	if err != nil {
		return 0, fmt.Errorf("insert charge of account %s: %w", mailbox.credential.AccountPublicID, err)
	}
	if !isNew {
		return outcomeAlreadySaved, nil
	}
	if err := u.dependencies.Unparsed.Forget(ctx, read.email.MessageID); err != nil {
		return 0, fmt.Errorf("forget unparsed email of account %s: %w", mailbox.credential.AccountPublicID, err)
	}
	return outcomeNew, nil
}

func (u SyncAll) saveUnparsed(ctx context.Context, mailbox bankMailbox, read domain.BankEmail) (emailOutcome, error) {
	if err := u.dependencies.Unparsed.Save(ctx, domain.UnparsedFromEmail(read)); err != nil {
		return 0, fmt.Errorf("save unparsed email of account %s: %w", mailbox.credential.AccountPublicID, err)
	}
	return outcomeUnparsed, nil
}
