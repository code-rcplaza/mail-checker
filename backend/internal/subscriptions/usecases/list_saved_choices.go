package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// SavedChoices is everything the user has here, as stored: their charges and
// what they decided, named or typed. It is what the copy of their data carries.
type SavedChoices struct {
	Charges   []shared.Charge
	Decisions map[domain.GroupKey]domain.Decision
	Aliases   map[domain.GroupKey]string
	Manuals   []domain.ManualSubscription
	Overrides map[domain.GroupKey]domain.LineOverride
}

type SavedChoicesSources struct {
	Charges   ports.ChargeReader
	Decisions ports.SubscriptionDecisions
	Aliases   ports.LineAliases
	Manuals   ports.ManualSubscriptions
	Overrides ports.LineOverrides
}

// ListSavedChoices reads what the user stored, without grouping nor detecting.
type ListSavedChoices struct {
	sources SavedChoicesSources
}

func NewListSavedChoices(sources SavedChoicesSources) ListSavedChoices {
	return ListSavedChoices{sources: sources}
}

// everCharged is the start of every charge the user could have.
var everCharged = time.Unix(0, 0).UTC()

func (u ListSavedChoices) Execute(ctx context.Context, userID int64) (SavedChoices, error) {
	var saved SavedChoices
	var err error
	if saved.Charges, err = u.sources.Charges.ListSince(ctx, userID, everCharged); err != nil {
		return SavedChoices{}, fmt.Errorf("read charges: %w", err)
	}
	if saved.Decisions, err = u.sources.Decisions.List(ctx, userID); err != nil {
		return SavedChoices{}, fmt.Errorf("read decisions: %w", err)
	}
	if saved.Aliases, err = u.sources.Aliases.List(ctx, userID); err != nil {
		return SavedChoices{}, fmt.Errorf("read aliases: %w", err)
	}
	if saved.Manuals, err = u.sources.Manuals.List(ctx, userID); err != nil {
		return SavedChoices{}, fmt.Errorf("read manual subscriptions: %w", err)
	}
	if saved.Overrides, err = u.sources.Overrides.List(ctx, userID); err != nil {
		return SavedChoices{}, fmt.Errorf("read overrides: %w", err)
	}
	return saved, nil
}
