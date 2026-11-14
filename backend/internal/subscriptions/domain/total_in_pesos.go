package domain

import shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"

// Total is the monthly spend in pesos.
type Total struct {
	Amount        shared.Money
	IsApproximate bool // includes dollars converted with an exchange rate
	IsIncomplete  bool // dollars were left out because there is no exchange rate
}

// TotalInPesos adds the monthly spend of every currency in pesos. Without an
// exchange rate (nil) the dollars are left out and the total says so, instead of
// showing a number the user would trust.
func TotalInPesos(spend map[shared.Currency]shared.Money, rate *shared.ExchangeRate) Total {
	total := Total{Amount: spend[shared.CurrencyCLP]}
	dollars := spend[shared.CurrencyUSD]
	if dollars == 0 {
		return total
	}
	if rate == nil {
		total.IsIncomplete = true
		return total
	}
	total.Amount += dollarsToPesos(dollars, *rate)
	total.IsApproximate = true
	return total
}
