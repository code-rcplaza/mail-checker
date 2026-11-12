package usecases

import (
	"context"
	"errors"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

// RemoveStatement deletes a statement and its charges; reimporting a billing
// date goes through here first.
type RemoveStatement struct {
	statements ports.StatementRepository
}

func NewRemoveStatement(statements ports.StatementRepository) RemoveStatement {
	return RemoveStatement{statements: statements}
}

// Execute adds the i18n key to the repository's plain ErrNotFound, so adapters
// stay free of user-facing keys.
func (u RemoveStatement) Execute(ctx context.Context, publicID string) error {
	err := u.statements.Delete(ctx, publicID)
	if errors.Is(err, shared.ErrNotFound) {
		return shared.NotFound(domain.KeyNotFound)
	}
	if err != nil {
		return fmt.Errorf("delete statement: %w", err)
	}
	return nil
}
