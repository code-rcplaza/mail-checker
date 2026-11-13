package bci

import (
	"strings"
	"time"
)

// statementDateLabel heads the billing date: "FECHA ESTADO DE CUENTA\t21/10/2026".
// The closing day changes from month to month, so it is always read.
const (
	statementDateLabel = "FECHA ESTADO DE CUENTA"
	billingDateLayout  = "02/01/2006"
)

// BillingDate is the first statement date of the file, at midnight in the
// parser's location; false when it is missing or unreadable.
func (p StatementParser) BillingDate(text string) (time.Time, bool) {
	for _, raw := range strings.Split(text, "\n") {
		label, value, hasValue := strings.Cut(raw, "\t")
		if !hasValue || strings.TrimSpace(label) != statementDateLabel {
			continue
		}
		date, err := time.ParseInLocation(billingDateLayout, strings.TrimSpace(value), p.location)
		return date, err == nil
	}
	return time.Time{}, false
}
