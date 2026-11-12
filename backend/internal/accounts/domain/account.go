package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Account is a mailbox the sync reads. The app password is only stored sealed.
type Account struct {
	ID           int64  // internal
	PublicID     string // UUID v7, the one exposed
	Email        string
	PasswordEnc  []byte
	CreatedAt    time.Time
	LastSyncAt   *time.Time
	LastErrorKey shared.MessageKey // empty when the last sync had no error
}
