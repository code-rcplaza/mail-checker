// Package ports declares what the statement use cases need from the outside.
package ports

import (
	"context"

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
}

type StatementRepository interface {
	// Create saves the statement and its charges atomically and assigns the ID.
	Create(ctx context.Context, statement *domain.Statement, charges []shared.Charge) error
	List(ctx context.Context) ([]domain.Statement, error)
	// Delete also removes the statement's charges; an unknown publicID is shared.ErrNotFound.
	Delete(ctx context.Context, publicID string) error
}
