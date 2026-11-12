package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

const possibleKey domain.GroupKey = "merchant:MUSIC APP:CLP"

// recordingDecisions keeps what was saved, or fails on Save when saveErr is set.
type recordingDecisions struct {
	saved   map[domain.GroupKey]domain.Decision
	saveErr error
}

func newRecordingDecisions() *recordingDecisions {
	return &recordingDecisions{saved: make(map[domain.GroupKey]domain.Decision)}
}

func (r *recordingDecisions) List(context.Context) (map[domain.GroupKey]domain.Decision, error) {
	return r.saved, nil
}

func (r *recordingDecisions) Save(_ context.Context, key domain.GroupKey, decision domain.Decision) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.saved[key] = decision
	return nil
}

// chargesWithPossibleAndActive: MUSIC APP is a current possible; Streaming Co is active.
func chargesWithPossibleAndActive() *fakeChargeReader {
	return &fakeChargeReader{charges: []shared.Charge{
		chargeDaysAgo("m1", "MUSIC APP", shared.SourceEmail, 399000, 40),
		chargeDaysAgo("m2", "MUSIC APP", shared.SourceEmail, 399000, 10),
		chargeDaysAgo("s1", "STREAMING CO", shared.SourceEmail, 599000, 40),
		chargeDaysAgo("s2", "STREAMING CO", shared.SourceEmail, 599000, 10),
	}}
}

func decide(decisions *recordingDecisions, key domain.GroupKey, action, cadence string) error {
	useCase := NewDecidePossibleSubscription(chargesWithPossibleAndActive(), fakeCatalog{}, decisions, func() time.Time { return testNow })
	return useCase.Execute(context.Background(), key, action, cadence)
}

func assertKeyedError(t *testing.T, err error, sentinel error, want shared.MessageKey) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is(err, %v) = false (err = %v), want true", sentinel, err)
	}
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) || keyed.Key != want {
		t.Errorf("key = %q, want %q", keyed.Key, want)
	}
}

func TestDecidePossibleSubscription_ConfirmSavesCadence(t *testing.T) {
	for _, cadence := range []domain.Cadence{domain.CadenceMonthly, domain.CadenceYearly} {
		t.Run(string(cadence), func(t *testing.T) {
			decisions := newRecordingDecisions()

			if err := decide(decisions, possibleKey, "confirm", string(cadence)); err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}

			want := domain.Decision{Kind: domain.DecisionConfirmed, Cadence: cadence}
			if got, isSaved := decisions.saved[possibleKey]; !isSaved || got != want {
				t.Errorf("saved[%s] = %+v (saved = %v), want %+v", possibleKey, got, isSaved, want)
			}
		})
	}
}

func TestDecidePossibleSubscription_ConfirmWithoutValidCadenceIsRejected(t *testing.T) {
	for _, cadence := range []string{"", "weekly"} {
		t.Run("cadence="+cadence, func(t *testing.T) {
			decisions := newRecordingDecisions()

			err := decide(decisions, possibleKey, "confirm", cadence)

			assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCadenceRequired)
			if len(decisions.saved) != 0 {
				t.Errorf("saved = %v, want nothing saved", decisions.saved)
			}
		})
	}
}

func TestDecidePossibleSubscription_DismissSavesDismissed(t *testing.T) {
	decisions := newRecordingDecisions()

	if err := decide(decisions, possibleKey, "dismiss", ""); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	want := domain.Decision{Kind: domain.DecisionDismissed}
	if got, isSaved := decisions.saved[possibleKey]; !isSaved || got != want {
		t.Errorf("saved[%s] = %+v (saved = %v), want %+v", possibleKey, got, isSaved, want)
	}
}

func TestDecidePossibleSubscription_KeyThatIsNotACurrentPossibleIsNotFound(t *testing.T) {
	keys := []domain.GroupKey{"service:Streaming Co", "merchant:NEVER SEEN:CLP"}
	for _, key := range keys {
		t.Run(string(key), func(t *testing.T) {
			decisions := newRecordingDecisions()

			err := decide(decisions, key, "dismiss", "")

			assertKeyedError(t, err, shared.ErrNotFound, domain.KeyPossibleNotFound)
			if len(decisions.saved) != 0 {
				t.Errorf("saved = %v, want nothing saved", decisions.saved)
			}
		})
	}
}

func TestDecidePossibleSubscription_UnknownDecisionIsRejected(t *testing.T) {
	err := decide(newRecordingDecisions(), possibleKey, "maybe", "")

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyDecisionUnknown)
}

func TestDecidePossibleSubscription_WrapsSaveError(t *testing.T) {
	errDatabase := errors.New("database is locked")
	decisions := newRecordingDecisions()
	decisions.saveErr = errDatabase

	err := decide(decisions, possibleKey, "dismiss", "")

	if !errors.Is(err, errDatabase) {
		t.Fatalf("errors.Is(err, errDatabase) = false (err = %v), want true", err)
	}
	if err.Error() == errDatabase.Error() {
		t.Errorf("err = %q, want it wrapped with context", err)
	}
}
