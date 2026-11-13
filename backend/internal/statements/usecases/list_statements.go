package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

type ListStatements struct {
	statements ports.StatementRepository
}

func NewListStatements(statements ports.StatementRepository) ListStatements {
	return ListStatements{statements: statements}
}

func (u ListStatements) Execute(ctx context.Context, userID int64) ([]domain.Statement, error) {
	statements, err := u.statements.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list statements: %w", err)
	}
	return statements, nil
}
