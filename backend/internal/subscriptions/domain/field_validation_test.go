package domain

import (
	"errors"
	"fmt"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func fieldsOf(t *testing.T, err error) string {
	t.Helper()
	var keyed shared.KeyedError
	if err == nil {
		return "none"
	}
	if !errors.As(err, &keyed) || keyed.Key != shared.KeyValidationFailed {
		t.Fatalf("error = %v, want a validation failure with fields", err)
	}
	return fmt.Sprint(keyed.Fields)
}

func TestValidateManualFields(t *testing.T) {
	cases := []struct {
		name  string
		input ManualSubscriptionInput
		want  string
	}{
		{"a valid manual", ManualSubscriptionInput{Name: "Gimnasio", Amount: 2500000, Currency: "CLP", BillingDay: 5}, "none"},
		{"every field wrong", ManualSubscriptionInput{Name: " ", Amount: 0, Currency: "EUR", BillingDay: 32},
			fmt.Sprint(map[string]shared.MessageKey{"amount": KeyManualAmountNotPositive, "billingDay": KeyManualDayOutOfRange, "currency": shared.KeyCurrencyUnknown, "name": KeyManualNameMissing})},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := fieldsOf(t, ValidateManualFields(tc.input)); got != tc.want {
				t.Errorf("fields = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestValidateOverrideFields(t *testing.T) {
	zero, day := shared.Money(0), 32
	amount, goodDay := shared.Money(900000), 31
	cases := []struct {
		name   string
		amount *shared.Money
		day    *int
		want   string
	}{
		{"both valid", &amount, &goodDay, "none"},
		{"neither given", nil, nil, "none"},
		{"both wrong", &zero, &day, fmt.Sprint(map[string]shared.MessageKey{"amount": KeyOverrideAmountNotPositive, "billingDay": KeyOverrideDayOutOfRange})},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := fieldsOf(t, ValidateOverrideFields(tc.amount, tc.day)); got != tc.want {
				t.Errorf("fields = %s, want %s", got, tc.want)
			}
		})
	}
}
