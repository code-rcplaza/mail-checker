package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// SyncResult is the outcome of syncing one account.
type SyncResult struct {
	AccountPublicID string
	Email           string
	Found           int
	New             int
	Unparsed        int
	ErrorKey        shared.MessageKey // empty when the account synced fine
}
