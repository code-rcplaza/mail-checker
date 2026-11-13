package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// GetSummary is the unified view: every charge of the period counted once,
// grouped and sorted into subscriptions, possible, inactive and other charges,
// each line with the name the user gave it.
type GetSummary struct {
	sources SummarySources
	now     func() time.Time
}

// SummarySources are the ports the summary reads, named so a new one does not
// shift the others.
type SummarySources struct {
	Charges   ports.ChargeReader
	Catalog   ports.ServiceCatalog
	Decisions ports.SubscriptionDecisions
	Aliases   ports.LineAliases
	Rates     ports.ExchangeRates
	KeptRises ports.PriceAlertDecisions
	Manuals   ports.ManualSubscriptions
	Overrides ports.LineOverrides
}

func NewGetSummary(sources SummarySources, now func() time.Time) GetSummary {
	return GetSummary{sources: sources, now: now}
}

func (u GetSummary) Execute(ctx context.Context, userID int64, months int) (domain.Summary, error) {
	if err := domain.ValidatePeriodMonths(months); err != nil {
		return domain.Summary{}, err
	}
	now := u.now()
	groups, err := readGroupsSince(ctx, userID, u.sources.Charges, u.sources.Catalog, domain.PeriodStart(now, months))
	if err != nil {
		return domain.Summary{}, err
	}
	decisions, err := u.readDecisions(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	aliases, err := u.readAliases(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	rates, err := u.readExchangeRates(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	manuals, err := u.readManuals(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	overrides, err := u.readOverrides(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	summary := domain.BuildSummary(groups, decisions, u.sources.Catalog.LooksLikeSubscription, now)
	summary = domain.ApplyLineOverrides(summary, overrides, now)
	summary = domain.AddManualSubscriptions(summary, manuals, now)
	keptRises, err := u.readKeptRises(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	summary = domain.ApplyExchangeRates(summary, rates)
	summary = domain.HideKeptRises(summary, keptRises)
	return domain.ApplyLineAliases(summary, aliases), nil
}

func (u GetSummary) readKeptRises(ctx context.Context, userID int64) (map[domain.GroupKey]shared.Money, error) {
	kept, err := u.sources.KeptRises.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read kept rises: %w", err)
	}
	return kept, nil
}

func (u GetSummary) readExchangeRates(ctx context.Context, userID int64) ([]shared.ExchangeRate, error) {
	rates, err := u.sources.Rates.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read exchange rates: %w", err)
	}
	return rates, nil
}

func (u GetSummary) readDecisions(ctx context.Context, userID int64) (map[domain.GroupKey]domain.Decision, error) {
	decisions, err := u.sources.Decisions.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read subscription decisions: %w", err)
	}
	return decisions, nil
}

func (u GetSummary) readAliases(ctx context.Context, userID int64) (map[domain.GroupKey]string, error) {
	aliases, err := u.sources.Aliases.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read line aliases: %w", err)
	}
	return aliases, nil
}

func (u GetSummary) readManuals(ctx context.Context, userID int64) ([]domain.ManualSubscription, error) {
	manuals, err := u.sources.Manuals.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read manual subscriptions: %w", err)
	}
	return manuals, nil
}

func (u GetSummary) readOverrides(ctx context.Context, userID int64) (map[domain.GroupKey]domain.LineOverride, error) {
	overrides, err := u.sources.Overrides.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read line overrides: %w", err)
	}
	return overrides, nil
}
