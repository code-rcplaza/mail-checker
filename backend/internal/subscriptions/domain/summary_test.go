package domain

import (
	"maps"
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func withKey(group ChargeGroup, key GroupKey) ChargeGroup {
	group.Key = key
	return group
}

var (
	activeStreaming   = withKey(groupOf(true, "STREAMING CO", sameAmounts(599000, 3), 70, 40, 10), "service:Streaming Co")
	inactiveStreaming = withKey(groupOf(true, "VIDEO CO", sameAmounts(799000, 3), 110, 80, 50), "service:Video Co")
	possibleMusic     = withKey(groupOf(false, "MUSIC APP", sameAmounts(1200000, 2), 40, 10), "merchant:MUSIC APP:CLP")
	cornerShop        = withKey(groupOf(false, "CORNER SHOP", sameAmounts(800000, 1), 10), "merchant:CORNER SHOP:CLP")
)

func keysOf(subscriptions []Subscription) []GroupKey {
	keys := make([]GroupKey, 0, len(subscriptions))
	for _, s := range subscriptions {
		keys = append(keys, s.Key)
	}
	return keys
}

func merchantsOf(charges []shared.Charge) []string {
	merchants := make([]string, 0, len(charges))
	for _, c := range charges {
		merchants = append(merchants, c.Merchant)
	}
	return merchants
}

func buildSummaryOf(groups []ChargeGroup, decisions map[GroupKey]Decision) Summary {
	return BuildSummary(groups, decisions, looksLikeSubscriptionFake, detectionNow)
}

func TestBuildSummary_ActiveKnownServiceInSubscriptionsAndInactiveInInactive(t *testing.T) {
	got := buildSummaryOf([]ChargeGroup{activeStreaming, inactiveStreaming}, nil)

	want := []Subscription{{
		Key: activeStreaming.Key, Status: StatusActive, Cadence: CadenceMonthly,
		Amount: 599000, Currency: shared.CurrencyCLP, Charges: activeStreaming.Charges,
		NextCharge: day(2026, 12, 3),
	}}
	if !reflect.DeepEqual(got.Subscriptions, want) {
		t.Errorf("subscriptions = %+v, want %+v", got.Subscriptions, want)
	}
	if keys := keysOf(got.Inactive); !reflect.DeepEqual(keys, []GroupKey{inactiveStreaming.Key}) {
		t.Errorf("inactive = %v, want [%s]", keys, inactiveStreaming.Key)
	}
}

func TestBuildSummary_UndecidedPossibleInPossibleAndNotInSpend(t *testing.T) {
	got := buildSummaryOf([]ChargeGroup{activeStreaming, possibleMusic}, nil)

	if keys := keysOf(got.Possible); !reflect.DeepEqual(keys, []GroupKey{possibleMusic.Key}) {
		t.Errorf("possible = %v, want [%s]", keys, possibleMusic.Key)
	}
	want := map[shared.Currency]shared.Money{shared.CurrencyCLP: 599000}
	if !maps.Equal(got.MonthlySpend, want) {
		t.Errorf("monthly spend = %v, want %v (only the active one)", got.MonthlySpend, want)
	}
}

func TestBuildSummary_ConfirmedPossibleInSubscriptionsWithUserCadenceAndInSpend(t *testing.T) {
	decisions := map[GroupKey]Decision{possibleMusic.Key: {Kind: DecisionConfirmed, Cadence: CadenceYearly}}

	got := buildSummaryOf([]ChargeGroup{possibleMusic}, decisions)

	if len(got.Subscriptions) != 1 || got.Subscriptions[0].Key != possibleMusic.Key {
		t.Fatalf("subscriptions = %v, want [%s]", keysOf(got.Subscriptions), possibleMusic.Key)
	}
	if cadence := got.Subscriptions[0].Cadence; cadence != CadenceYearly {
		t.Errorf("cadence = %q, want the user's %q", cadence, CadenceYearly)
	}
	want := map[shared.Currency]shared.Money{shared.CurrencyCLP: 100000}
	if !maps.Equal(got.MonthlySpend, want) {
		t.Errorf("monthly spend = %v, want %v (12.000 / 12)", got.MonthlySpend, want)
	}
	if len(got.Possible) != 0 {
		t.Errorf("possible = %v, want none", keysOf(got.Possible))
	}
}

func TestBuildSummary_DismissedPossibleChargesInOtherCharges(t *testing.T) {
	decisions := map[GroupKey]Decision{possibleMusic.Key: {Kind: DecisionDismissed}}

	got := buildSummaryOf([]ChargeGroup{possibleMusic}, decisions)

	if merchants := merchantsOf(got.OtherCharges); !reflect.DeepEqual(merchants, []string{"MUSIC APP", "MUSIC APP"}) {
		t.Errorf("other charges = %v, want both MUSIC APP charges", merchants)
	}
	if len(got.Possible) != 0 {
		t.Errorf("possible = %v, want none", keysOf(got.Possible))
	}
}

func TestBuildSummary_NonSubscriptionChargesInOtherCharges(t *testing.T) {
	got := buildSummaryOf([]ChargeGroup{cornerShop}, nil)

	if merchants := merchantsOf(got.OtherCharges); !reflect.DeepEqual(merchants, []string{"CORNER SHOP"}) {
		t.Errorf("other charges = %v, want [CORNER SHOP]", merchants)
	}
}

func TestBuildSummary_MonthlySpendOfActiveSubscriptions(t *testing.T) {
	activeMusic := withKey(groupOf(true, "MUSIC CO", sameAmounts(399000, 3), 70, 40, 10), "service:Music Co")

	got := buildSummaryOf([]ChargeGroup{activeStreaming, activeMusic}, nil)

	want := map[shared.Currency]shared.Money{shared.CurrencyCLP: 998000}
	if !maps.Equal(got.MonthlySpend, want) {
		t.Errorf("monthly spend = %v, want %v", got.MonthlySpend, want)
	}
}
