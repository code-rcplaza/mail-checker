package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var (
	_ ports.SubscriptionDecisions = Decisions{}
	_ ports.LineAliases           = Aliases{}
	_ ports.LineOverrides         = Overrides{}
	_ ports.ManualSubscriptions   = Manuals{}
	_ ports.HiddenLines           = HiddenLines{}
	_ ports.Cancellations         = Cancellations{}
)

// trickyKey has every character a line key can carry.
const trickyKey domain.GroupKey = "merchant:SPOTIFY P#12 AB:CLP||d05"

var decidedAt = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

// twoOwners opens a database with Ana and Beto.
func twoOwners(t *testing.T) (*sql.DB, int64, int64) {
	t.Helper()
	db := sqlitetest.OpenMigrated(t)
	users := sqlitetest.Users{DB: db}
	return db, users.Insert(t, "ana@example.com"), users.Insert(t, "beto@example.com")
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDecisions(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	decisions := NewDecisions(db, func() time.Time { return decidedAt })
	ctx := context.Background()
	confirmed := domain.Decision{Kind: domain.DecisionConfirmed, Cadence: domain.CadenceYearly}
	dismissed := domain.Decision{Kind: domain.DecisionDismissed}
	must(t, decisions.Save(ctx, anaID, domain.LineDecision{Key: trickyKey, Decision: dismissed}))
	must(t, decisions.Save(ctx, anaID, domain.LineDecision{Key: trickyKey, Decision: confirmed}))
	must(t, decisions.Save(ctx, anaID, domain.LineDecision{Key: "other", Decision: dismissed}))
	must(t, decisions.Save(ctx, betoID, domain.LineDecision{Key: trickyKey, Decision: dismissed}))

	got, err := decisions.List(ctx, anaID)
	if want := map[domain.GroupKey]domain.Decision{trickyKey: confirmed, "other": dismissed}; err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("Ana's decisions = %v, %v; want the last one per key %v", got, err, want)
	}

	must(t, decisions.Delete(ctx, anaID, "other"))
	must(t, decisions.Delete(ctx, anaID, "never-decided"))
	got, _ = decisions.List(ctx, anaID)
	betos, _ := decisions.List(ctx, betoID)
	if len(got) != 1 || betos[trickyKey] != dismissed {
		t.Errorf("after delete: Ana = %v, Beto = %v; want Ana with one and Beto untouched", got, betos)
	}
}

func TestAliases(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	aliases := NewAliases(db)
	ctx := context.Background()
	must(t, aliases.Save(ctx, anaID, domain.LineAlias{Key: trickyKey, Alias: "Padres"}))
	must(t, aliases.Save(ctx, anaID, domain.LineAlias{Key: trickyKey, Alias: "Música de mamá"}))
	must(t, aliases.Save(ctx, betoID, domain.LineAlias{Key: trickyKey, Alias: "Mío"}))

	got, err := aliases.List(ctx, anaID)
	if want := map[domain.GroupKey]string{trickyKey: "Música de mamá"}; err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("Ana's aliases = %v, %v; want %v", got, err, want)
	}

	must(t, aliases.Delete(ctx, anaID, trickyKey))
	must(t, aliases.Delete(ctx, anaID, trickyKey))
	got, _ = aliases.List(ctx, anaID)
	betos, _ := aliases.List(ctx, betoID)
	if len(got) != 0 || betos[trickyKey] != "Mío" {
		t.Errorf("after delete: Ana = %v, Beto = %v; want Ana empty and Beto untouched", got, betos)
	}
}

func TestOverrides(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	overrides := NewOverrides(db)
	ctx := context.Background()
	onlyAmount := domain.LineOverride{Amount: 900000}
	onlyDay := domain.LineOverride{BillingDay: 31}
	must(t, overrides.Save(ctx, anaID, domain.OverriddenLine{Key: trickyKey, Override: onlyAmount}))
	must(t, overrides.Save(ctx, anaID, domain.OverriddenLine{Key: "day", Override: onlyDay}))
	must(t, overrides.Save(ctx, betoID, domain.OverriddenLine{Key: trickyKey, Override: domain.LineOverride{Amount: 1, BillingDay: 1}}))

	got, err := overrides.List(ctx, anaID)
	if want := map[domain.GroupKey]domain.LineOverride{trickyKey: onlyAmount, "day": onlyDay}; err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("Ana's overrides = %v, %v; want %v", got, err, want)
	}

	must(t, overrides.Delete(ctx, anaID, trickyKey))
	must(t, overrides.Delete(ctx, anaID, "never-set"))
	got, _ = overrides.List(ctx, anaID)
	if _, isThere := got[trickyKey]; isThere || len(got) != 1 {
		t.Errorf("after delete: %v, want only the day override", got)
	}
}

