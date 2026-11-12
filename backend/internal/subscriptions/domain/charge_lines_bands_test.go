package domain

import (
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var lineMonths = []time.Month{time.August, time.September, time.October}

func TestChainLines_DifferentAmountsNearTheMonthTurnDoNotMix(t *testing.T) {
	var charges []shared.Charge
	for _, month := range lineMonths {
		charges = append(charges,
			lineCharge("day1-"+month.String(), 1299000, shared.CurrencyCLP, month, 1),
			lineCharge("day28-"+month.String(), 799000, shared.CurrencyCLP, month, 28))
	}

	assertLines(t, charges, [][]string{
		{"day1-August", "day1-September", "day1-October"},
		{"day28-August", "day28-September", "day28-October"},
	}, []string{})
}

func TestChainLines_SuccessiveBandsJoinAsPlanChange(t *testing.T) {
	charges := []shared.Charge{
		lineCharge("aug", 1299000, shared.CurrencyCLP, time.August, 1),
		lineCharge("sep", 1299000, shared.CurrencyCLP, time.September, 1),
		lineCharge("oct", 1799000, shared.CurrencyCLP, time.October, 1),
	}

	assertLines(t, charges, [][]string{{"aug", "sep", "oct"}}, []string{})
}

func TestChainLines_CurrencyChangeIsAPlanChange(t *testing.T) {
	charges := []shared.Charge{
		lineCharge("clp-1", 599000, shared.CurrencyCLP, time.July, 1),
		lineCharge("clp-2", 599000, shared.CurrencyCLP, time.August, 1),
		lineCharge("usd-1", 799, shared.CurrencyUSD, time.September, 1),
		lineCharge("usd-2", 799, shared.CurrencyUSD, time.October, 1),
	}

	assertLines(t, charges, [][]string{{"clp-1", "clp-2", "usd-1", "usd-2"}}, []string{})
}

func TestChainLines_BandsThatOverlapInTimeStayApart(t *testing.T) {
	charges := []shared.Charge{
		lineCharge("mine-aug", 1299000, shared.CurrencyCLP, time.August, 5),
		lineCharge("mine-sep", 1299000, shared.CurrencyCLP, time.September, 5),
		lineCharge("parents-sep", 799000, shared.CurrencyCLP, time.September, 20),
		lineCharge("mine-oct", 1299000, shared.CurrencyCLP, time.October, 5),
		lineCharge("parents-oct", 799000, shared.CurrencyCLP, time.October, 20),
	}

	assertLines(t, charges, [][]string{
		{"mine-aug", "mine-sep", "mine-oct"},
		{"parents-sep", "parents-oct"},
	}, []string{})
}
