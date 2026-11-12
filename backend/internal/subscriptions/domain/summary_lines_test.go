package domain

import (
	"maps"
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Scenarios from the audit with real statements; the charges are synthetic.

var (
	netflix = Service{Name: "Netflix"}
	claude  = Service{Name: "Claude"}
	apple   = Service{Name: "App Store", IsUmbrella: true}
	google  = Service{Name: "Google", IsUmbrella: true}
)

func identifyScenarioFake(merchant string) (Service, bool) {
	switch merchant {
	case "NETFLIX":
		return netflix, true
	case "ANTHROPIC CLAUDE":
		return claude, true
	case "APPLE COM BILL":
		return apple, true
	case "GOOGLE YOUTUBE":
		return google, true
	}
	return Service{}, false
}

// scenarioCharges bills amount on the same day of each month; without an
// account the charge comes from a statement.
func scenarioCharges(prefix, merchant string, account *int64, amount shared.Money, currency shared.Currency, day int, months ...time.Month) []shared.Charge {
	charges := make([]shared.Charge, 0, len(months))
	for _, month := range months {
		c := lineCharge(prefix+"-"+month.String(), amount, currency, month, day)
		c.Merchant = merchant
		c.AccountID = account
		if account == nil {
			c.Source = shared.SourceStatement
		}
		charges = append(charges, c)
	}
	return charges
}

func summaryOfCharges(charges []shared.Charge, decisions map[GroupKey]Decision) Summary {
	groups := GroupCharges(charges, identifyScenarioFake, normalizeFake)
	return BuildSummary(groups, decisions, looksLikeSubscriptionFake, detectionNow)
}

func assertTwoActiveOf(t *testing.T, got Summary, parent GroupKey, spend shared.Money) {
	t.Helper()
	if len(got.Subscriptions) != 2 || len(got.Possible) != 0 {
		t.Fatalf("subscriptions = %v, possible = %v; want two subscriptions and no possible", keysOf(got.Subscriptions), keysOf(got.Possible))
	}
	for _, s := range got.Subscriptions {
		if s.Status != StatusActive || s.Parent != parent {
			t.Errorf("%s: status = %q, parent = %q; want active under %q", s.Key, s.Status, s.Parent, parent)
		}
	}
	if want := map[shared.Currency]shared.Money{shared.CurrencyCLP: spend}; !maps.Equal(got.MonthlySpend, want) {
		t.Errorf("monthly spend = %v, want %v", got.MonthlySpend, want)
	}
}

var lastMonths = []time.Month{time.August, time.September, time.October}

func TestBuildSummary_TwoNetflixAreTwoActiveSubscriptions(t *testing.T) {
	cases := []struct {
		name           string
		second         *int64
		first, another shared.Money
	}{
		{"same price, same account", &accountA, 1299000, 1299000},
		{"same price, different account", &accountB, 1299000, 1299000},
		{"different price, same account", &accountA, 1299000, 799000},
		{"different price, different account", &accountB, 1299000, 799000},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			charges := slices.Concat(
				scenarioCharges("mine", "NETFLIX", &accountA, tc.first, shared.CurrencyCLP, 1, lastMonths...),
				scenarioCharges("parents", "NETFLIX", tc.second, tc.another, shared.CurrencyCLP, 3, lastMonths...),
			)

			got := summaryOfCharges(charges, nil)

			assertTwoActiveOf(t, got, "service:Netflix", tc.first+tc.another)
		})
	}
}

func TestBuildSummary_TwoAppleOfTheSamePriceAreTwoActiveSubscriptions(t *testing.T) {
	charges := slices.Concat(
		scenarioCharges("music", "APPLE COM BILL", nil, 599000, shared.CurrencyCLP, 5, lastMonths...),
		scenarioCharges("cloud", "APPLE COM BILL", nil, 599000, shared.CurrencyCLP, 9, lastMonths...),
	)

	got := summaryOfCharges(charges, nil)

	assertTwoActiveOf(t, got, "umbrella:App Store", 1198000)
}

func TestBuildSummary_TwoPoliciesOnTheSameDayAreTwoPossible(t *testing.T) {
	charges := slices.Concat(
		scenarioCharges("home", "SEGUROS CLC", nil, 1850000, shared.CurrencyCLP, 10, lastMonths...),
		scenarioCharges("life", "SEGUROS CLC", nil, 3200000, shared.CurrencyCLP, 10, lastMonths...),
	)

	got := summaryOfCharges(charges, nil)

	want := []GroupKey{"merchant:SEGUROS CLC:CLP||d10", "merchant:SEGUROS CLC:CLP||d10#2"}
	if keys := keysOf(got.Possible); !reflect.DeepEqual(keys, want) {
		t.Errorf("possible = %v, want %v", keys, want)
	}
	if len(got.OtherCharges) != 0 {
		t.Errorf("other charges = %v, want none", merchantsOf(got.OtherCharges))
	}
}

