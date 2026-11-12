package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// GetSummary is the unified view: every charge of the period counted once,
// grouped and sorted into subscriptions, possible, inactive and other charges.
type GetSummary struct {
	charges   ports.ChargeReader
	catalog   ports.ServiceCatalog
	decisions ports.SubscriptionDecisions
	now       func() time.Time
}

func NewGetSummary(
	charges ports.ChargeReader,
	catalog ports.ServiceCatalog,
	decisions ports.SubscriptionDecisions,
	now func() time.Time,
) GetSummary {
	return GetSummary{charges: charges, catalog: catalog, decisions: decisions, now: now}
}

func (u GetSummary) Execute(ctx context.Context, months int) (domain.Summary, error) {
	if err := domain.ValidatePeriodMonths(months); err != nil {
		return domain.Summary{}, err
	}
	now := u.now()
	charges, err := readChargesSince(ctx, u.charges, domain.PeriodStart(now, months))
	if err != nil {
		return domain.Summary{}, err
	}
	decisions, err := u.readDecisions(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	unique := domain.DropStatementChargesMatchingEmail(charges, u.catalog.Identify, u.catalog.Normalize)
	groups := domain.GroupCharges(unique, u.catalog.Identify, u.catalog.Normalize)
	return domain.BuildSummary(groups, decisions, u.catalog.LooksLikeSubscription, now), nil
}

func (u GetSummary) readDecisions(ctx context.Context) (map[domain.GroupKey]domain.Decision, error) {
	decisions, err := u.decisions.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("read subscription decisions: %w", err)
	}
	return decisions, nil
}