func TestManuals(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	manuals := NewManuals(db)
	ctx := context.Background()
	gym := domain.ManualSubscription{PublicID: "manual-gym", Name: "Gimnasio", Amount: 2500000, Currency: shared.CurrencyCLP, BillingDay: 5}
	must(t, manuals.Create(ctx, anaID, gym))
	must(t, manuals.Create(ctx, betoID, domain.ManualSubscription{PublicID: "manual-beto", Name: "Club", Amount: 100, Currency: shared.CurrencyUSD, BillingDay: 1}))

	got, err := manuals.List(ctx, anaID)
	if err != nil || !reflect.DeepEqual(got, []domain.ManualSubscription{gym}) {
		t.Fatalf("Ana's manuals = %+v, %v; want only %+v", got, err, gym)
	}

	edited := gym
	edited.Name, edited.Amount, edited.Currency, edited.BillingDay = "Gimnasio y piscina", 3200000, shared.CurrencyCLP, 10
	must(t, manuals.Update(ctx, anaID, edited))
	got, _ = manuals.List(ctx, anaID)
	if !reflect.DeepEqual(got, []domain.ManualSubscription{edited}) {
		t.Errorf("after update = %+v, want %+v with the same public id", got, edited)
	}

	must(t, manuals.Delete(ctx, anaID, gym.PublicID))
	got, _ = manuals.List(ctx, anaID)
	if len(got) != 0 {
		t.Errorf("after delete = %+v, want none", got)
	}
}

func TestManuals_OfAnotherUserOrUnknownAreNotFound(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	manuals := NewManuals(db)
	ctx := context.Background()
	gym := domain.ManualSubscription{PublicID: "manual-gym", Name: "Gimnasio", Amount: 2500000, Currency: shared.CurrencyCLP, BillingDay: 5}
	must(t, manuals.Create(ctx, anaID, gym))
	unknown := gym
	unknown.PublicID = "never-created"

	for name, err := range map[string]error{
		"update by another user": manuals.Update(ctx, betoID, gym),
		"update of an unknown":   manuals.Update(ctx, anaID, unknown),
		"delete by another user": manuals.Delete(ctx, betoID, gym.PublicID),
		"delete of an unknown":   manuals.Delete(ctx, anaID, unknown.PublicID),
	} {
		if !errors.Is(err, shared.ErrNotFound) {
			t.Errorf("%s error = %v, want ErrNotFound", name, err)
		}
	}
	if got, _ := manuals.List(ctx, anaID); !reflect.DeepEqual(got, []domain.ManualSubscription{gym}) {
		t.Errorf("Ana's manual = %+v, want it untouched", got)
	}
}

func TestHiddenLines(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	hidden := NewHiddenLines(db)
	ctx := context.Background()
	for range 2 {
		must(t, hidden.Hide(ctx, anaID, domain.HiddenLine{Key: trickyKey, HiddenAt: decidedAt}))
	}
	must(t, hidden.Hide(ctx, betoID, domain.HiddenLine{Key: "beto-line", HiddenAt: decidedAt}))

	got, err := hidden.List(ctx, anaID)

	if err != nil || !reflect.DeepEqual(got, []domain.GroupKey{trickyKey}) {
		t.Errorf("Ana's hidden lines = %v, %v; want only %q once", got, err, trickyKey)
	}
}

func TestCancellations(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	cancellations := NewCancellations(db)
	ctx := context.Background()
	must(t, cancellations.Mark(ctx, anaID, domain.CancellationMark{Key: trickyKey, MarkedAt: decidedAt}))
	must(t, cancellations.Mark(ctx, anaID, domain.CancellationMark{Key: trickyKey, MarkedAt: decidedAt.AddDate(0, 0, 3)}))
	must(t, cancellations.Mark(ctx, betoID, domain.CancellationMark{Key: "beto-line", MarkedAt: decidedAt}))

	got, err := cancellations.List(ctx, anaID)

	if want := map[domain.GroupKey]time.Time{trickyKey: decidedAt}; err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("Ana's marks = %v, %v; want %v (the first date)", got, err, want)
	}
}

func TestPriceAlertDecisions_KeepTheHistoryOfEachLine(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	decisions := NewPriceAlertDecisions(db)
	ctx := context.Background()
	for _, amount := range []shared.Money{1050000, 1100000, 1050000} {
		must(t, decisions.Save(ctx, anaID, domain.KeptRise{Key: trickyKey, NewAmount: amount}))
	}
	must(t, decisions.Save(ctx, betoID, domain.KeptRise{Key: trickyKey, NewAmount: 999}))

	got, err := decisions.List(ctx, anaID)

	if want := map[domain.GroupKey][]shared.Money{trickyKey: {1050000, 1100000}}; err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("Ana's kept amounts = %v, %v; want both, oldest first, once each: %v", got, err, want)
	}
}
