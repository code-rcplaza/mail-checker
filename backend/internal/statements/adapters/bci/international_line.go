package bci

import (
	"regexp"
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

// internationalPurchase is what an international line shows: the origin
// amount comes without its currency, next to the amount in dollars.
type internationalPurchase struct {
	merchant string
	date     time.Time
	origin   shared.Money
	dollars  shared.Money
}

// internationalOperation is the first column of a purchase: reference and
// date, "3009 00000000000000000000002 29/09/26". Payments and debt transfers
// start with a short code and the date in its own column, so they never match.
var internationalOperation = regexp.MustCompile(`^\d{4} ?\d{20,} (\d{2}/\d{2}/\d{2})$`)

// The columns after the operation: descriptor, city, country, origin amount
// and dollar amount. City or country may be missing, so amounts are read
// from the end.
const minInternationalColumns = 4

func (p StatementParser) readInternationalLine(raw, cardLast4 string) (shared.ParsedCharge, lineOutcome) {
	purchase, outcome := p.readInternationalColumns(raw)
	if outcome != linePurchase {
		return shared.ParsedCharge{}, outcome
	}
	amount, currency := chargedAmount(purchase)
	return shared.ParsedCharge{
		Amount:    amount,
		Currency:  currency,
		Merchant:  purchase.merchant,
		Card:      CardCredit,
		CardLast4: cardLast4,
		Date:      purchase.date,
	}, linePurchase
}

func (p StatementParser) readInternationalColumns(raw string) (internationalPurchase, lineOutcome) {
	columns := strings.Split(raw, "\t")
	match := internationalOperation.FindStringSubmatch(columns[0])
	if match == nil {
		return internationalPurchase{}, lineSkipped
	}
	if len(columns) < minInternationalColumns {
		return internationalPurchase{}, lineUnreadable
	}
	date, isDateRead := p.readDate(match[1])
	origin, isOriginRead := parseForeignAmount(columns[len(columns)-2])
	dollars, isDollarsRead := parseForeignAmount(columns[len(columns)-1])
	if !isDateRead || !isOriginRead || !isDollarsRead {
		return internationalPurchase{}, lineUnreadable
	}
	if dollars < 0 {
		return internationalPurchase{}, lineSkipped // a refund
	}
	return internationalPurchase{merchant: strings.TrimSpace(columns[1]), date: date, origin: origin, dollars: dollars}, linePurchase
}

// chargedAmount keeps a charge in pesos when its origin was pesos, so it is
// exact; anything else is kept in dollars, the only other currency stored.
func chargedAmount(purchase internationalPurchase) (shared.Money, shared.Currency) {
	if domain.IsChargedInPesos(purchase.origin, purchase.dollars) {
		return purchase.origin, shared.CurrencyCLP
	}
	return purchase.dollars, shared.CurrencyUSD
}
