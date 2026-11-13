package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// KeepPriceRise stores that the user keeps a line after its rise, so the
// notice for that amount does not show again. In the MVP "keep" is the only
// answer to a rise.
type KeepPriceRise struct {
	charges   ports.ChargeReader
	catalog   ports.ServiceCatalog
	keptRises ports.PriceAlertDecisions
	now       func() time.Time
}

func NewKeepPriceRise(
	charges ports.ChargeReader,
	catalog ports.ServiceCatalog,
	keptRises ports.PriceAlertDecisions,
	now func() time.Time,
) KeepPriceRise {
	return KeepPriceRise{charges: charges, catalog: catalog, keptRises: keptRises, now: now}
}

func (u KeepPriceRise) Execute(ctx context.Context, userID int64, key domain.GroupKey) error {
	rise, err := u.findRise(ctx, userID, key)
	if err != nil {
		return err
	}
	return u.saveKeptRise(ctx, userID, rise)
}

// findRise looks over the widest period, like the other line decisions.
func (u KeepPriceRise) findRise(ctx context.Context, userID int64, key domain.GroupKey) (domain.Rise, error) {
	groups, err := readGroupsSince(ctx, userID, u.charges, u.catalog, domain.PeriodStart(u.now(), domain.MaxPeriodMonths))
	if err != nil {
		return domain.Rise{}, err
	}
	rise, hasRise := domain.LineRise(groups, key)
	if !hasRise {
		return domain.Rise{}, shared.NotFound(domain.KeyRiseNotFound)
	}
	return rise, nil
}

func (u KeepPriceRise) saveKeptRise(ctx context.Context, userID int64, rise domain.Rise) error {
	if err := u.keptRises.Save(ctx, userID, rise.Key, rise.After); err != nil {
		return fmt.Errorf("save kept rise for %s: %w", rise.Key, err)
	}
	return nil
}
