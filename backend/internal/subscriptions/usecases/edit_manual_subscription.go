package usecases

import (
	"context"
	"errors"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// ManualSubscriptionEdit is what the user typed again for an existing manual.
type ManualSubscriptionEdit struct {
	PublicID string
	Input    domain.ManualSubscriptionInput
}

// EditManualSubscription replaces what the user typed, with the same
// validation as creating it; the public id stays.
type EditManualSubscription struct {
	manuals ports.ManualSubscriptions
}

func NewEditManualSubscription(manuals ports.ManualSubscriptions) EditManualSubscription {
	return EditManualSubscription{manuals: manuals}
}

func (u EditManualSubscription) Execute(ctx context.Context, userID int64, edit ManualSubscriptionEdit) error {
	manual, err := domain.NewManualSubscription(edit.PublicID, edit.Input)
	if err != nil {
		return err
	}
	return manualNotFoundOr(u.manuals.Update(ctx, userID, manual), "update manual subscription")
}

// RemoveManualSubscription deletes a manual subscription of the user.
type RemoveManualSubscription struct {
	manuals ports.ManualSubscriptions
}

func NewRemoveManualSubscription(manuals ports.ManualSubscriptions) RemoveManualSubscription {
	return RemoveManualSubscription{manuals: manuals}
}

func (u RemoveManualSubscription) Execute(ctx context.Context, userID int64, publicID string) error {
	return manualNotFoundOr(u.manuals.Delete(ctx, userID, publicID), "delete manual subscription")
}

// manualNotFoundOr adds the i18n key to the repository's plain ErrNotFound,
// so adapters stay free of user-facing keys; other errors get context.
func manualNotFoundOr(err error, action string) error {
	if errors.Is(err, shared.ErrNotFound) {
		return shared.NotFound(domain.KeyManualNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	return nil
}
