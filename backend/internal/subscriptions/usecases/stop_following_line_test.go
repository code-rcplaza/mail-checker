package usecases

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.HiddenLines = hiddenLinesByUser{}

// hiddenLinesByUser keeps one mark per line, like the UNIQUE of the table.
type hiddenLinesByUser map[int64]map[domain.GroupKey]time.Time

func (h hiddenLinesByUser) List(_ context.Context, userID int64) ([]domain.GroupKey, error) {
	keys := make([]domain.GroupKey, 0, len(h[userID]))
	for key := range h[userID] {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys, nil
}

func (h hiddenLinesByUser) Hide(_ context.Context, userID int64, hidden domain.HiddenLine) error {
	if h[userID] == nil {
		h[userID] = make(map[domain.GroupKey]time.Time)
	}
	if _, isHidden := h[userID][hidden.Key]; isHidden {
		return nil
	}
	h[userID][hidden.Key] = hidden.HiddenAt
	return nil
}

func TestStopFollowingLine(t *testing.T) {
	earlier := testNow.AddDate(0, 0, -3)
	cases := []struct {
		name    string
		userID  int64
		key     domain.GroupKey
		hidden  hiddenLinesByUser
		wantErr error
		want    map[domain.GroupKey]time.Time
	}{
		{"an existing line is hidden by its key", userOne, streamingLine, hiddenLinesByUser{}, nil, map[domain.GroupKey]time.Time{streamingLine: testNow}},
		{"a hidden line hidden again keeps its mark", userOne, streamingLine, hiddenLinesByUser{userOne: {streamingLine: earlier}}, nil, map[domain.GroupKey]time.Time{streamingLine: earlier}},
		{"a key that is not a line is not found", userOne, "merchant:NEVER SEEN:CLP||d01", hiddenLinesByUser{}, shared.ErrNotFound, nil},
		{"the line of another user is not found", userTwo, streamingLine, hiddenLinesByUser{}, shared.ErrNotFound, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			charges := chargesByUser{userOne: streamingOnce()}

			err := NewStopFollowingLine(linesOf(charges), tc.hidden).Execute(context.Background(), tc.userID, tc.key)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				assertKeyedError(t, err, shared.ErrNotFound, domain.KeyLineNotFound)
			}
			if got := tc.hidden[tc.userID]; !reflect.DeepEqual(got, tc.want) {
				t.Errorf("hidden lines of user %d = %v, want %v", tc.userID, got, tc.want)
			}
		})
	}
}

func TestGetSummary_LeavesOutTheHiddenLinesOfTheUser(t *testing.T) {
	charges := chargesByUser{userOne: streamingOnce(), userTwo: streamingOnce()}
	hidden := hiddenLinesByUser{userOne: {streamingLine: testNow}}
	sources := summarySourcesOf(charges)
	sources.Hidden = hidden

	mine, err := NewGetSummary(sources).Execute(context.Background(), userOne, 6)
	if err != nil || len(mine.Possible) != 0 {
		t.Errorf("summary of user %d = %v possible, %v; want the hidden line left out", userOne, keysOfLines(mine.Possible), err)
	}
	theirs, err := NewGetSummary(sources).Execute(context.Background(), userTwo, 6)
	if err != nil || len(theirs.Possible) != 1 {
		t.Errorf("summary of user %d = %v possible, %v; want their line, which they did not hide", userTwo, keysOfLines(theirs.Possible), err)
	}
}

func keysOfLines(lines []domain.Subscription) []domain.GroupKey {
	keys := make([]domain.GroupKey, 0, len(lines))
	for _, line := range lines {
		keys = append(keys, line.Key)
	}
	return keys
}
