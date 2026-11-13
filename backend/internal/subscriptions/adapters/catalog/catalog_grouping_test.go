package catalog

import (
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.ServiceCatalog = Catalog{}

func statementCharge(merchant string, amount shared.Money, month time.Month, day int) shared.Charge {
	return shared.Charge{
		Bank:       "bci",
		Source:     shared.SourceStatement,
		ExternalID: merchant + "-" + month.String(),
		OccurredAt: time.Date(2026, month, day, 12, 0, 0, 0, time.UTC),
		Amount:     amount,
		Currency:   shared.CurrencyCLP,
		Merchant:   merchant,
	}
}

func monthlyStatementCharges(merchant string, amount shared.Money, day int) []shared.Charge {
	return []shared.Charge{
		statementCharge(merchant, amount, time.August, day),
		statementCharge(merchant, amount, time.September, day),
		statementCharge(merchant, amount, time.October, day),
	}
}

func groupWithCatalog(charges []shared.Charge) []domain.ChargeGroup {
	catalog := Catalog{}
	return domain.GroupCharges(charges, catalog.Identify, catalog.Normalize)
}

func lineCountByParent(groups []domain.ChargeGroup) map[domain.GroupKey]int {
	counts := make(map[domain.GroupKey]int)
	for _, group := range groups {
		if !group.IsOneOff {
			counts[group.Parent]++
		}
	}
	return counts
}

func TestCatalogGroupsRealDescriptorsIntoLines(t *testing.T) {
	cases := []struct {
		name    string
		charges []shared.Charge
		want    map[domain.GroupKey]int
	}{
		{
			"Netflix descriptors share one line",
			[]shared.Charge{
				statementCharge("NETFLIX.COM 866-579-7172", 1299000, time.August, 5),
				statementCharge("NETFLIX LAS CONDES", 1299000, time.September, 5),
				statementCharge("NETFLIX", 1299000, time.October, 5),
			},
			map[domain.GroupKey]int{"service:Netflix": 1},
		},
		{
			"two Google amounts are two lines of the same umbrella",
			slices.Concat(
				monthlyStatementCharges("DL*GOOGLE YOUTUBE", 2170000, 3),
				monthlyStatementCharges("DL*GOOGLE YOUTUBE", 1100000, 17),
			),
			map[domain.GroupKey]int{"umbrella:Google": 2},
		},
		{
			"processor prefix does not change the unknown merchant",
			monthlyStatementCharges("PAYU *UBER TRIP", 450000, 10),
			map[domain.GroupKey]int{"merchant:UBER TRIP:CLP": 1},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := lineCountByParent(groupWithCatalog(tc.charges))

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("lines by parent = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCatalogKeepsSingleKnownChargeAsALine(t *testing.T) {
	charges := []shared.Charge{statementCharge("MICROSOFT. SUBSCRIPTION", 899000, time.October, 12)}

	got := groupWithCatalog(charges)

	if len(got) != 1 {
		t.Fatalf("groups = %d, want 1", len(got))
	}
	group := got[0]
	if group.Parent != "umbrella:Microsoft" || group.IsOneOff || !group.IsKnownService {
		t.Errorf("group = {Parent: %s, IsOneOff: %v, IsKnownService: %v}, want a known line of umbrella:Microsoft",
			group.Parent, group.IsOneOff, group.IsKnownService)
	}
}
