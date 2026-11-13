package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

const (
	userOne int64 = 1
	userTwo int64 = 2
)

// The doubles below keep each user's rows apart, like the real repositories.
type chargesByUser map[int64][]shared.Charge

func (c chargesByUser) ListSince(_ context.Context, userID int64, _ time.Time) ([]shared.Charge, error) {
	return c[userID], nil
}

type decisionsByUser map[int64]map[domain.GroupKey]domain.Decision

func (d decisionsByUser) List(_ context.Context, userID int64) (map[domain.GroupKey]domain.Decision, error) {
	return d[userID], nil
}

func (d decisionsByUser) Save(_ context.Context, userID int64, key domain.GroupKey, decision domain.Decision) error {
	if d[userID] == nil {
		d[userID] = make(map[domain.GroupKey]domain.Decision)
	}
	d[userID][key] = decision
	return nil
}

type aliasesByUser map[int64]map[domain.GroupKey]string

func (a aliasesByUser) List(_ context.Context, userID int64) (map[domain.GroupKey]string, error) {
	return a[userID], nil
}

func (a aliasesByUser) Save(_ context.Context, userID int64, key domain.GroupKey, alias string) error {
	if a[userID] == nil {
		a[userID] = make(map[domain.GroupKey]string)
	}
	a[userID][key] = alias
	return nil
}

func (a aliasesByUser) Delete(_ context.Context, userID int64, key domain.GroupKey) error {
	delete(a[userID], key)
	return nil
}

// A known service charged once: a possible line the user can decide and name.
const streamingLine domain.GroupKey = "service:Streaming Co||d02"

func streamingOnce() []shared.Charge {
	return []shared.Charge{chargeDaysAgo("s1", "STREAMING CO", shared.SourceStatement, 599000, 10)}
}

func fixedNow() time.Time { return testNow }

func TestDecisionsByUser_ADecisionIsNotSeenByAnotherUser(t *testing.T) {
	decisions := decisionsByUser{}
	confirmed := domain.Decision{Kind: domain.DecisionConfirmed, Cadence: domain.CadenceMonthly}

	if err := decisions.Save(context.Background(), userOne, streamingLine, confirmed); err != nil {
		t.Fatalf("Save error = %v, want nil", err)
	}

	theirs, err := decisions.List(context.Background(), userTwo)
	if err != nil || len(theirs) != 0 {
		t.Errorf("List of the other user = (%v, %v), want nothing", theirs, err)
	}
}

func TestGetSummary_ReadsOnlyTheDataOfTheUser(t *testing.T) {
	charges := chargesByUser{userOne: streamingOnce(), userTwo: streamingOnce()}
	decisions := decisionsByUser{userTwo: {streamingLine: {Kind: domain.DecisionDismissed}}}
	aliases := aliasesByUser{userTwo: {streamingLine: "Theirs"}}
	useCase := NewGetSummary(SummarySources{Charges: charges, Catalog: fakeCatalog{}, Decisions: decisions, Aliases: aliases, Rates: ratesByUser{}, KeptRises: keptRisesByUser{}, Manuals: manualsByUser{}, Overrides: overridesByUser{}}, fixedNow)

	got, err := useCase.Execute(context.Background(), userOne, 6)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if len(got.Possible) != 1 || got.Possible[0].Decision != nil || got.Possible[0].Alias != "" {
		t.Errorf("possible = %+v, want the undecided unnamed line of user %d", got.Possible, userOne)
	}
}

func TestDecidePossibleSubscription_LineOfAnotherUserIsNotFound(t *testing.T) {
	charges := chargesByUser{userTwo: streamingOnce()}
	decisions := decisionsByUser{}
	useCase := NewDecidePossibleSubscription(charges, fakeCatalog{}, decisions, fixedNow)

	err := useCase.Execute(context.Background(), userOne, streamingLine, "confirm", "monthly")

	if !errors.Is(err, shared.ErrNotFound) || len(decisions) != 0 {
		t.Errorf("Execute = %v with decisions %v, want ErrNotFound and nothing saved", err, decisions)
	}
}

func TestNameLine_LineOfAnotherUserIsNotFound(t *testing.T) {
	charges := chargesByUser{userTwo: streamingOnce()}
	aliases := aliasesByUser{}
	useCase := NewNameLine(charges, fakeCatalog{}, aliases, fixedNow)

	err := useCase.Execute(context.Background(), userOne, streamingLine, "Mine")

	assertKeyedError(t, err, shared.ErrNotFound, domain.KeyLineNotFound)
	if len(aliases) != 0 {
		t.Errorf("aliases = %v, want nothing saved", aliases)
	}
}
