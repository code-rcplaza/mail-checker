package catalog

import (
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

func TestIdentifyKnownServices(t *testing.T) {
	cases := []struct {
		name     string
		merchant string
		want     domain.Service
	}{
		{"Netflix with phone", "NETFLIX.COM 866-579-7172", domain.Service{Name: "Netflix"}},
		{"Netflix alone", "NETFLIX", domain.Service{Name: "Netflix"}},
		{"Disney+", "Disney Plus Burbank", domain.Service{Name: "Disney+"}},
		{"Spotify", "Spotify P1234ABCD Stockholm", domain.Service{Name: "Spotify"}},
		{"Claude", "ANTHROPIC* CLAUDE SUB", domain.Service{Name: "Claude"}},
		{"HostGator", "HOSTGATOR LAS CONDES", domain.Service{Name: "HostGator"}},
		{"playit.gg", "PLAYIT.GG* PLAYIT.GG/C", domain.Service{Name: "playit.gg"}},
		{"PedidosYa Plus", "PedidosYa*Plus CL", domain.Service{Name: "PedidosYa Plus"}},
		{"App Store is umbrella", "APPLE.COM/BILL 866-712-7753", domain.Service{Name: "App Store", IsUmbrella: true}},
		{"Google is umbrella", "DL*GOOGLE YOUTUBE", domain.Service{Name: "Google", IsUmbrella: true}},
		{"Microsoft is umbrella", "MICROSOFT. SUBSCRIPTION", domain.Service{Name: "Microsoft", IsUmbrella: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, isKnown := Catalog{}.Identify(tc.merchant)

			if !isKnown {
				t.Fatalf("Identify(%q) did not find a service, want %+v", tc.merchant, tc.want)
			}
			if got != tc.want {
				t.Errorf("Identify(%q) = %+v, want %+v", tc.merchant, got, tc.want)
			}
		})
	}
}

func TestIdentifyRejectsUnknownMerchants(t *testing.T) {
	cases := []struct {
		name     string
		merchant string
	}{
		{"pattern in the middle of a word", "NETFLIXER"},
		{"pattern in the middle of the descriptor", "MY NETFLIX"},
		{"PedidosYa restaurant", "PedidosYa*KFC Vicun"},
		{"PedidosYa tip", "PedidosYa*Propina"},
		{"EA is left out on purpose", "EA *ELECTRONIC ARTS"},
		{"insurance is left out on purpose", "SEGUROS CLC S.A."},
		{"My Mini Factory is left out on purpose", "MY MINI FACTORY LTD"},
		{"Nanonoble is left out on purpose", "NANONOBLE PTE. LTD."},
		{"processor prefix with unknown merchant", "PAYU *UBER TRIP"},
		{"Anthropic API is not Claude", "ANTHROPIC* API"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, isKnown := Catalog{}.Identify(tc.merchant)

			if isKnown {
				t.Errorf("Identify(%q) = %+v, want no service", tc.merchant, got)
			}
		})
	}
}
