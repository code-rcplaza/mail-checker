package domain

import (
	"fmt"
	"maps"
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// GroupKey identifies a group of charges and stays stable between
// recalculations, so user decisions can be stored against it.
type GroupKey string

// IdentifyServiceFunc and NormalizeMerchantFunc come from the catalog port;
// domain receives them as functions so it does not import ports.
type (
	IdentifyServiceFunc   func(merchant string) (Service, bool)
	NormalizeMerchantFunc func(merchant string) string
)

type ChargeGroup struct {
	Key     GroupKey
	Charges []shared.Charge // ascending by date
}

// GroupCharges puts together the charges of the same subscription: a known
// service by service, an umbrella service by exact amount and currency, and
// an unknown merchant by normalized name and currency. Groups come out
// ordered by key so the result is deterministic.
func GroupCharges(
	charges []shared.Charge,
	identify IdentifyServiceFunc,
	normalize NormalizeMerchantFunc,
) []ChargeGroup {
	chronological := sortByDate(charges)
	chargesByKey := collectByGroupKey(chronological, identify, normalize)
	return toGroupsSortedByKey(chargesByKey)
}

// sortByDate returns a sorted copy; the caller's slice is left untouched.
func sortByDate(charges []shared.Charge) []shared.Charge {
	sorted := slices.Clone(charges)
	slices.SortStableFunc(sorted, func(a, b shared.Charge) int {
		return a.OccurredAt.Compare(b.OccurredAt)
	})
	return sorted
}

func collectByGroupKey(
	charges []shared.Charge,
	identify IdentifyServiceFunc,
	normalize NormalizeMerchantFunc,
) map[GroupKey][]shared.Charge {
	chargesByKey := make(map[GroupKey][]shared.Charge)
	for _, c := range charges {
		key := groupKeyOf(c, identify, normalize)
		chargesByKey[key] = append(chargesByKey[key], c)
	}
	return chargesByKey
}

func groupKeyOf(
	c shared.Charge,
	identify IdentifyServiceFunc,
	normalize NormalizeMerchantFunc,
) GroupKey {
	service, isKnown := identify(c.Merchant)
	if !isKnown {
		return merchantKey(normalize(c.Merchant), c.Currency)
	}
	if service.IsUmbrella {
		return umbrellaKey(service, c.Amount, c.Currency)
	}
	return serviceKey(service)
}

func serviceKey(service Service) GroupKey {
	return GroupKey("service:" + service.Name)
}

func umbrellaKey(service Service, amount shared.Money, currency shared.Currency) GroupKey {
	return GroupKey(fmt.Sprintf("umbrella:%s:%d:%s", service.Name, amount, currency))
}

func merchantKey(normalizedMerchant string, currency shared.Currency) GroupKey {
	return GroupKey(fmt.Sprintf("merchant:%s:%s", normalizedMerchant, currency))
}

func toGroupsSortedByKey(chargesByKey map[GroupKey][]shared.Charge) []ChargeGroup {
	groups := make([]ChargeGroup, 0, len(chargesByKey))
	for _, key := range slices.Sorted(maps.Keys(chargesByKey)) {
		groups = append(groups, ChargeGroup{Key: key, Charges: chargesByKey[key]})
	}
	return groups
}
