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
	rate := latestRateOrNil(rates)
	summary.ExchangeRate = rate
	summary.Total = TotalInPesos(summary.MonthlySpend, rate)
	return changeEveryLine(summary, func(line Subscription) Subscription {
		line.InPesos = pesosAmountOf(line, rate)
		return line
	})
}

// latestRateOrNil is nil when the user's statements showed no rate.
func latestRateOrNil(rates []shared.ExchangeRate) *shared.ExchangeRate {
	rate, hasRate := LatestExchangeRate(rates)
	if !hasRate {
		return nil
	}
	return &rate
}

// pesosAmountOf is nil for dollars without a rate; charges are only ever
// stored in pesos or dollars.
func pesosAmountOf(line Subscription, rate *shared.ExchangeRate) *PesosAmount {
	switch {
	case line.Currency == shared.CurrencyCLP:
		return &PesosAmount{Amount: line.Amount}
	case line.Currency == shared.CurrencyUSD && rate != nil:
		return &PesosAmount{Amount: dollarsToPesos(line.Amount, *rate), IsApproximate: true}
	}
	return nil
}
