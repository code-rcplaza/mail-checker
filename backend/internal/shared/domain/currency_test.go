package domain

import (
	"errors"
	"testing"
)

func TestParseCurrency_AcceptsCLPAndUSD(t *testing.T) {
	cases := []struct {
		raw  string
		want Currency
	}{
		{"CLP", CurrencyCLP},
		{"USD", CurrencyUSD},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := ParseCurrency(tc.raw)
			if err != nil {
				t.Fatalf("ParseCurrency(%q) error = %v, want nil", tc.raw, err)
			}
			if got != tc.want {
				t.Errorf("ParseCurrency(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

func TestParseCurrency_RejectsUnknownWithCurrencyUnknownKey(t *testing.T) {
	_, err := ParseCurrency("EUR")

	assertInvalidWithKey(t, err, "common.currency_unknown")
}

// assertInvalidWithKey checks that err is ErrInvalid and carries the given i18n key.
func assertInvalidWithKey(t *testing.T, err error, want MessageKey) {
	t.Helper()
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false (err = %v), want true", err)
	}
	var keyed KeyedError
	if !errors.As(err, &keyed) {
		t.Fatalf("errors.As(err, &KeyedError) = false, want true")
	}
	if keyed.Key != want {
		t.Errorf("Key = %q, want %q", keyed.Key, want)
	}
}
