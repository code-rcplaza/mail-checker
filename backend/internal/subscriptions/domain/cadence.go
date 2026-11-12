package domain

import (
	"slices"
	"time"
)

type Cadence string

const (
	CadenceNone    Cadence = ""
	CadenceMonthly Cadence = "monthly"
	CadenceYearly  Cadence = "yearly"
)

// Recurrence rules. A monthly interval may span up to MaxMonthsPerInterval
// months because some purchase emails never arrive.
const (
	MonthlyIntervalMinDays = 24
	MonthlyIntervalMaxDays = 38
	MaxMonthsPerInterval   = 3
	YearlyIntervalMinDays  = 330
	YearlyIntervalMaxDays  = 400
)

// DetectCadence tells whether the dates repeat monthly or yearly. Every
// interval must fit the same cadence; with fewer than two dates there is none.
func DetectCadence(dates []time.Time) Cadence {
	intervals := intervalsInDays(sortDates(dates))
	if len(intervals) == 0 {
		return CadenceNone
	}
	if allIntervals(intervals, isMonthlyInterval) {
		return CadenceMonthly
	}
	if allIntervals(intervals, isYearlyInterval) {
		return CadenceYearly
	}
	return CadenceNone
}

// sortDates returns a sorted copy; the caller's slice is left untouched.
func sortDates(dates []time.Time) []time.Time {
	sorted := slices.Clone(dates)
	slices.SortFunc(sorted, time.Time.Compare)
	return sorted
}

func intervalsInDays(sorted []time.Time) []int {
	intervals := make([]int, 0, max(len(sorted)-1, 0))
	for i := 1; i < len(sorted); i++ {
		intervals = append(intervals, daysBetween(sorted[i-1], sorted[i]))
	}
	return intervals
}

// daysBetween truncates to whole days, so 30 days and 23 hours count as 30.
func daysBetween(earlier, later time.Time) int {
	return int(later.Sub(earlier).Hours() / 24)
}

func allIntervals(intervals []int, fits func(days int) bool) bool {
	for _, days := range intervals {
		if !fits(days) {
			return false
		}
	}
	return true
}

// isMonthlyInterval accepts k months of 24–38 days each, with k up to MaxMonthsPerInterval.
func isMonthlyInterval(days int) bool {
	for months := 1; months <= MaxMonthsPerInterval; months++ {
		if days >= months*MonthlyIntervalMinDays && days <= months*MonthlyIntervalMaxDays {
			return true
		}
	}
	return false
}

func isYearlyInterval(days int) bool {
	return days >= YearlyIntervalMinDays && days <= YearlyIntervalMaxDays
}
