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
	UserID      int64 // the owner; every read filters by it
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
	CardLast4   string // last four card digits; empty when the bank does not show them
}

// NewCharge returns c only if it holds the charge invariants, so no other
// layer has to revalidate amounts, dates or merchants.
func NewCharge(c Charge) (Charge, error) {
	if err := validateUserPresent(c.UserID); err != nil {
		return Charge{}, err
	}
	if err := validateAmountNotNegative(c.Amount); err != nil {
		return Charge{}, err
	}
	if err := validateDatePresent(c.OccurredAt); err != nil {
		return Charge{}, err
	}
	if err := validateMerchantPresent(c.Merchant); err != nil {
		return Charge{}, err
	}
	if err := validateCardLast4(c.CardLast4); err != nil {
		return Charge{}, err
	}
	return c, nil
}

// validateUserPresent: a charge without an owner could be shown to anyone.
func validateUserPresent(userID int64) error {
	if userID <= 0 {
		return Invalid(KeyUserMissing, nil)
	}
	return nil
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

// cardLast4Length is how many card digits a bank shows; the full number is never stored.
const cardLast4Length = 4

func validateCardLast4(last4 string) error {
	if last4 == "" || isFourDigits(last4) {
		return nil
	}
	return Invalid(KeyCardLastDigitsInvalid, nil)
}

func isFourDigits(text string) bool {
	if len(text) != cardLast4Length {
		return false
	}
	return !strings.ContainsFunc(text, func(r rune) bool { return r < '0' || r > '9' })
}

func validateMerchantPresent(merchant string) error {
	if strings.TrimSpace(merchant) == "" {
		return Invalid(KeyMerchantMissing, nil)
	}
	return nil
}
