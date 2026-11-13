package domain

import (
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var (
	emailDate  = time.Date(2026, 10, 6, 9, 30, 0, 0, time.UTC)
	parsedDate = time.Date(2026, 10, 5, 21, 15, 0, 0, time.UTC)
)

func purchaseEmail() RawEmail {
	return RawEmail{
		MessageID:  "<msg-1@example.com>",
		Subject:    "Compra con tu tarjeta",
		From:       "avisos@bank.example",
		ReceivedAt: emailDate,
		Body:       "Compra por $5.990 en STREAMING CO",
	}
}

func parsedPurchase(date time.Time) shared.ParsedCharge {
	return shared.ParsedCharge{Amount: 599000, Currency: shared.CurrencyCLP, Merchant: "STREAMING CO", Card: "credito", Date: date}
}

func TestChargeFromEmail_WithoutParsedDateUsesEmailValues(t *testing.T) {
	got, err := ChargeFromEmail(parsedPurchase(time.Time{}), purchaseEmail(), 1, 7, "bci")

	if err != nil {
		t.Fatalf("ChargeFromEmail error = %v, want nil", err)
	}
	if got.Source != shared.SourceEmail {
		t.Errorf("Source = %q, want %q", got.Source, shared.SourceEmail)
	}
	if got.ExternalID != "<msg-1@example.com>" {
		t.Errorf("ExternalID = %q, want the Message-ID", got.ExternalID)
	}
	if got.AccountID == nil || *got.AccountID != 7 {
		t.Errorf("AccountID = %v, want 7", got.AccountID)
	}
	if got.Bank != "bci" {
		t.Errorf("Bank = %q, want bci", got.Bank)
	}
	if !got.OccurredAt.Equal(emailDate) {
		t.Errorf("OccurredAt = %v, want the email date %v", got.OccurredAt, emailDate)
	}
}

func TestChargeFromEmail_KeepsParsedDate(t *testing.T) {
	got, err := ChargeFromEmail(parsedPurchase(parsedDate), purchaseEmail(), 1, 7, "bci")

	if err != nil {
		t.Fatalf("ChargeFromEmail error = %v, want nil", err)
	}
	if !got.OccurredAt.Equal(parsedDate) {
		t.Errorf("OccurredAt = %v, want the parsed date %v", got.OccurredAt, parsedDate)
	}
}
