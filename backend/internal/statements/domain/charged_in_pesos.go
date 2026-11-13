package domain

import shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"

// pesosPerDollarMinRatio separates a charge in pesos from one in dollars or
// euros: the international statement shows the origin amount without its
// currency, and only pesos are worth hundreds of times the US$ amount.
const pesosPerDollarMinRatio = 100

// IsChargedInPesos tells whether the origin amount of an international charge
// is in pesos. An origin of exactly pesosPerDollarMinRatio times the dollars
// counts as pesos. Multiplying instead of dividing keeps US$ 0,00 safe: any
// origin above zero is then pesos, and zero against zero is not.
func IsChargedInPesos(origin, dollars shared.Money) bool {
	return origin > 0 && origin >= dollars*pesosPerDollarMinRatio
}
