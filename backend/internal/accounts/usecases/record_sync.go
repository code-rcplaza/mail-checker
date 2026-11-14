package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// RecordSync keeps when each account was last synced and why it failed, if it did.
type RecordSync struct {
	accounts ports.AccountRepository
}

func NewRecordSync(accounts ports.AccountRepository) RecordSync {
	return RecordSync{accounts: accounts}
}

// Execute ignores an account removed while the sync ran: one missing account
// must not stop the sync of the others.
func (u RecordSync) Execute(ctx context.Context, record domain.SyncRecord) error {
	err := u.accounts.MarkSynced(ctx, record)
	if errors.Is(err, shared.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("record account sync: %w", err)
	}
	return nil
}
