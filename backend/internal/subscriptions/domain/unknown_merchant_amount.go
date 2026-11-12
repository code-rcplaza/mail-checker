package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// MaxPlanChangesForConsistentAmount: an unknown merchant may change plan once;
// an amount that keeps jumping is ordinary shopping, not a subscription.
const MaxPlanChangesForConsistentAmount = 1

type UnknownMerchantAmount struct {
	IsConsistent bool
	Current      shared.Money // amount of the last charge
	Currency     shared.Currency
	PlanChanges  []PlanChange
}

// AnalyzeUnknownMerchantAmount compares each charge with the previous one to
// tell whether the amount is stable enough to be a subscription.
func AnalyzeUnknownMerchantAmount(charges []shared.Charge) UnknownMerchantAmount {
	chronological := sortByDate(charges)
	if len(chronological) == 0 {
		return UnknownMerchantAmount{}
	}
	last := chronological[len(chronological)-1]
	planChanges := findPlanChanges(chronological)
	return UnknownMerchantAmount{
		IsConsistent: len(planChanges) <= MaxPlanChangesForConsistentAmount,
		Current:      last.Amount,
		Currency:     last.Currency,
		PlanChanges:  planChanges,
	}
}
