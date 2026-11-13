package bci

import (
	"regexp"
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// nationalOperation is the column with date, reference and descriptor:
// "26/09/26 2809 00000006 NETFLIX.COM". The place may come before it in its
// own column, and sometimes glued after the descriptor; it is kept, since
// the catalog recognises services by how the descriptor starts.
var nationalOperation = regexp.MustCompile(`^(\d{2}/\d{2}/\d{2}) \d{4} \d{8} (\S.*)$`)

// installmentMarker is in every installment purchase, which is out of scope:
// a subscription is never paid in installments.
const installmentMarker = "TASA INT"

const statementDateLayout = "02/01/06"

func (p StatementParser) readNationalLine(raw, cardLast4 string) (shared.ParsedCharge, lineOutcome) {
	columns := strings.Split(raw, "\t")
	operationAt := indexOfOperation(columns)
	if operationAt < 0 || isInstallment(raw) {
		return shared.ParsedCharge{}, lineSkipped
	}
	match := nationalOperation.FindStringSubmatch(columns[operationAt])
	date, isDateRead := p.readDate(match[1])
	amount, isAmountRead := firstCLPAmountAfter(columns, operationAt)
	if !isDateRead || !isAmountRead {
		return shared.ParsedCharge{}, lineUnreadable
	}
	if amount < 0 {
		return shared.ParsedCharge{}, lineSkipped // a payment or a refund
	}
	return shared.ParsedCharge{
		Amount:    amount,
		Currency:  shared.CurrencyCLP,
		Merchant:  strings.TrimSpace(match[2]),
		Card:      CardCredit,
		CardLast4: cardLast4,
		Date:      date,
	}, linePurchase
}

func indexOfOperation(columns []string) int {
	for i, column := range columns {
		if nationalOperation.MatchString(column) {
			return i
		}
	}
	return -1
}

func isInstallment(raw string) bool {
	return strings.Contains(raw, installmentMarker)
}

// firstCLPAmountAfter reads the operation amount, the first "$" column; a
// missing or unreadable one makes the line unreadable.
func firstCLPAmountAfter(columns []string, operationAt int) (shared.Money, bool) {
	for _, column := range columns[operationAt+1:] {
		if strings.HasPrefix(strings.TrimSpace(column), "$") {
			return parseCLPAmount(column)
		}
	}
	return 0, false
}

func (p StatementParser) readDate(raw string) (time.Time, bool) {
	date, err := time.ParseInLocation(statementDateLayout, raw, p.location)
	if err != nil {
		return time.Time{}, false
	}
	return date, true
}
