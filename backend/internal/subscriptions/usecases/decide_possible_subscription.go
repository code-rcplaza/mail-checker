package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// DecidePossibleSubscription stores the user's confirmation (with the cadence
// they chose) or dismissal of a possible subscription.
type DecidePossibleSubscription struct {
	charges   ports.ChargeReader
	catalog   ports.ServiceCatalog
	decisions ports.SubscriptionDecisions
	now       func() time.Time
}

func NewDecidePossibleSubscription(
	charges ports.ChargeReader,
	catalog ports.ServiceCatalog,
	decisions ports.SubscriptionDecisions,
	now func() time.Time,
) DecidePossibleSubscription {
	return DecidePossibleSubscription{charges: charges, catalog: catalog, decisions: decisions, now: now}
}

func (u DecidePossibleSubscription) Execute(ctx context.Context, key domain.GroupKey, action, cadence string) error {
	decision, err := domain.NewDecision(action, cadence)
	if err != nil {
		return err
	}
	if err := u.ensureCurrentPossible(ctx, key); err != nil {
		return err
	}
	return u.saveDecision(ctx, key, decision)
}

// ensureCurrentPossible looks over the widest period, so a possible visible
// in any period the user picked can be decided.
func (u DecidePossibleSubscription) ensureCurrentPossible(ctx context.Context, key domain.GroupKey) error {
	now := u.now()
	groups, err := readGroupsSince(ctx, u.charges, u.catalog, domain.PeriodStart(now, domain.MaxPeriodMonths))
	if err != nil {
		return err
	}
	if !domain.IsCurrentPossible(groups, key, u.catalog.LooksLikeSubscription, now) {
		return shared.NotFound(domain.KeyPossibleNotFound)
	}
	return nil
}

func (u DecidePossibleSubscription) saveDecision(ctx context.Context, key domain.GroupKey, decision domain.Decision) error {
	if err := u.decisions.Save(ctx, key, decision); err != nil {
		return fmt.Errorf("save decision for %s: %w", key, err)
	}
	return nil
}
