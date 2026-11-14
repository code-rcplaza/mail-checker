package domain

import (
	"slices"
	"time"
)

// HiddenLine is a line the user stopped following, and since when.
type HiddenLine struct {
	Key      GroupKey
	HiddenAt time.Time
}

// HideLines takes out the lines the user stopped following and recounts what
// depends on them; the given summary is left untouched.
func HideLines(summary Summary, hidden []GroupKey) Summary {
	isHidden := func(key GroupKey) bool { return slices.Contains(hidden, key) }
	summary = removeLines(summary, func(line Subscription) bool { return isHidden(line.Key) })
	summary.Rises = slices.DeleteFunc(slices.Clone(summary.Rises), func(rise Rise) bool { return isHidden(rise.Key) })
	return recountSpend(summary)
}
