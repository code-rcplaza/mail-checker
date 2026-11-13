package domain

import (
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// CountChargedThisMonth counts the subscriptions with a charge in the
// calendar month of now. The month is taken in the time zone now carries, so
// the caller decides whose calendar it is.
func CountChargedThisMonth(subscriptions []Subscription, now time.Time) int {
	count := 0
	for _, s := range subscriptions {
		if countsAsSubscription(s) && hasChargeInMonthOf(s.Charges, now) {
			count++
		}
	}
	return count
}

// CountRemainingThisMonth counts the subscriptions whose next charge falls
// between tomorrow and the end of the month. A next charge dated today is
// left out: it may already be charged and not be in a statement yet.
func CountRemainingThisMonth(subscriptions []Subscription, now time.Time) int {
	count := 0
	for _, s := range subscriptions {
		if countsAsSubscription(s) && isLaterThisMonth(s.NextCharge, now) {
			count++
		}
	}
	return count
}

func countsAsSubscription(s Subscription) bool {
	return sectionOf(s) == sectionSubscriptions
}

func hasChargeInMonthOf(charges []shared.Charge, now time.Time) bool {
	return slices.ContainsFunc(charges, func(c shared.Charge) bool {
		return isSameMonth(c.OccurredAt.In(now.Location()), now)
	})
}

func isLaterThisMonth(nextCharge, now time.Time) bool {
	if nextCharge.IsZero() {
		return false
	}
	local := nextCharge.In(now.Location())
	return isSameMonth(local, now) && local.Day() > now.Day()
}

func isSameMonth(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month()
}
