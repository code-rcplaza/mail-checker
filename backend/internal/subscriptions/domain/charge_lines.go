package domain

import (
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Ideal lengths used to pick between lines that both accept a charge.
const (
	daysPerMonth = 30
	daysPerYear  = 365
)

// ChainLines splits the charges of one origin into lines, one per
// subscription. Amount comes before date: charges are first split into amount
// bands, so two subscriptions of different prices never mix even when their
// days are close to the month turn. Each band is chained by date, a line that
// starts right after another ended joins it as a plan change, and one-off
// charges are returned apart.
func ChainLines(charges []shared.Charge) (lines [][]shared.Charge, other []shared.Charge) {
	var chained [][]shared.Charge
	for _, band := range splitByAmountBand(sortByDate(charges)) {
		chained = append(chained, chainChronologically(band)...)
	}
	joined := joinPlanChanges(sortLinesByFirstCharge(chained))
	return separateOneOffCharges(joined)
}

// splitByAmountBand puts each charge with the band of the same currency whose
// last amount is within PlanChangeThresholdPercent and closest.
func splitByAmountBand(chronological []shared.Charge) [][]shared.Charge {
	var bands [][]shared.Charge
	for _, c := range chronological {
		best, isFound := closestBandFor(bands, c)
		if !isFound {
			bands = append(bands, []shared.Charge{c})
			continue
		}
		bands[best] = append(bands[best], c)
	}
	return bands
}

func closestBandFor(bands [][]shared.Charge, c shared.Charge) (int, bool) {
	best, isFound := 0, false
	for i, band := range bands {
		last := lastOf(band)
		if last.Currency != c.Currency || isAmountJump(last.Amount, c.Amount) {
			continue
		}
		if !isFound || compareAmountDistance(last, lastOf(bands[best]), c) < 0 {
			best, isFound = i, true
		}
	}
	return best, isFound
}

// joinPlanChanges appends to a line the line that starts a valid interval
// after it ended; lines that overlap in time are separate subscriptions.
func joinPlanChanges(sortedLines [][]shared.Charge) [][]shared.Charge {
	var joined [][]shared.Charge
	for _, line := range sortedLines {
		previous, isFound := lineEndingRightBefore(joined, line[0])
		if !isFound {
			joined = append(joined, line)
			continue
		}
		joined[previous] = append(joined[previous], line...)
	}
	return joined
}

func lineEndingRightBefore(lines [][]shared.Charge, first shared.Charge) (int, bool) {
	for i, line := range lines {
		if lastOf(line).OccurredAt.Before(first.OccurredAt) && continuesLine(line, first) {
			return i, true
		}
	}
	return 0, false
}

func sortLinesByFirstCharge(lines [][]shared.Charge) [][]shared.Charge {
	sorted := slices.Clone(lines)
	slices.SortStableFunc(sorted, func(a, b []shared.Charge) int {
		return a[0].OccurredAt.Compare(b[0].OccurredAt)
	})
	return sorted
}

func chainChronologically(chronological []shared.Charge) [][]shared.Charge {
	var lines [][]shared.Charge
	for _, c := range chronological {
		best, isFound := bestLineFor(lines, c)
		if !isFound {
			lines = append(lines, []shared.Charge{c})
			continue
		}
		lines[best] = append(lines[best], c)
	}
	return lines
}

func bestLineFor(lines [][]shared.Charge, c shared.Charge) (int, bool) {
	best, isFound := 0, false
	for i, line := range lines {
		if !continuesLine(line, c) {
			continue
		}
		if !isFound || fitsBetter(line, lines[best], c) {
			best, isFound = i, true
		}
	}
	return best, isFound
}

func continuesLine(line []shared.Charge, c shared.Charge) bool {
	days := daysBetween(lastOf(line).OccurredAt, c.OccurredAt)
	return isMonthlyInterval(days) || isYearlyInterval(days)
}

// fitsBetter compares by amount first and by expected date on a tie.
func fitsBetter(candidate, current []shared.Charge, c shared.Charge) bool {
	byAmount := compareAmountDistance(lastOf(candidate), lastOf(current), c)
	if byAmount != 0 {
		return byAmount < 0
	}
	return daysOffCadence(lastOf(candidate), c) < daysOffCadence(lastOf(current), c)
}

// compareAmountDistance compares relative distances without floats:
// |a-c|/a against |b-c|/b. Another currency is always farther.
func compareAmountDistance(a, b, c shared.Charge) int {
	aOtherCurrency, bOtherCurrency := a.Currency != c.Currency, b.Currency != c.Currency
	if aOtherCurrency != bOtherCurrency {
		if aOtherCurrency {
			return 1
		}
		return -1
	}
	left := absMoney(a.Amount-c.Amount) * b.Amount
	right := absMoney(b.Amount-c.Amount) * a.Amount
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	}
	return 0
}

// daysOffCadence is how far the interval is from a whole number of months or a year.
func daysOffCadence(last, c shared.Charge) int {
	days := daysBetween(last.OccurredAt, c.OccurredAt)
	off := absInt(days - daysPerYear)
	for months := 1; months <= MaxMonthsPerInterval; months++ {
		off = min(off, absInt(days-months*daysPerMonth))
	}
	return off
}

func separateOneOffCharges(chained [][]shared.Charge) (lines [][]shared.Charge, other []shared.Charge) {
	other = []shared.Charge{}
	for _, line := range chained {
		if isOneOff(line, chained) {
			other = append(other, line[0])
			continue
		}
		lines = append(lines, line)
	}
	return lines, other
}

// isOneOff: a single charge that did not recur although the origin was billed
// again a full monthly interval later.
func isOneOff(line []shared.Charge, chained [][]shared.Charge) bool {
	if len(line) != 1 {
		return false
	}
	for _, otherLine := range chained {
		if len(otherLine) > 1 && daysBetween(line[0].OccurredAt, lastOf(otherLine).OccurredAt) > MonthlyIntervalMaxDays {
			return true
		}
	}
	return false
}

func lastOf(line []shared.Charge) shared.Charge {
	return line[len(line)-1]
}

func absMoney(m shared.Money) shared.Money {
	return max(m, -m)
}

func absInt(n int) int {
	return max(n, -n)
}
