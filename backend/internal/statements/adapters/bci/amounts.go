package bci

import (
	"regexp"
	"strconv"
	"strings"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// clpAmount is a peso amount as the national statement prints it: "$5.990",
// "$-480.512", "$ 4.200.000"; dots group thousands and there are no cents.
var clpAmount = regexp.MustCompile(`^\$ ?(-?)(\d{1,3}(?:\.\d{3})*)$`)

const hundredthsPerPeso = 100

// foreignAmount is an amount of the international statement: "20,00",
// "5.990,00", "US$ 100,00", "US$-55,40"; it already comes in hundredths.
var foreignAmount = regexp.MustCompile(`^(?:US\$ ?)?(-?)(\d{1,3}(?:\.\d{3})*),(\d{2})$`)

func parseCLPAmount(raw string) (shared.Money, bool) {
	match := clpAmount.FindStringSubmatch(strings.TrimSpace(raw))
	if match == nil {
		return 0, false
	}
	pesos, err := strconv.ParseInt(strings.ReplaceAll(match[2], ".", ""), 10, 64)
	if err != nil {
		return 0, false
	}
	if match[1] == "-" {
		pesos = -pesos
	}
	return shared.Money(pesos * hundredthsPerPeso), true
}

func parseForeignAmount(raw string) (shared.Money, bool) {
	match := foreignAmount.FindStringSubmatch(strings.TrimSpace(raw))
	if match == nil {
		return 0, false
	}
	hundredths, err := strconv.ParseInt(strings.ReplaceAll(match[2], ".", "")+match[3], 10, 64)
	if err != nil {
		return 0, false
	}
	if match[1] == "-" {
		hundredths = -hundredths
	}
	return shared.Money(hundredths), true
}
