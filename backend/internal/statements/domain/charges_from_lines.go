package domain

import (
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// ChargesFromLines turns the parsed lines of a statement into charges. A line
// that is not a valid charge is reported by its number (the front words it)
// and the others still convert. Lines are numbered from 1, in the order the
// parser returned them.
func ChargesFromLines(lines []shared.ParsedCharge, bank shared.Bank, publicID string, billingDate time.Time) (charges []shared.Charge, badLines []int) {
	for i, parsed := range lines {
		lineNumber := i + 1
		charge, err := parsed.ToCharge(statementOrigin(bank, publicID, lineNumber, billingDate))
		if err != nil {
			badLines = append(badLines, lineNumber)
			continue
		}
		charges = append(charges, charge)
	}
	return charges, badLines
}

// statementOrigin builds the external id from the public id, because the
// internal id does not exist until the statement is saved.
func statementOrigin(bank shared.Bank, publicID string, lineNumber int, billingDate time.Time) shared.ChargeOrigin {
	return shared.ChargeOrigin{
		Bank:         bank,
		Source:       shared.SourceStatement,
		ExternalID:   fmt.Sprintf("%s:%d", publicID, lineNumber),
		FallbackDate: billingDate,
	}
}
