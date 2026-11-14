package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
)

func serveHealth(t *testing.T, check HealthCheck) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	NewRouter(discardLogger(), check).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	return recorder
}

func TestHealth_OpenDatabaseIsOK(t *testing.T) {
	db := sqlitetest.OpenMigrated(t)

	got := serveHealth(t, db.PingContext)

	if got.Code != http.StatusOK || strings.TrimSpace(got.Body.String()) != `{"status":"ok"}` {
		t.Errorf("health = %d %s, want 200 {\"status\":\"ok\"}", got.Code, got.Body.String())
	}
}

func TestHealth_DatabaseThatDoesNotAnswerIsUnavailable(t *testing.T) {
	closed := sqlitetest.OpenMigrated(t)
	if err := closed.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}
	cases := []struct {
		name  string
		check HealthCheck
	}{
		{"a closed database", closed.PingContext},
		{"a database that hangs past the limit", func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			started := time.Now()

			got := serveHealth(t, tc.check)

			if body := strings.TrimSpace(got.Body.String()); got.Code != http.StatusServiceUnavailable || body != `{"error":{"key":"common.unavailable"}}` {
				t.Errorf("health = %d %s, want 503 with only common.unavailable", got.Code, body)
			}
			if elapsed := time.Since(started); elapsed > healthTimeout+time.Second {
				t.Errorf("health took %v, want at most about %v", elapsed, healthTimeout)
			}
		})
	}
}

func TestHealth_IsPublic(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	request.Header.Set("Origin", "https://elsewhere.example")

	NewRouter(discardLogger(), AlwaysHealthy).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("health without cookie from another origin = %d, want 200", recorder.Code)
	}
}
