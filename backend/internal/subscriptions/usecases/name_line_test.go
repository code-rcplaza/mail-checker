package usecases

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type recordingAliases struct {
	aliases   map[domain.GroupKey]string
	saved     map[domain.GroupKey]string
	deleted   []domain.GroupKey
	listErr   error
	saveErr   error
	deleteErr error
}

func newRecordingAliases() *recordingAliases {
	return &recordingAliases{saved: make(map[domain.GroupKey]string)}
}

func (r *recordingAliases) List(context.Context, int64) (map[domain.GroupKey]string, error) {
	return r.aliases, r.listErr
}

func (r *recordingAliases) Save(_ context.Context, _ int64, alias domain.LineAlias) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.saved[alias.Key] = alias.Alias
	return nil
}

func (r *recordingAliases) Delete(_ context.Context, _ int64, key domain.GroupKey) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	r.deleted = append(r.deleted, key)
	return nil
}

const (
	activeLineKey   domain.GroupKey = "service:Streaming Co||d02"
	inactiveLineKey domain.GroupKey = "merchant:OLD APP:CLP||d26"
)

// chargesWithEveryKindOfLine: Streaming Co is active, MUSIC APP possible and OLD APP inactive.
func chargesWithEveryKindOfLine() *fakeChargeReader {
	reader := chargesWithPossibleAndActive()
	reader.charges = append(reader.charges,
		chargeDaysAgo("o1", "OLD APP", shared.SourceEmail, 299000, 200),
		chargeDaysAgo("o2", "OLD APP", shared.SourceEmail, 299000, 170),
	)
	return reader
}

func nameLine(aliases *recordingAliases, key domain.GroupKey, alias string) error {
	useCase := NewNameLine(linesOf(chargesWithEveryKindOfLine()), aliases)
	return useCase.Execute(context.Background(), userOne, domain.LineAlias{Key: key, Alias: alias})
}

func TestNameLine_SavesTheTrimmedAliasOfAnyLine(t *testing.T) {
	for _, key := range []domain.GroupKey{activeLineKey, possibleKey, inactiveLineKey} {
		t.Run(string(key), func(t *testing.T) {
			aliases := newRecordingAliases()

			err := nameLine(aliases, key, "  Padres ")

			if err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}
			if want := map[domain.GroupKey]string{key: "Padres"}; !reflect.DeepEqual(aliases.saved, want) {
				t.Errorf("saved = %v, want %v", aliases.saved, want)
			}
		})
	}
}

func TestNameLine_BlankAliasDeletesIt(t *testing.T) {
	aliases := newRecordingAliases()

	err := nameLine(aliases, activeLineKey, "   ")

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if !reflect.DeepEqual(aliases.deleted, []domain.GroupKey{activeLineKey}) || len(aliases.saved) != 0 {
		t.Errorf("deleted = %v, saved = %v; want only %s deleted", aliases.deleted, aliases.saved, activeLineKey)
	}
}

func TestNameLine_TooLongAliasIsRejectedAndNotSaved(t *testing.T) {
	aliases := newRecordingAliases()

	err := nameLine(aliases, activeLineKey, strings.Repeat("a", 41))

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyAliasTooLong)
	if len(aliases.saved) != 0 || len(aliases.deleted) != 0 {
		t.Errorf("saved = %v, deleted = %v; want nothing stored", aliases.saved, aliases.deleted)
	}
}

func TestNameLine_KeyThatIsNotALineIsNotFound(t *testing.T) {
	for _, key := range []domain.GroupKey{"service:Streaming Co", "merchant:NEVER SEEN:CLP||d01"} {
		t.Run(string(key), func(t *testing.T) {
			aliases := newRecordingAliases()

			err := nameLine(aliases, key, "Padres")

			assertKeyedError(t, err, shared.ErrNotFound, domain.KeyLineNotFound)
			if len(aliases.saved) != 0 {
				t.Errorf("saved = %v, want nothing saved", aliases.saved)
			}
		})
	}
}

func TestNameLine_WrapsAliasErrors(t *testing.T) {
	errDatabase := errors.New("database is locked")
	cases := []struct {
		name    string
		aliases *recordingAliases
		alias   string
	}{
		{"save fails", &recordingAliases{saved: map[domain.GroupKey]string{}, saveErr: errDatabase}, "Padres"},
		{"delete fails", &recordingAliases{saved: map[domain.GroupKey]string{}, deleteErr: errDatabase}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := nameLine(tc.aliases, activeLineKey, tc.alias)

			if !errors.Is(err, errDatabase) {
				t.Fatalf("errors.Is(err, errDatabase) = false (err = %v), want true", err)
			}
			if err.Error() == errDatabase.Error() {
				t.Errorf("err = %q, want it wrapped with context", err)
			}
		})
	}
}

func TestGetSummary_SubscriptionsCarryTheAliasOfTheirLine(t *testing.T) {
	aliases := newRecordingAliases()
	aliases.aliases = map[domain.GroupKey]string{activeLineKey: "Padres"}
	sources := summarySourcesOf(chargesWithPossibleAndActive())
	sources.Aliases = aliases
	useCase := NewGetSummary(sources)

	got, err := useCase.Execute(context.Background(), userOne, 6)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if len(got.Subscriptions) != 1 || got.Subscriptions[0].Alias != "Padres" {
		t.Errorf("subscriptions = %+v, want the active line named Padres", got.Subscriptions)
	}
	if len(got.Possible) != 1 || got.Possible[0].Alias != "" {
		t.Errorf("possible = %+v, want the possible line without alias", got.Possible)
	}
}

func TestGetSummary_WrapsAliasReadError(t *testing.T) {
	errDatabase := errors.New("database is locked")
	aliases := newRecordingAliases()
	aliases.listErr = errDatabase
	sources := summarySourcesOf(chargesWithPossibleAndActive())
	sources.Aliases = aliases
	useCase := NewGetSummary(sources)

	_, err := useCase.Execute(context.Background(), userOne, 6)

	if !errors.Is(err, errDatabase) {
		t.Fatalf("errors.Is(err, errDatabase) = false (err = %v), want true", err)
	}
	if err.Error() == errDatabase.Error() {
		t.Errorf("err = %q, want it wrapped with context", err)
	}
}
