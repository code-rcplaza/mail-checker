package usecases

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

var _ ports.ExchangeRateRecorder = &fakeExchangeRates{}

// fakeExchangeRates keeps one rate per user and date, like the real table.
type fakeExchangeRates struct {
	byUser    map[int64][]shared.ExchangeRate
	recordErr error
}

func (f *fakeExchangeRates) Record(_ context.Context, userID int64, rate shared.ExchangeRate) error {
	if f.recordErr != nil {
		return f.recordErr
	}
	if f.byUser == nil {
		f.byUser = make(map[int64][]shared.ExchangeRate)
	}
	for _, saved := range f.byUser[userID] {
		if saved.Date.Equal(rate.Date) {
			return nil
		}
	}
	f.byUser[userID] = append(f.byUser[userID], rate)
	return nil
}

func TestFakeExchangeRates_KeepsOneRatePerUserAndDate(t *testing.T) {
	rates := &fakeExchangeRates{}
	rate := shared.ExchangeRate{PesosPerDollar: 95000, Date: billingDate}

	for _, userID := range []int64{userOne, userOne, userTwo} {
		if err := rates.Record(context.Background(), userID, rate); err != nil {
			t.Fatalf("Record error = %v, want nil", err)
		}
	}

	want := map[int64][]shared.ExchangeRate{userOne: {rate}, userTwo: {rate}}
	if !reflect.DeepEqual(rates.byUser, want) {
		t.Errorf("rates = %v, want one per user", rates.byUser)
	}
}

func TestImportStatement_RecordsTheExchangeRateOfTheTransferPair(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.transferPair = &domain.TransferPair{Pesos: 9500000, Dollars: 10000, Date: chargeDate}

	importAs(t, fixture, userOne)

	want := []shared.ExchangeRate{{PesosPerDollar: 95000, Date: billingDate}}
	if got := fixture.rates.byUser[userOne]; !reflect.DeepEqual(got, want) {
		t.Errorf("recorded rates = %v, want %v (dated on the billing date)", got, want)
	}
}

func TestImportStatement_WithoutTransferPairRecordsNoRate(t *testing.T) {
	fixture := newImportFixture()

	importAs(t, fixture, userOne)

	if len(fixture.rates.byUser) != 0 {
		t.Errorf("recorded rates = %v, want none", fixture.rates.byUser)
	}
}

func TestImportStatement_RateFailureSavesNoStatement(t *testing.T) {
	cause := errors.New("boom")
	fixture := newImportFixture()
	fixture.bci.transferPair = &domain.TransferPair{Pesos: 9500000, Dollars: 10000, Date: chargeDate}
	fixture.rates.recordErr = cause

	_, err := fixture.useCase.Execute(context.Background(), userOne, bciInput())

	if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), "record exchange rate: ") {
		t.Errorf("error = %v, want \"record exchange rate\" wrapping %v", err, cause)
	}
	if len(fixture.statements.created) != 0 {
		t.Errorf("Create called %d times, want 0", len(fixture.statements.created))
	}
}
