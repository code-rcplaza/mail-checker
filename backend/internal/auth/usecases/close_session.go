package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

type CloseSession struct {
	sessions ports.Sessions
}

func NewCloseSession(sessions ports.Sessions) CloseSession {
	return CloseSession{sessions: sessions}
}

// Execute closes only this device's session; closing an unknown one is not an error.
func (u CloseSession) Execute(ctx context.Context, token string) error {
	if err := u.sessions.Delete(ctx, domain.HashSecret(token)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
