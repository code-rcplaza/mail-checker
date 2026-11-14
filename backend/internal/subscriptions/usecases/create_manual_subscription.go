package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// CreateManualSubscription adds a subscription that is not in the user's
// statements. The public id is generated here, before saving, so the result
// can be returned as it was stored.
type CreateManualSubscription struct {
	manuals     ports.ManualSubscriptions
	newPublicID func() string
}

func NewCreateManualSubscription(manuals ports.ManualSubscriptions, newPublicID func() string) CreateManualSubscription {
	return CreateManualSubscription{manuals: manuals, newPublicID: newPublicID}
}

func (u CreateManualSubscription) Execute(ctx context.Context, userID int64, input domain.ManualSubscriptionInput) (domain.ManualSubscription, error) {
	manual, err := domain.NewManualSubscription(u.newPublicID(), input)
	if err != nil {
		return domain.ManualSubscription{}, err
	}
	if err := u.manuals.Create(ctx, userID, manual); err != nil {
		return domain.ManualSubscription{}, fmt.Errorf("create manual subscription: %w", err)
	}
	return manual, nil
}
