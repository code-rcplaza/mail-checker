package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type Summary struct {
	Subscriptions []Subscription // active + confirmed possible; they add to spend
	Possible      []Subscription // undecided
	Inactive      []Subscription
	OtherCharges  []shared.Charge // not a subscription + dismissed possible
	MonthlySpend  map[shared.Currency]shared.Money
}

// BuildSummary puts each group in its section, honoring the user's decisions
// on possible subscriptions, and adds up the monthly spend.
func BuildSummary(
	groups []ChargeGroup,
	decisions map[GroupKey]Decision,
	looksLikeSubscription LooksLikeSubscriptionFunc,
	now time.Time,
) Summary {
	subscriptions := buildSubscriptions(groups, decisions, looksLikeSubscription, now)
	summary := sortIntoSections(subscriptions)
	summary.MonthlySpend = MonthlySpend(summary.Subscriptions)
	return summary
}

func buildSubscriptions(
	groups []ChargeGroup,
	decisions map[GroupKey]Decision,
	looksLikeSubscription LooksLikeSubscriptionFunc,
	now time.Time,
) []Subscription {
	subscriptions := make([]Subscription, 0, len(groups))
	for _, group := range groups {
		subscriptions = append(subscriptions, buildSubscription(group, decisions, looksLikeSubscription, now))
	}
	return subscriptions
}

func buildSubscription(
	group ChargeGroup,
	decisions map[GroupKey]Decision,
	looksLikeSubscription LooksLikeSubscriptionFunc,
	now time.Time,
) Subscription {
	detection := DetectSubscription(group, looksLikeSubscription, now)
	amount, currency := currentAmount(group)
	subscription := Subscription{
		Key:      group.Key,
		Status:   detection.Status,
		Cadence:  detection.Cadence,
		Amount:   amount,
		Currency: currency,
		Charges:  group.Charges,
	}
	return applyDecision(subscription, decisions)
}

func currentAmount(group ChargeGroup) (shared.Money, shared.Currency) {
	if group.IsKnownService {
		known := AnalyzeKnownServiceAmount(group.Charges)
		return known.Current, known.Currency
	}
	unknown := AnalyzeUnknownMerchantAmount(group.Charges)
	return unknown.Current, unknown.Currency
}

// applyDecision only touches possible subscriptions: a decision answers a
// possible one, and a confirmation brings the cadence the user picked.
func applyDecision(s Subscription, decisions map[GroupKey]Decision) Subscription {
	decision, isDecided := decisions[s.Key]
	if s.Status != StatusPossible || !isDecided {
		return s
	}
	s.Decision = &decision
	if decision.Kind == DecisionConfirmed {
		s.Cadence = decision.Cadence
	}
	return s
}

func sortIntoSections(subscriptions []Subscription) Summary {
	var summary Summary
	for _, s := range subscriptions {
		switch sectionOf(s) {
		case sectionSubscriptions:
			summary.Subscriptions = append(summary.Subscriptions, s)
		case sectionPossible:
			summary.Possible = append(summary.Possible, s)
		case sectionInactive:
			summary.Inactive = append(summary.Inactive, s)
		case sectionOtherCharges:
			summary.OtherCharges = append(summary.OtherCharges, s.Charges...)
		}
	}
	return summary
}

type section int

const (
	sectionSubscriptions section = iota
	sectionPossible
	sectionInactive
	sectionOtherCharges
)

func sectionOf(s Subscription) section {
	switch s.Status {
	case StatusActive:
		return sectionSubscriptions
	case StatusInactive:
		return sectionInactive
	case StatusPossible:
		return sectionOfPossible(s.Decision)
	}
	return sectionOtherCharges
}

func sectionOfPossible(decision *Decision) section {
	if decision == nil {
		return sectionPossible
	}
	if decision.Kind == DecisionConfirmed {
		return sectionSubscriptions
	}
	return sectionOtherCharges
}
