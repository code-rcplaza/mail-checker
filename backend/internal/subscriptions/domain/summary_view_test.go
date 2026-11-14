package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestParentName(t *testing.T) {
	cases := map[GroupKey]string{
		"service:Netflix":         "Netflix",
		"umbrella:App Store":      "App Store",
		"merchant:OLD APP:CLP":    "OLD APP",
		"merchant:A:B CORP:USD":   "A:B CORP",
		"manual:0193f1a2-0000-70": "",
	}
	for parent, want := range cases {
		if got := ParentName(parent); got != want {
			t.Errorf("ParentName(%q) = %q, want %q", parent, got, want)
		}
	}
}

func TestRisesWithImpact_UsesTheCadenceOfEachLine(t *testing.T) {
	summary := Summary{
		Subscriptions: []Subscription{{Key: "m", Cadence: CadenceMonthly}, {Key: "y", Cadence: CadenceYearly}},
		Rises:         []Rise{{Key: "m", Before: 1000000, After: 1050000, Currency: shared.CurrencyCLP}, {Key: "y", Before: 12000000, After: 13000000, Currency: shared.CurrencyCLP}},
	}

	got := RisesWithImpact(summary)

	if len(got) != 2 || got[0].Impact != (Impact{Monthly: 50000, Yearly: 600000}) || got[1].Impact != (Impact{Monthly: 83333, Yearly: 1000000}) {
		t.Errorf("RisesWithImpact = %+v, want each rise with the impact of its cadence", got)
	}
}

func TestDefaultPeriodMonthsIsValid(t *testing.T) {
	if err := ValidatePeriodMonths(DefaultPeriodMonths); err != nil || DefaultPeriodMonths != 12 {
		t.Errorf("DefaultPeriodMonths = %d (%v), want a valid 12", DefaultPeriodMonths, err)
	}
}
