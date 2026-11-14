package http

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const unknownLine = "merchant:NEVER SEEN:CLP||d01"

const lineKeyRequired = `"fields":{"lineKey":"subscriptions.line_key_required"}`

func TestPostDecision(t *testing.T) {
	cases := []struct {
		name       string
		body       func(seededKeys) any
		wantStatus int
		wantBody   string
	}{
		{"confirm with cadence", func(k seededKeys) any {
			return map[string]string{"groupKey": k.possible, "decision": "confirm", "cadence": "monthly"}
		}, http.StatusNoContent, ""},
		{"confirm without cadence", func(k seededKeys) any {
			return map[string]string{"groupKey": k.possible, "decision": "confirm"}
		}, http.StatusBadRequest, "subscriptions.cadence_required"},
		{"dismiss", func(k seededKeys) any {
			return map[string]string{"groupKey": k.possible, "decision": "dismiss"}
		}, http.StatusNoContent, ""},
		{"a key that is not a possible", func(k seededKeys) any {
			return map[string]string{"groupKey": k.active, "decision": "dismiss"}
		}, http.StatusNotFound, "subscriptions.possible_not_found"},
		{"no group key", func(seededKeys) any {
			return map[string]string{"decision": "dismiss"}
		}, http.StatusBadRequest, `"fields":{"groupKey":"subscriptions.line_key_required"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			server.seedLines(t)
			keys := server.seededKeys(t)

			got := server.send(t, http.MethodPost, "/api/subscriptions/decisions", jsonBody(t, tc.body(keys)))

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
}

func TestPostDecision_ConfirmedLineIsASubscription(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)

	server.send(t, http.MethodPost, "/api/subscriptions/decisions", jsonBody(t, map[string]string{"groupKey": keys.possible, "decision": "confirm", "cadence": "monthly"}))

	if view := server.summary(t); len(view.Subscriptions) != 2 || len(view.Possible) != 0 {
		t.Errorf("summary = %+v, want the confirmed line among the subscriptions", view)
	}
}

func TestDeleteDecision(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)
	server.send(t, http.MethodPost, "/api/subscriptions/decisions", jsonBody(t, map[string]string{"groupKey": keys.possible, "decision": "dismiss"}))
	target := "/api/subscriptions/decisions?groupKey=" + url.QueryEscape(keys.possible)

	assertResponse(t, server.send(t, http.MethodDelete, target, ""), http.StatusNoContent, "")
	if _, isPossible := server.find(t, keys.possible); !isPossible {
		t.Errorf("line %q is not back after undoing its decision", keys.possible)
	}
	assertResponse(t, server.send(t, http.MethodDelete, target, ""), http.StatusNoContent, "")
}

func TestDeleteDecision_KeyIsRequired(t *testing.T) {
	for _, target := range []string{"/api/subscriptions/decisions", "/api/subscriptions/decisions?groupKey=", "/api/subscriptions/decisions?groupKey=%20"} {
		t.Run(target, func(t *testing.T) {
			got := newTestServer(t).send(t, http.MethodDelete, target, "")

			assertResponse(t, got, http.StatusBadRequest, `"fields":{"groupKey":"subscriptions.line_key_required"}`)
		})
	}
}

func TestPutAlias(t *testing.T) {
	cases := []struct {
		name       string
		lineKey    func(seededKeys) string
		alias      string
		wantStatus int
		wantBody   string
	}{
		{"a valid alias", func(k seededKeys) string { return k.active }, "Padres", http.StatusNoContent, ""},
		{"an empty alias", func(k seededKeys) string { return k.active }, "", http.StatusNoContent, ""},
		{"a too long alias", func(k seededKeys) string { return k.active }, strings.Repeat("a", 41), http.StatusBadRequest, "subscriptions.alias_too_long"},
		{"a line that does not exist", func(seededKeys) string { return unknownLine }, "Padres", http.StatusNotFound, "subscriptions.line_not_found"},
		{"no line key", func(seededKeys) string { return "" }, "Padres", http.StatusBadRequest, lineKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			server.seedLines(t)
			keys := server.seededKeys(t)

			got := server.send(t, http.MethodPut, "/api/subscriptions/aliases", jsonBody(t, map[string]string{"lineKey": tc.lineKey(keys), "alias": tc.alias}))

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
}

func TestPutAlias_EmptyAliasRemovesIt(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)
	putAlias := func(alias string) {
		server.send(t, http.MethodPut, "/api/subscriptions/aliases", jsonBody(t, map[string]string{"lineKey": keys.active, "alias": alias}))
	}

	putAlias("Padres")
	if line, _ := server.find(t, keys.active); line.Alias != "Padres" {
		t.Errorf("alias = %q, want Padres", line.Alias)
	}
	putAlias("")
	if line, _ := server.find(t, keys.active); line.Alias != "" {
		t.Errorf("alias = %q, want it removed", line.Alias)
	}
}

func TestPutOverride(t *testing.T) {
	cases := []struct {
		name       string
		body       func(seededKeys) map[string]any
		wantStatus int
		wantBody   string
	}{
		{"a valid override", func(k seededKeys) map[string]any {
			return map[string]any{"lineKey": k.active, "amount": 900000, "billingDay": 15}
		}, http.StatusNoContent, ""},
		{"no amount nor day", func(k seededKeys) map[string]any {
			return map[string]any{"lineKey": k.active}
		}, http.StatusNoContent, ""},
		{"a zero amount and a day out of range", func(k seededKeys) map[string]any {
			return map[string]any{"lineKey": k.active, "amount": 0, "billingDay": 32}
		}, http.StatusBadRequest, `"fields":{"amount":"subscriptions.override_amount_not_positive","billingDay":"subscriptions.override_day_out_of_range"}`},
		{"a negative amount", func(k seededKeys) map[string]any {
			return map[string]any{"lineKey": k.active, "amount": -1}
		}, http.StatusBadRequest, `"fields":{"amount":"subscriptions.override_amount_not_positive"}`},
		{"a line that does not exist", func(seededKeys) map[string]any {
			return map[string]any{"lineKey": unknownLine, "amount": 900000}
		}, http.StatusNotFound, "subscriptions.line_not_found"},
		{"no line key", func(seededKeys) map[string]any {
			return map[string]any{"amount": 900000}
		}, http.StatusBadRequest, lineKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			server.seedLines(t)
			keys := server.seededKeys(t)

			got := server.send(t, http.MethodPut, "/api/subscriptions/overrides", jsonBody(t, tc.body(keys)))

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
}

func TestPutOverride_EmptyOverrideRestoresTheDetectedAmount(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)
	put := func(body map[string]any) {
		server.send(t, http.MethodPut, "/api/subscriptions/overrides", jsonBody(t, body))
	}

	put(map[string]any{"lineKey": keys.active, "amount": 900000})
	if line, _ := server.find(t, keys.active); line.Amount.Amount != 900000 {
		t.Errorf("amount = %d, want the overridden 900000", line.Amount.Amount)
	}
	put(map[string]any{"lineKey": keys.active})
	if line, _ := server.find(t, keys.active); line.Amount.Amount != 1050000 {
		t.Errorf("amount = %d, want the detected 1050000 back", line.Amount.Amount)
	}
}

func TestPostPriceAlertDecision(t *testing.T) {
	cases := []struct {
		name       string
		body       func(seededKeys) string
		wantStatus int
		wantBody   string
	}{
		{"a line with a rise", func(k seededKeys) string { return `{"lineKey":"` + k.active + `"}` }, http.StatusNoContent, ""},
		{"a line without a rise", func(k seededKeys) string { return `{"lineKey":"` + k.possible + `"}` }, http.StatusNotFound, "subscriptions.rise_not_found"},
		{"later is not an answer", func(k seededKeys) string { return `{"lineKey":"` + k.active + `","later":true}` }, http.StatusBadRequest, "common.invalid_json"},
		{"no line key", func(seededKeys) string { return `{}` }, http.StatusBadRequest, lineKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			server.seedLines(t)
			keys := server.seededKeys(t)

			got := server.send(t, http.MethodPost, "/api/subscriptions/price-alerts/decision", tc.body(keys))

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
}

func TestPostPriceAlertDecision_KeptRiseLeavesTheSummary(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)
	if rises := server.summary(t).Rises; len(rises) != 1 {
		t.Fatalf("rises = %v, want the seeded rise", rises)
	}

	server.send(t, http.MethodPost, "/api/subscriptions/price-alerts/decision", jsonBody(t, map[string]string{"lineKey": keys.active}))

	if rises := server.summary(t).Rises; len(rises) != 0 {
		t.Errorf("rises = %v, want none after keeping it", rises)
	}
}

func TestGetCancelGuide(t *testing.T) {
	cases := []struct {
		name       string
		query      func(seededKeys) string
		wantStatus int
		wantBody   string
	}{
		{"a service with a link", func(k seededKeys) string { return "?lineKey=" + url.QueryEscape(k.active) }, http.StatusOK, `{"url":"` + streamingCancelURL + `"}`},
		{"a merchant without a link", func(k seededKeys) string { return "?lineKey=" + url.QueryEscape(k.possible) }, http.StatusOK, `{"url":""}`},
		{"a line that does not exist", func(seededKeys) string { return "?lineKey=" + url.QueryEscape(unknownLine) }, http.StatusNotFound, "subscriptions.line_not_found"},
		{"no line key", func(seededKeys) string { return "" }, http.StatusBadRequest, lineKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			server.seedLines(t)
			keys := server.seededKeys(t)

			got := server.send(t, http.MethodGet, "/api/subscriptions/cancel-guide"+tc.query(keys), "")

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
}

func TestPostCancellation(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)
	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{"an existing line", jsonBody(t, map[string]string{"lineKey": keys.active}), http.StatusNoContent, ""},
		{"a line that does not exist", jsonBody(t, map[string]string{"lineKey": unknownLine}), http.StatusNotFound, "subscriptions.line_not_found"},
		{"no line key", `{}`, http.StatusBadRequest, lineKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := server.send(t, http.MethodPost, "/api/subscriptions/cancellations", tc.body)

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
	if line, _ := server.find(t, keys.active); line.MarkedForCancellationAt == nil {
		t.Errorf("line %q has no mark after marking it", keys.active)
	}
}

func TestDeleteLine(t *testing.T) {
	server := newTestServer(t)
	server.seedLines(t)
	keys := server.seededKeys(t)
	cases := []struct {
		name       string
		query      string
		wantStatus int
		wantBody   string
	}{
		{"an existing line", "?lineKey=" + url.QueryEscape(keys.active), http.StatusNoContent, ""},
		{"a line already hidden", "?lineKey=" + url.QueryEscape(keys.active), http.StatusNoContent, ""},
		{"a line that does not exist", "?lineKey=" + url.QueryEscape(unknownLine), http.StatusNotFound, "subscriptions.line_not_found"},
		{"no line key", "", http.StatusBadRequest, lineKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := server.send(t, http.MethodDelete, "/api/subscriptions/line"+tc.query, "")

			assertResponse(t, got, tc.wantStatus, tc.wantBody)
		})
	}
	if _, isShown := server.find(t, keys.active); isShown {
		t.Errorf("line %q still shows after hiding it", keys.active)
	}
}
