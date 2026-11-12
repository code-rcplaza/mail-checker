package domain

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestParseLineAlias(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"trims spaces", "  Padres ", "Padres"},
		{"blank clears", "   ", ""},
		{"empty clears", "", ""},
		{"40 runes with accents fit", strings.Repeat("á", MaxLineAliasRunes), strings.Repeat("á", MaxLineAliasRunes)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseLineAlias(tc.raw)

			if err != nil || got != tc.want {
				t.Errorf("ParseLineAlias(%q) = %q, %v; want %q, nil", tc.raw, got, err, tc.want)
			}
		})
	}
}

func TestParseLineAlias_RejectsMoreThan40Runes(t *testing.T) {
	_, err := ParseLineAlias(strings.Repeat("a", MaxLineAliasRunes+1))

	if !errors.Is(err, shared.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false (err = %v), want true", err)
	}
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) || keyed.Key != KeyAliasTooLong {
		t.Errorf("key = %q, want %q", keyed.Key, KeyAliasTooLong)
	}
	if MaxLineAliasRunes != 40 || KeyAliasTooLong != "subscriptions.alias_too_long" || KeyLineNotFound != "subscriptions.line_not_found" {
		t.Errorf("max = %d, keys = %q %q", MaxLineAliasRunes, KeyAliasTooLong, KeyLineNotFound)
	}
}

func TestIsLine_OnlyLinesOfTheGroups(t *testing.T) {
	groups := []ChargeGroup{
		{Key: "service:Claude||d01", Parent: "service:Claude"},
		{Key: "service:Claude|one-off", Parent: "service:Claude", IsOneOff: true},
	}
	cases := map[GroupKey]bool{
		"service:Claude||d01":    true,
		"service:Claude|one-off": false,
		"service:Claude":         false,
	}
	for key, want := range cases {
		if got := IsLine(groups, key); got != want {
			t.Errorf("IsLine(%q) = %v, want %v", key, got, want)
		}
	}
}

func TestApplyLineAliases_EverySectionGetsTheAliasOfItsKey(t *testing.T) {
	summary := Summary{
		Subscriptions: []Subscription{{Key: "service:Netflix||d01"}, {Key: "service:Netflix||d03"}},
		Possible:      []Subscription{{Key: "merchant:SEGUROS CLC:CLP||d10"}},
		Inactive:      []Subscription{{Key: "service:Video Co||d05"}},
	}
	aliases := map[GroupKey]string{
		"service:Netflix||d03":          "Padres",
		"merchant:SEGUROS CLC:CLP||d10": "Casa",
		"service:Video Co||d05":         "Hijos",
	}

	got := ApplyLineAliases(summary, aliases)

	want := Summary{
		Subscriptions: []Subscription{{Key: "service:Netflix||d01"}, {Key: "service:Netflix||d03", Alias: "Padres"}},
		Possible:      []Subscription{{Key: "merchant:SEGUROS CLC:CLP||d10", Alias: "Casa"}},
		Inactive:      []Subscription{{Key: "service:Video Co||d05", Alias: "Hijos"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("summary = %+v, want %+v", got, want)
	}
	if summary.Subscriptions[1].Alias != "" {
		t.Errorf("ApplyLineAliases changed its input")
	}
}
