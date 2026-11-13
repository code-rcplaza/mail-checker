package usecases

import (
	"context"
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.ManualSubscriptions = manualsByUser{}

// manualsByUser keeps each user's manual subscriptions apart.
type manualsByUser map[int64][]domain.ManualSubscription

func (m manualsByUser) List(_ context.Context, userID int64) ([]domain.ManualSubscription, error) {
	return m[userID], nil
}

func (m manualsByUser) Create(_ context.Context, userID int64, manual domain.ManualSubscription) error {
	m[userID] = append(m[userID], manual)
	return nil
}

func createManual(manuals manualsByUser, userID int64, input ManualSubscriptionInput) (domain.ManualSubscription, error) {
	useCase := NewCreateManualSubscription(manuals, func() string { return "manual-1" })
	return useCase.Execute(context.Background(), userID, input)
}

var gymInput = ManualSubscriptionInput{Name: "Gimnasio", Amount: 500000, Currency: "CLP", BillingDay: 15}

func TestCreateManualSubscription_SavesItWithAGeneratedPublicID(t *testing.T) {
	manuals := manualsByUser{}

	got, err := createManual(manuals, userOne, gymInput)

	want := domain.ManualSubscription{PublicID: "manual-1", Name: "Gimnasio", Amount: 500000, Currency: shared.CurrencyCLP, BillingDay: 15}
	if err != nil || got != want {
		t.Fatalf("Execute = (%+v, %v), want (%+v, nil)", got, err, want)
	}
	if saved := manuals[userOne]; len(saved) != 1 || saved[0] != want {
		t.Errorf("saved = %+v, want [%+v]", saved, want)
	}
}

func TestCreateManualSubscription_InvalidDataSavesNothing(t *testing.T) {
	manuals := manualsByUser{}
	invalid := gymInput
	invalid.BillingDay = 32

	_, err := createManual(manuals, userOne, invalid)

	if !errors.Is(err, shared.ErrInvalid) || len(manuals) != 0 {
		t.Errorf("Execute error = %v with saved %v, want ErrInvalid and nothing saved", err, manuals)
	}
}

func TestGetSummary_ShowsOnlyTheManualSubscriptionsOfTheUser(t *testing.T) {
	manuals := manualsByUser{}
	if _, err := createManual(manuals, userTwo, gymInput); err != nil {
		t.Fatalf("create error = %v", err)
	}
	useCase := NewGetSummary(SummarySources{Charges: chargesByUser{}, Catalog: fakeCatalog{}, Decisions: decisionsByUser{}, Aliases: aliasesByUser{}, Rates: ratesByUser{}, KeptRises: keptRisesByUser{}, Manuals: manuals, Overrides: overridesByUser{}}, fixedNow)

	mine, errMine := useCase.Execute(context.Background(), userOne, 6)
	theirs, errTheirs := useCase.Execute(context.Background(), userTwo, 6)

	if errMine != nil || errTheirs != nil {
		t.Fatalf("Execute errors = %v, %v", errMine, errTheirs)
	}
	if len(mine.Subscriptions) != 0 || len(theirs.Subscriptions) != 1 || theirs.Total.Amount != 500000 {
		t.Errorf("mine = %+v, theirs = %+v; want the manual only for its owner, adding to the total", mine.Subscriptions, theirs)
	}
}
