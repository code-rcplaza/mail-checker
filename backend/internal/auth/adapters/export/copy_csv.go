package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	subscriptionsdomain "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// The CSV is one table to open in a spreadsheet: amounts with two decimals and
// dates as YYYY-MM-DD, the way a spreadsheet reads them.
var csvHeader = []string{"name", "alias", "amount", "currency", "cadence", "nextCharge"}

// formulaStarts are the first characters a spreadsheet runs as a formula.
const formulaStarts = "=+-@\t\r"

// EncodeCSV writes one row per subscription of the summary, under a header.
func EncodeCSV(subscriptions []subscriptionsdomain.Subscription) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	rows := append([][]string{csvHeader}, subscriptionRows(subscriptions)...)
	if err := writer.WriteAll(rows); err != nil {
		return nil, fmt.Errorf("encode csv copy: %w", err)
	}
	return buffer.Bytes(), nil
}

func subscriptionRows(subscriptions []subscriptionsdomain.Subscription) [][]string {
	rows := make([][]string, 0, len(subscriptions))
	for _, line := range subscriptions {
		rows = append(rows, []string{
			neutralizeFormula(subscriptionsdomain.ParentName(line.Parent)),
			neutralizeFormula(line.Alias),
			decimalOf(line.Amount),
			string(line.Currency),
			string(line.Cadence),
			dateOf(line.NextCharge),
		})
	}
	return rows
}

// neutralizeFormula prefixes a quote so the spreadsheet shows the text instead of running it.
func neutralizeFormula(field string) string {
	if field != "" && strings.ContainsRune(formulaStarts, rune(field[0])) {
		return "'" + field
	}
	return field
}

func decimalOf(amount shared.Money) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}

func dateOf(at time.Time) string {
	if at.IsZero() {
		return ""
	}
	return at.UTC().Format(time.DateOnly)
}
