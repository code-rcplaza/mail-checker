package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// RemoveAccount deletes a mailbox and its charges; changing its password goes
// through here and AddAccount.
type RemoveAccount struct {
	accounts ports.AccountRepository
}

func NewRemoveAccount(accounts ports.AccountRepository) RemoveAccount {
	return RemoveAccount{accounts: accounts}
}

// Execute adds the i18n key to the repository's plain ErrNotFound, so adapters
// stay free of user-facing keys.
func (u RemoveAccount) Execute(ctx context.Context, publicID string) error {
	err := u.accounts.Delete(ctx, publicID)
	if errors.Is(err, shared.ErrNotFound) {
		return shared.NotFound(domain.KeyNotFound)
	}
	if err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	return nil
}
