package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// GetSummary is the unified view: every charge of the period counted once,
// grouped and sorted into subscriptions, possible, inactive and other charges,
// each line with the name the user gave it.
type GetSummary struct {
	charges   ports.ChargeReader
	catalog   ports.ServiceCatalog
	decisions ports.SubscriptionDecisions
	aliases   ports.LineAliases
	now       func() time.Time
}

func NewGetSummary(
	charges ports.ChargeReader,
	catalog ports.ServiceCatalog,
	decisions ports.SubscriptionDecisions,
	aliases ports.LineAliases,
	now func() time.Time,
) GetSummary {
	return GetSummary{charges: charges, catalog: catalog, decisions: decisions, aliases: aliases, now: now}
}

func (u GetSummary) Execute(ctx context.Context, months int) (domain.Summary, error) {
	if err := domain.ValidatePeriodMonths(months); err != nil {
		return domain.Summary{}, err
	}
	now := u.now()
	groups, err := readGroupsSince(ctx, u.charges, u.catalog, domain.PeriodStart(now, months))
	if err != nil {
		return domain.Summary{}, err
	}
	decisions, err := u.readDecisions(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	aliases, err := u.readAliases(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	summary := domain.BuildSummary(groups, decisions, u.catalog.LooksLikeSubscription, now)
	return domain.ApplyLineAliases(summary, aliases), nil
}

func (u GetSummary) readDecisions(ctx context.Context) (map[domain.GroupKey]domain.Decision, error) {
	decisions, err := u.decisions.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("read subscription decisions: %w", err)
	}
	return decisions, nil
}

func (u GetSummary) readAliases(ctx context.Context) (map[domain.GroupKey]string, error) {
	aliases, err := u.aliases.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("read line aliases: %w", err)
	}
	return aliases, nil
}
