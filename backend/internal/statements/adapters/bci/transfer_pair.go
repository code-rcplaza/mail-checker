package bci

import (
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

// The international balance moves to the national statement: in pesos as a
// fee line of the national one, and in dollars in the payment information of
// the international one.
const (
	nationalTransferDescriptor = "TRASPASO DEUDA INTERNACIONAL"
	internationalTransferLabel = "TRASPASO DEUDA NACIONAL"
)

// TransferPair needs both sides and some dollars: a month without
// international debt shows US$ 0,00 and gives no rate.
func (p StatementParser) TransferPair(text string) (domain.TransferPair, bool) {
	lines := strings.Split(text, "\n")
	pesos, date, hasPesos := p.readNationalTransfer(lines)
	dollars, hasDollars := readInternationalTransfer(lines)
	if !hasPesos || !hasDollars || pesos <= 0 || dollars <= 0 {
		return domain.TransferPair{}, false
	}
	return domain.TransferPair{Pesos: pesos, Dollars: dollars, Date: date}, true
}

func (p StatementParser) readNationalTransfer(lines []string) (shared.Money, time.Time, bool) {
	for _, raw := range lines {
		columns := strings.Split(raw, "\t")
		operationAt := indexOfOperation(columns)
		if operationAt < 0 {
			continue
		}
		match := nationalOperation.FindStringSubmatch(columns[operationAt])
		if !strings.HasPrefix(match[2], nationalTransferDescriptor) {
			continue
		}
		date, isDateRead := p.readDate(match[1])
		pesos, isAmountRead := firstCLPAmountAfter(columns, operationAt)
		return pesos, date, isDateRead && isAmountRead
	}
	return 0, time.Time{}, false
}

func readInternationalTransfer(lines []string) (shared.Money, bool) {
	for _, raw := range lines {
		columns := strings.Split(raw, "\t")
		if len(columns) < 2 || strings.TrimSpace(columns[0]) != internationalTransferLabel {
			continue
		}
		return parseForeignAmount(columns[1])
	}
	return 0, false
}
