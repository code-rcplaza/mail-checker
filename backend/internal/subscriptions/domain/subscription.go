package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Subscription is a detected group with its current amount and, for a
// possible one, the user's decision.
type Subscription struct {
	Key        GroupKey // the line
	Parent     GroupKey // the service, umbrella or merchant the line belongs to
	Alias      string   // the user's name for the line; empty if none
	Status     Status
	Cadence    Cadence
	Amount     shared.Money
	Currency   shared.Currency
	Decision   *Decision    // nil while undecided
	InPesos    *PesosAmount // nil when it is in dollars and there is no exchange rate
	NextCharge time.Time    // zero while the cadence is unknown or the line is inactive
	Charges    []shared.Charge
}
