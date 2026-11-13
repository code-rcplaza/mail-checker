package bci

import "strings"

// section is the part of the statement a line belongs to; only some of them
// hold purchases.
type section int

const (
	sectionNone section = iota
	sectionNationalHeader
	sectionNationalPurchases // "1. TOTAL OPERACIONES": purchases, payments and installments
	sectionNationalOther     // products, fees, interest and debt transfers
	sectionInternational
)

// lineOutcome tells what reading a line gave.
type lineOutcome int

const (
	lineSkipped lineOutcome = iota
	linePurchase
	lineUnreadable
)

// Section headings, without spaces. A page break repeats "2. PERÍODO ACTUAL",
// which does not end the purchases.
const (
	purchasesHeading      = "1.TOTALOPERACIONES"
	productsHeading       = "PRODUCTOSOSERVICIOSVOLUNTARIAMENTE"
	feesHeading           = "CARGOS,COMISIONES"
	installmentsInfoTitle = "INFORMACIONCOMPRASENCUOTAS"
)

func nextSection(current section, raw string) section {
	compact := withoutSpaces(raw)
	switch {
	case strings.Contains(compact, internationalTitle):
		return sectionInternational
	case strings.Contains(compact, nationalTitle):
		return sectionNationalHeader
	case strings.HasPrefix(compact, purchasesHeading):
		return sectionNationalPurchases
	case current == sectionNationalPurchases && isNationalOtherHeading(compact):
		return sectionNationalOther
	}
	return current
}

func isNationalOtherHeading(compact string) bool {
	return strings.Contains(compact, productsHeading) ||
		strings.Contains(compact, feesHeading) ||
		strings.Contains(compact, installmentsInfoTitle)
}
