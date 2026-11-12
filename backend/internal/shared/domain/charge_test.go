package domain

import (
	"reflect"
	"testing"
	"time"
)

func validCharge() Charge {
	accountID := int64(7)
	return Charge{
		Bank:       "bci",
		Source:     SourceEmail,
		ExternalID: "<msg-1@example.com>",
		AccountID:  &accountID,
		OccurredAt: time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC),
		Amount:     599000,
		Currency:   CurrencyCLP,
		Merchant:   "STREAMING CO",
		Card:       "credito",
	}
}

func TestNewCharge_ValidDataReturnsChargeWithThoseValues(t *testing.T) {
	want := validCharge()

	got, err := NewCharge(want)

	if err != nil {
		t.Fatalf("NewCharge error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewCharge = %+v, want %+v", got, want)
	}
}

func TestNewCharge_RejectsInvalidData(t *testing.T) {
	cases := []struct {
		name   string
		modify func(c Charge) Charge
		want   MessageKey
	}{
		{"negative amount", func(c Charge) Charge { c.Amount = -1; return c }, "common.amount_negative"},
		{"zero date", func(c Charge) Charge { c.OccurredAt = time.Time{}; return c }, "common.date_missing"},
		{"empty merchant", func(c Charge) Charge { c.Merchant = ""; return c }, "common.merchant_missing"},
		{"blank merchant", func(c Charge) Charge { c.Merchant = "  \t "; return c }, "common.merchant_missing"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCharge(tc.modify(validCharge()))

			assertInvalidWithKey(t, err, tc.want)
		})
	}
}
