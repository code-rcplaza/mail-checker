package domain

import "strings"

// DefaultPeriodMonths is the period the summary reads when none is asked for.
const DefaultPeriodMonths = 12

// ParentName is the readable name of a line's parent: the service or
// umbrella name, or the merchant without its currency. A manual line has none:
// its name is its alias.
func ParentName(parent GroupKey) string {
	kind, rest, _ := strings.Cut(string(parent), ":")
	switch kind {
	case "service", "umbrella":
		return rest
	case "merchant":
		merchant, _ := cutLast(rest, ":")
		return merchant
	default:
		return ""
	}
}

func cutLast(text, separator string) (string, string) {
	at := strings.LastIndex(text, separator)
	if at < 0 {
		return text, ""
	}
	return text[:at], text[at+len(separator):]
}

// RisesWithImpact gives each rise of the summary with what it costs, by the
// cadence of its line.
func RisesWithImpact(summary Summary) []RiseWithImpact {
	rises := make([]RiseWithImpact, 0, len(summary.Rises))
	for _, rise := range summary.Rises {
		rises = append(rises, RiseWithImpact{Rise: rise, Impact: RiseImpact(rise, cadenceOfLine(summary.Subscriptions, rise.Key))})
	}
	return rises
}
