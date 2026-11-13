package domain

import (
	"slices"
	"time"
)

// NextCharge is the last charge plus one interval of the cadence, on the
// typical day of the line; without a cadence there is none. A day the month
// does not have (31, or 30 in February) lands on the last day of the month.
// The date is not moved forward when it already passed: a late charge is
// shown as it is.
func NextCharge(lastCharge time.Time, typicalDay int, cadence Cadence) (time.Time, bool) {
	switch cadence {
	case CadenceMonthly:
		return nextMonthlyCharge(lastCharge, typicalDay), true
	case CadenceYearly:
		return dayOfMonth(lastCharge.Year()+1, lastCharge.Month(), typicalDay, lastCharge.Location()), true
	}
	return time.Time{}, false
}

// nextMonthlyCharge picks the typical day closest to a month after the last
// charge, so a charge a few days late or early does not shift the month.
func nextMonthlyCharge(lastCharge time.Time, typicalDay int) time.Time {
	expected := lastCharge.AddDate(0, 0, daysPerMonth)
	candidates := []time.Time{
		dayOfMonth(expected.Year(), expected.Month()-1, typicalDay, expected.Location()),
		dayOfMonth(expected.Year(), expected.Month(), typicalDay, expected.Location()),
		dayOfMonth(expected.Year(), expected.Month()+1, typicalDay, expected.Location()),
	}
	return slices.MinFunc(candidates, func(a, b time.Time) int {
		return int(a.Sub(expected).Abs() - b.Sub(expected).Abs())
	})
}

// dayOfMonth clamps the day to the last one of the month; a month outside
// 1–12 rolls into the neighbouring year.
func dayOfMonth(year int, month time.Month, day int, location *time.Location) time.Time {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, location).Day()
	return time.Date(year, month, min(day, lastDay), 0, 0, 0, 0, location)
}
