package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// PlanChangeThresholdPercent: a jump above this between consecutive recurring
// charges is a plan change; at or below it is ordinary price variation.
const PlanChangeThresholdPercent = 25

type PlanChange struct {
	From         shared.Money
	To           shared.Money
	FromCurrency shared.Currency
	ToCurrency   shared.Currency
	At           time.Time
}

type KnownServiceAmount struct {
	Current          shared.Money // amount of the last recurring charge
	Currency         shared.Currency
	PlanChanges      []PlanChange
	RecurringCharges []shared.Charge
	OtherCharges     []shared.Charge // off-pattern charges, e.g. API credits; outside monthly spend
}

// AnalyzeKnownServiceAmount works out the current amount and plan changes of
// a known service, keeping one-off charges of the same service apart.
func AnalyzeKnownServiceAmount(charges []shared.Charge) KnownServiceAmount {
	recurring, other := splitRecurringFromOther(sortByDate(charges))
	if len(recurring) == 0 {
		return KnownServiceAmount{}
	}
	last := recurring[len(recurring)-1]
	return KnownServiceAmount{
		Current:          last.Amount,
		Currency:         last.Currency,
		PlanChanges:      findPlanChanges(recurring),
		RecurringCharges: recurring,
		OtherCharges:     other,
	}
}

func findPlanChanges(recurring []shared.Charge) []PlanChange {
	var changes []PlanChange
	for i := 1; i < len(recurring); i++ {
		previous, current := recurring[i-1], recurring[i]
		if !isPlanChange(previous, current) {
			continue
		}
		changes = append(changes, PlanChange{
			From:         previous.Amount,
			To:           current.Amount,
			FromCurrency: previous.Currency,
			ToCurrency:   current.Currency,
			At:           current.OccurredAt,
		})
	}
	return changes
}

// isPlanChange: amounts in different currencies cannot be compared, so a
// currency change is a plan change by itself.
func isPlanChange(previous, current shared.Charge) bool {
	return previous.Currency != current.Currency || isAmountJump(previous.Amount, current.Amount)
}

// isAmountJump compares in integers to avoid floats: |to-from| / from > 25 %.
func isAmountJump(from, to shared.Money) bool {
	difference := to - from
	if difference < 0 {
		difference = -difference
	}
	return difference*100 > from*PlanChangeThresholdPercent
}
