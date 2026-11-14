package domain

import (
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// netflixAndAppStoreRising bills Netflix on netflixDay, so its key is
// "service:Netflix||d<netflixDay>".
func netflixAndAppStoreRising(netflixDay int) []shared.Charge {
	return slices.Concat(
		risingCharges("netflix", "NETFLIX", 1299000, 1349000, netflixDay),
		risingCharges("apple", "APPLE COM BILL", 469000, 499000, 9),
	)
}

func assertOnlyAppStoreLeft(t *testing.T, got Summary, hiddenNetflix GroupKey) {
	t.Helper()
	if keys := keysOf(got.Subscriptions); !slices.Equal(keys, []GroupKey{"umbrella:App Store||d09"}) {
		t.Errorf("subscriptions = %v, want only the App Store line", keys)
	}
	if spend := got.MonthlySpend[shared.CurrencyCLP]; spend != 499000 {
		t.Errorf("monthly spend = %d, want 499000 (without Netflix)", spend)
	}
	if slices.ContainsFunc(got.Rises, func(r Rise) bool { return r.Key == hiddenNetflix }) {
		t.Errorf("rises = %+v, want no rise of the hidden line", got.Rises)
	}
}

func TestHideLines_HiddenLineDoesNotShowNorAdd(t *testing.T) {
	const hiddenNetflix GroupKey = "service:Netflix||d25"
	summary := summaryOfCharges(netflixAndAppStoreRising(25), nil)

	got := HideLines(summary, []GroupKey{hiddenNetflix})

	assertOnlyAppStoreLeft(t, got, hiddenNetflix)
	if got.RemainingThisMonth != 0 {
		t.Errorf("remaining this month = %d, want 0 (Netflix on the 25th is hidden)", got.RemainingThisMonth)
	}
	if !slices.Contains(keysOf(summary.Subscriptions), hiddenNetflix) {
		t.Errorf("the given summary changed")
	}
}

// The November charge is the next statement: it joins the same line, which
// keeps its key and so its mark.
func TestHideLines_HiddenLineStaysHiddenAfterANewCharge(t *testing.T) {
	const hiddenNetflix GroupKey = "service:Netflix||d05"
	newCharge := scenarioCharges("netflix", "NETFLIX", nil, 1349000, shared.CurrencyCLP, 5, time.November)
	summary := summaryOfCharges(slices.Concat(netflixAndAppStoreRising(5), newCharge), nil)
	if summary.ChargedThisMonth != 1 {
		t.Fatalf("before hiding: charged this month = %d, want 1 (the new Netflix charge)", summary.ChargedThisMonth)
	}

	got := HideLines(summary, []GroupKey{hiddenNetflix})

	assertOnlyAppStoreLeft(t, got, hiddenNetflix)
	if got.ChargedThisMonth != 0 {
		t.Errorf("charged this month = %d, want 0 (the November charge is of the hidden line)", got.ChargedThisMonth)
	}
}
