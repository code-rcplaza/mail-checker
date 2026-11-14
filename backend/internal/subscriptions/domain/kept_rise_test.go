package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestHideKeptRises_UsesEveryKeptAmountOfTheLine(t *testing.T) {
	const netflix GroupKey = "service:Netflix||d25"
	kept := map[GroupKey][]shared.Money{netflix: {1050000, 1100000}}
	cases := []struct {
		name      string
		after     shared.Money
		wantShown bool
	}{
		{"the last kept amount is hidden", 1100000, false},
		{"an earlier kept amount is hidden too", 1050000, false},
		{"a new amount shows", 1150000, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			summary := Summary{Rises: []Rise{{Key: netflix, Before: 1000000, After: tc.after}}}

			got := HideKeptRises(summary, kept)

			if shown := len(got.Rises) == 1; shown != tc.wantShown {
				t.Errorf("rise shown = %v, want %v", shown, tc.wantShown)
			}
			if len(summary.Rises) != 1 {
				t.Errorf("the given summary was changed")
			}
		})
	}
}