func TestBuildSummary_GoogleWithTwoAmountsIsTwoActiveSubscriptions(t *testing.T) {
	charges := slices.Concat(
		scenarioCharges("family", "GOOGLE YOUTUBE", nil, 2170000, shared.CurrencyCLP, 12, lastMonths...),
		scenarioCharges("storage", "GOOGLE YOUTUBE", nil, 1100000, shared.CurrencyCLP, 20, lastMonths...),
	)

	got := summaryOfCharges(charges, nil)

	assertTwoActiveOf(t, got, "umbrella:Google", 3270000)
}

func TestBuildSummary_NetflixPlanChangeOnTheSameAccountIsOneSubscription(t *testing.T) {
	charges := slices.Concat(
		scenarioCharges("basic", "NETFLIX", &accountA, 1299000, shared.CurrencyCLP, 1, time.June, time.July, time.August),
		scenarioCharges("premium", "NETFLIX", &accountA, 1799000, shared.CurrencyCLP, 1, time.September, time.October),
	)

	got := summaryOfCharges(charges, nil)

	if len(got.Subscriptions) != 1 {
		t.Fatalf("subscriptions = %v, want one", keysOf(got.Subscriptions))
	}
	s := got.Subscriptions[0]
	if s.Status != StatusActive || s.Amount != 1799000 {
		t.Errorf("status = %q, amount = %d; want active with 1799000", s.Status, s.Amount)
	}
	if changes := AnalyzeKnownServiceAmount(s.Charges).PlanChanges; len(changes) != 1 {
		t.Errorf("plan changes = %+v, want one", changes)
	}
}

func TestBuildSummary_OneOffCreditOfAServiceGoesToOtherCharges(t *testing.T) {
	credit := scenarioCharges("credit", "ANTHROPIC CLAUDE", &accountA, 500, shared.CurrencyUSD, 15, time.August)
	charges := append(scenarioCharges("plan", "ANTHROPIC CLAUDE", &accountA, 2000, shared.CurrencyUSD, 1, lastMonths...), credit...)

	got := summaryOfCharges(charges, nil)

	if len(got.Subscriptions) != 1 || got.Subscriptions[0].Amount != 2000 {
		t.Fatalf("subscriptions = %+v, want one active line of USD 20", got.Subscriptions)
	}
	if len(got.Possible) != 0 {
		t.Errorf("possible = %v, want none", keysOf(got.Possible))
	}
	if ids := externalIDs(got.OtherCharges); !reflect.DeepEqual(ids, []string{"credit-August"}) {
		t.Errorf("other charges = %v, want the USD 5 credit", ids)
	}
}

func TestBuildSummary_SubscriptionCarriesParentAndLineKey(t *testing.T) {
	charges := scenarioCharges("mine", "NETFLIX", &accountA, 1299000, shared.CurrencyCLP, 1, lastMonths...)

	got := summaryOfCharges(charges, nil)

	if len(got.Subscriptions) != 1 {
		t.Fatalf("subscriptions = %v, want one", keysOf(got.Subscriptions))
	}
	s := got.Subscriptions[0]
	if s.Key != "service:Netflix||d01" || s.Parent != "service:Netflix" {
		t.Errorf("key = %q, parent = %q; want service:Netflix||d01 under service:Netflix", s.Key, s.Parent)
	}
}

func TestBuildSummary_DecisionOnALineAppliesOnlyToThatLine(t *testing.T) {
	charges := slices.Concat(
		scenarioCharges("home", "SEGUROS CLC", nil, 1850000, shared.CurrencyCLP, 10, lastMonths...),
		scenarioCharges("life", "SEGUROS CLC", nil, 3200000, shared.CurrencyCLP, 10, lastMonths...),
	)
	decisions := map[GroupKey]Decision{"merchant:SEGUROS CLC:CLP||d10": {Kind: DecisionConfirmed, Cadence: CadenceMonthly}}

	got := summaryOfCharges(charges, decisions)

	if keys := keysOf(got.Subscriptions); !reflect.DeepEqual(keys, []GroupKey{"merchant:SEGUROS CLC:CLP||d10"}) {
		t.Errorf("subscriptions = %v, want only the confirmed line", keys)
	}
	if keys := keysOf(got.Possible); !reflect.DeepEqual(keys, []GroupKey{"merchant:SEGUROS CLC:CLP||d10#2"}) {
		t.Errorf("possible = %v, want the other line still undecided", keys)
	}
}

func TestIsCurrentPossible_OneOffGroupIsNotAPossibleLine(t *testing.T) {
	credit := scenarioCharges("credit", "ANTHROPIC CLAUDE", &accountA, 500, shared.CurrencyUSD, 15, time.August)
	charges := append(scenarioCharges("plan", "ANTHROPIC CLAUDE", &accountA, 2000, shared.CurrencyUSD, 1, lastMonths...), credit...)
	groups := GroupCharges(charges, identifyScenarioFake, normalizeFake)

	if IsCurrentPossible(groups, "service:Claude|one-off", looksLikeSubscriptionFake, detectionNow) {
		t.Errorf("one-off group accepted as a possible subscription, want rejected like in the summary")
	}
}
