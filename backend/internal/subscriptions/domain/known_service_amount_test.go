package domain

import (
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func dayOf(day int) time.Time {
	return time.Date(2026, 1, day, 12, 0, 0, 0, time.UTC)
}

func TestAnalyzeKnownServiceAmount_PlanChangeFrom20To100(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 2000, shared.CurrencyUSD, 1),
		charge("b", "STREAMING CO", 2000, shared.CurrencyUSD, 31),
		charge("c", "STREAMING CO", 10000, shared.CurrencyUSD, 61),
		charge("d", "STREAMING CO", 10000, shared.CurrencyUSD, 91),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if got.Current != 10000 || got.Currency != shared.CurrencyUSD {
		t.Errorf("current = %d %s, want 10000 USD", got.Current, got.Currency)
	}
	want := []PlanChange{{From: 2000, To: 10000, FromCurrency: shared.CurrencyUSD, ToCurrency: shared.CurrencyUSD, At: dayOf(61)}}
	if !reflect.DeepEqual(got.PlanChanges, want) {
		t.Errorf("plan changes = %+v, want %+v", got.PlanChanges, want)
	}
}

func TestAnalyzeKnownServiceAmount_VariationWithin25PercentIsNotAPlanChange(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 659000, shared.CurrencyCLP, 1),
		charge("b", "STREAMING CO", 659000, shared.CurrencyCLP, 31),
		charge("c", "STREAMING CO", 719000, shared.CurrencyCLP, 61),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if got.Current != 719000 {
		t.Errorf("current = %d, want 719000", got.Current)
	}
	if len(got.PlanChanges) != 0 {
		t.Errorf("plan changes = %+v, want none", got.PlanChanges)
	}
}

func TestAnalyzeKnownServiceAmount_ExtraChargeGoesToOtherCharges(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 2000, shared.CurrencyUSD, 1),
		charge("b", "STREAMING CO", 2000, shared.CurrencyUSD, 31),
		charge("extra", "STREAMING CO", 500, shared.CurrencyUSD, 41),
		charge("c", "STREAMING CO", 2000, shared.CurrencyUSD, 61),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if got.Current != 2000 {
		t.Errorf("current = %d, want 2000", got.Current)
	}
	if ids := externalIDs(got.OtherCharges); !reflect.DeepEqual(ids, []string{"extra"}) {
		t.Errorf("other charges = %v, want [extra]", ids)
	}
	if ids := externalIDs(got.RecurringCharges); !reflect.DeepEqual(ids, []string{"a", "b", "c"}) {
		t.Errorf("recurring charges = %v, want [a b c]", ids)
	}
	if len(got.PlanChanges) != 0 {
		t.Errorf("plan changes = %+v, want none", got.PlanChanges)
	}
}

// Story 002-12: the recurring pattern is decided by backing, not only by order.

func TestAnalyzeKnownServiceAmount_ExtraBeforeFirstRecurringChargeGoesToOtherCharges(t *testing.T) {
	charges := []shared.Charge{
		charge("extra", "STREAMING CO", 500, shared.CurrencyUSD, 1),
		charge("a", "STREAMING CO", 2000, shared.CurrencyUSD, 11),
		charge("b", "STREAMING CO", 2000, shared.CurrencyUSD, 41),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if ids := externalIDs(got.OtherCharges); !reflect.DeepEqual(ids, []string{"extra"}) {
		t.Errorf("other charges = %v, want [extra]", ids)
	}
	if got.Current != 2000 {
		t.Errorf("current = %d, want 2000", got.Current)
	}
	if len(got.PlanChanges) != 0 {
		t.Errorf("plan changes = %+v, want none", got.PlanChanges)
	}
}

func TestAnalyzeKnownServiceAmount_ExtraInYearlyServiceGoesToOtherCharges(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 12000, shared.CurrencyUSD, 1),
		charge("extra", "STREAMING CO", 500, shared.CurrencyUSD, 91),
		charge("b", "STREAMING CO", 12000, shared.CurrencyUSD, 366),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if ids := externalIDs(got.OtherCharges); !reflect.DeepEqual(ids, []string{"extra"}) {
		t.Errorf("other charges = %v, want [extra]", ids)
	}
	if ids := externalIDs(got.RecurringCharges); !reflect.DeepEqual(ids, []string{"a", "b"}) {
		t.Errorf("recurring charges = %v, want [a b]", ids)
	}
}

func TestAnalyzeKnownServiceAmount_RecentPlanChangeWithoutPairStaysRecurring(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 2000, shared.CurrencyUSD, 1),
		charge("b", "STREAMING CO", 2000, shared.CurrencyUSD, 31),
		charge("c", "STREAMING CO", 10000, shared.CurrencyUSD, 61),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if ids := externalIDs(got.RecurringCharges); !reflect.DeepEqual(ids, []string{"a", "b", "c"}) {
		t.Errorf("recurring charges = %v, want [a b c]", ids)
	}
	if got.Current != 10000 {
		t.Errorf("current = %d, want 10000", got.Current)
	}
	want := []PlanChange{{From: 2000, To: 10000, FromCurrency: shared.CurrencyUSD, ToCurrency: shared.CurrencyUSD, At: dayOf(61)}}
	if !reflect.DeepEqual(got.PlanChanges, want) {
		t.Errorf("plan changes = %+v, want %+v", got.PlanChanges, want)
	}
}

func TestAnalyzeKnownServiceAmount_PlanDowngradeKeepsFirstChargeRecurring(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 10000, shared.CurrencyUSD, 1),
		charge("b", "STREAMING CO", 2000, shared.CurrencyUSD, 31),
		charge("c", "STREAMING CO", 2000, shared.CurrencyUSD, 61),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if ids := externalIDs(got.RecurringCharges); !reflect.DeepEqual(ids, []string{"a", "b", "c"}) {
		t.Errorf("recurring charges = %v, want [a b c]", ids)
	}
	want := []PlanChange{{From: 10000, To: 2000, FromCurrency: shared.CurrencyUSD, ToCurrency: shared.CurrencyUSD, At: dayOf(31)}}
	if !reflect.DeepEqual(got.PlanChanges, want) {
		t.Errorf("plan changes = %+v, want %+v", got.PlanChanges, want)
	}
}

func TestAnalyzeKnownServiceAmount_CurrencyChangeIsPlanChange(t *testing.T) {
	charges := []shared.Charge{
		charge("a", "STREAMING CO", 599000, shared.CurrencyCLP, 1),
		charge("b", "STREAMING CO", 599000, shared.CurrencyCLP, 31),
		charge("c", "STRMNG*CO", 799, shared.CurrencyUSD, 61),
		charge("d", "STRMNG*CO", 799, shared.CurrencyUSD, 91),
	}

	got := AnalyzeKnownServiceAmount(charges)

	if got.Current != 799 || got.Currency != shared.CurrencyUSD {
		t.Errorf("current = %d %s, want 799 USD", got.Current, got.Currency)
	}
	want := []PlanChange{{From: 599000, To: 799, FromCurrency: shared.CurrencyCLP, ToCurrency: shared.CurrencyUSD, At: dayOf(61)}}
	if !reflect.DeepEqual(got.PlanChanges, want) {
		t.Errorf("plan changes = %+v, want %+v", got.PlanChanges, want)
	}
}
