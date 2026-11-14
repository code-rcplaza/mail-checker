package http

import (
	"encoding/json"
	"net/http"
	"testing"
)

func validManual() map[string]any {
	return map[string]any{"name": "Gimnasio", "amount": 2500000, "currency": "CLP", "billingDay": 15}
}

func TestPostManual_CreatesItWithItsPublicID(t *testing.T) {
	server := newTestServer(t)

	got := server.send(t, http.MethodPost, "/api/subscriptions/manual", jsonBody(t, validManual()))

	var created struct {
		ID string `json:"id"`
	}
	if got.Code != http.StatusCreated || json.Unmarshal(got.Body.Bytes(), &created) != nil || created.ID == "" {
		t.Fatalf("response = %d %s, want 201 with the id", got.Code, got.Body.String())
	}
	if line, isShown := server.find(t, "manual:"+created.ID); !isShown || line.Alias != "Gimnasio" {
		t.Errorf("summary line of manual:%s = %+v (shown %v), want it named Gimnasio", created.ID, line, isShown)
	}
}

func TestPostManual_InvalidFields(t *testing.T) {
	got := newTestServer(t).send(t, http.MethodPost, "/api/subscriptions/manual", jsonBody(t, map[string]any{"name": " ", "amount": 0, "currency": "EUR", "billingDay": 0}))

	assertResponse(t, got, http.StatusBadRequest,
		`{"error":{"key":"common.validation_failed","fields":{"amount":"subscriptions.manual_amount_not_positive","billingDay":"subscriptions.manual_day_out_of_range","currency":"common.currency_unknown","name":"subscriptions.manual_name_missing"}}}`)
}

func TestPutAndDeleteManual(t *testing.T) {
	server := newTestServer(t)
	created := server.send(t, http.MethodPost, "/api/subscriptions/manual", jsonBody(t, validManual()))
	var own struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &own); err != nil {
		t.Fatalf("decode created manual: %v (%s)", err, created.Body.String())
	}
	theirs := insertManualOf(t, server.db, "beto@example.com")
	edited := validManual()
	edited["name"] = "Gimnasio nuevo"
	cases := []struct {
		name       string
		method     string
		publicID   string
		body       map[string]any
		wantStatus int
	}{
		{"edit own", http.MethodPut, own.ID, edited, http.StatusNoContent},
		{"edit own with invalid fields", http.MethodPut, own.ID, map[string]any{"name": "", "amount": 1, "currency": "CLP", "billingDay": 1}, http.StatusBadRequest},
		{"edit another user's", http.MethodPut, theirs, edited, http.StatusNotFound},
		{"edit one that does not exist", http.MethodPut, "0190a000-0000-7000-8000-000000000404", edited, http.StatusNotFound},
		{"edit with a malformed id", http.MethodPut, "not-an-id", edited, http.StatusNotFound},
		{"delete another user's", http.MethodDelete, theirs, nil, http.StatusNotFound},
		{"delete with a malformed id", http.MethodDelete, "not-an-id", nil, http.StatusNotFound},
		{"delete own", http.MethodDelete, own.ID, nil, http.StatusNoContent},
		{"delete own again", http.MethodDelete, own.ID, nil, http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := ""
			if tc.body != nil {
				body = jsonBody(t, tc.body)
			}

			got := server.send(t, tc.method, "/api/subscriptions/manual/"+tc.publicID, body)

			assertResponse(t, got, tc.wantStatus, "")
		})
	}
	if _, isShown := server.find(t, "manual:"+own.ID); isShown {
		t.Errorf("manual %s still shows after deleting it", own.ID)
	}
}
