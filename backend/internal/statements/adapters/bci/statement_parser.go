// Package bci reads BCI credit card statements, from the text the PDF
// extractor writes: one line per text line, columns separated by tabs.
package bci

import (
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

const (
	Bank       shared.Bank = "bci"
	CardCredit shared.Card = "credito"
)

// Markers are compared without spaces, so a statement whose words came out
// glued still matches. The titles are the CMF format every bank shares; the
// blocking phone is what makes it BCI's.
const (
	nationalTitle      = "ESTADODECUENTANACIONALDETARJETADECRÉDITO"
	internationalTitle = "ESTADODECUENTAINTERNACIONALDETARJETADECRÉDITO"
	bciBlockingPhone   = "800201090PARASOLICITARELBLOQUEO"
)

// StatementParser reads the national and international statements of one
// BCI file as a single list of purchases. Dates are midnight in location.
type StatementParser struct {
	location *time.Location
}

func NewStatementParser(location *time.Location) StatementParser {
	return StatementParser{location: location}
}

func (StatementParser) Bank() shared.Bank { return Bank }

func (StatementParser) Detect(text string) bool {
	compact := withoutSpaces(text)
	hasTitle := strings.Contains(compact, nationalTitle) || strings.Contains(compact, internationalTitle)
	return hasTitle && strings.Contains(compact, bciBlockingPhone)
}

// Parse returns the purchases with their line number. Payments, interest,
// fees, debt transfers and installments are skipped; a line shaped like a
// purchase that cannot be read is reported by number only.
func (p StatementParser) Parse(text string) (lines []domain.StatementLine, bad []int) {
	cardLast4 := readCardLast4(text)
	current := sectionNone
	for index, raw := range strings.Split(text, "\n") {
		current = nextSection(current, raw)
		number := index + 1
		charge, outcome := p.readLine(current, raw, cardLast4)
		switch outcome {
		case linePurchase:
			lines = append(lines, domain.StatementLine{Number: number, Charge: charge})
		case lineUnreadable:
			bad = append(bad, number)
		}
	}
	return lines, bad
}

func (p StatementParser) readLine(current section, raw, cardLast4 string) (shared.ParsedCharge, lineOutcome) {
	switch current {
	case sectionNationalPurchases:
		return p.readNationalLine(raw, cardLast4)
	case sectionInternational:
		return p.readInternationalLine(raw, cardLast4)
	}
	return shared.ParsedCharge{}, lineSkipped
}

func withoutSpaces(text string) string {
	return strings.Join(strings.Fields(text), "")
}
