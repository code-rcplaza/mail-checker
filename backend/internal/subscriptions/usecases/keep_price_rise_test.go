package usecases

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.PriceAlertDecisions = keptRisesByUser{}

// keptRisesByUser keeps, per user and line, every amount whose rise was kept, like the UNIQUE of the table.
type keptRisesByUser map[int64]map[domain.GroupKey][]shared.Money

func (k keptRisesByUser) List(_ context.Context, userID int64) (map[domain.GroupKey][]shared.Money, error) {
	return k[userID], nil
}

func (k keptRisesByUser) Save(_ context.Context, userID int64, kept domain.KeptRise) error {
	if k[userID] == nil {
		k[userID] = make(map[domain.GroupKey][]shared.Money)
	}
	if !slices.Contains(k[userID][kept.Key], kept.NewAmount) {
		k[userID][kept.Key] = append(k[userID][kept.Key], kept.NewAmount)
	}
	return nil
}

func TestKeptRisesByUser_SavingTwiceKeepsOne(t *testing.T) {
	kept := keptRisesByUser{}

	for range 2 {
		if err := kept.Save(context.Background(), userOne, domain.KeptRise{Key: streamingLine, NewAmount: 1050000}); err != nil {
			t.Fatalf("Save error = %v, want nil", err)
		}
	}

	got, _ := kept.List(context.Background(), userOne)
	if !reflect.DeepEqual(got, map[domain.GroupKey][]shared.Money{streamingLine: {1050000}}) {
		t.Errorf("kept = %v, want one entry for the line", got)
	}
}

// risingStreaming is a monthly line of a known service whose last charge went up.
func risingStreaming() []shared.Charge {
	return []shared.Charge{
		chargeDaysAgo("s1", "STREAMING CO", shared.SourceStatement, 1000000, 70),
		chargeDaysAgo("s2", "STREAMING CO", shared.SourceStatement, 1000000, 40),
		chargeDaysAgo("s3", "STREAMING CO", shared.SourceStatement, 1050000, 10),
	}
}

const risingLine domain.GroupKey = "service:Streaming Co||d03"

func TestKeepPriceRise(t *testing.T) {
	cases := []struct {
		name    string
		charges chargesByUser
		wantErr error
		want    keptRisesByUser
	}{
		{"a line with a rise is kept at its new amount", chargesByUser{userOne: risingStreaming()}, nil, keptRisesByUser{userOne: {risingLine: {1050000}}}},
		{"a line without a rise is not found", chargesByUser{userOne: risingStreaming()[:2]}, shared.ErrNotFound, keptRisesByUser{}},
		{"the line of another user is not found", chargesByUser{userTwo: risingStreaming()}, shared.ErrNotFound, keptRisesByUser{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kept := keptRisesByUser{}

			err := NewKeepPriceRise(linesOf(tc.charges), kept).Execute(context.Background(), userOne, risingLine)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				assertKeyedError(t, err, shared.ErrNotFound, domain.KeyRiseNotFound)
			}
			if !reflect.DeepEqual(kept, tc.want) {
				t.Errorf("kept = %v, want %v", kept, tc.want)
			}
		})
	}
}

func TestGetSummary_HidesTheKeptRiseUntilANewOne(t *testing.T) {
	newRise := []shared.Charge{
		chargeDaysAgo("s1", "STREAMING CO", shared.SourceStatement, 1000000, 70),
		chargeDaysAgo("s2", "STREAMING CO", shared.SourceStatement, 1050000, 40),
		chargeDaysAgo("s3", "STREAMING CO", shared.SourceStatement, 1100000, 10),
	}
	cases := []struct {
		name      string
		charges   []shared.Charge
		wantRises int
	}{
		{"the kept rise does not show", risingStreaming(), 0},
		{"a new higher charge shows a new rise", newRise, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kept := keptRisesByUser{userOne: {risingLine: {1050000}}}
			sources := summarySourcesOf(chargesByUser{userOne: tc.charges})
			sources.KeptRises = kept
			useCase := NewGetSummary(sources)

			got, err := useCase.Execute(context.Background(), userOne, 6)

			if err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}
			if len(got.Rises) != tc.wantRises {
				t.Errorf("rises = %+v, want %d", got.Rises, tc.wantRises)
			}
		})
	}
}
