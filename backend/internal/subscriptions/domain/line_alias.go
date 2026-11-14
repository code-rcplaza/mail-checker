package domain

import (
	"strings"
	"unicode/utf8"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// MaxLineAliasRunes counts characters, not bytes, so accents cost one each.
const MaxLineAliasRunes = 40

// LineAlias is the user's name for a line.
type LineAlias struct {
	Key   GroupKey
	Alias string
}

const (
	KeyAliasTooLong shared.MessageKey = "subscriptions.alias_too_long"
	KeyLineNotFound shared.MessageKey = "subscriptions.line_not_found"
)

// ParseLineAlias trims the user's name for a line; an empty result means the
// alias is cleared.
func ParseLineAlias(raw string) (string, error) {
	alias := strings.TrimSpace(raw)
	if utf8.RuneCountInString(alias) > MaxLineAliasRunes {
		return "", shared.Invalid(KeyAliasTooLong, map[string]any{"max": MaxLineAliasRunes})
	}
	return alias, nil
}

// IsLine tells whether key names a line of the groups; one-off charges are
// not a line.
func IsLine(groups []ChargeGroup, key GroupKey) bool {
	_, isLine := FindLine(groups, key)
	return isLine
}

// FindLine is the group of the line key names; one-off charges are not a line.
func FindLine(groups []ChargeGroup, key GroupKey) (ChargeGroup, bool) {
	group, isFound := findGroup(groups, key)
	if !isFound || group.IsOneOff {
		return ChargeGroup{}, false
	}
	return group, true
}

// ApplyLineAliases returns the summary with each line carrying the alias of
// its key; a line without one keeps its own (a manual's name). The given
// summary is left untouched.
func ApplyLineAliases(summary Summary, aliases map[GroupKey]string) Summary {
	return changeEveryLine(summary, func(line Subscription) Subscription {
		if alias, hasAlias := aliases[line.Key]; hasAlias {
			line.Alias = alias
		}
		return line
	})
}
