package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// SyncRecord is when an account was synced and why it failed, if it did.
type SyncRecord struct {
	AccountID int64
	At        time.Time
	ErrorKey  shared.MessageKey // empty when the sync had no error
}

// SyncResult is the outcome of syncing one account.
type SyncResult struct {
	AccountPublicID string
	Email           string
	Found           int
	New             int
	Unparsed        int
	ErrorKey        shared.MessageKey // empty when the account synced fine
}
