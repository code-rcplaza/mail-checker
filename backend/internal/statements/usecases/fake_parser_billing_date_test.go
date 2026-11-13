package usecases

import (
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

func TestFakeStatementParser_BillingDate(t *testing.T) {
	cases := []struct {
		name     string
		parser   ports.BankStatementParser
		want     time.Time
		wantDate bool
	}{
		{"parser with a date returns it", &fakeStatementParser{bank: "bci", billingDate: billingDate}, billingDate, true},
		{"parser without a date has none", &fakeStatementParser{bank: "bci"}, time.Time{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasDate := tc.parser.BillingDate("statement text")

			if hasDate != tc.wantDate || !got.Equal(tc.want) {
				t.Errorf("BillingDate = (%v, %v), want (%v, %v)", got, hasDate, tc.want, tc.wantDate)
			}
		})
	}
}

func TestFakeStatementParser_TransferPair(t *testing.T) {
	pair := domain.TransferPair{Pesos: 9500000, Dollars: 10000, Date: billingDate}
	cases := []struct {
		name     string
		parser   ports.BankStatementParser
		want     domain.TransferPair
		wantPair bool
	}{
		{"a bank without the pair has none", &fakeStatementParser{bank: "other-bank"}, domain.TransferPair{}, false},
		{"a parser with the pair gives pesos, dollars and date", &fakeStatementParser{bank: "bci", transferPair: &pair}, pair, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasPair := tc.parser.TransferPair("statement text")

			if hasPair != tc.wantPair || got != tc.want {
				t.Errorf("TransferPair = (%+v, %v), want (%+v, %v)", got, hasPair, tc.want, tc.wantPair)
			}
		})
	}
}
