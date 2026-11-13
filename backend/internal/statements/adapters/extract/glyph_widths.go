package extract

import (
	"strings"
	"unicode/utf8"
)

// The bank PDFs give no glyph widths, so the space between two glyphs is
// measured against the width the previous glyph should have. Statements use
// Arial (or Arial Narrow), whose widths match Helvetica's, in thousandths of
// the font size. A glyph missing from the table counts as a digit.
var helveticaWidths = map[rune]int{
	' ': 278, '!': 278, '"': 355, '#': 556, '$': 556, '%': 889, '&': 667, '\'': 191,
	'(': 333, ')': 333, '*': 389, '+': 584, ',': 278, '-': 333, '.': 278, '/': 278,
	':': 278, ';': 278, '<': 584, '=': 584, '>': 584, '?': 556, '@': 1015, '_': 556,
	'A': 667, 'B': 667, 'C': 722, 'D': 722, 'E': 667, 'F': 611, 'G': 778, 'H': 722,
	'I': 278, 'J': 500, 'K': 667, 'L': 556, 'M': 833, 'N': 722, 'O': 778, 'P': 667,
	'Q': 778, 'R': 722, 'S': 667, 'T': 611, 'U': 722, 'V': 667, 'W': 944, 'X': 667,
	'Y': 667, 'Z': 611,
	'a': 556, 'b': 556, 'c': 500, 'd': 556, 'e': 556, 'f': 278, 'g': 556, 'h': 556,
	'i': 222, 'j': 222, 'k': 500, 'l': 222, 'm': 833, 'n': 556, 'o': 556, 'p': 556,
	'q': 556, 'r': 333, 's': 500, 't': 278, 'u': 556, 'v': 500, 'w': 722, 'x': 500,
	'y': 500, 'z': 500,
	'Á': 667, 'É': 667, 'Í': 278, 'Ó': 778, 'Ú': 722, 'Ñ': 722, 'Ü': 722, 'º': 365, '°': 400,
	'á': 556, 'é': 556, 'í': 278, 'ó': 556, 'ú': 556, 'ñ': 556, 'ü': 556,
}

const (
	defaultGlyphWidth = 556  // a digit
	narrowFontScale   = 0.82 // Arial Narrow is Arial condensed to 82 %
	thousandths       = 1000.0
)

// glyphAdvance is how far a glyph should move the pen at its font and size.
func glyphAdvance(glyph, font string, fontSize float64) float64 {
	return float64(widthOf(glyph)) / thousandths * fontSize * fontScale(font)
}

func widthOf(glyph string) int {
	first, _ := utf8.DecodeRuneInString(glyph)
	if width, isKnown := helveticaWidths[first]; isKnown {
		return width
	}
	return defaultGlyphWidth
}

func fontScale(font string) float64 {
	if strings.Contains(font, "Narrow") {
		return narrowFontScale
	}
	return 1
}
