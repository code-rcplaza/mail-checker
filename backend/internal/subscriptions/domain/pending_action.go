package domain

import (
	"cmp"
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// PendingActionKind is what the user has to resolve on a line.
type PendingActionKind string

const (
	PendingActionRise     PendingActionKind = "rise"     // keep the subscription after a rise
	PendingActionPossible PendingActionKind = "possible" // confirm or dismiss a possible subscription
)

type PendingAction struct {
	Kind PendingActionKind
	Key  GroupKey
}

// PendingActions lists what the user still has to resolve: first the rises
// not yet kept, the largest yearly impact in pesos first, then the possible
// subscriptions not yet decided. Kept rises and decided or inactive lines are
// already out of those sections of the summary.
func PendingActions(summary Summary) []PendingAction {
	var actions []PendingAction
	for _, rise := range sortRisesByYearlyImpact(summary) {
		actions = append(actions, PendingAction{Kind: PendingActionRise, Key: rise.Key})
	}
	for _, possible := range summary.Possible {
		actions = append(actions, PendingAction{Kind: PendingActionPossible, Key: possible.Key})
	}
	return actions
}

// rankedRise is a rise and what it costs a year in pesos; a rise in dollars
// without an exchange rate has no known cost and goes after the known ones.
type rankedRise struct {
	rise          Rise
	yearlyInPesos shared.Money
	isKnown       bool
}

func sortRisesByYearlyImpact(summary Summary) []Rise {
	ranked := make([]rankedRise, 0, len(summary.Rises))
	for _, rise := range summary.Rises {
		ranked = append(ranked, rankRise(rise, summary))
	}
	slices.SortStableFunc(ranked, compareByYearlyImpact)
	rises := make([]Rise, 0, len(ranked))
	for _, r := range ranked {
		rises = append(rises, r.rise)
	}
	return rises
}

func rankRise(rise Rise, summary Summary) rankedRise {
	yearly := RiseImpact(rise, cadenceOfLine(summary.Subscriptions, rise.Key)).Yearly
	switch {
	case rise.Currency == shared.CurrencyCLP:
		return rankedRise{rise: rise, yearlyInPesos: yearly, isKnown: true}
	case rise.Currency == shared.CurrencyUSD && summary.ExchangeRate != nil:
		return rankedRise{rise: rise, yearlyInPesos: dollarsToPesos(yearly, *summary.ExchangeRate), isKnown: true}
	}
	return rankedRise{rise: rise}
}

// compareByYearlyImpact puts known costs first, the largest first.
func compareByYearlyImpact(a, b rankedRise) int {
	return cmp.Or(compareKnownFirst(a, b), cmp.Compare(b.yearlyInPesos, a.yearlyInPesos))
}

func compareKnownFirst(a, b rankedRise) int {
	if a.isKnown == b.isKnown {
		return 0
	}
	if a.isKnown {
		return -1
	}
	return 1
}

// cadenceOfLine is CadenceNone for a key that is not among the lines.
func cadenceOfLine(lines []Subscription, key GroupKey) Cadence {
	for _, line := range lines {
		if line.Key == key {
			return line.Cadence
		}
	}
	return CadenceNone
}
