package usecases

import (
	"context"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// DecisionInput is the user's answer on a possible line: the action and,
// when confirming, the cadence they chose.
type DecisionInput struct {
	Key     domain.GroupKey
	Action  string
	Cadence string
}

// DecidePossibleSubscription stores the user's confirmation (with the cadence
// they chose) or dismissal of a possible subscription.
type DecidePossibleSubscription struct {
	lines     ChargeLines
	decisions ports.SubscriptionDecisions
}

func NewDecidePossibleSubscription(lines ChargeLines, decisions ports.SubscriptionDecisions) DecidePossibleSubscription {
	return DecidePossibleSubscription{lines: lines, decisions: decisions}
}

func (u DecidePossibleSubscription) Execute(ctx context.Context, userID int64, input DecisionInput) error {
	decision, err := domain.NewDecision(input.Action, input.Cadence)
	if err != nil {
		return err
	}
	if err := u.ensureCurrentPossible(ctx, userID, input.Key); err != nil {
		return err
	}
	return u.saveDecision(ctx, userID, domain.LineDecision{Key: input.Key, Decision: decision})
}

// ensureCurrentPossible looks over the widest period, so a possible visible
// in any period the user picked can be decided.
func (u DecidePossibleSubscription) ensureCurrentPossible(ctx context.Context, userID int64, key domain.GroupKey) error {
	groups, err := u.lines.readAllGroups(ctx, userID)
	if err != nil {
		return err
	}
	if !u.lines.detector(u.lines.Now()).IsCurrentPossible(groups, key) {
		return shared.NotFound(domain.KeyPossibleNotFound)
	}
	return nil
}

func (u DecidePossibleSubscription) saveDecision(ctx context.Context, userID int64, decision domain.LineDecision) error {
	if err := u.decisions.Save(ctx, userID, decision); err != nil {
		return fmt.Errorf("save decision for %s: %w", decision.Key, err)
	}
	return nil
}
