package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
)

// ListUnparsed shows the emails no parser could read, so the user knows which
// purchase notices are being missed.
type ListUnparsed struct {
	unparsed ports.UnparsedRepository
}

func NewListUnparsed(unparsed ports.UnparsedRepository) ListUnparsed {
	return ListUnparsed{unparsed: unparsed}
}

func (u ListUnparsed) Execute(ctx context.Context, limit int) ([]domain.UnparsedEmail, error) {
	if err := domain.ValidateUnparsedLimit(limit); err != nil {
		return nil, err
	}
	return u.readUnparsed(ctx, limit)
}

func (u ListUnparsed) readUnparsed(ctx context.Context, limit int) ([]domain.UnparsedEmail, error) {
	emails, err := u.unparsed.List(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list unparsed emails: %w", err)
	}
	return emails, nil
}
