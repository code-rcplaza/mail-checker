package usecases

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

type OpenSessionDependencies struct {
	Verify   VerifyCode
	Sessions ports.Sessions
	Random   io.Reader
	Now      func() time.Time
}

type OpenSession struct {
	dependencies OpenSessionDependencies
}

func NewOpenSession(dependencies OpenSessionDependencies) OpenSession {
	return OpenSession{dependencies: dependencies}
}

// Execute returns the token in clear, the only time it exists outside the browser.
func (u OpenSession) Execute(ctx context.Context, entry CodeEntry) (string, error) {
	user, err := u.dependencies.Verify.Execute(ctx, entry)
	if err != nil {
		return "", err
	}
	return u.createSession(ctx, user.ID)
}

func (u OpenSession) createSession(ctx context.Context, userID int64) (string, error) {
	opened, err := domain.NewSession(u.dependencies.Random, domain.SessionOpening{UserID: userID, Now: u.dependencies.Now()})
	if err != nil {
		return "", err
	}
	if err := u.dependencies.Sessions.Create(ctx, opened.Session); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	return opened.Token, nil
}
