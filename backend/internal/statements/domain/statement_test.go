package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestStatementKeys(t *testing.T) {
	cases := []struct {
		key  shared.MessageKey
		want shared.MessageKey
	}{
		{KeyBankNotDetected, "statements.bank_not_detected"},
		{KeyBankUnknown, "statements.bank_unknown"},
		{KeyAlreadyImported, "statements.already_imported"},
		{KeyNoCharges, "statements.no_charges"},
		{KeyBillingDateRequired, "statements.billing_date_required"},
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
