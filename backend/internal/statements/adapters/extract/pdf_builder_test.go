package extract

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Test PDFs are built here instead of committing a binary: the repo is public
// and a real statement must never be in it. Like the real bank PDFs, every
// glyph is its own text operation, so the extractor has to rebuild words and
// columns from positions.

const (
	builderFont     = "Helvetica"
	builderFontSize = 9
	builderLineStep = 12
	builderTopY     = 760
)

// pdfCell is text that starts at a column position on its line.
type pdfCell struct {
	x    float64
	text string
}

// pdfPage holds lines from top to bottom.
type pdfPage [][]pdfCell

func buildPDF(pages ...pdfPage) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		pagesObject(len(pages)),
		"<< /Type /Font /Subtype /Type1 /BaseFont /" + builderFont + " >>",
	}
	for i, page := range pages {
		content := contentStream(page)
		contentID := 4 + 2*i + 1
		objects = append(objects,
			fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 3 0 R >> >> /Contents %d 0 R >>", contentID),
			fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content),
		)
	}
	return writeObjects(objects)
}

func pagesObject(pageCount int) string {
	kids := make([]string, 0, pageCount)
	for i := range pageCount {
		kids = append(kids, fmt.Sprintf("%d 0 R", 4+2*i))
	}
	return fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), pageCount)
}

func contentStream(page pdfPage) string {
	var stream strings.Builder
	for row, cells := range page {
		y := float64(builderTopY - row*builderLineStep)
		for _, cell := range cells {
			writeGlyphs(&stream, cell, y)
		}
	}
	return stream.String()
}

// writeGlyphs emits one operation per glyph; a space only moves the position.
func writeGlyphs(stream *strings.Builder, cell pdfCell, y float64) {
	x := cell.x
	for _, glyph := range cell.text {
		if glyph != ' ' {
			fmt.Fprintf(stream, "BT /F1 %d Tf 1 0 0 1 %.1f %.1f Tm (%s) Tj ET\n", builderFontSize, x, y, escapePDFText(string(glyph)))
		}
		x += glyphAdvance(pdf.Text{S: string(glyph), Font: builderFont, FontSize: builderFontSize})
	}
}

func escapePDFText(text string) string {
	return strings.NewReplacer(`\`, `\\`, "(", `\(`, ")", `\)`).Replace(text)
}

func writeObjects(objects []string) []byte {
	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n")
	offsets := make([]int, 0, len(objects))
	for i, object := range objects {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", i+1, object)
	}
	xrefAt := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets {
		fmt.Fprintf(&out, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefAt)
	return out.Bytes()
}
