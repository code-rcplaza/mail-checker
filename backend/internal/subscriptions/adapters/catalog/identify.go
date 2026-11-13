package catalog

import (
	"strings"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// Catalog recognises known services from the merchant descriptor of a charge.
type Catalog struct{}

// serviceRule ties a normalized descriptor pattern to the service it names.
type serviceRule struct {
	pattern string
	service domain.Service
}

// serviceRules only holds descriptors verified in real charges: calling
// something Netflix when it is not is worse than not recognising it, because
// an unknown merchant still shows up as possible through its cadence.
var serviceRules = []serviceRule{
	{"NETFLIX", domain.Service{Name: "Netflix"}},
	{"DISNEY PLUS", domain.Service{Name: "Disney+"}},
	{"SPOTIFY", domain.Service{Name: "Spotify"}},
	{"ANTHROPIC CLAUDE", domain.Service{Name: "Claude"}},
	{"HOSTGATOR", domain.Service{Name: "HostGator"}},
	{"PLAYIT GG", domain.Service{Name: "playit.gg"}},
	{"PEDIDOSYA PLUS", domain.Service{Name: "PedidosYa Plus"}},
	{"APPLE COM BILL", domain.Service{Name: "App Store", IsUmbrella: true}},
	{"GOOGLE", domain.Service{Name: "Google", IsUmbrella: true}},
	{"MICROSOFT", domain.Service{Name: "Microsoft", IsUmbrella: true}},
}

// Identify finds the service whose pattern starts the normalized descriptor.
func (Catalog) Identify(merchant string) (domain.Service, bool) {
	normalized := normalizeMerchant(merchant)
	for _, rule := range serviceRules {
		if startsWithWords(normalized, rule.pattern) {
			return rule.service, true
		}
	}
	return domain.Service{}, false
}

// startsWithWords reports whether text begins with the pattern and the
// pattern ends at a word boundary, so "NETFLIX" does not match "NETFLIXER".
func startsWithWords(text, pattern string) bool {
	rest, hasPattern := strings.CutPrefix(text, pattern)
	if !hasPattern {
		return false
	}
	return rest == "" || strings.HasPrefix(rest, " ")
}
