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
	List(ctx context.Context, email string) ([]time.Time, error)
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
