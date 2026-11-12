package domain

import (
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestParseFileKind_AcceptsKnownKinds(t *testing.T) {
	cases := []struct {
		raw  string
		want FileKind
	}{
		{"pdf", FileKindPDF},
		{"xlsx", FileKindXLSX},
		{"csv", FileKindCSV},
		{"text", FileKindText},
	}
	for _, tc := range cases {
		got, err := ParseFileKind(tc.raw)

		if err != nil || got != tc.want {
			t.Errorf("ParseFileKind(%q) = %q, %v; want %q, nil", tc.raw, got, err, tc.want)
		}
	}
}

func TestParseFileKind_RejectsUnknownKind(t *testing.T) {
	_, err := ParseFileKind("doc")

	if !errors.Is(err, shared.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false (err = %v), want true", err)
	}
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) || keyed.Key != "statements.file_kind_unknown" {
		t.Errorf("key = %q, want statements.file_kind_unknown", keyed.Key)
	}
}

func TestStatementKeys(t *testing.T) {
	cases := []struct {
		key  shared.MessageKey
		want shared.MessageKey
	}{
		{KeyBankNotDetected, "statements.bank_not_detected"},
		{KeyBankUnknown, "statements.bank_unknown"},
		{KeyAlreadyImported, "statements.already_imported"},
		{KeyNoCharges, "statements.no_charges"},
		{KeyBillingDateMissing, "statements.billing_date_missing"},
		{KeyFileEmpty, "statements.file_empty"},
		{KeyNotFound, "statements.not_found"},
		{KeyFileKindUnknown, "statements.file_kind_unknown"},
	}
	for _, tc := range cases {
		if tc.key != tc.want {
			t.Errorf("key = %q, want %q", tc.key, tc.want)
		}
	}
}
