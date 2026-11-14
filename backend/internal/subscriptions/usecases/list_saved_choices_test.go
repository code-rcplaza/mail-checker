package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// failingCharges fails every read, to check the error comes back with context.
type failingCharges struct{}

func (failingCharges) ListSince(context.Context, int64, time.Time) ([]shared.Charge, error) {
	return nil, errDatabaseLocked
}

var errDatabaseLocked = errors.New("database is locked")

func savedSourcesOf(owner int64) SavedChoicesSources {
	confirmed := domain.Decision{Kind: domain.DecisionConfirmed, Cadence: domain.CadenceMonthly}
	return SavedChoicesSources{
		Charges:   chargesByUser{owner: streamingOnce()},
		Decisions: decisionsByUser{owner: {streamingLine: confirmed}},
		Aliases:   aliasesByUser{owner: {streamingLine: "Padres"}},
		Manuals:   manualsByUser{owner: {{PublicID: "m1", Name: "Gimnasio", Amount: 2500000, Currency: shared.CurrencyCLP, BillingDay: 5}}},
		Overrides: overridesByUser{owner: {streamingLine: {Amount: 900000}}},
	}
}

func TestListSavedChoices_ReadsEverythingOfTheUser(t *testing.T) {
	sources := savedSourcesOf(userOne)

	got, err := NewListSavedChoices(sources).Execute(context.Background(), userOne)

	want := SavedChoices{
		Charges:   streamingOnce(),
		Decisions: sources.Decisions.(decisionsByUser)[userOne],
		Aliases:   sources.Aliases.(aliasesByUser)[userOne],
		Manuals:   sources.Manuals.(manualsByUser)[userOne],
		Overrides: sources.Overrides.(overridesByUser)[userOne],
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("Execute = %+v, %v; want %+v", got, err, want)
	}
}

func TestListSavedChoices_AnotherUserGetsNothing(t *testing.T) {
	got, err := NewListSavedChoices(savedSourcesOf(userOne)).Execute(context.Background(), userTwo)

	if err != nil || len(got.Charges)+len(got.Decisions)+len(got.Aliases)+len(got.Manuals)+len(got.Overrides) != 0 {
		t.Errorf("Execute = %+v, %v; want nothing of another user", got, err)
	}
}

func TestListSavedChoices_WrapsReadErrors(t *testing.T) {
	sources := savedSourcesOf(userOne)
	sources.Charges = failingCharges{}

	_, err := NewListSavedChoices(sources).Execute(context.Background(), userOne)

	if !errors.Is(err, errDatabaseLocked) || err.Error() == errDatabaseLocked.Error() {
		t.Errorf("err = %v, want the read error wrapped with context", err)
	}
}
