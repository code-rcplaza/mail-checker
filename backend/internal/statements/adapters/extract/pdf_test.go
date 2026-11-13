package extract

import (
	"errors"
	"strings"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

var _ ports.FileExtractor = PDFExtractor{}

// Lines shaped like an anonymized statement: a label and its value in two
// columns, and a purchase with reference, descriptor and amounts.
var (
	statementDateLine = []pdfCell{{40, "FECHA ESTADO DE CUENTA"}, {300, "21/10/2026"}}
	purchaseLine      = []pdfCell{{40, "26/09/26 280900000001 NETFLIX.COM"}, {330, "$12.990"}, {400, "$12.990"}}
	secondPageLine    = []pdfCell{{40, "TOTAL TARJETA XXXXXXXXXXXX0000"}, {330, "$12.990"}}
)

func TestExtract_ReadsTheTextOfAPDF(t *testing.T) {
	cases := []struct {
		name  string
		pages []pdfPage
		want  string
	}{
		{
			"lines in page order with words and columns",
			[]pdfPage{{statementDateLine, purchaseLine}},
			"FECHA ESTADO DE CUENTA\t21/10/2026\n26/09/26 280900000001 NETFLIX.COM\t$12.990\t$12.990",
		},
		{
			"pages in order",
			[]pdfPage{{statementDateLine}, {secondPageLine}},
			"FECHA ESTADO DE CUENTA\t21/10/2026\nTOTAL TARJETA XXXXXXXXXXXX0000\t$12.990",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			text, kind, err := PDFExtractor{}.Extract("cartola.pdf", buildPDF(tc.pages...))

			if err != nil {
				t.Fatalf("Extract error = %v, want nil", err)
			}
			if kind != domain.FileKindPDF {
				t.Errorf("kind = %q, want pdf", kind)
			}
			if text != tc.want {
				t.Errorf("text =\n%q\nwant\n%q", text, tc.want)
			}
		})
	}
}

func TestExtract_RejectsWhatIsNotAPDF(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want shared.MessageKey
	}{
		{"a text file", []byte("FECHA ESTADO DE CUENTA 21/10/2026"), domain.KeyFileKindUnknown},
		{"an empty file", nil, domain.KeyFileEmpty},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := PDFExtractor{}.Extract("cartola.pdf", tc.data)

			var keyed shared.KeyedError
			if !errors.As(err, &keyed) || !errors.Is(err, shared.ErrInvalid) || keyed.Key != tc.want {
				t.Errorf("Extract error = %v, want ErrInvalid with %q", err, tc.want)
			}
		})
	}
}

func TestExtract_DamagedPDFIsAnErrorWithContext(t *testing.T) {
	whole := buildPDF(pdfPage{statementDateLine, purchaseLine})
	cases := []struct {
		name string
		data []byte
	}{
		{"truncated", whole[:len(whole)/2]},
		{"header only", []byte("%PDF-1.4\n")},
		{"garbage after the header", append([]byte("%PDF-1.4\n"), strings.Repeat("\x00\xff obj stream", 50)...)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := PDFExtractor{}.Extract("cartola.pdf", tc.data)

			if err == nil || !strings.HasPrefix(err.Error(), "read pdf: ") {
				t.Errorf("Extract error = %v, want one starting with \"read pdf: \"", err)
			}
		})
	}
}

func TestPanicAsError_TurnsALibraryPanicIntoAnError(t *testing.T) {
	_, err := panicAsError(func() (string, error) { panic("malformed xref") })

	if err == nil || !strings.Contains(err.Error(), "malformed xref") {
		t.Errorf("error = %v, want one that carries the panic", err)
	}
}
