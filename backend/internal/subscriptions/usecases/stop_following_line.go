package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// StopFollowingLine hides a line: it no longer shows nor adds to the spend.
type StopFollowingLine struct {
	lines  ChargeLines
	hidden ports.HiddenLines
}

func NewStopFollowingLine(lines ChargeLines, hidden ports.HiddenLines) StopFollowingLine {
	return StopFollowingLine{lines: lines, hidden: hidden}
}

func (u StopFollowingLine) Execute(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.lines.ensureLine(ctx, userID, key); err != nil {
		return err
	}
	return u.hideLine(ctx, userID, domain.HiddenLine{Key: key, HiddenAt: u.lines.Now()})
}

func (u StopFollowingLine) hideLine(ctx context.Context, userID int64, hidden domain.HiddenLine) error {
	if err := u.hidden.Hide(ctx, userID, hidden); err != nil {
		return fmt.Errorf("hide line %s: %w", hidden.Key, err)
	}
	return nil
}
