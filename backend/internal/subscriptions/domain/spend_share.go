package domain

import (
	"cmp"
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	fullPercent = 100
	// SpendRingSegments is how many services the ring shows on their own;
	// the rest go together in "Others" so it stays readable.
	SpendRingSegments = 5
)

// SpendShare is what one service (the parent of its lines) takes of the
// monthly spend in pesos.
type SpendShare struct {
	Parent  GroupKey
	Amount  shared.Money // monthly, in pesos
	Percent int
	// IsOthers marks the segment that joins the services past the largest
	// ones; it has no parent.
	IsOthers bool
}

// TopSpendShares keeps the largest services and joins the rest in one
// "Others" segment; it expects the shares sorted largest first, as
// SpendByService gives them.
func TopSpendShares(shares []SpendShare) []SpendShare {
	if len(shares) <= SpendRingSegments {
		return slices.Clone(shares)
	}
	top := slices.Clone(shares[:SpendRingSegments])
	return append(top, othersShare(shares[SpendRingSegments:]))
}

func othersShare(rest []SpendShare) SpendShare {
	others := SpendShare{IsOthers: true}
	for _, share := range rest {
		others.Amount += share.Amount
		others.Percent += share.Percent
	}
	return others
}

// SpendByService splits the monthly spend in pesos by service, largest first.
// Lines in dollars without an exchange rate are left out, like in the total.
func SpendByService(subscriptions []Subscription) []SpendShare {
	shares := sortLargestFirst(monthlyPesosByParent(subscriptions))
	total := totalOfShares(shares)
	if total == 0 {
		return nil
	}
	return withPercents(shares, total)
}

func monthlyPesosByParent(subscriptions []Subscription) []SpendShare {
	byParent := make(map[GroupKey]shared.Money)
	for _, s := range subscriptions {
		cadence, isCounted := spendCadence(s)
		if !isCounted || s.InPesos == nil {
			continue
		}
		byParent[s.Parent] += monthlyAmount(s.InPesos.Amount, cadence)
	}
	shares := make([]SpendShare, 0, len(byParent))
	for parent, amount := range byParent {
		shares = append(shares, SpendShare{Parent: parent, Amount: amount})
	}
	return shares
}

// sortLargestFirst breaks ties by parent so the order does not change between reads.
func sortLargestFirst(shares []SpendShare) []SpendShare {
	sorted := slices.Clone(shares)
	slices.SortFunc(sorted, func(a, b SpendShare) int {
		return cmp.Or(cmp.Compare(b.Amount, a.Amount), cmp.Compare(a.Parent, b.Parent))
	})
	return sorted
}

func totalOfShares(shares []SpendShare) shared.Money {
	var total shared.Money
	for _, share := range shares {
		total += share.Amount
	}
	return total
}

// withPercents truncates each percent and gives what is missing to reach 100
// to the largest share, which expects the shares sorted largest first.
func withPercents(shares []SpendShare, total shared.Money) []SpendShare {
	withPercent := make([]SpendShare, 0, len(shares))
	assigned := 0
	for _, share := range shares {
		share.Percent = int(share.Amount * fullPercent / total)
		assigned += share.Percent
		withPercent = append(withPercent, share)
	}
	withPercent[0].Percent += fullPercent - assigned
	return withPercent
}
