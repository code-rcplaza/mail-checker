// Package extract turns an uploaded statement file into text.
package extract

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/ledongthuc/pdf"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

// The bank PDFs draw one glyph per operation, so words and columns are
// rebuilt from positions: what counts is the room left after the previous
// glyph's own width, measured in font sizes.
const (
	sameLineTolerance  = 2.0  // points of vertical drift still on the same line
	wordGapFontSizes   = 0.12 // a space is about 0.23 font sizes wide
	columnGapFontSizes = 1.0  // a wider gap starts a new column, written as a tab
)

var pdfHeader = []byte("%PDF-")

// PDFExtractor reads the text of a PDF statement: one line per text line of
// the page, words separated by a space and columns by a tab.
type PDFExtractor struct{}

func (PDFExtractor) Extract(_ string, data []byte) (string, domain.FileKind, error) {
	if err := validatePDF(data); err != nil {
		return "", "", err
	}
	text, err := panicAsError(func() (string, error) { return readPDFText(data) })
	if err != nil {
		return "", "", fmt.Errorf("read pdf: %w", err)
	}
	return text, domain.FileKindPDF, nil
}

func validatePDF(data []byte) error {
	if len(data) == 0 {
		return shared.Invalid(domain.KeyFileEmpty, nil)
	}
	if !bytes.HasPrefix(data, pdfHeader) {
		return shared.Invalid(domain.KeyFileKindUnknown, nil)
	}
	return nil
}

// panicAsError contains the PDF library, which panics on some malformed
// files instead of returning an error.
func panicAsError(read func() (string, error)) (text string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("malformed pdf: %v", recovered)
		}
	}()
	return read()
}

func readPDFText(data []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	var lines []string
	for number := 1; number <= reader.NumPage(); number++ {
		page := reader.Page(number)
		if page.V.IsNull() {
			continue
		}
		lines = append(lines, linesOfGlyphs(sortTopDown(page.Content().Text))...)
	}
	return strings.Join(lines, "\n"), nil
}

// sortTopDown orders glyphs by line from the top of the page, then left to
// right; the slice the library returned is left untouched.
func sortTopDown(glyphs []pdf.Text) []pdf.Text {
	sorted := slices.Clone(glyphs)
	slices.SortStableFunc(sorted, func(a, b pdf.Text) int {
		if !isSameLine(a, b) {
			return compareFloats(b.Y, a.Y)
		}
		return compareFloats(a.X, b.X)
	})
	return sorted
}

func linesOfGlyphs(sorted []pdf.Text) []string {
	var lines []string
	var line strings.Builder
	for i, glyph := range sorted {
		if i > 0 && !isSameLine(sorted[i-1], glyph) {
			lines = appendNonBlank(lines, line.String())
			line.Reset()
		}
		if i > 0 && isSameLine(sorted[i-1], glyph) {
			line.WriteString(separatorBetween(sorted[i-1], glyph))
		}
		line.WriteString(glyph.S)
	}
	return appendNonBlank(lines, line.String())
}

func separatorBetween(previous, glyph pdf.Text) string {
	gap := glyph.X - previous.X - glyphAdvance(previous.S, previous.Font, previous.FontSize)
	switch {
	case gap > glyph.FontSize*columnGapFontSizes:
		return "\t"
	case gap > glyph.FontSize*wordGapFontSizes:
		return " "
	}
	return ""
}

func isSameLine(a, b pdf.Text) bool {
	return math.Abs(a.Y-b.Y) <= sameLineTolerance
}

func appendNonBlank(lines []string, line string) []string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return lines
	}
	return append(lines, trimmed)
}

func compareFloats(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}
