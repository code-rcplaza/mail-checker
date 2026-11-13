package catalog

import "testing"

func TestNormalizeMerchant(t *testing.T) {
	cases := []struct {
		name     string
		merchant string
		want     string
	}{
		{"punctuation becomes a space", "C. Verde Mall", "C VERDE MALL"},
		{"accents and enye are removed", "Café Ñuñoa", "CAFE NUNOA"},
		{"dLocal prefix", "DL*GOOGLE YOUTUBE", "GOOGLE YOUTUBE"},
		{"PayU prefix", "PAYU *UBER TRIP", "UBER TRIP"},
		{"Mercado Pago short prefix", "MP *MERCADO LIB", "MERCADO LIB"},
		{"Mercado Pago long prefix", "MERCADOPAGO*ISTEC", "ISTEC"},
		{"DP prefix", "DP *FALABELLA.C", "FALABELLA C"},
		{"star after the merchant is kept as a word break", "PedidosYa*Plus", "PEDIDOSYA PLUS"},
		{"merchant before the star is not a prefix: Anthropic", "ANTHROPIC* CLAUDE SUB", "ANTHROPIC CLAUDE SUB"},
		{"merchant before the star is not a prefix: EA", "EA *ELECTRONIC ARTS", "EA ELECTRONIC ARTS"},
		{"punctuation collapses repeated spaces", "APPLE.COM/BILL   866-712-7753", "APPLE COM BILL 866 712 7753"},
		{"prefix ignores case and surrounding spaces", "  dl*google  ", "GOOGLE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeMerchant(tc.merchant)

			if got != tc.want {
				t.Errorf("normalizeMerchant(%q) = %q, want %q", tc.merchant, got, tc.want)
			}
		})
	}
}
