package usecases

import (
	"context"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// readGroupsSince reads the charges, counts each purchase once and groups
// them. The summary and the decision use case share it, so a possible the
// user sees is the same one they can decide on.
func readGroupsSince(ctx context.Context, reader ports.ChargeReader, catalog ports.ServiceCatalog, since time.Time) ([]domain.ChargeGroup, error) {
	charges, err := readChargesSince(ctx, reader, since)
	if err != nil {
		return nil, err
	}
	unique := domain.DropStatementChargesMatchingEmail(charges, catalog.Identify, catalog.Normalize)
	return domain.GroupCharges(unique, catalog.Identify, catalog.Normalize), nil
}
