package bci

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestReadInternationalColumns(t *testing.T) {
	raw := "3009 00000000000000000000002 29/09/26\tANTHROPIC* CLAUDE SUB\tANTHROPIC.COM\tUS\t20,00\t20,00"

	parser := newParser(t)

	got, outcome := parser.readInternationalColumns(raw)

	want := internationalPurchase{
		merchant: "ANTHROPIC* CLAUDE SUB",
		date:     time.Date(2026, time.September, 29, 0, 0, 0, 0, parser.location),
		origin:   2000,
		dollars:  2000,
	}
	if outcome != linePurchase || got != want {
		t.Errorf("readInternationalColumns = (%+v, %v), want (%+v, purchase)", got, outcome, want)
	}
}

func TestParseForeignAmount(t *testing.T) {
	cases := []struct {
		raw  string
		want shared.Money
	}{
		{"20,00", 2000},
		{"5.990,00", 599000},
		{"-30,00", -3000},
		{"US$ 100,00", 10000},
		{"US$0,00", 0},
		{"US$-55,40", -5540},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, isRead := parseForeignAmount(tc.raw)

			if !isRead || got != tc.want {
				t.Errorf("parseForeignAmount(%q) = (%d, %v), want (%d, true)", tc.raw, got, isRead, tc.want)
			}
		})
	}
}

func internationalOnly(lines []lineSummary) []lineSummary {
	return slices.DeleteFunc(slices.Clone(lines), func(l lineSummary) bool { return l.number <= nationalSectionEnd })
}

func TestParse_ReadsTheInternationalPurchasesInTheirCurrency(t *testing.T) {
	lines, bad := newParser(t).Parse(readFixture(t, "statement_october.txt"))

	want := []lineSummary{
		{106, "PLAYIT.GG* PLAYIT.GG/C", 300, shared.CurrencyUSD, "2026-10-01"},
		{107, "ANTHROPIC* CLAUDE SUB", 2000, shared.CurrencyUSD, "2026-09-29"},
		{108, "TIENDA EJEMPLO LTD", 600, shared.CurrencyUSD, "2026-10-01"},
		{110, "COMERCIO EJEMPLO PTE. LTD.", 2000, shared.CurrencyUSD, "2026-09-28"},
		{111, "APPLE.COM/BILL", 599000, shared.CurrencyCLP, "2026-09-20"},
	}
	if got := internationalOnly(summarize(lines)); !reflect.DeepEqual(got, want) {
		t.Errorf("international lines =\n%v\nwant\n%v", got, want)
	}
	if len(bad) != 0 {
		t.Errorf("bad lines = %v, want none: payments and transfers are skipped, not bad", bad)
	}
}

func internationalStatement(lines ...string) string {
	return strings.Join(append([]string{
		"ESTADO DE CUENTA INTERNACIONAL DE TARJETA DE CRÉDITO",
		"Nº DE TARJETA DE CRÉDITO\tXXXX XXXX XXXX 1234",
	}, lines...), "\n")
}

func TestParse_InternationalLineInEurosIsKeptInDollars(t *testing.T) {
	text := internationalStatement("0000 00000000000000000000001 05/10/26\tTIENDA EUROPA\tMADRID\tES\t15,00\t16,20")

	lines, _ := newParser(t).Parse(text)

	if len(lines) != 1 || lines[0].Charge.Currency != shared.CurrencyUSD || lines[0].Charge.Amount != 1620 {
		t.Errorf("lines = %+v, want one charge of US$ 16,20", lines)
	}
}

func TestParse_UnreadableInternationalLineGoesToBadLines(t *testing.T) {
	text := internationalStatement(
		"0000 00000000000000000000001 31/02/26\tSECRETO UNO\tCITY\tUS\t3,00\t3,00",
		"0000 00000000000000000000002 01/10/26\tSECRETO DOS\tCITY\tUS\t3,00\ttres",
		"0000 00000000000000000000003 01/10/26\tSECRETO TRES",
		"0610\t06/10/26\tMONTO CANCELADO\tCL\t0,00\t-30,00",
	)

	lines, bad := newParser(t).Parse(text)

	if len(lines) != 0 || !reflect.DeepEqual(bad, []int{3, 4, 5}) {
		t.Errorf("lines = %+v, bad = %v; want no charges and bad [3 4 5]", lines, bad)
	}
}
