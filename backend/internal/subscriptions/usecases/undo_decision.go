package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// UndoDecision forgets the user's confirmation or dismissal of a line, so it
// goes back to possible. Undoing a line without a decision changes nothing.
type UndoDecision struct {
	decisions ports.SubscriptionDecisions
}

func NewUndoDecision(decisions ports.SubscriptionDecisions) UndoDecision {
	return UndoDecision{decisions: decisions}
}

func (u UndoDecision) Execute(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.decisions.Delete(ctx, userID, key); err != nil {
		return fmt.Errorf("delete decision for %s: %w", key, err)
	}
	return nil
}
