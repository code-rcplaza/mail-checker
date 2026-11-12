package domain

import (
	"errors"
	"fmt"
	"testing"
)

const keyBankNotDetected MessageKey = "statements.bank_not_detected"

func TestInvalid_WrapsErrInvalidAndKeepsKeyAndParams(t *testing.T) {
	err := Invalid(keyBankNotDetected, map[string]any{"count": 2})

	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, want true")
	}
	var keyed KeyedError
	if !errors.As(err, &keyed) {
		t.Fatalf("errors.As(err, &KeyedError) = false, want true")
	}
	if keyed.Key != keyBankNotDetected {
		t.Errorf("Key = %q, want %q", keyed.Key, keyBankNotDetected)
	}
	if keyed.Params["count"] != 2 {
		t.Errorf("Params[count] = %v, want 2", keyed.Params["count"])
	}
}

func TestNotFoundAndConflict_WrapSentinelAndKeepKey(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		sentinel error
		key      MessageKey
	}{
		{"NotFound", NotFound("accounts.not_found"), ErrNotFound, "accounts.not_found"},
		{"Conflict", Conflict("accounts.email_taken"), ErrConflict, "accounts.email_taken"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !errors.Is(tc.err, tc.sentinel) {
				t.Fatalf("errors.Is(err, sentinel) = false, want true")
			}
			var keyed KeyedError
			if !errors.As(tc.err, &keyed) {
				t.Fatalf("errors.As(err, &KeyedError) = false, want true")
			}
			if keyed.Key != tc.key {
				t.Errorf("Key = %q, want %q", keyed.Key, tc.key)
			}
		})
	}
}

func TestKeyedError_SurvivesFmtErrorfWrapping(t *testing.T) {
	wrapped := fmt.Errorf("import statement: %w", Invalid(keyBankNotDetected, nil))

	if !errors.Is(wrapped, ErrInvalid) {
		t.Fatalf("errors.Is(wrapped, ErrInvalid) = false, want true")
	}
	var keyed KeyedError
	if !errors.As(wrapped, &keyed) {
		t.Fatalf("errors.As(wrapped, &KeyedError) = false, want true")
	}
	if keyed.Key != keyBankNotDetected {
		t.Errorf("Key = %q, want %q", keyed.Key, keyBankNotDetected)
	}
}

func TestSentinels_AreDistinct(t *testing.T) {
	sentinels := map[string]error{
		"ErrInvalid":  ErrInvalid,
		"ErrNotFound": ErrNotFound,
		"ErrConflict": ErrConflict,
	}
	for nameA, a := range sentinels {
		for nameB, b := range sentinels {
			if nameA == nameB {
				continue
			}
			if errors.Is(a, b) {
				t.Errorf("errors.Is(%s, %s) = true, want false", nameA, nameB)
			}
		}
	}
}

func TestInvalid_CopiesParams(t *testing.T) {
	params := map[string]any{"count": 2}
	err := Invalid(keyBankNotDetected, params)

	params["count"] = 99

	var keyed KeyedError
	if !errors.As(err, &keyed) {
		t.Fatalf("errors.As(err, &KeyedError) = false, want true")
	}
	if keyed.Params["count"] != 2 {
		t.Errorf("Params[count] = %v after mutating caller map, want 2", keyed.Params["count"])
	}
}
