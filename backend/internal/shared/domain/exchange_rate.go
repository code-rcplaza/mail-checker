package domain

import "time"

// ExchangeRate is what one US dollar cost in pesos on a date.
type ExchangeRate struct {
	PesosPerDollar Money
	Date           time.Time
}

// hundredthsPerUnit is how many hundredths a whole peso or dollar has.
const hundredthsPerUnit = 100

// DollarsPaidInPesos is an amount of dollars and the pesos paid for it.
type DollarsPaidInPesos struct {
	Dollars Money
	Pesos   Money
}

// DeriveExchangeRate gives the rate the bank really charged on date, from the
// pesos paid for an amount of dollars. Zero dollars give no rate.
func DeriveExchangeRate(paid DollarsPaidInPesos, date time.Time) (ExchangeRate, bool, error) {
	if paid.Pesos < 0 || paid.Dollars < 0 {
		return ExchangeRate{}, false, Invalid(KeyAmountNegative, nil)
	}
	if paid.Dollars == 0 {
		return ExchangeRate{}, false, nil
	}
	return ExchangeRate{PesosPerDollar: pesosPerDollar(paid), Date: date}, true, nil
}

// pesosPerDollar rounds half up to hundredths; both amounts are non-negative.
func pesosPerDollar(paid DollarsPaidInPesos) Money {
	return (paid.Pesos*hundredthsPerUnit + paid.Dollars/2) / paid.Dollars
}
