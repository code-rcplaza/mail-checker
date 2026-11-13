package domain

import shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"

// PesosAmount is a line's amount in pesos; approximate when it was converted
// from dollars, so the front shows it with "≈".
type PesosAmount struct {
	Amount        shared.Money
	IsApproximate bool
}

// ApplyExchangeRates sets the total in pesos and each line's equivalent with
// the most recent rate. Without one, the total leaves the dollars out and
// says so, and a line in dollars shows only its original amount.
func ApplyExchangeRates(summary Summary, rates []shared.ExchangeRate) Summary {
	rate, hasRate := LatestExchangeRate(rates)
	summary.Total = TotalInPesos(summary.MonthlySpend, rate, hasRate)
	summary.Subscriptions = withPesosAmounts(summary.Subscriptions, rate, hasRate)
	summary.Possible = withPesosAmounts(summary.Possible, rate, hasRate)
	summary.Inactive = withPesosAmounts(summary.Inactive, rate, hasRate)
	return summary
}

// withPesosAmounts returns new lines; the caller's slice is left untouched.
func withPesosAmounts(lines []Subscription, rate shared.ExchangeRate, hasRate bool) []Subscription {
	if lines == nil {
		return nil
	}
	converted := make([]Subscription, 0, len(lines))
	for _, line := range lines {
		line.InPesos = pesosAmountOf(line, rate, hasRate)
		converted = append(converted, line)
	}
	return converted
}

// pesosAmountOf is nil for dollars without a rate; charges are only ever
// stored in pesos or dollars.
func pesosAmountOf(line Subscription, rate shared.ExchangeRate, hasRate bool) *PesosAmount {
	switch {
	case line.Currency == shared.CurrencyCLP:
		return &PesosAmount{Amount: line.Amount}
	case line.Currency == shared.CurrencyUSD && hasRate:
		return &PesosAmount{Amount: dollarsToPesos(line.Amount, rate), IsApproximate: true}
	}
	return nil
}
