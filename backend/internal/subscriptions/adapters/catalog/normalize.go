package catalog

import (
	"strings"
	"unicode"
)

// paymentProcessorPrefixes is a closed list seen in real statements. Any
// other "X*" is not stripped: in "PEDIDOSYA*PLUS" or "ANTHROPIC* CLAUDE" what
// comes before the star is the merchant.
var paymentProcessorPrefixes = []string{
	"PAYU *",
	"DL*",
	"MP *",
	"MERCADOPAGO*",
	"DP *",
}

// accentFolds maps the uppercase Spanish letters with a mark to their plain
// letter; it avoids a dependency for six letters.
var accentFolds = strings.NewReplacer(
	"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U", "Ü", "U", "Ñ", "N",
)

func (Catalog) Normalize(merchant string) string {
	return normalizeMerchant(merchant)
}

// normalizeMerchant gives the descriptor a stable form so charges of the same
// merchant group together however the bank wrote it.
func normalizeMerchant(merchant string) string {
	upper := strings.ToUpper(strings.TrimSpace(merchant))
	withoutPrefix := stripPaymentProcessorPrefix(upper)
	return collapsePunctuation(accentFolds.Replace(withoutPrefix))
}

func stripPaymentProcessorPrefix(upperMerchant string) string {
	for _, prefix := range paymentProcessorPrefixes {
		if rest, isPrefixed := strings.CutPrefix(upperMerchant, prefix); isPrefixed {
			return rest
		}
	}
	return upperMerchant
}

// collapsePunctuation turns every run of non letters and digits into one space.
func collapsePunctuation(text string) string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return strings.Join(words, " ")
}
