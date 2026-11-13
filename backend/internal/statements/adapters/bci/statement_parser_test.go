package bci

import (
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

var _ ports.BankStatementParser = StatementParser{}

func santiago(t *testing.T) *time.Location {
	t.Helper()
	location, err := time.LoadLocation("America/Santiago")
	if err != nil {
		t.Fatalf("load America/Santiago: %v", err)
	}
	return location
}

func newParser(t *testing.T) StatementParser {
	t.Helper()
	return NewStatementParser(santiago(t))
}

// readFixture loads an anonymized statement, as the PDF extractor writes it.
func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}

func withoutLinesContaining(text, needle string) string {
	lines := strings.Split(text, "\n")
	kept := slices.DeleteFunc(slices.Clone(lines), func(line string) bool { return strings.Contains(line, needle) })
	return strings.Join(kept, "\n")
}

func TestDetect(t *testing.T) {
	october := readFixture(t, "statement_october.txt")
	cases := []struct {
		name string
		text string
		want bool
	}{
		{"a BCI statement", october, true},
		{"another BCI statement", readFixture(t, "statement_september.txt"), true},
		{"the same CMF format without the BCI blocking phone", withoutLinesContaining(october, "800 201 090"), false},
		{"the BCI phone in another document", "Para bloquear su tarjeta llame al\n800 201 090 PARA SOLICITAR EL BLOQUEO", false},
		{"another document", "BOLETA ELECTRÓNICA\nTOTAL $12.990", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := newParser(t).Detect(tc.text); got != tc.want {
				t.Errorf("Detect = %v, want %v", got, tc.want)
			}
		})
	}
}

// lineSummary keeps what a test compares of a parsed line.
type lineSummary struct {
	number   int
	merchant string
	amount   shared.Money
	currency shared.Currency
	day      string
}

func summarize(lines []domain.StatementLine) []lineSummary {
	summaries := make([]lineSummary, 0, len(lines))
	for _, line := range lines {
		summaries = append(summaries, lineSummary{
			number:   line.Number,
			merchant: line.Charge.Merchant,
			amount:   line.Charge.Amount,
			currency: line.Charge.Currency,
			day:      line.Charge.Date.Format(time.DateOnly),
		})
	}
	return summaries
}

func nationalOnly(lines []lineSummary) []lineSummary {
	return slices.DeleteFunc(slices.Clone(lines), func(l lineSummary) bool { return l.number > nationalSectionEnd })
}

// nationalSectionEnd is the last line of the national statement in the October fixture.
const nationalSectionEnd = 87

func TestParse_ReadsTheNationalPurchases(t *testing.T) {
	lines, bad := newParser(t).Parse(readFixture(t, "statement_october.txt"))

	want := []lineSummary{
		{30, "PAYU *UBER TRIP", 400000, shared.CurrencyCLP, "2026-09-19"},
		{31, "PedidosYa*RESTAURANTE", 900000, shared.CurrencyCLP, "2026-09-22"},
		{32, "PedidosYa*Propina", 50000, shared.CurrencyCLP, "2026-09-22"},
		{33, "NETFLIX.COM", 1299000, shared.CurrencyCLP, "2026-09-26"},
		{34, "Disney Plus", 1599000, shared.CurrencyCLP, "2026-09-27"},
		{49, "PAYU *UBER TRIP", 300000, shared.CurrencyCLP, "2026-09-28"},
		{50, "APPLE.COM/BILL", 469000, shared.CurrencyCLP, "2026-10-01"},
		{51, "PedidosYa*Plus", 199500, shared.CurrencyCLP, "2026-10-03"},
		{52, "EA *ELECTRONIC ARTS help.ea.com", 1390000, shared.CurrencyCLP, "2026-10-04"},
		{53, "SEGUROS EJEMPLO S.A.", 1000000, shared.CurrencyCLP, "2026-10-06"},
		{54, "MICROSOFT. SUBSCRIPTIONLAS CONDES", 539900, shared.CurrencyCLP, "2026-10-12"},
		{55, "HOSTGATOR", 870000, shared.CurrencyCLP, "2026-10-14"},
		{56, "DL*GOOGLE YOUTUBE", 2170000, shared.CurrencyCLP, "2026-10-15"},
		{57, "Spotify P00000000A", 825000, shared.CurrencyCLP, "2026-10-15"},
		{58, "DL*GOOGLE YOUTUBE", 1100000, shared.CurrencyCLP, "2026-10-16"},
	}
	if got := nationalOnly(summarize(lines)); !reflect.DeepEqual(got, want) {
		t.Errorf("national lines =\n%v\nwant\n%v", got, want)
	}
	if len(bad) != 0 {
		t.Errorf("bad lines = %v, want none", bad)
	}
}

