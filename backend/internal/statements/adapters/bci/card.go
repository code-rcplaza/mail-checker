package bci

import (
	"regexp"
	"strings"
)

// cardNumberLine is "Nº DE TARJETA DE CRÉDITO\tXXXX XXXX XXXX 1234": only
// the last four digits are shown, and only those are kept.
var cardNumberLine = regexp.MustCompile(`TARJETA DE CRÉDITO\t[X ]+(\d{4})\s*$`)

// readCardLast4 takes the first card number of the file; both statements of
// a file are for the same card. Empty when it is not shown.
func readCardLast4(text string) string {
	for _, raw := range strings.Split(text, "\n") {
		if match := cardNumberLine.FindStringSubmatch(raw); match != nil {
			return match[1]
		}
	}
	return ""
}
