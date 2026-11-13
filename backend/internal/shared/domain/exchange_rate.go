package domain

import "time"

// ExchangeRate is what one US dollar cost in pesos on a date.
type ExchangeRate struct {
	PesosPerDollar Money
	Date           time.Time
}

// hundredthsPerUnit is how many hundredths a whole peso or dollar has.
const hundredthsPerUnit = 100

// DeriveExchangeRate gives the rate the bank really charged, from the pesos
// paid for an amount of dollars. Zero dollars give no rate.
func DeriveExchangeRate(pesos, dollars Money, date time.Time) (ExchangeRate, bool, error) {
	if pesos < 0 || dollars < 0 {
		return ExchangeRate{}, false, Invalid(KeyAmountNegative, nil)
	}
	if dollars == 0 {
		return ExchangeRate{}, false, nil
	}
	return ExchangeRate{PesosPerDollar: pesosPerDollar(pesos, dollars), Date: date}, true, nil
}

// pesosPerDollar rounds half up to hundredths; both amounts are non-negative.
func pesosPerDollar(pesos, dollars Money) Money {
	return (pesos*hundredthsPerUnit + dollars/2) / dollars
}
