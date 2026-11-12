package domain

import (
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// musicAppCharges builds monthly CLP charges of an unknown merchant with the given amounts.
func musicAppCharges(amounts ...shared.Money) []shared.Charge {
	charges := make([]shared.Charge, 0, len(amounts))
	for i, amount := range amounts {
		charges = append(charges, charge("m", "MUSIC APP", amount, shared.CurrencyCLP, 1+30*i))
	}
	return charges
}

func TestAnalyzeUnknownMerchantAmount_SmallVariationIsConsistentWithoutPlanChange(t *testing.T) {
	got := AnalyzeUnknownMerchantAmount(musicAppCharges(2000, 2000, 2100))

	if !got.IsConsistent {
		t.Errorf("IsConsistent = false, want true")
	}
	if got.Current != 2100 || got.Currency != shared.CurrencyCLP {
		t.Errorf("current = %d %s, want 2100 CLP", got.Current, got.Currency)
	}
	if len(got.PlanChanges) != 0 {
		t.Errorf("plan changes = %+v, want none", got.PlanChanges)
	}
}

func TestAnalyzeUnknownMerchantAmount_OneJumpIsConsistentWithPlanChange(t *testing.T) {
	got := AnalyzeUnknownMerchantAmount(musicAppCharges(2000, 2000, 10000, 10000))

	if !got.IsConsistent {
		t.Errorf("IsConsistent = false, want true")
	}
	if got.Current != 10000 {
		t.Errorf("current = %d, want 10000", got.Current)
	}
	want := []PlanChange{{From: 2000, To: 10000, FromCurrency: shared.CurrencyCLP, ToCurrency: shared.CurrencyCLP, At: dayOf(61)}}
	if !reflect.DeepEqual(got.PlanChanges, want) {
		t.Errorf("plan changes = %+v, want %+v", got.PlanChanges, want)
	}
}

func TestAnalyzeUnknownMerchantAmount_AmountChangingOftenIsNotConsistent(t *testing.T) {
	// The near miss has a single jump, so the case cannot pass just because
	// IsConsistent defaults to false.
	if !AnalyzeUnknownMerchantAmount(musicAppCharges(800000, 1450000, 1450000)).IsConsistent {
		t.Fatalf("near miss is not consistent; the case would prove nothing")
	}

	got := AnalyzeUnknownMerchantAmount(musicAppCharges(800000, 1450000, 520000))

	if got.IsConsistent {
		t.Errorf("IsConsistent = true, want false")
	}
}
