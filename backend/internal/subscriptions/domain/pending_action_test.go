package domain

import (
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestPendingActions(t *testing.T) {
	summary := Summary{
		Possible: []Subscription{{Key: "service:Music||d03", Status: StatusPossible}, {Key: "merchant:GYM:CLP||d01", Status: StatusPossible}},
		Inactive: []Subscription{{Key: "merchant:OLD APP:CLP||d26", Status: StatusInactive}},
		Rises:    []Rise{{Key: "service:Netflix||d25", Before: 1299000, After: 1349000, Currency: shared.CurrencyCLP}},
	}
	keptNetflix := map[GroupKey][]shared.Money{"service:Netflix||d25": {1349000}}
	possibleActions := []PendingAction{
		{Kind: PendingActionPossible, Key: "service:Music||d03"},
		{Kind: PendingActionPossible, Key: "merchant:GYM:CLP||d01"},
	}
	cases := []struct {
		name    string
		summary Summary
		want    []PendingAction
	}{
		{"current possible and a rise ask for an action each; an inactive line does not", summary, append([]PendingAction{{Kind: PendingActionRise, Key: "service:Netflix||d25"}}, possibleActions...)},
		{"a rise already kept does not ask for an action", HideKeptRises(summary, keptNetflix), possibleActions},
		{"nothing pending has no actions", Summary{}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PendingActions(tc.summary)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("PendingActions = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestPendingActions_Order(t *testing.T) {
	monthly := func(key GroupKey, currency shared.Currency) Subscription {
		return Subscription{Key: key, Status: StatusActive, Cadence: CadenceMonthly, Currency: currency}
	}
	yearly := Subscription{Key: "service:Canva||d01", Status: StatusActive, Cadence: CadenceYearly, Currency: shared.CurrencyCLP}
	rate := &shared.ExchangeRate{PesosPerDollar: 95000}
	lines := []Subscription{monthly("service:Netflix||d25", shared.CurrencyCLP), monthly("service:Claude||d20", shared.CurrencyUSD), yearly}
	// Netflix +CLP 500 a month is CLP 6.000 a year; Canva +CLP 3.000 once a year;
	// Claude +US$ 1 a month is US$ 12 a year, CLP 11.400 at 950.
	rises := []Rise{
		{Key: "service:Canva||d01", Before: 9000000, After: 9300000, Currency: shared.CurrencyCLP},
		{Key: "service:Netflix||d25", Before: 1299000, After: 1349000, Currency: shared.CurrencyCLP},
		{Key: "service:Claude||d20", Before: 2000, After: 2100, Currency: shared.CurrencyUSD},
	}
	possible := []Subscription{{Key: "merchant:GYM:CLP||d01", Status: StatusPossible}}
	riseAction := func(key GroupKey) PendingAction { return PendingAction{Kind: PendingActionRise, Key: key} }
	possibleAction := PendingAction{Kind: PendingActionPossible, Key: "merchant:GYM:CLP||d01"}
	cases := []struct {
		name    string
		summary Summary
		want    []PendingAction
	}{
		{
			"rises go first, the largest yearly impact in pesos first, then the possible",
			Summary{Subscriptions: lines, Possible: possible, Rises: rises, ExchangeRate: rate},
			[]PendingAction{riseAction("service:Claude||d20"), riseAction("service:Netflix||d25"), riseAction("service:Canva||d01"), possibleAction},
		},
		{
			"without an exchange rate a rise in dollars goes after the ones in pesos",
			Summary{Subscriptions: lines, Possible: possible, Rises: rises},
			[]PendingAction{riseAction("service:Netflix||d25"), riseAction("service:Canva||d01"), riseAction("service:Claude||d20"), possibleAction},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PendingActions(tc.summary)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("PendingActions = %+v, want %+v", got, tc.want)
			}
		})
	}
}
