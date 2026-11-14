// Package ports declares what the account use cases need from the outside.
package ports

import (
	"context"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
)

type AccountRepository interface {
	// Create assigns the ID; an email already registered is ErrConflict with accounts.email_taken.
	Create(ctx context.Context, account *domain.Account) error
	List(ctx context.Context) ([]domain.Account, error)
	// Delete also removes the account's charges; an unknown publicID is shared.ErrNotFound.
	Delete(ctx context.Context, publicID string) error
	// MarkSynced records the last sync; an unknown id is shared.ErrNotFound.
	MarkSynced(ctx context.Context, record domain.SyncRecord) error
}

// SecretCipher seals app passwords at rest.
type SecretCipher interface {
	Encrypt(plain string) ([]byte, error)
	Decrypt(sealed []byte) (string, error)
}
