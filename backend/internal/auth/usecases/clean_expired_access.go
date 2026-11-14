package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

// CleanExpiredAccess removes what no longer serves; an expired code or session
// is already rejected, so this only keeps the tables small.
type CleanExpiredAccess struct {
	expired ports.ExpiredAccess
	now     func() time.Time
}

func NewCleanExpiredAccess(expired ports.ExpiredAccess, now func() time.Time) CleanExpiredAccess {
	return CleanExpiredAccess{expired: expired, now: now}
}

func (u CleanExpiredAccess) Execute(ctx context.Context) error {
	if err := u.expired.DeleteExpired(ctx, domain.CleanupCutoffAt(u.now())); err != nil {
		return fmt.Errorf("clean expired access: %w", err)
	}
	return nil
}
