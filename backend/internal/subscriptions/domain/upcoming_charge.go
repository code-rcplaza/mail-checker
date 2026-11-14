package domain

import (
	"cmp"
	"slices"
	"time"
)

// UpcomingHorizonDays is how far ahead the list of upcoming charges looks.
const UpcomingHorizonDays = 30

// UpcomingCharges lists the subscriptions (manual ones included) charging
// from today to the horizon, nearest first. A next charge dated today is in:
// it may not be charged yet.
func UpcomingCharges(subscriptions []Subscription, now time.Time) []Subscription {
	today := startOfDay(now)
	horizon := dayRange{from: today, to: today.AddDate(0, 0, UpcomingHorizonDays)}
	upcoming := slices.DeleteFunc(slices.Clone(subscriptions), func(s Subscription) bool {
		return !horizon.contains(s.NextCharge)
	})
	slices.SortStableFunc(upcoming, compareByNextChargeThenName)
	return upcoming
}

func startOfDay(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// dayRange runs from one day to another, both included.
type dayRange struct {
	from time.Time
	to   time.Time
}

// contains leaves out a zero moment, which is before any day.
func (r dayRange) contains(moment time.Time) bool {
	return !moment.Before(r.from) && !moment.After(r.to)
}

func compareByNextChargeThenName(a, b Subscription) int {
	return cmp.Or(a.NextCharge.Compare(b.NextCharge), cmp.Compare(lineName(a), lineName(b)))
}

// lineName is the user's name for the line or, without one, its parent key;
// the summary has no other name for a detected line.
func lineName(s Subscription) string {
	if s.Alias != "" {
		return s.Alias
	}
	return string(s.Parent)
}
