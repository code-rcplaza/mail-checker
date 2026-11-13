package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// ratesByUser keeps each user's rates apart, like the real repository.
type ratesByUser struct {
	rates map[int64][]shared.ExchangeRate
	err   error
}

func (r ratesByUser) List(_ context.Context, userID int64) ([]shared.ExchangeRate, error) {
	return r.rates[userID], r.err
}

// dollarsAndPesos is one active subscription in each currency.
func dollarsAndPesos() chargesByUser {
	return chargesByUser{userOne: {
		chargeInCurrency("d1", "STREAMING CO", 2000, shared.CurrencyUSD, 40),
		chargeInCurrency("d2", "STREAMING CO", 2000, shared.CurrencyUSD, 10),
		chargeInCurrency("p1", "STRMNG*CO", 1000000, shared.CurrencyCLP, 35),
		chargeInCurrency("p2", "STRMNG*CO", 1000000, shared.CurrencyCLP, 5),
	}}
}

func chargeInCurrency(externalID, merchant string, amount shared.Money, currency shared.Currency, daysAgo int) shared.Charge {
	c := chargeDaysAgo(externalID, merchant, shared.SourceStatement, amount, daysAgo)
	c.Currency = currency
	return c
}

func summaryWithRates(t *testing.T, rates ratesByUser) (domain.Summary, error) {
	t.Helper()
	useCase := NewGetSummary(SummarySources{Charges: dollarsAndPesos(), Catalog: fakeCatalog{}, Decisions: decisionsByUser{}, Aliases: aliasesByUser{}, Rates: rates, KeptRises: keptRisesByUser{}, Manuals: manualsByUser{}, Overrides: overridesByUser{}}, fixedNow)
	return useCase.Execute(context.Background(), userOne, 6)
}

func TestGetSummary_TotalUsesTheLatestExchangeRateOfTheUser(t *testing.T) {
	rates := ratesByUser{rates: map[int64][]shared.ExchangeRate{
		userOne: {
			{PesosPerDollar: 90000, Date: time.Date(2026, 9, 17, 0, 0, 0, 0, time.UTC)},
			{PesosPerDollar: 95000, Date: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)},
		},
		userTwo: {{PesosPerDollar: 100000, Date: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}},
	}}

	got, err := summaryWithRates(t, rates)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	want := domain.Total{Amount: 2900000, IsApproximate: true}
	if got.Total != want {
		t.Errorf("total = %+v, want %+v", got.Total, want)
	}
}

func TestGetSummary_WithoutExchangeRateTheTotalIsIncomplete(t *testing.T) {
	got, err := summaryWithRates(t, ratesByUser{})

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	want := domain.Total{Amount: 1000000, IsIncomplete: true}
	if got.Total != want {
		t.Errorf("total = %+v, want %+v", got.Total, want)
	}
}

func TestGetSummary_ExchangeRateErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")

	_, err := summaryWithRates(t, ratesByUser{err: cause})

	if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), "read exchange rates: ") {
		t.Errorf("error = %v, want \"read exchange rates\" wrapping %v", err, cause)
	}
}
