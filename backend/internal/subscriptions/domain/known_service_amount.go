package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// PlanChangeThresholdPercent: a jump above this between consecutive recurring
// charges is a plan change; at or below it is ordinary price variation.
const PlanChangeThresholdPercent = 25

type PlanChange struct {
	From shared.Money
	To   shared.Money
	At   time.Time
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

// splitRecurringFromOther treats a charge arriving sooner than a monthly
// interval after the last recurring one as off-pattern: it cannot be the
// next bill.
func splitRecurringFromOther(chronological []shared.Charge) (recurring, other []shared.Charge) {
	for _, c := range chronological {
		if len(recurring) > 0 && arrivesTooSoon(recurring[len(recurring)-1], c) {
			other = append(other, c)
			continue
		}
		recurring = append(recurring, c)
	}
	return recurring, other
}

func arrivesTooSoon(previous, next shared.Charge) bool {
	return daysBetween(previous.OccurredAt, next.OccurredAt) < MonthlyIntervalMinDays
}

func findPlanChanges(recurring []shared.Charge) []PlanChange {
	var changes []PlanChange
	for i := 1; i < len(recurring); i++ {
		previous, current := recurring[i-1], recurring[i]
		if !isPlanChange(previous.Amount, current.Amount) {
			continue
		}
		changes = append(changes, PlanChange{From: previous.Amount, To: current.Amount, At: current.OccurredAt})
	}
	return changes
}

// isPlanChange compares in integers to avoid floats: |to-from| / from > 25 %.
func isPlanChange(from, to shared.Money) bool {
	difference := to - from
	if difference < 0 {
		difference = -difference
	}
	return difference*100 > from*PlanChangeThresholdPercent
}
