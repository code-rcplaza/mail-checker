package domain

import (
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// lineCharge is a charge of one origin on the given date of 2026.
func lineCharge(externalID string, amount shared.Money, currency shared.Currency, month time.Month, day int) shared.Charge {
	return shared.Charge{
		Bank:       "bci",
		Source:     shared.SourceEmail,
		ExternalID: externalID,
		OccurredAt: time.Date(2026, month, day, 12, 0, 0, 0, time.UTC),
		Amount:     amount,
		Currency:   currency,
		Merchant:   "STREAMING CO",
	}
}

func lineIDs(lines [][]shared.Charge) [][]string {
	ids := make([][]string, 0, len(lines))
	for _, line := range lines {
		ids = append(ids, externalIDs(line))
	}
	return ids
}

func assertLines(t *testing.T, charges []shared.Charge, wantLines [][]string, wantOther []string) {
	t.Helper()
	lines, other := ChainLines(charges)
	if got := lineIDs(lines); !reflect.DeepEqual(got, wantLines) {
		t.Errorf("lines = %v, want %v", got, wantLines)
	}
	if got := externalIDs(other); !reflect.DeepEqual(got, wantOther) {
		t.Errorf("other = %v, want %v", got, wantOther)
	}
}

func TestChainLines_TwoNetflixWithDifferentPricesAndDays(t *testing.T) {
	var charges []shared.Charge
	for _, month := range []time.Month{time.August, time.September, time.October} {
		charges = append(charges,
			lineCharge("mine-"+month.String(), 1299000, shared.CurrencyCLP, month, 25),
			lineCharge("parents-"+month.String(), 799000, shared.CurrencyCLP, month, 9))
	}

	assertLines(t, charges, [][]string{
		{"parents-August", "parents-September", "parents-October"},
		{"mine-August", "mine-September", "mine-October"},
	}, []string{})
}

func TestChainLines_TwoAppleSubscriptionsWithTheSamePrice(t *testing.T) {
	var charges []shared.Charge
	for _, month := range []time.Month{time.August, time.September, time.October} {
		charges = append(charges,
			lineCharge("day5-"+month.String(), 599000, shared.CurrencyCLP, month, 5),
			lineCharge("day9-"+month.String(), 599000, shared.CurrencyCLP, month, 9))
	}

	assertLines(t, charges, [][]string{
		{"day5-August", "day5-September", "day5-October"},
		{"day9-August", "day9-September", "day9-October"},
	}, []string{})
}

func TestChainLines_TwoIdenticalChargesOnTheSameDay(t *testing.T) {
	var charges []shared.Charge
	for _, month := range []time.Month{time.August, time.September, time.October} {
		charges = append(charges,
			lineCharge("first-"+month.String(), 1471500, shared.CurrencyCLP, month, 4),
			lineCharge("second-"+month.String(), 1471500, shared.CurrencyCLP, month, 4))
	}

	assertLines(t, charges, [][]string{
		{"first-August", "first-September", "first-October"},
		{"second-August", "second-September", "second-October"},
	}, []string{})
}

func TestChainLines_PlanChangeOnTheExpectedDateStaysInTheLine(t *testing.T) {
	charges := []shared.Charge{
		lineCharge("aug", 1299000, shared.CurrencyCLP, time.August, 1),
		lineCharge("sep", 1299000, shared.CurrencyCLP, time.September, 1),
		lineCharge("oct", 1799000, shared.CurrencyCLP, time.October, 1),
	}

	assertLines(t, charges, [][]string{{"aug", "sep", "oct"}}, []string{})
}

func TestChainLines_OneOffChargeBetweenMonthlyChargesIsOther(t *testing.T) {
	charges := []shared.Charge{
		lineCharge("aug", 2000, shared.CurrencyUSD, time.August, 1),
		lineCharge("credits", 500, shared.CurrencyUSD, time.August, 10),
		lineCharge("sep", 2000, shared.CurrencyUSD, time.September, 1),
		lineCharge("oct", 2000, shared.CurrencyUSD, time.October, 1),
	}

	assertLines(t, charges, [][]string{{"aug", "sep", "oct"}}, []string{"credits"})
}

func TestChainLines_SingleChargeIsALine(t *testing.T) {
	charges := []shared.Charge{lineCharge("only", 1299000, shared.CurrencyCLP, time.August, 1)}

	assertLines(t, charges, [][]string{{"only"}}, []string{})
}

func TestChainLines_YearlyChargesAreOneLine(t *testing.T) {
	first := lineCharge("2026", 12000, shared.CurrencyUSD, time.January, 10)
	second := first
	second.ExternalID = "2027"
	second.OccurredAt = first.OccurredAt.AddDate(0, 0, 365)

	assertLines(t, []shared.Charge{first, second}, [][]string{{"2026", "2027"}}, []string{})
}
