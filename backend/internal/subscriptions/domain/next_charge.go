package domain

import (
	"slices"
	"time"
)

// LineCycle says how often a line charges and when: the day for a monthly
// line, the day and month for a yearly one.
type LineCycle struct {
	Cadence Cadence
	Day     int
	Month   time.Month // zero unless the line is yearly
}

// NextCharge is the last charge plus one interval of the cadence, on the
// cycle's day; without a cadence there is none. A day the month does not
// have (31, or 30 in February) lands on the last day of the month. The date
// is not moved forward when it already passed: a late charge is shown as it is.
func NextCharge(lastCharge time.Time, cycle LineCycle) (time.Time, bool) {
	switch cycle.Cadence {
	case CadenceMonthly:
		return nextMonthlyCharge(lastCharge, cycle.Day), true
	case CadenceYearly:
		return monthOf(lastCharge).plus(monthsPerYear).day(cycle.Day), true
	}
	return time.Time{}, false
}

// nextMonthlyCharge picks the typical day closest to a month after the last
// charge, so a charge a few days late or early does not shift the month.
func nextMonthlyCharge(lastCharge time.Time, typicalDay int) time.Time {
	expected := lastCharge.AddDate(0, 0, daysPerMonth)
	candidates := []time.Time{
		monthOf(expected).plus(-1).day(typicalDay),
		monthOf(expected).day(typicalDay),
		monthOf(expected).plus(1).day(typicalDay),
	}
	return slices.MinFunc(candidates, func(a, b time.Time) int {
		return int(a.Sub(expected).Abs() - b.Sub(expected).Abs())
	})
}

// calendarMonth is a month of a year in a time zone.
type calendarMonth struct {
	year     int
	month    time.Month
	location *time.Location
}

func monthOf(moment time.Time) calendarMonth {
	return calendarMonth{year: moment.Year(), month: moment.Month(), location: moment.Location()}
}

// plus moves by whole months, rolling into the neighbouring years.
func (m calendarMonth) plus(months int) calendarMonth {
	first := time.Date(m.year, m.month+time.Month(months), 1, 0, 0, 0, 0, m.location)
	return monthOf(first)
}

// day clamps the day to the last one of the month.
func (m calendarMonth) day(day int) time.Time {
	lastDay := time.Date(m.year, m.month+1, 0, 0, 0, 0, 0, m.location).Day()
	return time.Date(m.year, m.month, min(day, lastDay), 0, 0, 0, 0, m.location)
}
