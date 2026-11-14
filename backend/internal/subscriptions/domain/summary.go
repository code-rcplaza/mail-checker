package domain

import (
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type Summary struct {
	Subscriptions []Subscription // active + confirmed possible; they add to spend
	Possible      []Subscription // undecided
	Inactive      []Subscription
	OtherCharges  []shared.Charge // not a subscription + dismissed possible
	MonthlySpend  map[shared.Currency]shared.Money
	// Charges of the calendar month of "now", in the time zone "now" carries.
	ChargedThisMonth   int    // subscriptions already charged
	RemainingThisMonth int    // subscriptions whose next charge is still to come
	Rises              []Rise // subscriptions whose last charge went up
	Total              Total  // the monthly spend in pesos; set with ApplyExchangeRates
	// ExchangeRate is the rate used for every amount in pesos; nil without one.
	ExchangeRate *shared.ExchangeRate
	// AsOf is the moment the summary describes; its time zone is the user's
	// calendar for "this month".
	AsOf time.Time
}

// SummaryInput is what a summary is built from: the grouped charges, the
// user's decisions on possible subscriptions and the detector, whose Now is
// the moment the summary describes.
type SummaryInput struct {
	Groups    []ChargeGroup
	Decisions map[GroupKey]Decision
	Detector  SubscriptionDetector
}

// BuildSummary puts each line in its section, honoring the user's decisions
// on possible subscriptions, sends one-off charges to other charges and adds
// up the monthly spend.
func BuildSummary(input SummaryInput) Summary {
	lines, oneOff := splitOneOffGroups(input.Groups)
	summary := sortIntoSections(buildSubscriptions(lines, input))
	summary.OtherCharges = append(summary.OtherCharges, chargesOfGroups(oneOff)...)
	summary.Rises = risesOf(summary.Subscriptions)
	summary.AsOf = input.Detector.Now
	return recountSpend(summary)
}

// recountSpend recalculates what depends on the subscriptions section: the
// monthly spend and this month's charged and remaining counts.
func recountSpend(summary Summary) Summary {
	summary.MonthlySpend = MonthlySpend(summary.Subscriptions)
	summary.ChargedThisMonth = CountChargedThisMonth(summary.Subscriptions, summary.AsOf)
	summary.RemainingThisMonth = CountRemainingThisMonth(summary.Subscriptions, summary.AsOf)
	return summary
}

// changeEveryLine applies change to each line of the three sections that hold
// lines; the given summary is left untouched.
func changeEveryLine(summary Summary, change func(Subscription) Subscription) Summary {
	summary.Subscriptions = changeEach(summary.Subscriptions, change)
	summary.Possible = changeEach(summary.Possible, change)
	summary.Inactive = changeEach(summary.Inactive, change)
	return summary
}

func changeEach(lines []Subscription, change func(Subscription) Subscription) []Subscription {
	if lines == nil {
		return nil
	}
	changed := make([]Subscription, 0, len(lines))
	for _, line := range lines {
		changed = append(changed, change(line))
	}
	return changed
}

// removeLines takes out of the three sections the lines isRemoved picks; the
// given summary is left untouched.
func removeLines(summary Summary, isRemoved func(Subscription) bool) Summary {
	summary.Subscriptions = slices.DeleteFunc(slices.Clone(summary.Subscriptions), isRemoved)
	summary.Possible = slices.DeleteFunc(slices.Clone(summary.Possible), isRemoved)
	summary.Inactive = slices.DeleteFunc(slices.Clone(summary.Inactive), isRemoved)
	return summary
}

func splitOneOffGroups(groups []ChargeGroup) (lines, oneOff []ChargeGroup) {
	for _, group := range groups {
		if group.IsOneOff {
			oneOff = append(oneOff, group)
			continue
		}
		lines = append(lines, group)
	}
	return lines, oneOff
}

func chargesOfGroups(groups []ChargeGroup) []shared.Charge {
	var charges []shared.Charge
	for _, group := range groups {
		charges = append(charges, group.Charges...)
	}
	return charges
}

func buildSubscriptions(lines []ChargeGroup, input SummaryInput) []Subscription {
	subscriptions := make([]Subscription, 0, len(lines))
	for _, group := range lines {
		subscriptions = append(subscriptions, buildSubscription(group, input))
	}
	return subscriptions
}

func buildSubscription(group ChargeGroup, input SummaryInput) Subscription {
	detection := input.Detector.Detect(group)
	amount, currency := currentAmount(group)
	subscription := Subscription{
		Key:      group.Key,
		Parent:   group.Parent,
		Status:   detection.Status,
		Cadence:  detection.Cadence,
		Amount:   amount,
		Currency: currency,
		Charges:  group.Charges,
	}
	return withNextCharge(applyDecision(subscription, input.Decisions))
}

// withNextCharge runs after the decision: a confirmed possible only then has
// the cadence the user picked. An inactive line is not expected to charge again.
func withNextCharge(s Subscription) Subscription {
	if s.Status == StatusInactive || len(s.Charges) == 0 {
		return s
	}
	lastCharge := lastOf(s.Charges).OccurredAt
	nextCharge, hasNextCharge := NextCharge(lastCharge, LineCycle{Cadence: s.Cadence, Day: medianDayOfMonth(s.Charges)})
	if hasNextCharge {
		s.NextCharge = nextCharge
	}
	return s
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
