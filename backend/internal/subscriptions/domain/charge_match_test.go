package domain

import (
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func statementCharge(externalID, merchant string, amount shared.Money, day int) shared.Charge {
	c := charge(externalID, merchant, amount, shared.CurrencyCLP, day)
	c.Source = shared.SourceStatement
	return c
}

func emailCharge(externalID, merchant string, amount shared.Money, day int) shared.Charge {
	return charge(externalID, merchant, amount, shared.CurrencyCLP, day)
}

func TestChargesMatch_SameBankServiceAndAmountWithin20Days(t *testing.T) {
	statement := statementCharge("s", "STRMNG*CO", 599000, 21)
	email := emailCharge("e", "STREAMING CO", 599000, 1)

	if !ChargesMatch(statement, email, identifyFake, normalizeFake) {
		t.Errorf("ChargesMatch at 20 days = false, want true")
	}
}

func TestChargesMatch_DoesNotMatch(t *testing.T) {
	cases := []struct {
		name   string
		modify func(email shared.Charge) shared.Charge
	}{
		{"36 days apart", func(e shared.Charge) shared.Charge { e.OccurredAt = e.OccurredAt.AddDate(0, 0, -16); return e }},
		{"different bank", func(e shared.Charge) shared.Charge { e.Bank = "falabella"; return e }},
		{"different amount", func(e shared.Charge) shared.Charge { e.Amount = 799000; return e }},
		{"different service", func(e shared.Charge) shared.Charge { e.Merchant = "MUSIC APP"; return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			statement := statementCharge("s", "STRMNG*CO", 599000, 21)
			email := emailCharge("e", "STREAMING CO", 599000, 1)
			if !ChargesMatch(statement, email, identifyFake, normalizeFake) {
				t.Fatalf("baseline pair does not match; the case would prove nothing")
			}

			if ChargesMatch(statement, tc.modify(email), identifyFake, normalizeFake) {
				t.Errorf("ChargesMatch = true, want false")
			}
		})
	}
}

func externalIDs(charges []shared.Charge) []string {
	ids := make([]string, 0, len(charges))
	for _, c := range charges {
		ids = append(ids, c.ExternalID)
	}
	return ids
}

func TestDropStatementChargesMatchingEmail(t *testing.T) {
	cases := []struct {
		name    string
		charges []shared.Charge
		want    []string
	}{
		{
			name: "keeps the email charge and the unpaired statement charge",
			charges: []shared.Charge{
				emailCharge("email", "STREAMING CO", 599000, 1),
				statementCharge("paired", "STRMNG*CO", 599000, 3),
				statementCharge("unpaired", "MUSIC APP", 399000, 3),
			},
			want: []string{"email", "unpaired"},
		},
		{
			name: "pairs one to one: one email drops only one statement charge",
			charges: []shared.Charge{
				emailCharge("email", "STREAMING CO", 599000, 1),
				statementCharge("first", "STRMNG*CO", 599000, 3),
				statementCharge("second", "STRMNG*CO", 599000, 3),
			},
			want: []string{"email", "second"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := externalIDs(DropStatementChargesMatchingEmail(tc.charges, identifyFake, normalizeFake))

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("charges = %v, want %v", got, tc.want)
			}
		})
	}
}
