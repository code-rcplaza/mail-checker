// Package ports declares what the subscriptions use cases need from the outside.
package ports

import (
	"context"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type ChargeReader interface {
	ListSince(ctx context.Context, since time.Time) ([]shared.Charge, error)
}

type ServiceCatalog interface {
	Identify(merchant string) (domain.Service, bool)
	Normalize(merchant string) string
	LooksLikeSubscription(merchant string) bool
}

type SubscriptionDecisions interface {
	List(ctx context.Context) (map[domain.GroupKey]domain.Decision, error)
}
