package domain

import (
	"errors"
	"testing"
)

func TestParsedChargeToCharge_CopiesCardLast4(t *testing.T) {
	parsed := validParsedCharge()
	parsed.CardLast4 = "4821"

	got, err := parsed.ToCharge(validOrigin())

	if err != nil {
		t.Fatalf("ToCharge error = %v, want nil", err)
	}
	if got.CardLast4 != "4821" {
		t.Errorf("CardLast4 = %q, want 4821", got.CardLast4)
	}
}

func TestParsedChargeToCharge_EmptyCardLast4IsAllowed(t *testing.T) {
	got, err := validParsedCharge().ToCharge(validOrigin())

	if err != nil {
		t.Fatalf("ToCharge error = %v, want nil", err)
	}
	if got.CardLast4 != "" {
		t.Errorf("CardLast4 = %q, want empty", got.CardLast4)
	}
}

func TestCardLast4_RejectsAnythingButFourDigits(t *testing.T) {
	for _, last4 := range []string{"482", "48211", "48a1", " 4821"} {
		t.Run(last4, func(t *testing.T) {
			charge := validCharge()
			charge.CardLast4 = last4
			_, chargeErr := NewCharge(charge)

			parsed := validParsedCharge()
			parsed.CardLast4 = last4
			_, parsedErr := parsed.ToCharge(validOrigin())

			for name, err := range map[string]error{"NewCharge": chargeErr, "ToCharge": parsedErr} {
				var keyed KeyedError
				if !errors.Is(err, ErrInvalid) || !errors.As(err, &keyed) || keyed.Key != "common.card_last_digits_invalid" {
					t.Errorf("%s(%q) error = %v, want invalid with common.card_last_digits_invalid", name, last4, err)
				}
			}
		})
	}
}

func TestCardLastDigitsKey_FollowsFeatureReasonFormat(t *testing.T) {
	if KeyCardLastDigitsInvalid != "common.card_last_digits_invalid" || !featureReasonFormat.MatchString(string(KeyCardLastDigitsInvalid)) {
		t.Errorf("key = %q, want common.card_last_digits_invalid in <feature>.<reason> format", KeyCardLastDigitsInvalid)
	}
}
