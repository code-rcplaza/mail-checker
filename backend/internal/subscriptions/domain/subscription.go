package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Subscription is a detected group with its current amount and, for a
// possible one, the user's decision.
type Subscription struct {
	Key      GroupKey // the line
	Parent   GroupKey // the service, umbrella or merchant the line belongs to
	Status   Status
	Cadence  Cadence
	Amount   shared.Money
	Currency shared.Currency
	Decision *Decision // nil while undecided
	Charges  []shared.Charge
}
