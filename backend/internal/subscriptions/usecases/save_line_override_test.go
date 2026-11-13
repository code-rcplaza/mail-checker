package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.LineOverrides = overridesByUser{}

type overridesByUser map[int64]map[domain.GroupKey]domain.LineOverride

func (o overridesByUser) List(_ context.Context, userID int64) (map[domain.GroupKey]domain.LineOverride, error) {
	return o[userID], nil
}

func (o overridesByUser) Save(_ context.Context, userID int64, key domain.GroupKey, override domain.LineOverride) error {
	if o[userID] == nil {
		o[userID] = make(map[domain.GroupKey]domain.LineOverride)
	}
	o[userID][key] = override
	return nil
}

func (o overridesByUser) Delete(_ context.Context, userID int64, key domain.GroupKey) error {
	delete(o[userID], key)
	return nil
}

func amountOf(m shared.Money) *shared.Money { return &m }

func TestSaveLineOverride(t *testing.T) {
	existing := domain.LineOverride{Amount: 700000}
	cases := []struct {
		name    string
		userID  int64
		amount  *shared.Money
		wantErr error
		want    map[domain.GroupKey]domain.LineOverride
	}{
		{"a valid override is saved by line key", userOne, amountOf(900000), nil, map[domain.GroupKey]domain.LineOverride{streamingLine: {Amount: 900000}}},
		{"an empty override removes the saved one", userOne, nil, nil, map[domain.GroupKey]domain.LineOverride{}},
		{"an invalid override changes nothing", userOne, amountOf(0), shared.ErrInvalid, map[domain.GroupKey]domain.LineOverride{streamingLine: existing}},
		{"the line of another user is not found", userTwo, amountOf(900000), shared.ErrNotFound, map[domain.GroupKey]domain.LineOverride{streamingLine: existing}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			overrides := overridesByUser{userOne: {streamingLine: existing}}
			charges := chargesByUser{userOne: streamingOnce()}

			err := NewSaveLineOverride(charges, fakeCatalog{}, overrides, fixedNow).Execute(context.Background(), tc.userID, streamingLine, tc.amount, nil)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if errors.Is(tc.wantErr, shared.ErrNotFound) {
				assertKeyedError(t, err, shared.ErrNotFound, domain.KeyLineNotFound)
			}
			if !reflect.DeepEqual(overrides[userOne], tc.want) {
				t.Errorf("overrides = %v, want %v", overrides[userOne], tc.want)
			}
		})
	}
}

func TestGetSummary_ShowsTheOverriddenAmount(t *testing.T) {
	overrides := overridesByUser{userOne: {"service:Streaming Co||d02": {Amount: 900000}}}
	sources := SummarySources{
		Charges: chargesByUser{userOne: {
			chargeDaysAgo("s1", "STREAMING CO", shared.SourceStatement, 599000, 40),
			chargeDaysAgo("s2", "STREAMING CO", shared.SourceStatement, 599000, 10),
		}},
		Catalog: fakeCatalog{}, Decisions: decisionsByUser{}, Aliases: aliasesByUser{}, Rates: ratesByUser{},
		KeptRises: keptRisesByUser{}, Manuals: manualsByUser{}, Overrides: overrides,
	}

	got, err := NewGetSummary(sources, fixedNow).Execute(context.Background(), userOne, 6)

	if err != nil || len(got.Subscriptions) != 1 || got.Subscriptions[0].Amount != 900000 || got.Total.Amount != 900000 {
		t.Errorf("summary = %+v, %v; want the overridden amount in the line and the total", got, err)
	}
}
