package domain

type DecisionKind string

const (
	DecisionConfirmed DecisionKind = "confirmed"
	DecisionDismissed DecisionKind = "dismissed"
)

// Decision is the user's answer about a possible subscription. When
// confirming, the user picks the cadence and the app does not infer it again.
type Decision struct {
	Kind    DecisionKind
	Cadence Cadence // only when confirmed
}

// isConfirmed treats a missing decision as undecided.
func isConfirmed(d *Decision) bool {
	return d != nil && d.Kind == DecisionConfirmed
}
