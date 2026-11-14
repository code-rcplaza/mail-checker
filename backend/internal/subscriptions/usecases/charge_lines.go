package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// ChargeLines reads a user's charges and groups them into lines. Every use
// case that reads or checks lines composes it, so a line the user sees in the
// summary is the same one they can decide on, name, override or hide.
type ChargeLines struct {
	Charges ports.ChargeReader
	Catalog ports.ServiceCatalog
	Now     func() time.Time
}

// readGroupsSince reads the charges, counts each purchase once and groups them.
func (l ChargeLines) readGroupsSince(ctx context.Context, userID int64, since time.Time) ([]domain.ChargeGroup, error) {
	charges, err := l.readChargesSince(ctx, userID, since)
	if err != nil {
		return nil, err
	}
	unique := domain.DropStatementChargesMatchingEmail(charges, l.merchantCatalog())
	return domain.GroupCharges(unique, l.merchantCatalog()), nil
}

// readAllGroups looks over the widest period, so any line the user sees in
// the period they picked can be acted on.
func (l ChargeLines) readAllGroups(ctx context.Context, userID int64) ([]domain.ChargeGroup, error) {
	return l.readGroupsSince(ctx, userID, domain.PeriodStart(l.Now(), domain.MaxPeriodMonths))
}

// ensureLine fails with not found for a key that is not a line of the user;
// a line of another user is not among their groups.
func (l ChargeLines) ensureLine(ctx context.Context, userID int64, key domain.GroupKey) error {
	groups, err := l.readAllGroups(ctx, userID)
	if err != nil {
		return err
	}
	if !domain.IsLine(groups, key) {
		return shared.NotFound(domain.KeyLineNotFound)
	}
	return nil
}

func (l ChargeLines) readChargesSince(ctx context.Context, userID int64, since time.Time) ([]shared.Charge, error) {
	charges, err := l.Charges.ListSince(ctx, userID, since)
	if err != nil {
		return nil, fmt.Errorf("read charges since %s: %w", since.Format(time.DateOnly), err)
	}
	return charges, nil
}

func (l ChargeLines) merchantCatalog() domain.MerchantCatalog {
	return domain.MerchantCatalog{Identify: l.Catalog.Identify, Normalize: l.Catalog.Normalize}
}

// detector classifies as of the given moment, so one use case run sees one "now".
func (l ChargeLines) detector(now time.Time) domain.SubscriptionDetector {
	return domain.SubscriptionDetector{LooksLikeSubscription: l.Catalog.LooksLikeSubscription, Now: now}
}
