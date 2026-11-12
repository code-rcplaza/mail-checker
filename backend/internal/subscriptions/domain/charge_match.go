package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// MatchWindow is how far apart a statement line and its purchase email can
// be: the statement often books the charge days after the purchase.
const MatchWindow = 35 * 24 * time.Hour

// ChargesMatch tells whether a statement charge and an email charge are the
// same purchase: same bank, same service (the identity used for grouping),
// same amount and currency, within MatchWindow.
func ChargesMatch(statement, email shared.Charge, identify IdentifyServiceFunc, normalize NormalizeMerchantFunc) bool {
	return statement.Bank == email.Bank &&
		statement.Amount == email.Amount &&
		statement.Currency == email.Currency &&
		isWithinMatchWindow(statement.OccurredAt, email.OccurredAt) &&
		groupKeyOf(statement, identify, normalize) == groupKeyOf(email, identify, normalize)
}

func isWithinMatchWindow(a, b time.Time) bool {
	return a.Sub(b).Abs() <= MatchWindow
}

// DropStatementChargesMatchingEmail removes the statement charges already
// counted through an email, pairing one to one, so each purchase counts once.
// It returns a new slice in the original order.
func DropStatementChargesMatchingEmail(charges []shared.Charge, identify IdentifyServiceFunc, normalize NormalizeMerchantFunc) []shared.Charge {
	duplicates := findStatementDuplicates(charges, identify, normalize)
	return withoutIndexes(charges, duplicates)
}

// findStatementDuplicates returns the indexes of statement charges paired
// with an email charge; each email charge pairs at most once.
func findStatementDuplicates(charges []shared.Charge, identify IdentifyServiceFunc, normalize NormalizeMerchantFunc) map[int]bool {
	duplicates := make(map[int]bool)
	pairedEmails := make(map[int]bool)
	for i, statement := range charges {
		if statement.Source != shared.SourceStatement {
			continue
		}
		emailIndex, isFound := findUnpairedEmailMatch(statement, charges, pairedEmails, identify, normalize)
		if !isFound {
			continue
		}
		pairedEmails[emailIndex] = true
		duplicates[i] = true
	}
	return duplicates
}

func findUnpairedEmailMatch(
	statement shared.Charge,
	charges []shared.Charge,
	pairedEmails map[int]bool,
	identify IdentifyServiceFunc,
	normalize NormalizeMerchantFunc,
) (int, bool) {
	for i, email := range charges {
		if email.Source != shared.SourceEmail || pairedEmails[i] {
			continue
		}
		if ChargesMatch(statement, email, identify, normalize) {
			return i, true
		}
	}
	return 0, false
}

func withoutIndexes(charges []shared.Charge, excluded map[int]bool) []shared.Charge {
	kept := make([]shared.Charge, 0, len(charges)-len(excluded))
	for i, c := range charges {
		if !excluded[i] {
			kept = append(kept, c)
		}
	}
	return kept
}
