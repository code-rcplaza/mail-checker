package usecases

import (
	"context"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// ReconcileCharges is the differentiated view: it crosses statement and
// email charges with the same matching rule as GetSummary.
type ReconcileCharges struct {
	charges ports.ChargeReader
	catalog ports.ServiceCatalog
	now     func() time.Time
}

func NewReconcileCharges(charges ports.ChargeReader, catalog ports.ServiceCatalog, now func() time.Time) ReconcileCharges {
	return ReconcileCharges{charges: charges, catalog: catalog, now: now}
}

func (u ReconcileCharges) Execute(ctx context.Context, userID int64, months int) (domain.Reconciliation, error) {
	if err := domain.ValidatePeriodMonths(months); err != nil {
		return domain.Reconciliation{}, err
	}
	charges, err := readChargesSince(ctx, userID, u.charges, domain.PeriodStart(u.now(), months))
	if err != nil {
		return domain.Reconciliation{}, err
	}
	return domain.BuildReconciliation(charges, u.catalog.Identify, u.catalog.Normalize), nil
}
