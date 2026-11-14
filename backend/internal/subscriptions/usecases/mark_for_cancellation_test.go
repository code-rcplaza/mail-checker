package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.Cancellations = cancellationsByUser{}

// cancellationsByUser keeps the first date of each mark, like the UNIQUE of the table.
type cancellationsByUser map[int64]map[domain.GroupKey]time.Time

func (c cancellationsByUser) List(_ context.Context, userID int64) (map[domain.GroupKey]time.Time, error) {
	return c[userID], nil
}

func (c cancellationsByUser) Mark(_ context.Context, userID int64, mark domain.CancellationMark) error {
	if c[userID] == nil {
		c[userID] = make(map[domain.GroupKey]time.Time)
	}
	if _, isMarked := c[userID][mark.Key]; isMarked {
		return nil
	}
	c[userID][mark.Key] = mark.MarkedAt
	return nil
}

func TestMarkForCancellation(t *testing.T) {
	earlier := testNow.AddDate(0, 0, -3)
	cases := []struct {
		name          string
		userID        int64
		key           domain.GroupKey
		cancellations cancellationsByUser
		wantErr       error
		want          map[domain.GroupKey]time.Time
	}{
		{"an existing line is marked with the date of now", userOne, streamingLine, cancellationsByUser{}, nil, map[domain.GroupKey]time.Time{streamingLine: testNow}},
		{"a marked line marked again keeps its date", userOne, streamingLine, cancellationsByUser{userOne: {streamingLine: earlier}}, nil, map[domain.GroupKey]time.Time{streamingLine: earlier}},
		{"a key that is not a line is not found", userOne, "merchant:NEVER SEEN:CLP||d01", cancellationsByUser{}, shared.ErrNotFound, nil},
		{"the line of another user is not found", userTwo, streamingLine, cancellationsByUser{}, shared.ErrNotFound, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			charges := chargesByUser{userOne: streamingOnce()}

			err := NewMarkForCancellation(linesOf(charges), tc.cancellations).Execute(context.Background(), tc.userID, tc.key)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				assertKeyedError(t, err, shared.ErrNotFound, domain.KeyLineNotFound)
			}
			if got := tc.cancellations[tc.userID]; !reflect.DeepEqual(got, tc.want) {
				t.Errorf("marks of user %d = %v, want %v", tc.userID, got, tc.want)
			}
		})
	}
}

func TestGetSummary_ShowsTheMarkOnlyOnTheMarkedLine(t *testing.T) {
	markedAt := testNow.AddDate(0, 0, -2)
	sources := summarySourcesOf(chargesWithPossibleAndActive())
	sources.Cancellations = cancellationsByUser{userOne: {activeLineKey: markedAt}}

	got, err := NewGetSummary(sources).Execute(context.Background(), userOne, 6)

	if err != nil || len(got.Subscriptions) != 1 || !got.Subscriptions[0].MarkedForCancellationAt.Equal(markedAt) {
		t.Errorf("subscriptions = %+v, %v; want the active line marked at %v", got.Subscriptions, err, markedAt)
	}
	if len(got.Possible) != 1 || !got.Possible[0].MarkedForCancellationAt.IsZero() {
		t.Errorf("possible = %+v, want the unmarked line without a mark", got.Possible)
	}
}
