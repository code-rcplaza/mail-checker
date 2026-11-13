package usecases

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

func (m manualsByUser) Update(_ context.Context, userID int64, manual domain.ManualSubscription) error {
	index := slices.IndexFunc(m[userID], func(saved domain.ManualSubscription) bool { return saved.PublicID == manual.PublicID })
	if index < 0 {
		return shared.ErrNotFound
	}
	m[userID][index] = manual
	return nil
}

func (m manualsByUser) Delete(_ context.Context, userID int64, publicID string) error {
	index := slices.IndexFunc(m[userID], func(saved domain.ManualSubscription) bool { return saved.PublicID == publicID })
	if index < 0 {
		return shared.ErrNotFound
	}
	m[userID] = slices.Delete(slices.Clone(m[userID]), index, index+1)
	return nil
}

func savedGym() domain.ManualSubscription {
	return domain.ManualSubscription{PublicID: "manual-1", Name: "Gimnasio", Amount: 500000, Currency: shared.CurrencyCLP, BillingDay: 15}
}

func TestEditManualSubscription(t *testing.T) {
	edited := ManualSubscriptionInput{Name: "Gimnasio anual", Amount: 600000, Currency: "CLP", BillingDay: 20}
	invalid := ManualSubscriptionInput{Name: "", Amount: 600000, Currency: "CLP", BillingDay: 20}
	cases := []struct {
		name     string
		userID   int64
		publicID string
		input    ManualSubscriptionInput
		wantErr  error
		want     domain.ManualSubscription
	}{
		{"valid data are saved under the same public id", userOne, "manual-1", edited,
			nil, domain.ManualSubscription{PublicID: "manual-1", Name: "Gimnasio anual", Amount: 600000, Currency: shared.CurrencyCLP, BillingDay: 20}},
		{"invalid data change nothing", userOne, "manual-1", invalid, shared.ErrInvalid, savedGym()},
		{"an unknown public id is not found", userOne, "manual-404", edited, shared.ErrNotFound, savedGym()},
		{"the manual of another user is not found", userTwo, "manual-1", edited, shared.ErrNotFound, savedGym()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manuals := manualsByUser{userOne: {savedGym()}}

			err := NewEditManualSubscription(manuals).Execute(context.Background(), tc.userID, tc.publicID, tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if errors.Is(tc.wantErr, shared.ErrNotFound) {
				assertKeyedError(t, err, shared.ErrNotFound, domain.KeyManualNotFound)
			}
			if !reflect.DeepEqual(manuals[userOne], []domain.ManualSubscription{tc.want}) {
				t.Errorf("saved = %+v, want [%+v]", manuals[userOne], tc.want)
			}
		})
	}
}

func TestRemoveManualSubscription(t *testing.T) {
	cases := []struct {
		name     string
		userID   int64
		publicID string
		wantErr  error
		wantLeft int
	}{
		{"an existing manual leaves the summary", userOne, "manual-1", nil, 0},
		{"an unknown public id is not found", userOne, "manual-404", shared.ErrNotFound, 1},
		{"the manual of another user is not found", userTwo, "manual-1", shared.ErrNotFound, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manuals := manualsByUser{userOne: {savedGym()}}

			err := NewRemoveManualSubscription(manuals).Execute(context.Background(), tc.userID, tc.publicID)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if errors.Is(tc.wantErr, shared.ErrNotFound) {
				assertKeyedError(t, err, shared.ErrNotFound, domain.KeyManualNotFound)
			}
			summary, _ := NewGetSummary(SummarySources{Charges: chargesByUser{}, Catalog: fakeCatalog{}, Decisions: decisionsByUser{}, Aliases: aliasesByUser{}, Rates: ratesByUser{}, KeptRises: keptRisesByUser{}, Manuals: manuals, Overrides: overridesByUser{}}, fixedNow).Execute(context.Background(), userOne, 6)
			if len(summary.Subscriptions) != tc.wantLeft {
				t.Errorf("subscriptions in the summary = %d, want %d", len(summary.Subscriptions), tc.wantLeft)
			}
		})
	}
}
