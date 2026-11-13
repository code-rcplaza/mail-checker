package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// SaveLineOverride stores the amount or billing day the user set for a line;
// an empty override removes it. Real charges are never changed.
type SaveLineOverride struct {
	charges   ports.ChargeReader
	catalog   ports.ServiceCatalog
	overrides ports.LineOverrides
	now       func() time.Time
}

func NewSaveLineOverride(
	charges ports.ChargeReader,
	catalog ports.ServiceCatalog,
	overrides ports.LineOverrides,
	now func() time.Time,
) SaveLineOverride {
	return SaveLineOverride{charges: charges, catalog: catalog, overrides: overrides, now: now}
}

func (u SaveLineOverride) Execute(ctx context.Context, userID int64, key domain.GroupKey, amount *shared.Money, billingDay *int) error {
	override, err := domain.NewLineOverride(amount, billingDay)
	if err != nil {
		return err
	}
	if err := ensureLineExists(ctx, userID, u.charges, u.catalog, u.now(), key); err != nil {
		return err
	}
	return u.storeOverride(ctx, userID, key, override)
}

func (u SaveLineOverride) storeOverride(ctx context.Context, userID int64, key domain.GroupKey, override domain.LineOverride) error {
	if override.IsEmpty() {
		return u.deleteOverride(ctx, userID, key)
	}
	if err := u.overrides.Save(ctx, userID, key, override); err != nil {
		return fmt.Errorf("save override for %s: %w", key, err)
	}
	return nil
}

func (u SaveLineOverride) deleteOverride(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.overrides.Delete(ctx, userID, key); err != nil {
		return fmt.Errorf("delete override for %s: %w", key, err)
	}
	return nil
}
