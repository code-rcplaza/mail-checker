package domain

import "testing"

func TestParseAppPassword_RemovesSpaces(t *testing.T) {
	got, err := ParseAppPassword("abcd efgh ijkl mnop")

	if err != nil || got != "abcdefghijklmnop" {
		t.Errorf("ParseAppPassword = %q, %v; want abcdefghijklmnop, nil", got, err)
	}
}

func TestParseAppPassword_RejectsWrongLengthOrNonLetters(t *testing.T) {
	cases := map[string]string{
		"15 letters": "abcdefghijklmno",
		"17 letters": "abcdefghijklmnopq",
		"digit":      "abcdefghijklmno1",
		"symbol":     "abcdefghijklmno-",
		"empty":      "",
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := ParseAppPassword(raw)

			if got != "" {
				t.Errorf("ParseAppPassword(%q) = %q, want empty", raw, got)
			}
			assertInvalidWithKey(t, err, "accounts.app_password_invalid")
		})
	}
}
