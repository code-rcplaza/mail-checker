package domain

import (
	"fmt"
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// StatementLine is a charge the parser read, with its line number in the
// statement text, so external ids and bad lines point to the same place.
type StatementLine struct {
	Number int
	Charge shared.ParsedCharge
}

// StatementSource is where the charges of a statement come from: whose it
// is, which bank, its public id and its billing date.
type StatementSource struct {
	UserID      int64
	Bank        shared.Bank
	PublicID    string
	BillingDate time.Time
}

// ChargesFromLines turns the parsed lines of a statement into charges. A line
// that is not a valid charge is reported by its number (the front words it)
// and the others still convert.
func ChargesFromLines(lines []StatementLine, source StatementSource) (charges []shared.Charge, badLines []int) {
	for _, line := range lines {
		charge, err := line.Charge.ToCharge(statementOrigin(source, line.Number))
		if err != nil {
			badLines = append(badLines, line.Number)
			continue
		}
		charges = append(charges, charge)
	}
	return charges, badLines
}

// statementOrigin builds the external id from the public id, because the
// internal id does not exist until the statement is saved.
func statementOrigin(source StatementSource, lineNumber int) shared.ChargeOrigin {
	return shared.ChargeOrigin{
		UserID:       source.UserID,
		Bank:         source.Bank,
		Source:       shared.SourceStatement,
		ExternalID:   fmt.Sprintf("%s:%d", source.PublicID, lineNumber),
		FallbackDate: source.BillingDate,
	}
}

// MergeBadLines joins the lines the parser could not read with those that
// were not valid charges, in a new slice ordered by line number.
func MergeBadLines(unread, invalid []int) []int {
	merged := slices.Concat(unread, invalid)
	slices.Sort(merged)
	return merged
}
