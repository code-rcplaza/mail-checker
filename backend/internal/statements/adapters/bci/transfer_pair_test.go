package bci

import (
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

func TestTransferPair(t *testing.T) {
	parser := newParser(t)
	october := readFixture(t, "statement_october.txt")
	cases := []struct {
		name     string
		text     string
		want     domain.TransferPair
		wantPair bool
	}{
		{
			"both transfers give pesos of the national and dollars of the international",
			october,
			domain.TransferPair{Pesos: 9500000, Dollars: 10000, Date: time.Date(2026, time.October, 5, 0, 0, 0, 0, parser.location)},
			true,
		},
		{"a statement without transfers has no pair", withoutLinesContaining(october, "TRASPASO"), domain.TransferPair{}, false},
		{"a transfer of US$ 0,00 has no pair", readFixture(t, "statement_september.txt"), domain.TransferPair{}, false},
		{"only the national transfer has no pair", withoutLinesContaining(october, "TRASPASO DEUDA NACIONAL"), domain.TransferPair{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, hasPair := parser.TransferPair(tc.text)

			if hasPair != tc.wantPair || got != tc.want {
				t.Errorf("TransferPair = (%+v, %v), want (%+v, %v)", got, hasPair, tc.want, tc.wantPair)
			}
		})
	}
}
