package domain

import (
	"errors"
	"time"
)

// RawEmail is a purchase email as the mail reader returns it.
type RawEmail struct {
	MessageID  string
	Subject    string
	From       string
	ReceivedAt time.Time
	Body       string
}

// MailQuery is a bank's filter for its purchase emails, from a date on.
type MailQuery struct {
	From    string
	Subject string
	Since   time.Time
}

// Credential holds a decrypted app password; it must never be logged.
type Credential struct {
	AccountID       int64
	AccountPublicID string
	Email           string
	Password        string
}

// ErrMailAuthRejected is what the mail reader returns when the server
// rejects the credentials, so the sync can tell it from other read errors.
var ErrMailAuthRejected = errors.New("mail authentication rejected")
