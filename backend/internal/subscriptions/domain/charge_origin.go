package domain

import (
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// origin is where charges were paid from. An empty card matches any card,
// since some notices and statements do not show the digits.
type origin struct {
	accountID *int64
	cardLast4 string
	charges   []shared.Charge
}

// SplitByOrigin separates the charges of one parent by where they were paid
// from: mail account and card. A statement charge has no account, so it joins
// the origin it chains with, or stays apart. Origins come out in order of
// their first charge; the caller's slice is left untouched.
func SplitByOrigin(charges []shared.Charge) [][]shared.Charge {
	chronological := sortByDate(charges)
	origins := collectByAccount(chronological)
	origins = joinChargesWithoutAccount(origins, chronological)
	return chargesOf(origins)
}

func collectByAccount(chronological []shared.Charge) []origin {
	var origins []origin
	for _, c := range chronological {
		if c.AccountID == nil {
			continue
		}
		origins = addToOrigin(origins, c, func(o origin) bool { return isSameAccount(o, c) })
	}
	return origins
}

func joinChargesWithoutAccount(origins []origin, chronological []shared.Charge) []origin {
	for _, c := range chronological {
		if c.AccountID != nil {
			continue
		}
		origins = addToOrigin(origins, c, func(o origin) bool { return chainsWithOrigin(o, c) })
	}
	return sortOriginsByFirstCharge(origins)
}

// addToOrigin puts c in the first origin that fits, or opens a new one.
func addToOrigin(origins []origin, c shared.Charge, fits func(origin) bool) []origin {
	for i, o := range origins {
		if !fits(o) || !isSameCard(o.cardLast4, c.CardLast4) {
			continue
		}
		origins[i] = withCharge(o, c)
		return origins
	}
	return append(origins, origin{accountID: c.AccountID, cardLast4: c.CardLast4, charges: []shared.Charge{c}})
}

func withCharge(o origin, c shared.Charge) origin {
	if o.cardLast4 == "" {
		o.cardLast4 = c.CardLast4
	}
	o.charges = append(o.charges, c)
	return o
}

func isSameAccount(o origin, c shared.Charge) bool {
	return o.accountID != nil && *o.accountID == *c.AccountID
}

// isSameCard treats a missing card as unknown, not as a different card.
func isSameCard(a, b string) bool {
	return a == "" || b == "" || a == b
}

func chainsWithOrigin(o origin, c shared.Charge) bool {
	for _, existing := range o.charges {
		earlier, later := existing, c
		if later.OccurredAt.Before(earlier.OccurredAt) {
			earlier, later = later, earlier
		}
		if chargesBackEachOther(earlier, later) {
			return true
		}
	}
	return false
}

func sortOriginsByFirstCharge(origins []origin) []origin {
	sorted := make([]origin, 0, len(origins))
	for _, o := range origins {
		o.charges = sortByDate(o.charges)
		sorted = append(sorted, o)
	}
	slices.SortStableFunc(sorted, func(a, b origin) int {
		return a.charges[0].OccurredAt.Compare(b.charges[0].OccurredAt)
	})
	return sorted
}

func chargesOf(origins []origin) [][]shared.Charge {
	split := make([][]shared.Charge, 0, len(origins))
	for _, o := range origins {
		split = append(split, o.charges)
	}
	return split
}
