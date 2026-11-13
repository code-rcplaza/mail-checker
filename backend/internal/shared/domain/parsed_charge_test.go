package domain

import (
	"reflect"
	"testing"
	"time"
)

func validParsedCharge() ParsedCharge {
	return ParsedCharge{
		Amount:   599000,
		Currency: CurrencyCLP,
		Merchant: "STREAMING CO",
		Card:     "credito",
		Date:     time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC),
	}
}

func validOrigin() ChargeOrigin {
	accountID := int64(7)
	return ChargeOrigin{
		UserID:       1,
		Bank:         "bci",
		Source:       SourceEmail,
		ExternalID:   "<msg-1@example.com>",
		AccountID:    &accountID,
		FallbackDate: time.Date(2026, 10, 6, 9, 30, 0, 0, time.UTC),
	}
}

func TestParsedChargeToCharge_CombinesParsedValuesAndOrigin(t *testing.T) {
	parsed, origin := validParsedCharge(), validOrigin()

	got, err := parsed.ToCharge(origin)

	if err != nil {
		t.Fatalf("ToCharge error = %v, want nil", err)
	}
	want := Charge{
		UserID:     origin.UserID,
		Bank:       origin.Bank,
		Source:     origin.Source,
		ExternalID: origin.ExternalID,
		AccountID:  origin.AccountID,
		OccurredAt: parsed.Date,
		Amount:     parsed.Amount,
		Currency:   parsed.Currency,
		Merchant:   parsed.Merchant,
		Card:       parsed.Card,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToCharge = %+v, want %+v", got, want)
	}
}

func TestParsedChargeToCharge_ZeroDateUsesOriginFallbackDate(t *testing.T) {
	parsed, origin := validParsedCharge(), validOrigin()
	parsed.Date = time.Time{}

	got, err := parsed.ToCharge(origin)

	if err != nil {
		t.Fatalf("ToCharge error = %v, want nil", err)
	}
	if !got.OccurredAt.Equal(origin.FallbackDate) {
		t.Errorf("OccurredAt = %v, want fallback %v", got.OccurredAt, origin.FallbackDate)
	}
}

func TestParsedChargeToCharge_RejectsMissingExternalID(t *testing.T) {
	for _, externalID := range []string{"", "   "} {
		t.Run("external id="+externalID, func(t *testing.T) {
			origin := validOrigin()
			origin.ExternalID = externalID

			_, err := validParsedCharge().ToCharge(origin)

			assertInvalidWithKey(t, err, "common.external_id_missing")
		})
	}
}

func TestParsedChargeToCharge_RejectsNegativeAmountLikeNewCharge(t *testing.T) {
	parsed := validParsedCharge()
	parsed.Amount = -1

	_, err := parsed.ToCharge(validOrigin())

	assertInvalidWithKey(t, err, "common.amount_negative")
}
