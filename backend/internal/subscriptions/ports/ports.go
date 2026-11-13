// Package ports declares what the subscriptions use cases need from the outside.
package ports

import (
	"context"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// Every port that reads or writes user data takes the user right after ctx
// and only ever touches that user's rows.
type ChargeReader interface {
	ListSince(ctx context.Context, userID int64, since time.Time) ([]shared.Charge, error)
}

type ServiceCatalog interface {
	Identify(merchant string) (domain.Service, bool)
	Normalize(merchant string) string
	LooksLikeSubscription(merchant string) bool
}

type SubscriptionDecisions interface {
	List(ctx context.Context, userID int64) (map[domain.GroupKey]domain.Decision, error)
	Save(ctx context.Context, userID int64, key domain.GroupKey, decision domain.Decision) error
	// Delete of a key without a decision is not an error.
	Delete(ctx context.Context, userID int64, key domain.GroupKey) error
}

// ExchangeRates lists the rates the user's statements showed.
type ExchangeRates interface {
	List(ctx context.Context, userID int64) ([]shared.ExchangeRate, error)
}

// PriceAlertDecisions keeps, per line, the new amount of a rise the user
// chose to keep, so its notice does not show again; saving again replaces it.
type PriceAlertDecisions interface {
	List(ctx context.Context, userID int64) (map[domain.GroupKey]shared.Money, error)
	Save(ctx context.Context, userID int64, key domain.GroupKey, newAmount shared.Money) error
}

// ManualSubscriptions are the ones the user added by hand.
type ManualSubscriptions interface {
	List(ctx context.Context, userID int64) ([]domain.ManualSubscription, error)
	Create(ctx context.Context, userID int64, manual domain.ManualSubscription) error
	// Update and Delete find the manual by public id; one that is unknown or
	// belongs to another user is shared.ErrNotFound.
	Update(ctx context.Context, userID int64, manual domain.ManualSubscription) error
	Delete(ctx context.Context, userID int64, publicID string) error
}

// LineOverrides keep the amount or billing day the user set for a line.
type LineOverrides interface {
	List(ctx context.Context, userID int64) (map[domain.GroupKey]domain.LineOverride, error)
	Save(ctx context.Context, userID int64, key domain.GroupKey, override domain.LineOverride) error
	Delete(ctx context.Context, userID int64, key domain.GroupKey) error
}

type LineAliases interface {
	List(ctx context.Context, userID int64) (map[domain.GroupKey]string, error)
	Save(ctx context.Context, userID int64, key domain.GroupKey, alias string) error
	Delete(ctx context.Context, userID int64, key domain.GroupKey) error
}
