package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

var testNow = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)

type fakeChargeReader struct {
	charges    []shared.Charge
	err        error
	askedSince time.Time
}

func (f *fakeChargeReader) ListSince(_ context.Context, _ int64, since time.Time) ([]shared.Charge, error) {
	f.askedSince = since
	return f.charges, f.err
}

// fakeCatalog knows one service under two descriptors.
type fakeCatalog struct{}

func (fakeCatalog) Identify(merchant string) (domain.Service, bool) {
	if merchant == "STREAMING CO" || merchant == "STRMNG*CO" {
		return domain.Service{Name: "Streaming Co"}, true
	}
	return domain.Service{}, false
}

func (fakeCatalog) Normalize(merchant string) string { return merchant }

func (fakeCatalog) LooksLikeSubscription(string) bool { return false }

type fakeDecisions struct {
	decisions map[domain.GroupKey]domain.Decision
	err       error
}

func (f fakeDecisions) List(context.Context, int64) (map[domain.GroupKey]domain.Decision, error) {
	return f.decisions, f.err
}

// Save is unused by GetSummary; it is here to satisfy the port.
func (fakeDecisions) Save(context.Context, int64, domain.GroupKey, domain.Decision) error { return nil }

// Delete is unused by GetSummary; it is here to satisfy the port.
func (fakeDecisions) Delete(context.Context, int64, domain.GroupKey) error { return nil }

func chargeDaysAgo(externalID, merchant string, source shared.Source, amount shared.Money, daysAgo int) shared.Charge {
	return shared.Charge{
		Bank:       "bci",
		Source:     source,
		ExternalID: externalID,
		OccurredAt: testNow.AddDate(0, 0, -daysAgo),
		Amount:     amount,
		Currency:   shared.CurrencyCLP,
		Merchant:   merchant,
		Card:       "credito",
	}
}

func getSummary(reader *fakeChargeReader, decisions fakeDecisions, months int) (domain.Summary, error) {
	useCase := NewGetSummary(SummarySources{Charges: reader, Catalog: fakeCatalog{}, Decisions: decisions, Aliases: newRecordingAliases(), Rates: ratesByUser{}, KeptRises: keptRisesByUser{}, Manuals: manualsByUser{}, Overrides: overridesByUser{}}, func() time.Time { return testNow })
	return useCase.Execute(context.Background(), userOne, months)
}

func TestGetSummary_ChargeInEmailAndStatementCountsOnce(t *testing.T) {
	reader := &fakeChargeReader{charges: []shared.Charge{
		chargeDaysAgo("email-1", "STREAMING CO", shared.SourceEmail, 599000, 40),
		chargeDaysAgo("email-2", "STREAMING CO", shared.SourceEmail, 599000, 10),
		chargeDaysAgo("statement-2", "STRMNG*CO", shared.SourceStatement, 599000, 8),
	}}

	got, err := getSummary(reader, fakeDecisions{}, 6)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if len(got.Subscriptions) != 1 {
		t.Fatalf("subscriptions = %d, want 1", len(got.Subscriptions))
	}
	if count := len(got.Subscriptions[0].Charges); count != 2 {
		t.Errorf("subscription charges = %d, want 2 (the statement line is the same purchase)", count)
	}
}

func TestGetSummary_ReadsChargesSinceNowMinusMonths(t *testing.T) {
	reader := &fakeChargeReader{}

	if _, err := getSummary(reader, fakeDecisions{}, 6); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	want := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)
	if !reader.askedSince.Equal(want) {
		t.Errorf("ListSince(since = %v), want %v", reader.askedSince, want)
	}
}

func TestGetSummary_ActiveServiceAndConfirmedPossibleAreSubscriptions(t *testing.T) {
	reader := &fakeChargeReader{charges: []shared.Charge{
		chargeDaysAgo("s1", "STREAMING CO", shared.SourceEmail, 599000, 40),
		chargeDaysAgo("s2", "STREAMING CO", shared.SourceEmail, 599000, 10),
		chargeDaysAgo("m1", "MUSIC APP", shared.SourceEmail, 399000, 40),
		chargeDaysAgo("m2", "MUSIC APP", shared.SourceEmail, 399000, 10),
	}}
	decisions := fakeDecisions{decisions: map[domain.GroupKey]domain.Decision{
		"merchant:MUSIC APP:CLP||d02": {Kind: domain.DecisionConfirmed, Cadence: domain.CadenceMonthly},
	}}

	got, err := getSummary(reader, decisions, 6)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	keys := make(map[domain.GroupKey]bool)
	for _, s := range got.Subscriptions {
		keys[s.Key] = true
	}
	if !keys["service:Streaming Co||d02"] || !keys["merchant:MUSIC APP:CLP||d02"] || len(keys) != 2 {
		t.Errorf("subscriptions = %v, want the active service and the confirmed possible", keys)
	}
}

func TestGetSummary_RejectsMonthsOutOfRange(t *testing.T) {
	for _, months := range []int{0, 25} {
		_, err := getSummary(&fakeChargeReader{}, fakeDecisions{}, months)

		if !errors.Is(err, shared.ErrInvalid) {
			t.Fatalf("months = %d: errors.Is(err, ErrInvalid) = false (err = %v), want true", months, err)
		}
		var keyed shared.KeyedError
		if !errors.As(err, &keyed) || keyed.Key != "subscriptions.months_out_of_range" {
			t.Errorf("months = %d: key = %q, want subscriptions.months_out_of_range", months, keyed.Key)
		}
	}
}

func TestGetSummary_WrapsPortErrors(t *testing.T) {
	errDatabase := errors.New("database is locked")
	cases := []struct {
		name      string
		reader    *fakeChargeReader
		decisions fakeDecisions
	}{
		{"charge reader fails", &fakeChargeReader{err: errDatabase}, fakeDecisions{}},
		{"decisions fail", &fakeChargeReader{}, fakeDecisions{err: errDatabase}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getSummary(tc.reader, tc.decisions, 6)

			if !errors.Is(err, errDatabase) {
				t.Errorf("errors.Is(err, errDatabase) = false (err = %v), want true", err)
			}
			if err != nil && err.Error() == errDatabase.Error() {
				t.Errorf("err = %q, want it wrapped with context", err)
			}
		})
	}
}
