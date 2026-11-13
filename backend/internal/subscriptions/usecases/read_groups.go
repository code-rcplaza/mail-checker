package usecases

import (
	"context"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// readGroupsSince reads the charges, counts each purchase once and groups
// them. The summary and the decision use case share it, so a possible the
// user sees is the same one they can decide on.
func readGroupsSince(ctx context.Context, userID int64, reader ports.ChargeReader, catalog ports.ServiceCatalog, since time.Time) ([]domain.ChargeGroup, error) {
	charges, err := readChargesSince(ctx, userID, reader, since)
	if err != nil {
		return nil, err
	}
	unique := domain.DropStatementChargesMatchingEmail(charges, catalog.Identify, catalog.Normalize)
	return domain.GroupCharges(unique, catalog.Identify, catalog.Normalize), nil
}

// ensureLineExists looks over the widest period, so any line the user sees in
// the period they picked can be named, overridden or hidden. A line of
// another user is not among their groups, so it is not found either.
func ensureLineExists(ctx context.Context, userID int64, reader ports.ChargeReader, catalog ports.ServiceCatalog, now time.Time, key domain.GroupKey) error {
	groups, err := readGroupsSince(ctx, userID, reader, catalog, domain.PeriodStart(now, domain.MaxPeriodMonths))
	if err != nil {
		return err
	}
	if !domain.IsLine(groups, key) {
		return shared.NotFound(domain.KeyLineNotFound)
	}
	return nil
}
