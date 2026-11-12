package domain

// Credential is an account with its app password decrypted, for the sync to
// log into the mailbox. It must never be logged.
type Credential struct {
	AccountID       int64
	AccountPublicID string
	Email           string
	Password        string
}
