package usecases

import (
	"context"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// LineOverrideInput is what the user sent for a line; nil means "not given".
type LineOverrideInput struct {
	Key        domain.GroupKey
	Amount     *shared.Money
	BillingDay *int
}

// SaveLineOverride stores the amount or billing day the user set for a line;
// an empty override removes it. Real charges are never changed.
type SaveLineOverride struct {
	lines     ChargeLines
	overrides ports.LineOverrides
}

func NewSaveLineOverride(lines ChargeLines, overrides ports.LineOverrides) SaveLineOverride {
	return SaveLineOverride{lines: lines, overrides: overrides}
}

func (u SaveLineOverride) Execute(ctx context.Context, userID int64, input LineOverrideInput) error {
	override, err := domain.NewLineOverride(input.Amount, input.BillingDay)
	if err != nil {
		return err
	}
	if err := u.lines.ensureLine(ctx, userID, input.Key); err != nil {
		return err
	}
	return u.storeOverride(ctx, userID, domain.OverriddenLine{Key: input.Key, Override: override})
}

func (u SaveLineOverride) storeOverride(ctx context.Context, userID int64, overridden domain.OverriddenLine) error {
	if overridden.Override.IsEmpty() {
		return u.deleteOverride(ctx, userID, overridden.Key)
	}
	if err := u.overrides.Save(ctx, userID, overridden); err != nil {
		return fmt.Errorf("save override for %s: %w", overridden.Key, err)
	}
	return nil
}

func (u SaveLineOverride) deleteOverride(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.overrides.Delete(ctx, userID, key); err != nil {
		return fmt.Errorf("delete override for %s: %w", key, err)
	}
	return nil
}
