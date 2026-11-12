package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/ports"
)

type ListAccounts struct {
	accounts ports.AccountRepository
}

func NewListAccounts(accounts ports.AccountRepository) ListAccounts {
	return ListAccounts{accounts: accounts}
}

func (u ListAccounts) Execute(ctx context.Context) ([]domain.Account, error) {
	accounts, err := u.accounts.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return accounts, nil
}