func TestParse_PurchaseCarriesCardAndDateAtMidnightInSantiago(t *testing.T) {
	lines, _ := newParser(t).Parse(readFixture(t, "statement_october.txt"))

	index := slices.IndexFunc(lines, func(l domain.StatementLine) bool { return l.Number == 33 })
	if index < 0 {
		t.Fatalf("line 33 was not parsed")
	}
	netflix := lines[index].Charge

	want := shared.ParsedCharge{
		Amount:    1299000,
		Currency:  shared.CurrencyCLP,
		Merchant:  "NETFLIX.COM",
		Card:      CardCredit,
		CardLast4: "1234",
		Date:      time.Date(2026, time.September, 26, 0, 0, 0, 0, santiago(t)),
	}
	if !reflect.DeepEqual(netflix, want) {
		t.Errorf("charge = %+v, want %+v", netflix, want)
	}
}

func TestParseCLPAmount(t *testing.T) {
	cases := []struct {
		raw  string
		want shared.Money
	}{
		{"$5.990", 599000},
		{"$480", 48000},
		{"$1.503.960", 150396000},
		{"$ 4.200.000", 420000000},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, isRead := parseCLPAmount(tc.raw)

			if !isRead || got != tc.want {
				t.Errorf("parseCLPAmount(%q) = (%d, %v), want (%d, true)", tc.raw, got, isRead, tc.want)
			}
		})
	}
}

func TestParse_SkipsWhatIsNotAPurchase(t *testing.T) {
	text := readFixture(t, "statement_october.txt")
	notPurchases := []string{"MONTO CANCELADO", "PAGO PAC", "TRASPASO DEUDA", "IMPUESTO", "COMISION", "INTERESES", "TASA INT."}

	lines, bad := newParser(t).Parse(text)

	for _, line := range lines {
		for _, word := range notPurchases {
			if strings.Contains(line.Charge.Merchant, word) || strings.Contains(textLine(text, line.Number), word) {
				t.Errorf("line %d (%s) parsed as a purchase", line.Number, line.Charge.Merchant)
			}
		}
	}
	if len(bad) != 0 {
		t.Errorf("bad lines = %v, want none: skipped lines are not bad", bad)
	}
}

func textLine(text string, number int) string {
	return strings.Split(text, "\n")[number-1]
}

func TestParse_UnreadablePurchaseGoesToBadLines(t *testing.T) {
	text := strings.Join([]string{
		"ESTADO DE CUENTA NACIONAL DE TARJETA DE CRÉDITO",
		"Nº DE TARJETA DE CRÉDITO\tXXXX XXXX XXXX 1234",
		"1.TOTAL OPERACIONES\t(B)\t$0",
		"31/02/26 0000 00000001 NETFLIX.COM\t$12.990\t$12.990\t01/01\t$12.990",
		"26/09/26 0000 00000002 SECRETO COMERCIO\t$12,99O\t$12.990\t01/01\t$12.990",
		"27/09/26 0000 00000003 HOSTGATOR\t$8.700\t$8.700\t01/01\t$8.700",
	}, "\n")

	lines, bad := newParser(t).Parse(text)

	if !reflect.DeepEqual(bad, []int{4, 5}) {
		t.Errorf("bad lines = %v, want [4 5]", bad)
	}
	for _, line := range lines {
		if strings.Contains(line.Charge.Merchant, "SECRETO") || line.Number != 6 {
			t.Errorf("parsed line %d %q, want only line 6", line.Number, line.Charge.Merchant)
		}
	}
}
