package domain

import (
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func monthlyCharges(prefix, merchant string, amount shared.Money, currency shared.Currency, day int, months ...time.Month) []shared.Charge {
	charges := make([]shared.Charge, 0, len(months))
	for _, month := range months {
		c := lineCharge(prefix+"-"+month.String(), amount, currency, month, day)
		c.Merchant = merchant
		c.AccountID = &accountA
		charges = append(charges, c)
	}
	return charges
}

func groupsByParent(groups []ChargeGroup) map[GroupKey][]GroupKey {
	byParent := make(map[GroupKey][]GroupKey)
	for _, g := range groups {
		byParent[g.Parent] = append(byParent[g.Parent], g.Key)
	}
	return byParent
}

func TestGroupCharges_UmbrellaParentHasOneLinePerSubscription(t *testing.T) {
	charges := slices.Concat(
		monthlyCharges("small", "APP STORE", 469000, shared.CurrencyCLP, 1, time.August, time.September, time.October),
		monthlyCharges("large", "APP STORE", 599000, shared.CurrencyCLP, 9, time.August, time.September, time.October),
	)

	got := groupsByParent(GroupCharges(charges, identifyFake, normalizeFake))

	want := map[GroupKey][]GroupKey{"umbrella:App Store": {"umbrella:App Store||d01", "umbrella:App Store||d09"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("groups by parent = %v, want %v", got, want)
	}
}

func TestGroupCharges_TwoSubscriptionsOfAServiceShareTheParent(t *testing.T) {
	charges := slices.Concat(
		monthlyCharges("mine", "STREAMING CO", 1299000, shared.CurrencyCLP, 25, time.August, time.September, time.October),
		monthlyCharges("parents", "STREAMING CO", 799000, shared.CurrencyCLP, 9, time.August, time.September, time.October),
	)

	got := GroupCharges(charges, identifyFake, normalizeFake)

	want := map[GroupKey][]GroupKey{"service:Streaming Co": {"service:Streaming Co||d09", "service:Streaming Co||d25"}}
	if byParent := groupsByParent(got); !reflect.DeepEqual(byParent, want) {
		t.Errorf("groups by parent = %v, want %v", byParent, want)
	}
	for _, g := range got {
		if !g.IsKnownService || len(g.Charges) != 3 {
			t.Errorf("group %s: known = %v, charges = %d; want known with 3 charges", g.Key, g.IsKnownService, len(g.Charges))
		}
	}
}

func TestGroupCharges_OneOffChargeGoesToItsParentsOtherCharges(t *testing.T) {
	credits := lineCharge("credits", 500, shared.CurrencyUSD, time.August, 10)
	credits.Merchant = "STREAMING CO"
	credits.AccountID = &accountA
	charges := append(monthlyCharges("plan", "STREAMING CO", 2000, shared.CurrencyUSD, 1, time.August, time.September, time.October), credits)

	got := externalIDsByKey(GroupCharges(charges, identifyFake, normalizeFake))

	want := map[GroupKey][]string{
		"service:Streaming Co||d01":    {"plan-August", "plan-September", "plan-October"},
		"service:Streaming Co|one-off": {"credits"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
}

func TestGroupCharges_LinesSortedByKeyAndInputUntouched(t *testing.T) {
	charges := slices.Concat(
		monthlyCharges("music", "MUSIC APP", 399000, shared.CurrencyCLP, 3, time.September, time.August),
		monthlyCharges("stream", "STREAMING CO", 599000, shared.CurrencyCLP, 5, time.September, time.August),
	)
	before := slices.Clone(charges)

	groups := GroupCharges(charges, identifyFake, normalizeFake)

	keys := make([]GroupKey, 0, len(groups))
	for _, g := range groups {
		keys = append(keys, g.Key)
	}
	if want := []GroupKey{"merchant:MUSIC APP:CLP||d03", "service:Streaming Co||d05"}; !reflect.DeepEqual(keys, want) {
		t.Errorf("keys = %v, want %v", keys, want)
	}
	if !reflect.DeepEqual(charges, before) {
		t.Errorf("GroupCharges changed its input")
	}
}
