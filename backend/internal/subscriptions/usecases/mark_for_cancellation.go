package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// MarkForCancellation reminds the user of a line they mean to cancel; the app
// cancels nothing.
type MarkForCancellation struct {
	lines         ChargeLines
	cancellations ports.Cancellations
}

func NewMarkForCancellation(lines ChargeLines, cancellations ports.Cancellations) MarkForCancellation {
	return MarkForCancellation{lines: lines, cancellations: cancellations}
}

func (u MarkForCancellation) Execute(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.lines.ensureLine(ctx, userID, key); err != nil {
		return err
	}
	return u.markLine(ctx, userID, domain.CancellationMark{Key: key, MarkedAt: u.lines.Now()})
}

func (u MarkForCancellation) markLine(ctx context.Context, userID int64, mark domain.CancellationMark) error {
	if err := u.cancellations.Mark(ctx, userID, mark); err != nil {
		return fmt.Errorf("mark line %s for cancellation: %w", mark.Key, err)
	}
	return nil
}
