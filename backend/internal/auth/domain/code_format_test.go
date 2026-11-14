package domain

import "testing"

func TestParseCode(t *testing.T) {
	cases := []struct {
		raw     string
		want    string
		wantErr bool
	}{
		{"042917", "042917", false},
		{"  042917 ", "042917", false},
		{"42917", "", true},
		{"0429170", "", true},
		{"04291a", "", true},
		{"", "", true},
		{"０４２９１７", "", true}, // full-width digits are not ASCII digits
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := ParseCode(tc.raw)

			if got != tc.want {
				t.Errorf("ParseCode(%q) = %q, want %q", tc.raw, got, tc.want)
			}
			if tc.wantErr {
				assertInvalidWithKey(t, err, KeyCodeFormatInvalid)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("ParseCode(%q) error = %v, want nil", tc.raw, err)
			}
		})
	}
}

func TestInitials(t *testing.T) {
	cases := map[string]string{
		"Ana Pérez":            "AP",
		"Ana":                  "A",
		"ana maría de la cruz": "AM",
		"  ana   pérez ":       "AP",
		"Óscar Ñúñez":          "ÓÑ",
		"🙂 Ana":                "🙂A",
		"":                     "",
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Initials(name); got != want {
				t.Errorf("Initials(%q) = %q, want %q", name, got, want)
			}
		})
	}
}

func TestRequireTermsAndParseName(t *testing.T) {
	assertInvalidWithKey(t, RequireTerms(false), KeyTermsRequired)
	if err := RequireTerms(true); err != nil {
		t.Errorf("RequireTerms(true) = %v, want nil", err)
	}
	if got, err := ParseName("  Ana "); got != "Ana" || err != nil {
		t.Errorf("ParseName = %q, %v; want Ana", got, err)
	}
}
