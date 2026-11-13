package catalog

import (
	"slices"
	"strings"
)

// subscriptionWords name a recurring charge in English or Spanish; they are
// compared against whole words of the normalized descriptor.
var subscriptionWords = []string{
	"SUBSCRIPTION",
	"MEMBERSHIP",
	"SUSCRIPCION",
	"MEMBRESIA",
	"PREMIUM",
	"PLAN",
}

// LooksLikeSubscription lets a single charge show up as possible without
// waiting for a second month.
func (Catalog) LooksLikeSubscription(merchant string) bool {
	words := strings.Fields(normalizeMerchant(merchant))
	return slices.ContainsFunc(words, isSubscriptionWord)
}

func isSubscriptionWord(word string) bool {
	return slices.Contains(subscriptionWords, word)
}
