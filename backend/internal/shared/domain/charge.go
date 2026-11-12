package domain

import (
	"strings"
	"time"
)

type Bank string

type Card string

type Source string

const (
	SourceEmail     Source = "email"
	SourceStatement Source = "statement"
)

// Charge is a stored card charge. It is passed by value; the optional
// account and statement links are pointers.
type Charge struct {
	ID          int64
	Bank        Bank
	Source      Source
	ExternalID  string // email: Message-ID; statement: "<statementPublicID>:<line>"
	AccountID   *int64
	StatementID *int64
	OccurredAt  time.Time
	Amount      Money
	Currency    Currency
	Merchant    string
	Card        Card
}

// NewCharge returns c only if it holds the charge invariants, so no other
// layer has to revalidate amounts, dates or merchants.
func NewCharge(c Charge) (Charge, error) {
	if err := validateAmountNotNegative(c.Amount); err != nil {
		return Charge{}, err
	}
	if err := validateDatePresent(c.OccurredAt); err != nil {
		return Charge{}, err
	}
	if err := validateMerchantPresent(c.Merchant); err != nil {
		return Charge{}, err
	}
	return c, nil
}

func validateAmountNotNegative(amount Money) error {
	if amount < 0 {
		return Invalid(KeyAmountNegative, nil)
	}
	return nil
}

func validateDatePresent(date time.Time) error {
	if date.IsZero() {
		return Invalid(KeyDateMissing, nil)
	}
	return nil
}

func validateMerchantPresent(merchant string) error {
	if strings.TrimSpace(merchant) == "" {
		return Invalid(KeyMerchantMissing, nil)
	}
	return nil
}
