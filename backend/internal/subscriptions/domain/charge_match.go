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
// same amount and currency, same card when both show it, within MatchWindow.
func ChargesMatch(statement, email shared.Charge, identify IdentifyServiceFunc, normalize NormalizeMerchantFunc) bool {
	return statement.Bank == email.Bank &&
		statement.Amount == email.Amount &&
		statement.Currency == email.Currency &&
		isSameCard(statement.CardLast4, email.CardLast4) &&
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
	emailByStatement := pairStatementsWithEmails(charges, identify, normalize)
	return withoutIndexes(charges, emailByStatement)
}

// pairStatementsWithEmails maps the index of each statement charge to the
// index of the email charge it pairs with; each email pairs at most once.
// Both views (summary and reconciliation) use this one pairing.
func pairStatementsWithEmails(charges []shared.Charge, identify IdentifyServiceFunc, normalize NormalizeMerchantFunc) map[int]int {
	emailByStatement := make(map[int]int)
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
		emailByStatement[i] = emailIndex
	}
	return emailByStatement
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

func withoutIndexes(charges []shared.Charge, excluded map[int]int) []shared.Charge {
	kept := make([]shared.Charge, 0, len(charges)-len(excluded))
	for i, c := range charges {
		if _, isExcluded := excluded[i]; !isExcluded {
			kept = append(kept, c)
		}
	}
	return kept
}
