package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// splitRecurringFromOther decides which charges of a known service follow its
// recurring pattern. A charge backed by another (same currency, similar
// amount, valid interval) is recurring. An unbacked one is recurring only if
// it sits at a cadence interval from its nearest backed neighbor, which keeps
// a recent plan change that has no pair yet. Without any backing it falls back
// to arrival time.
func splitRecurringFromOther(chronological []shared.Charge) (recurring, other []shared.Charge) {
	backed := findBackedIndexes(chronological)
	if len(backed) == 0 {
		return splitByArrivalTime(chronological)
	}
	pattern := backedPattern{charges: chronological, backed: backed, cadence: DetectCadence(datesAtIndexes(chronological, backed))}
	for i, c := range chronological {
		if backed[i] || pattern.isAtCadenceFromNearestBacked(i) {
			recurring = append(recurring, c)
			continue
		}
		other = append(other, c)
	}
	return recurring, other
}

func findBackedIndexes(charges []shared.Charge) map[int]bool {
	backed := make(map[int]bool)
	for i := range charges {
		for j := i + 1; j < len(charges); j++ {
			if !chargesBackEachOther(charges[i], charges[j]) {
				continue
			}
			backed[i], backed[j] = true, true
		}
	}
	return backed
}

func chargesBackEachOther(earlier, later shared.Charge) bool {
	days := daysBetween(earlier.OccurredAt, later.OccurredAt)
	return earlier.Currency == later.Currency &&
		!isAmountJump(earlier.Amount, later.Amount) &&
		(isMonthlyInterval(days) || isYearlyInterval(days))
}

func datesAtIndexes(charges []shared.Charge, indexes map[int]bool) []time.Time {
	dates := make([]time.Time, 0, len(indexes))
	for i, c := range charges {
		if indexes[i] {
			dates = append(dates, c.OccurredAt)
		}
	}
	return dates
}

// backedPattern is a known service's charges (chronological), which of them
// back each other and the cadence those backed ones follow.
type backedPattern struct {
	charges []shared.Charge
	backed  map[int]bool
	cadence Cadence
}

// isAtCadenceFromNearestBacked checks the nearest backed charge on each side.
func (p backedPattern) isAtCadenceFromNearestBacked(i int) bool {
	for _, neighbor := range p.nearestBackedNeighbors(i) {
		days := daysBetween(p.charges[neighbor].OccurredAt, p.charges[i].OccurredAt)
		if fitsCadence(days, p.cadence) {
			return true
		}
	}
	return false
}

func (p backedPattern) nearestBackedNeighbors(i int) []int {
	var neighbors []int
	for before := i - 1; before >= 0; before-- {
		if p.backed[before] {
			neighbors = append(neighbors, before)
			break
		}
	}
	for after := i + 1; after < len(p.charges); after++ {
		if p.backed[after] {
			neighbors = append(neighbors, after)
			break
		}
	}
	return neighbors
}

// fitsCadence takes the absolute interval, since the neighbor may come after.
func fitsCadence(days int, cadence Cadence) bool {
	absoluteDays := max(days, -days)
	switch cadence {
	case CadenceMonthly:
		return isMonthlyInterval(absoluteDays)
	case CadenceYearly:
		return isYearlyInterval(absoluteDays)
	}
	return false
}

// splitByArrivalTime treats a charge arriving sooner than a monthly interval
// after the last recurring one as off-pattern: it cannot be the next bill.
func splitByArrivalTime(chronological []shared.Charge) (recurring, other []shared.Charge) {
	for _, c := range chronological {
		if len(recurring) > 0 && arrivesTooSoon(recurring[len(recurring)-1], c) {
			other = append(other, c)
			continue
		}
		recurring = append(recurring, c)
	}
	return recurring, other
}

func arrivesTooSoon(previous, next shared.Charge) bool {
	return daysBetween(previous.OccurredAt, next.OccurredAt) < MonthlyIntervalMinDays
}
