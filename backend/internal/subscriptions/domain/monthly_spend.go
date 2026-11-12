package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const monthsPerYear = 12

// MonthlySpend adds what active subscriptions and confirmed possible ones
// cost per month, per currency. A yearly one counts as a twelfth.
func MonthlySpend(subscriptions []Subscription) map[shared.Currency]shared.Money {
	spend := make(map[shared.Currency]shared.Money)
	for _, s := range subscriptions {
		cadence, isCounted := spendCadence(s)
		if !isCounted {
			continue
		}
		spend[s.Currency] += monthlyAmount(s.Amount, cadence)
	}
	return spend
}

// spendCadence tells whether a subscription adds to spend and with which
// cadence: its own when active, the user's choice when a possible was confirmed.
func spendCadence(s Subscription) (Cadence, bool) {
	if s.Status == StatusActive {
		return s.Cadence, true
	}
	if s.Status == StatusPossible && isConfirmed(s.Decision) {
		return s.Decision.Cadence, true
	}
	return CadenceNone, false
}

func monthlyAmount(amount shared.Money, cadence Cadence) shared.Money {
	if cadence == CadenceYearly {
		return divideRoundingHalfUp(amount, monthsPerYear)
	}
	return amount
}

// divideRoundingHalfUp works for non-negative amounts, which Charge guarantees.
func divideRoundingHalfUp(amount shared.Money, divisor shared.Money) shared.Money {
	return (amount + divisor/2) / divisor
}
