package domain

type Currency string

const (
	CurrencyCLP Currency = "CLP"
	CurrencyUSD Currency = "USD"
)

func ParseCurrency(raw string) (Currency, error) {
	currency := Currency(raw)
	if currency != CurrencyCLP && currency != CurrencyUSD {
		return "", Invalid(KeyCurrencyUnknown, nil)
	}
	return currency, nil
}
