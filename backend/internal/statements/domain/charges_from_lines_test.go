package domain

import (
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var billingDate = time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)

func line(number int, merchant string, amount shared.Money, date time.Time) StatementLine {
	return StatementLine{
		Number: number,
		Charge: shared.ParsedCharge{Amount: amount, Currency: shared.CurrencyCLP, Merchant: merchant, Card: "credito", Date: date},
	}
}

func TestChargesFromLines_DatedLinesBecomeStatementCharges(t *testing.T) {
	first := time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)
	second := time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)
	lines := []StatementLine{line(1, "STREAMING CO", 599000, first), line(2, "MUSIC APP", 399000, second)}

	charges, badLines := ChargesFromLines(lines, 1, "bci", "stmt-1", billingDate)

	if len(badLines) != 0 {
		t.Errorf("bad lines = %v, want none", badLines)
	}
	if len(charges) != 2 {
		t.Fatalf("charges = %d, want 2", len(charges))
	}
	wantIDs := []string{"stmt-1:1", "stmt-1:2"}
	wantDates := []time.Time{first, second}
	for i, c := range charges {
		if c.Source != shared.SourceStatement || c.Bank != "bci" {
			t.Errorf("charge %d source/bank = %q/%q, want statement/bci", i, c.Source, c.Bank)
		}
		if c.ExternalID != wantIDs[i] {
			t.Errorf("charge %d external id = %q, want %q", i, c.ExternalID, wantIDs[i])
		}
		if !c.OccurredAt.Equal(wantDates[i]) {
			t.Errorf("charge %d date = %v, want its own %v", i, c.OccurredAt, wantDates[i])
		}
	}
}

func TestChargesFromLines_UndatedLineTakesBillingDate(t *testing.T) {
	charges, _ := ChargesFromLines([]StatementLine{line(1, "STREAMING CO", 599000, time.Time{})}, 1, "bci", "stmt-1", billingDate)

	if len(charges) != 1 || !charges[0].OccurredAt.Equal(billingDate) {
		t.Errorf("charges = %+v, want one dated %v", charges, billingDate)
	}
}

func TestChargesFromLines_InvalidLineIsReportedAndOthersConvert(t *testing.T) {
	date := time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)
	lines := []StatementLine{line(1, "STREAMING CO", 599000, date), line(2, "", 100000, date), line(3, "MUSIC APP", 399000, date)}

	charges, badLines := ChargesFromLines(lines, 1, "bci", "stmt-1", billingDate)

	if !reflect.DeepEqual(badLines, []int{2}) {
		t.Errorf("bad lines = %v, want [2]", badLines)
	}
	ids := make([]string, 0, len(charges))
	for _, c := range charges {
		ids = append(ids, c.ExternalID)
	}
	if !reflect.DeepEqual(ids, []string{"stmt-1:1", "stmt-1:3"}) {
		t.Errorf("charge ids = %v, want [stmt-1:1 stmt-1:3]", ids)
	}
}
