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
// each line with what the user chose for it.
type GetSummary struct {
	sources SummarySources
}

// SummarySources are the ports the summary reads, named so a new one does not
// shift the others.
type SummarySources struct {
	Lines         ChargeLines
	Decisions     ports.SubscriptionDecisions
	Aliases       ports.LineAliases
	Rates         ports.ExchangeRates
	KeptRises     ports.PriceAlertDecisions
	Manuals       ports.ManualSubscriptions
	Overrides     ports.LineOverrides
	Hidden        ports.HiddenLines
	Cancellations ports.Cancellations
}

// userChoices is everything the user decided or added on top of their charges.
type userChoices struct {
	decisions     map[domain.GroupKey]domain.Decision
	aliases       map[domain.GroupKey]string
	rates         []shared.ExchangeRate
	keptRises     map[domain.GroupKey][]shared.Money
	manuals       []domain.ManualSubscription
	overrides     map[domain.GroupKey]domain.LineOverride
	hidden        []domain.GroupKey
	cancellations map[domain.GroupKey]time.Time
}

func NewGetSummary(sources SummarySources) GetSummary {
	return GetSummary{sources: sources}
}

func (u GetSummary) Execute(ctx context.Context, userID int64, months int) (domain.Summary, error) {
	if err := domain.ValidatePeriodMonths(months); err != nil {
		return domain.Summary{}, err
	}
	now := u.sources.Lines.Now()
	groups, err := u.sources.Lines.readGroupsSince(ctx, userID, domain.PeriodStart(now, months))
	if err != nil {
		return domain.Summary{}, err
	}
	choices, err := u.readUserChoices(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	input := domain.SummaryInput{Groups: groups, Decisions: choices.decisions, Detector: u.sources.Lines.detector(now)}
	return composeSummary(input, choices), nil
}

// composeSummary applies the user's choices in order: what changes a line
// before what counts it, and the names last.
func composeSummary(input domain.SummaryInput, choices userChoices) domain.Summary {
	summary := domain.BuildSummary(input)
	summary = domain.ApplyLineOverrides(summary, choices.overrides)
	summary = domain.AddManualSubscriptions(summary, choices.manuals)
	summary = domain.HideLines(summary, choices.hidden)
	summary = domain.ApplyCancellations(summary, choices.cancellations)
	summary = domain.ApplyExchangeRates(summary, choices.rates)
	summary = domain.HideKeptRises(summary, choices.keptRises)
	return domain.ApplyLineAliases(summary, choices.aliases)
}

func (u GetSummary) readUserChoices(ctx context.Context, userID int64) (userChoices, error) {
	var choices userChoices
	var err error
	if choices.decisions, err = u.sources.Decisions.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read subscription decisions: %w", err)
	}
	if choices.aliases, err = u.sources.Aliases.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read line aliases: %w", err)
	}
	if choices.rates, err = u.sources.Rates.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read exchange rates: %w", err)
	}
	if choices.keptRises, err = u.sources.KeptRises.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read kept rises: %w", err)
	}
	if choices.manuals, err = u.sources.Manuals.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read manual subscriptions: %w", err)
	}
	if choices.overrides, err = u.sources.Overrides.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read line overrides: %w", err)
	}
	if choices.hidden, err = u.sources.Hidden.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read hidden lines: %w", err)
	}
	if choices.cancellations, err = u.sources.Cancellations.List(ctx, userID); err != nil {
		return userChoices{}, fmt.Errorf("read cancellations: %w", err)
	}
	return choices, nil
}
