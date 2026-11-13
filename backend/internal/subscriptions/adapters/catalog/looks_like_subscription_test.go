package catalog

import "testing"

func TestLooksLikeSubscription(t *testing.T) {
	cases := []struct {
		name     string
		merchant string
		want     bool
	}{
		{"SUBSCRIPTION as a word", "MICROSOFT. SUBSCRIPTION", true},
		{"MEMBERSHIP as a word", "GYM MEMBERSHIP", true},
		{"SUSCRIPCION with accent and lower case", "Suscripción Anual", true},
		{"MEMBRESIA with accent", "Membresía Club", true},
		{"PREMIUM as a word", "ACME Premium", true},
		{"PLAN as a word", "PLAN MOVIL", true},
		{"PLAN inside a longer word", "PLANETA", false},
		{"PREMIUM inside a longer word", "SUPREMIUM", false},
		{"plural is another word", "MEMBERSHIPS", false},
		{"no subscription word", "NETFLIX.COM", false},
		{"word glued with a dot", "MICROSOFT.SUBSCRIPTION", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Catalog{}.LooksLikeSubscription(tc.merchant)

			if got != tc.want {
				t.Errorf("LooksLikeSubscription(%q) = %v, want %v", tc.merchant, got, tc.want)
			}
		})
	}
}
