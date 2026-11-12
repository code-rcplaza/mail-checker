package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Subscription is a detected group with its current amount and, for a
// possible one, the user's decision.
type Subscription struct {
	Key      GroupKey
	Status   Status
	Cadence  Cadence
	Amount   shared.Money
	Currency shared.Currency
	Decision *Decision // nil while undecided
	Charges  []shared.Charge
}
