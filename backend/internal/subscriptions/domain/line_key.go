package domain

import (
	"fmt"
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// LineKeys gives each line of a parent a key that stays stable between reads,
// so decisions and aliases stored against it survive: parent, card digits and
// the line's typical billing day ("service:Netflix|4821|d01"). The typical day
// is the median, so a late charge or a shorter period does not move it. When
// two lines share a key, the one that started later gets "#2", "#3"… Keys
// never hold internal ids or emails, since they leave through the API. Keys
// come back in the order of lines.
func LineKeys(parent GroupKey, lines [][]shared.Charge) []GroupKey {
	keys := make([]GroupKey, len(lines))
	seen := make(map[GroupKey]int)
	for _, i := range indexesByFirstCharge(lines) {
		base := baseLineKey(parent, lines[i])
		seen[base]++
		keys[i] = withCollisionSuffix(base, seen[base])
	}
	return keys
}

func baseLineKey(parent GroupKey, line []shared.Charge) GroupKey {
	return GroupKey(fmt.Sprintf("%s|%s|d%02d", parent, cardLast4Of(line), medianDayOfMonth(line)))
}

func withCollisionSuffix(base GroupKey, occurrence int) GroupKey {
	if occurrence == 1 {
		return base
	}
	return GroupKey(fmt.Sprintf("%s#%d", base, occurrence))
}

// cardLast4Of takes the first card the line shows; some charges do not show one.
func cardLast4Of(line []shared.Charge) string {
	for _, c := range line {
		if c.CardLast4 != "" {
			return c.CardLast4
		}
	}
	return ""
}

// medianDayOfMonth takes the lower middle day when the count is even.
func medianDayOfMonth(line []shared.Charge) int {
	days := make([]int, 0, len(line))
	for _, c := range line {
		days = append(days, c.OccurredAt.Day())
	}
	slices.Sort(days)
	return days[(len(days)-1)/2]
}

func indexesByFirstCharge(lines [][]shared.Charge) []int {
	indexes := make([]int, len(lines))
	for i := range lines {
		indexes[i] = i
	}
	slices.SortStableFunc(indexes, func(a, b int) int {
		return firstChargeAt(lines[a]).Compare(firstChargeAt(lines[b]))
	})
	return indexes
}

func firstChargeAt(line []shared.Charge) time.Time {
	return slices.MinFunc(line, func(a, b shared.Charge) int { return a.OccurredAt.Compare(b.OccurredAt) }).OccurredAt
}
