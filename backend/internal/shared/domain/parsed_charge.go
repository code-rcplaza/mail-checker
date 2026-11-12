package domain

import (
	"strings"
	"time"
)

// ParsedCharge is what a bank parser extracts from a purchase email or a
// statement line.
type ParsedCharge struct {
	Amount   Money
	Currency Currency
	Merchant string
	Card     Card
	Date     time.Time // zero when a statement line has no date
}

// ChargeOrigin is where a parsed charge came from.
type ChargeOrigin struct {
	Bank         Bank
	Source       Source
	ExternalID   string
	AccountID    *int64
	FallbackDate time.Time // used when the parsed charge has no date
}

// ToCharge is the single conversion for email and statement charges. The
// external id is checked here; the rest of the invariants are NewCharge's.
func (p ParsedCharge) ToCharge(origin ChargeOrigin) (Charge, error) {
	if err := validateExternalIDPresent(origin.ExternalID); err != nil {
		return Charge{}, err
	}
	return NewCharge(Charge{
		Bank:       origin.Bank,
		Source:     origin.Source,
		ExternalID: origin.ExternalID,
		AccountID:  origin.AccountID,
		OccurredAt: occurredAtOrFallback(p.Date, origin.FallbackDate),
		Amount:     p.Amount,
		Currency:   p.Currency,
		Merchant:   p.Merchant,
		Card:       p.Card,
	})
}

func validateExternalIDPresent(externalID string) error {
	if strings.TrimSpace(externalID) == "" {
		return Invalid(KeyExternalIDMissing, nil)
	}
	return nil
}

func occurredAtOrFallback(parsed, fallback time.Time) time.Time {
	if parsed.IsZero() {
		return fallback
	}
	return parsed
}
