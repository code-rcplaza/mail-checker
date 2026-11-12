package domain

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var (
	streamingCo = Service{Name: "Streaming Co"}
	appStore    = Service{Name: "App Store", IsUmbrella: true}
)

// identifyFake stands in for the catalog: two descriptors map to the same
// known service and one to an umbrella service.
func identifyFake(merchant string) (Service, bool) {
	switch merchant {
	case "STREAMING CO", "STRMNG*CO":
		return streamingCo, true
	case "APP STORE":
		return appStore, true
	}
	return Service{}, false
}

// normalizeFake drops digits and surrounding spaces, as a stand-in for the catalog.
func normalizeFake(merchant string) string {
	withoutDigits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return -1
		}
		return r
	}, merchant)
	return strings.TrimSpace(withoutDigits)
}

func charge(externalID, merchant string, amount shared.Money, currency shared.Currency, day int) shared.Charge {
	return shared.Charge{
		Bank:       "bci",
		Source:     shared.SourceEmail,
		ExternalID: externalID,
		OccurredAt: time.Date(2026, 1, day, 12, 0, 0, 0, time.UTC),
		Amount:     amount,
		Currency:   currency,
		Merchant:   merchant,
		Card:       "credito",
	}
}

// externalIDsByKey reduces groups to their keys and charge ids, in order.
func externalIDsByKey(groups []ChargeGroup) map[GroupKey][]string {
	result := make(map[GroupKey][]string, len(groups))
	for _, group := range groups {
		for _, c := range group.Charges {
			result[group.Key] = append(result[group.Key], c.ExternalID)
		}
	}
	return result
}

func TestGroupCharges(t *testing.T) {
	cases := []struct {
		name    string
		charges []shared.Charge
		want    map[GroupKey][]string
	}{
		{
			name: "known service with different descriptors and currencies forms one group",
			charges: []shared.Charge{
				charge("a", "STREAMING CO", 599000, shared.CurrencyCLP, 1),
				charge("b", "STRMNG*CO", 799, shared.CurrencyUSD, 2),
			},
			want: map[GroupKey][]string{"service:Streaming Co": {"a", "b"}},
		},
		{
			name: "umbrella service splits by exact amount",
			charges: []shared.Charge{
				charge("a", "APP STORE", 299000, shared.CurrencyCLP, 1),
				charge("b", "APP STORE", 999000, shared.CurrencyCLP, 2),
			},
			want: map[GroupKey][]string{
				"umbrella:App Store:299000:CLP": {"a"},
				"umbrella:App Store:999000:CLP": {"b"},
			},
		},
		{
			name: "umbrella service with same amount splits by currency",
			charges: []shared.Charge{
				charge("a", "APP STORE", 299000, shared.CurrencyCLP, 1),
				charge("b", "APP STORE", 299000, shared.CurrencyUSD, 2),
			},
			want: map[GroupKey][]string{
				"umbrella:App Store:299000:CLP": {"a"},
				"umbrella:App Store:299000:USD": {"b"},
			},
		},
		{
			name: "unknown merchants join by normalized name and split by currency",
			charges: []shared.Charge{
				charge("a", "MUSIC APP 123", 399000, shared.CurrencyCLP, 1),
				charge("b", "MUSIC APP 456", 399000, shared.CurrencyCLP, 2),
				charge("c", "MUSIC APP 789", 499, shared.CurrencyUSD, 3),
			},
			want: map[GroupKey][]string{
				"merchant:MUSIC APP:CLP": {"a", "b"},
				"merchant:MUSIC APP:USD": {"c"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := externalIDsByKey(GroupCharges(tc.charges, identifyFake, normalizeFake))

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("groups = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGroupCharges_SortsEachGroupByDateAndKeepsInputUntouched(t *testing.T) {
	charges := []shared.Charge{
		charge("march", "STREAMING CO", 599000, shared.CurrencyCLP, 60),
		charge("january", "STRMNG*CO", 599000, shared.CurrencyCLP, 1),
		charge("february", "STREAMING CO", 599000, shared.CurrencyCLP, 31),
	}
	original := slices.Clone(charges)

	got := externalIDsByKey(GroupCharges(charges, identifyFake, normalizeFake))

	want := map[GroupKey][]string{"service:Streaming Co": {"january", "february", "march"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
	if !reflect.DeepEqual(charges, original) {
		t.Errorf("input was modified: %v, want %v", charges, original)
	}
}
