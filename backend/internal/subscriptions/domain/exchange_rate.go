package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// hundredthsPerDollar is how many hundredths a dollar has.
const hundredthsPerDollar = 100

// LatestExchangeRate picks the rate of the most recent date; false means
// there is none and a total in pesos cannot include dollars. Two rates of the
// same date tie towards the higher one, so the result does not depend on the
// order and the total is never understated.
func LatestExchangeRate(rates []shared.ExchangeRate) (shared.ExchangeRate, bool) {
	if len(rates) == 0 {
		return shared.ExchangeRate{}, false
	}
	latest := rates[0]
	for _, rate := range rates[1:] {
		if isMoreRecentRate(rate, latest) {
			latest = rate
		}
	}
	return latest, true
}

func isMoreRecentRate(rate, current shared.ExchangeRate) bool {
	if rate.Date.Equal(current.Date) {
		return rate.PesosPerDollar > current.PesosPerDollar
	}
	return rate.Date.After(current.Date)
}

// ConvertToPesos gives an amount in pesos and tells whether it is an
// approximation: dollars are converted with the rate, pesos stay as they are.
func ConvertToPesos(amount shared.Money, currency shared.Currency, rate shared.ExchangeRate) (pesos shared.Money, isApproximate bool, err error) {
	switch currency {
	case shared.CurrencyCLP:
		return amount, false, nil
	case shared.CurrencyUSD:
		return dollarsToPesos(amount, rate), true, nil
	}
	return 0, false, shared.Invalid(shared.KeyCurrencyUnknown, nil)
}

// dollarsToPesos rounds half up to hundredths; amounts are non-negative.
func dollarsToPesos(dollars shared.Money, rate shared.ExchangeRate) shared.Money {
	return (dollars*rate.PesosPerDollar + hundredthsPerDollar/2) / hundredthsPerDollar
}
