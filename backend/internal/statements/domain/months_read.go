package domain

import (
	"cmp"
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Month is a calendar month of a year.
type Month struct {
	Year  int
	Month time.Month
}

// MonthsRead tells which months the charges of a statement cover, each once
// and in order. The month is the one the user saw on the charge, so it is
// taken in their time zone and not in UTC.
func MonthsRead(charges []shared.Charge, location *time.Location) []Month {
	months := make([]Month, 0, len(charges))
	for _, charge := range charges {
		months = append(months, monthOf(charge.OccurredAt, location))
	}
	slices.SortFunc(months, compareMonths)
	return slices.Compact(months)
}

func monthOf(moment time.Time, location *time.Location) Month {
	local := moment.In(location)
	return Month{Year: local.Year(), Month: local.Month()}
}

func compareMonths(a, b Month) int {
	return cmp.Or(cmp.Compare(a.Year, b.Year), cmp.Compare(a.Month, b.Month))
}
