package usecases

import (
	"context"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// ReconcileCharges is the differentiated view: it crosses statement and
// email charges with the same matching rule as GetSummary.
type ReconcileCharges struct {
	lines ChargeLines
}

func NewReconcileCharges(lines ChargeLines) ReconcileCharges {
	return ReconcileCharges{lines: lines}
}

func (u ReconcileCharges) Execute(ctx context.Context, userID int64, months int) (domain.Reconciliation, error) {
	if err := domain.ValidatePeriodMonths(months); err != nil {
		return domain.Reconciliation{}, err
	}
	charges, err := u.lines.readChargesSince(ctx, userID, domain.PeriodStart(u.lines.Now(), months))
	if err != nil {
		return domain.Reconciliation{}, err
	}
	return domain.BuildReconciliation(charges, u.lines.merchantCatalog()), nil
}
