package usecases

import (
	"context"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// KeepPriceRise stores that the user keeps a line after its rise, so the
// notice for that amount does not show again. In the MVP "keep" is the only
// answer to a rise.
type KeepPriceRise struct {
	lines     ChargeLines
	keptRises ports.PriceAlertDecisions
}

func NewKeepPriceRise(lines ChargeLines, keptRises ports.PriceAlertDecisions) KeepPriceRise {
	return KeepPriceRise{lines: lines, keptRises: keptRises}
}

func (u KeepPriceRise) Execute(ctx context.Context, userID int64, key domain.GroupKey) error {
	rise, err := u.findRise(ctx, userID, key)
	if err != nil {
		return err
	}
	return u.saveKeptRise(ctx, userID, domain.KeptRise{Key: rise.Key, NewAmount: rise.After})
}

// findRise looks over the widest period, like the other line decisions.
func (u KeepPriceRise) findRise(ctx context.Context, userID int64, key domain.GroupKey) (domain.Rise, error) {
	groups, err := u.lines.readAllGroups(ctx, userID)
	if err != nil {
		return domain.Rise{}, err
	}
	rise, hasRise := domain.LineRise(groups, key)
	if !hasRise {
		return domain.Rise{}, shared.NotFound(domain.KeyRiseNotFound)
	}
	return rise, nil
}

func (u KeepPriceRise) saveKeptRise(ctx context.Context, userID int64, kept domain.KeptRise) error {
	if err := u.keptRises.Save(ctx, userID, kept); err != nil {
		return fmt.Errorf("save kept rise for %s: %w", kept.Key, err)
	}
	return nil
}
