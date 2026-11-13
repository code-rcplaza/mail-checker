package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

func (d decisionsByUser) Delete(_ context.Context, userID int64, key domain.GroupKey) error {
	delete(d[userID], key)
	return nil
}

func TestDecisionsByUser_Delete(t *testing.T) {
	cases := []struct {
		name  string
		saved map[domain.GroupKey]domain.Decision
	}{
		{"a saved decision is no longer listed", map[domain.GroupKey]domain.Decision{streamingLine: {Kind: domain.DecisionDismissed}}},
		{"a key without a decision does not fail", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			decisions := decisionsByUser{userOne: tc.saved}

			err := decisions.Delete(context.Background(), userOne, streamingLine)

			listed, _ := decisions.List(context.Background(), userOne)
			if _, isListed := listed[streamingLine]; err != nil || isListed {
				t.Errorf("Delete = %v, still listed = %v; want no error and not listed", err, isListed)
			}
		})
	}
}

func possibleKeysAfterUndo(t *testing.T, decisions decisionsByUser) []domain.GroupKey {
	t.Helper()
	charges := chargesByUser{userOne: streamingOnce()}
	if err := NewUndoDecision(decisions).Execute(context.Background(), userOne, streamingLine); err != nil {
		t.Fatalf("UndoDecision error = %v, want nil", err)
	}
	summary, err := NewGetSummary(SummarySources{Charges: charges, Catalog: fakeCatalog{}, Decisions: decisions, Aliases: aliasesByUser{}, Rates: ratesByUser{}, KeptRises: keptRisesByUser{}, Manuals: manualsByUser{}, Overrides: overridesByUser{}}, fixedNow).Execute(context.Background(), userOne, 6)
	if err != nil {
		t.Fatalf("GetSummary error = %v, want nil", err)
	}
	keys := make([]domain.GroupKey, 0, len(summary.Possible))
	for _, s := range summary.Possible {
		if s.Decision == nil {
			keys = append(keys, s.Key)
		}
	}
	return keys
}

func TestUndoDecision_LineGoesBackToPossible(t *testing.T) {
	cases := []struct {
		name     string
		decision *domain.Decision
	}{
		{"a confirmed line", &domain.Decision{Kind: domain.DecisionConfirmed, Cadence: domain.CadenceMonthly}},
		{"a dismissed line", &domain.Decision{Kind: domain.DecisionDismissed}},
		{"a line without a decision", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			decisions := decisionsByUser{userOne: {}}
			if tc.decision != nil {
				decisions[userOne][streamingLine] = *tc.decision
			}

			got := possibleKeysAfterUndo(t, decisions)

			if len(got) != 1 || got[0] != streamingLine {
				t.Errorf("undecided possible = %v, want [%s]", got, streamingLine)
			}
		})
	}
}

func TestUndoDecision_DoesNotTouchTheDecisionOfAnotherUser(t *testing.T) {
	theirs := domain.Decision{Kind: domain.DecisionDismissed}
	decisions := decisionsByUser{userTwo: {streamingLine: theirs}}

	if err := NewUndoDecision(decisions).Execute(context.Background(), userOne, streamingLine); err != nil {
		t.Fatalf("UndoDecision error = %v, want nil", err)
	}

	if got, isKept := decisions[userTwo][streamingLine]; !isKept || got != theirs {
		t.Errorf("decision of the other user = (%+v, %v), want it untouched", got, isKept)
	}
}

type failingDecisions struct {
	decisionsByUser
	err error
}

func (f failingDecisions) Delete(context.Context, int64, domain.GroupKey) error { return f.err }

func TestUndoDecision_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")

	err := NewUndoDecision(failingDecisions{err: cause}).Execute(context.Background(), userOne, streamingLine)

	if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), "delete decision for ") {
		t.Errorf("error = %v, want \"delete decision for …\" wrapping %v", err, cause)
	}
}
