package httpserver

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func serve(t *testing.T, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	NewRouter(discardLogger()).ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder
}

func TestRouter_HealthReturnsOK(t *testing.T) {
	res := serve(t, http.MethodGet, "/api/health")

	if res.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", res.Code)
	}
	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body map[string]string
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not JSON: %v", res.Body.String(), err)
	}
	if body["status"] != "ok" {
		t.Errorf("body = %v, want {\"status\":\"ok\"}", body)
	}
}

func TestRouter_UnknownAPIRouteReturnsRouteNotFound(t *testing.T) {
	res := serve(t, http.MethodGet, "/api/does-not-exist")

	if res.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.Code)
	}
	var body struct {
		Error struct {
			Key string `json:"key"`
		} `json:"error"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not JSON: %v", res.Body.String(), err)
	}
	if body.Error.Key != "common.route_not_found" {
		t.Errorf("error.key = %q, want common.route_not_found", body.Error.Key)
	}
}
