package domain

import (
	"reflect"
	"slices"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func activeInPesos(key, parent GroupKey, cadence Cadence, pesos shared.Money) Subscription {
	return Subscription{
		Key: key, Parent: parent, Status: StatusActive, Cadence: cadence,
		Amount: pesos, Currency: shared.CurrencyCLP, InPesos: &PesosAmount{Amount: pesos},
	}
}

func TestSpendByService(t *testing.T) {
	dollarsWithoutRate := Subscription{Key: "service:Claude||d20", Parent: "service:Claude", Status: StatusActive, Cadence: CadenceMonthly, Amount: 2000, Currency: shared.CurrencyUSD}
	cases := []struct {
		name  string
		lines []Subscription
		want  []SpendShare
	}{
		{
			"the lines of one service add up in one share",
			[]Subscription{
				activeInPesos("service:Netflix||d05", "service:Netflix", CadenceMonthly, 1000000),
				activeInPesos("service:Spotify||d10", "service:Spotify", CadenceMonthly, 500000),
				activeInPesos("service:Netflix||d20", "service:Netflix", CadenceMonthly, 500000),
			},
			[]SpendShare{{Parent: "service:Netflix", Amount: 1500000, Percent: 75}, {Parent: "service:Spotify", Amount: 500000, Percent: 25}},
		},
		{
			"the percents add up to 100 and the rest of the rounding goes to the largest",
			[]Subscription{
				activeInPesos("service:Spotify||d10", "service:Spotify", CadenceMonthly, 200000),
				activeInPesos("service:Netflix||d05", "service:Netflix", CadenceMonthly, 400000),
				activeInPesos("merchant:GYM:CLP||d01", "merchant:GYM:CLP", CadenceMonthly, 100000),
			},
			[]SpendShare{{Parent: "service:Netflix", Amount: 400000, Percent: 58}, {Parent: "service:Spotify", Amount: 200000, Percent: 28}, {Parent: "merchant:GYM:CLP", Amount: 100000, Percent: 14}},
		},
		{
			"a yearly line counts its monthly twelfth and lines without pesos are left out",
			[]Subscription{
				activeInPesos("service:Netflix||d05", "service:Netflix", CadenceMonthly, 100000),
				activeInPesos("service:Canva||d01", "service:Canva", CadenceYearly, 1200000),
				dollarsWithoutRate,
			},
			[]SpendShare{{Parent: "service:Canva", Amount: 100000, Percent: 50}, {Parent: "service:Netflix", Amount: 100000, Percent: 50}},
		},
		{"a zero spend has no shares", []Subscription{activeInPesos("service:Free||d01", "service:Free", CadenceMonthly, 0)}, nil},
		{"no subscriptions have no shares", nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SpendByService(tc.lines)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("SpendByService = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func sharesOf(percents ...int) []SpendShare {
	shares := make([]SpendShare, 0, len(percents))
	for i, percent := range percents {
		parent := GroupKey("service:S" + string(rune('A'+i)))
		shares = append(shares, SpendShare{Parent: parent, Amount: shared.Money(percent * 1000), Percent: percent})
	}
	return shares
}

func TestTopSpendShares(t *testing.T) {
	eight := sharesOf(30, 20, 15, 10, 10, 6, 5, 4)
	cases := []struct {
		name   string
		shares []SpendShare
		want   []SpendShare
	}{
		{"eight services are the five largest and Others", eight, append(sharesOf(30, 20, 15, 10, 10), SpendShare{Amount: 15000, Percent: 15, IsOthers: true})},
		{"four services are all shown without Others", sharesOf(40, 30, 20, 10), sharesOf(40, 30, 20, 10)},
		{"five services are all shown without Others", sharesOf(40, 20, 20, 10, 10), sharesOf(40, 20, 20, 10, 10)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := slices.Clone(tc.shares)

			got := TopSpendShares(tc.shares)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("TopSpendShares = %+v, want %+v", got, tc.want)
			}
			if !reflect.DeepEqual(tc.shares, original) {
				t.Errorf("the given shares changed")
			}
		})
	}
}
