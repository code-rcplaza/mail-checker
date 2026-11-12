package domain

import (
	"fmt"
	"slices"
	"strings"

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
	Key            GroupKey        // the line, or the parent's one-off charges
	Parent         GroupKey        // the service, umbrella or merchant the line belongs to
	IsKnownService bool            // known or umbrella service from the catalog
	IsOneOff       bool            // charges of the parent that belong to no line
	Charges        []shared.Charge // ascending by date
}

// oneOffSuffix marks the group holding a parent's charges that belong to no line.
const oneOffSuffix = "|one-off"

// GroupCharges puts each subscription in its own group, a line. Charges are
// first gathered by parent (a known or umbrella service by service, an
// unknown merchant by normalized name and currency), then split by origin and
// chained into lines; the one-off charges of a parent get a group of their
// own. Groups come out ordered by key so the result is deterministic.
func GroupCharges(
	charges []shared.Charge,
	identify IdentifyServiceFunc,
	normalize NormalizeMerchantFunc,
) []ChargeGroup {
	var groups []ChargeGroup
	for parent, parentCharges := range collectByParent(sortByDate(charges), identify, normalize) {
		_, isKnownService := identify(parentCharges[0].Merchant)
		groups = append(groups, groupsOfParent(parent, parentCharges, isKnownService)...)
	}
	return sortGroupsByKey(groups)
}

func collectByParent(
	chronological []shared.Charge,
	identify IdentifyServiceFunc,
	normalize NormalizeMerchantFunc,
) map[GroupKey][]shared.Charge {
	byParent := make(map[GroupKey][]shared.Charge)
	for _, c := range chronological {
		parent := groupKeyOf(c, identify, normalize)
		byParent[parent] = append(byParent[parent], c)
	}
	return byParent
}

func groupsOfParent(parent GroupKey, charges []shared.Charge, isKnownService bool) []ChargeGroup {
	lines, oneOff := linesOfParent(charges)
	groups := make([]ChargeGroup, 0, len(lines)+1)
	for i, key := range LineKeys(parent, lines) {
		groups = append(groups, ChargeGroup{Key: key, Parent: parent, IsKnownService: isKnownService, Charges: lines[i]})
	}
	if len(oneOff) > 0 {
		groups = append(groups, ChargeGroup{Key: parent + oneOffSuffix, Parent: parent, IsKnownService: isKnownService, IsOneOff: true, Charges: sortByDate(oneOff)})
	}
	return groups
}

func linesOfParent(charges []shared.Charge) (lines [][]shared.Charge, oneOff []shared.Charge) {
	for _, origin := range SplitByOrigin(charges) {
		originLines, originOneOff := ChainLines(origin)
		lines = append(lines, originLines...)
		oneOff = append(oneOff, originOneOff...)
	}
	return lines, oneOff
}

// sortByDate returns a sorted copy; the caller's slice is left untouched.
func sortByDate(charges []shared.Charge) []shared.Charge {
	sorted := slices.Clone(charges)
	slices.SortStableFunc(sorted, func(a, b shared.Charge) int {
		return a.OccurredAt.Compare(b.OccurredAt)
	})
	return sorted
}

// groupKeyOf gives the parent of a charge; ChargesMatch compares charges by it.
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
		return umbrellaKey(service)
	}
	return serviceKey(service)
}

func serviceKey(service Service) GroupKey {
	return GroupKey("service:" + service.Name)
}

func umbrellaKey(service Service) GroupKey {
	return GroupKey("umbrella:" + service.Name)
}

func merchantKey(normalizedMerchant string, currency shared.Currency) GroupKey {
	return GroupKey(fmt.Sprintf("merchant:%s:%s", normalizedMerchant, currency))
}

func sortGroupsByKey(groups []ChargeGroup) []ChargeGroup {
	sorted := slices.Clone(groups)
	slices.SortFunc(sorted, func(a, b ChargeGroup) int { return strings.Compare(string(a.Key), string(b.Key)) })
	return sorted
}
