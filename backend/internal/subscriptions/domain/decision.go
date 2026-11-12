package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type DecisionKind string

const (
	DecisionConfirmed DecisionKind = "confirmed"
	DecisionDismissed DecisionKind = "dismissed"
)

const (
	KeyCadenceRequired  shared.MessageKey = "subscriptions.cadence_required"
	KeyDecisionUnknown  shared.MessageKey = "subscriptions.decision_unknown"
	KeyPossibleNotFound shared.MessageKey = "subscriptions.possible_not_found"
)

// Decision is the user's answer about a possible subscription. When
// confirming, the user picks the cadence and the app does not infer it again.
type Decision struct {
	Kind    DecisionKind
	Cadence Cadence // only when confirmed
}

// The actions the user can take on a possible subscription.
const (
	ActionConfirm = "confirm"
	ActionDismiss = "dismiss"
)

// NewDecision validates the user's answer: confirming needs a monthly or
// yearly cadence; dismissing needs none.
func NewDecision(action, cadence string) (Decision, error) {
	switch action {
	case ActionConfirm:
		return newConfirmation(cadence)
	case ActionDismiss:
		return Decision{Kind: DecisionDismissed}, nil
	}
	return Decision{}, shared.Invalid(KeyDecisionUnknown, nil)
}

func newConfirmation(cadence string) (Decision, error) {
	chosen := Cadence(cadence)
	if chosen != CadenceMonthly && chosen != CadenceYearly {
		return Decision{}, shared.Invalid(KeyCadenceRequired, nil)
	}
	return Decision{Kind: DecisionConfirmed, Cadence: chosen}, nil
}

// IsCurrentPossible tells whether key names a group that is a possible
// subscription today, with the same detection the summary uses.
func IsCurrentPossible(groups []ChargeGroup, key GroupKey, looksLikeSubscription LooksLikeSubscriptionFunc, now time.Time) bool {
	group, isFound := findGroup(groups, key)
	if !isFound {
		return false
	}
	return DetectSubscription(group, looksLikeSubscription, now).Status == StatusPossible
}

func findGroup(groups []ChargeGroup, key GroupKey) (ChargeGroup, bool) {
	for _, group := range groups {
		if group.Key == key {
			return group, true
		}
	}
	return ChargeGroup{}, false
}

// isConfirmed treats a missing decision as undecided.
func isConfirmed(d *Decision) bool {
	return d != nil && d.Kind == DecisionConfirmed
}
