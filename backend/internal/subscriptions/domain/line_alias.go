package domain

import (
	"strings"
	"unicode/utf8"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// MaxLineAliasRunes counts characters, not bytes, so accents cost one each.
const MaxLineAliasRunes = 40

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
	group, isFound := findGroup(groups, key)
	return isFound && !group.IsOneOff
}

// ApplyLineAliases returns the summary with each line carrying the alias of
// its key; the given summary is left untouched.
func ApplyLineAliases(summary Summary, aliases map[GroupKey]string) Summary {
	summary.Subscriptions = withAliases(summary.Subscriptions, aliases)
	summary.Possible = withAliases(summary.Possible, aliases)
	summary.Inactive = withAliases(summary.Inactive, aliases)
	return summary
}

func withAliases(subscriptions []Subscription, aliases map[GroupKey]string) []Subscription {
	named := make([]Subscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		s.Alias = aliases[s.Key]
		named = append(named, s)
	}
	return named
}
