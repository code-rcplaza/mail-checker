// Package ports declares what the statement use cases need from the outside.
package ports

import (
	"context"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

// FileExtractor turns an uploaded file into text and tells its kind.
type FileExtractor interface {
	Extract(filename string, data []byte) (text string, kind domain.FileKind, err error)
}

// BankStatementParser knows how one bank's statement looks. Bad lines are
// line numbers in the text, never the text itself.
type BankStatementParser interface {
	Bank() shared.Bank
	Detect(text string) bool
	Parse(text string) (lines []domain.StatementLine, bad []int)
	// BillingDate is the billing date the file itself shows; false when it has none.
	BillingDate(text string) (time.Time, bool)
	// TransferPair is the debt moved between the national and the
	// international statement, in pesos and in dollars; false when the bank
	// does not show one or the file has none.
	TransferPair(text string) (domain.TransferPair, bool)
}

// ExchangeRateRecorder keeps the rate a statement showed, one per user and
// date: recording the same date again changes nothing.
type ExchangeRateRecorder interface {
	Record(ctx context.Context, userID int64, rate shared.ExchangeRate) error
}

// StatementRepository only ever touches the statements of the given user.
type StatementRepository interface {
	// Create saves the statement and its charges atomically and assigns the ID.
	// A statement is unique per user, bank and billing date.
	Create(ctx context.Context, userID int64, statement *domain.Statement, charges []shared.Charge) error
	List(ctx context.Context, userID int64) ([]domain.Statement, error)
	// Delete also removes the statement's charges; a publicID that is unknown
	// or belongs to another user is shared.ErrNotFound.
	Delete(ctx context.Context, userID int64, publicID string) error
}
