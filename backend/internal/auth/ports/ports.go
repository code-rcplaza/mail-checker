package ports

import (
	"context"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
)

type UserRepository interface {
	// Create assigns the ID.
	Create(ctx context.Context, user domain.User) (domain.User, error)
	// FindByEmail returns ErrNotFound when no user has that normalized email.
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	// FindByID returns ErrNotFound when no user has that ID.
	FindByID(ctx context.Context, id int64) (domain.User, error)
	// Delete removes the user and everything of theirs in one transaction;
	// ErrNotFound when no user has that ID.
	Delete(ctx context.Context, id int64) error
}

// LoginCodes keeps one code per user.
type LoginCodes interface {
	// Save replaces the previous code of the user.
	Save(ctx context.Context, code domain.LoginCode) error
	// Find returns ErrNotFound when the user has no code.
	Find(ctx context.Context, userID int64) (domain.LoginCode, error)
	// AddFailedAttempt adds one attempt atomically and returns the total.
	AddFailedAttempt(ctx context.Context, userID int64) (int, error)
	// Delete is idempotent.
	Delete(ctx context.Context, userID int64) error
}

type CodeSends interface {
	// ListSince returns the sends to email after since (strictly).
	ListSince(ctx context.Context, email string, since time.Time) ([]time.Time, error)
	Record(ctx context.Context, send domain.CodeSend) error
}

type Sessions interface {
	Create(ctx context.Context, session domain.Session) error
	// FindByTokenHash returns ErrNotFound when no session has that hash.
	FindByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error)
	// Delete is idempotent.
	Delete(ctx context.Context, tokenHash string) error
}

// CodeSender delivers the code: console in development, Resend in production.
type CodeSender interface {
	Send(ctx context.Context, email, code string) error
}

// ExpiredAccess deletes the codes and sessions that expired and the sends that
// no longer count for the limit.
type ExpiredAccess interface {
	DeleteExpired(ctx context.Context, cutoff domain.CleanupCutoff) error
}
