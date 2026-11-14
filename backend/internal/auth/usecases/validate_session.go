package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type ValidateSession struct {
	sessions ports.Sessions
	now      func() time.Time
}

func NewValidateSession(sessions ports.Sessions, now func() time.Time) ValidateSession {
	return ValidateSession{sessions: sessions, now: now}
}

// Execute does not tell an expired token from an unknown one: both are unauthorized.
func (u ValidateSession) Execute(ctx context.Context, token string) (int64, error) {
	session, err := u.sessions.FindByTokenHash(ctx, domain.HashSecret(token))
	if errors.Is(err, shared.ErrNotFound) {
		return 0, shared.Unauthorized(shared.KeyUnauthorized)
	}
	if err != nil {
		return 0, fmt.Errorf("find session: %w", err)
	}
	if session.IsExpiredAt(u.now()) {
		return 0, shared.Unauthorized(shared.KeyUnauthorized)
	}
	return session.UserID, nil
}
