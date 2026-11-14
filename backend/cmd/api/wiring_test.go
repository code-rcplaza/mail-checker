package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/mailer"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/bci"
)

func testRouter(t *testing.T) http.Handler {
	t.Helper()
	router, err := buildRouter(appDependencies{db: sqlitetest.OpenMigrated(t), logger: slog.New(slog.NewTextHandler(io.Discard, nil)), stdout: io.Discard})
	if err != nil {
		t.Fatalf("buildRouter error = %v, want nil", err)
	}
	return router
}

type route struct {
	method string
	path   string
}

// mvpRoutes is the API table of ARCHITECTURE.md; every one needs a session
// except health and auth.
var (
	publicRoutes = []route{
		{http.MethodPost, "/api/auth/register"}, {http.MethodPost, "/api/auth/login"},
		{http.MethodPost, "/api/auth/code/verify"}, {http.MethodPost, "/api/auth/code/resend"},
		{http.MethodPost, "/api/auth/logout"},
	}
	protectedRoutes = []route{
		{http.MethodGet, "/api/me"}, {http.MethodDelete, "/api/me"}, {http.MethodGet, "/api/export"},
		{http.MethodGet, "/api/banks"}, {http.MethodPost, "/api/statements"}, {http.MethodGet, "/api/summary"},
		{http.MethodPost, "/api/subscriptions/decisions"}, {http.MethodDelete, "/api/subscriptions/decisions"},
		{http.MethodPut, "/api/subscriptions/aliases"}, {http.MethodPut, "/api/subscriptions/overrides"},
		{http.MethodPost, "/api/subscriptions/manual"}, {http.MethodPut, "/api/subscriptions/manual/0190a000-0000-7000-8000-000000000001"},
		{http.MethodDelete, "/api/subscriptions/manual/0190a000-0000-7000-8000-000000000001"},
		{http.MethodDelete, "/api/subscriptions/line"}, {http.MethodPost, "/api/subscriptions/price-alerts/decision"},
		{http.MethodGet, "/api/subscriptions/cancel-guide"}, {http.MethodPost, "/api/subscriptions/cancellations"},
	}
)

func serveRoute(router http.Handler, r route) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(r.method, r.path, strings.NewReader("{}")))
	return recorder
}

func TestBuildRouter_EveryRouteExists(t *testing.T) {
	router := testRouter(t)
	for _, r := range append(append([]route{{http.MethodGet, "/api/health"}}, publicRoutes...), protectedRoutes...) {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := serveRoute(router, r); got.Code == http.StatusNotFound || got.Code == http.StatusMethodNotAllowed {
				t.Errorf("%s %s = %d %s, want the route to exist", r.method, r.path, got.Code, got.Body.String())
			}
		})
	}
}

func TestBuildRouter_OnlyHealthAndAuthGoWithoutSession(t *testing.T) {
	router := testRouter(t)
	for _, r := range protectedRoutes {
		if got := serveRoute(router, r); got.Code != http.StatusUnauthorized {
			t.Errorf("%s %s without session = %d, want 401", r.method, r.path, got.Code)
		}
	}
	for _, r := range append([]route{{http.MethodGet, "/api/health"}}, publicRoutes...) {
		if got := serveRoute(router, r); got.Code == http.StatusUnauthorized {
			t.Errorf("%s %s without session = 401, want it public", r.method, r.path)
		}
	}
}

func TestBuildRouter_PostMVPRoutesDoNotExist(t *testing.T) {
	router := testRouter(t)
	for _, r := range []route{{http.MethodGet, "/api/accounts"}, {http.MethodPost, "/api/accounts"}, {http.MethodPost, "/api/sync"}} {
		if got := serveRoute(router, r); got.Code != http.StatusNotFound {
			t.Errorf("%s %s = %d, want 404", r.method, r.path, got.Code)
		}
	}
}

func TestBankParsers_HasBCIOnce(t *testing.T) {
	parsers := bankParsers(time.UTC)

	if len(parsers) != 1 || parsers[0].Bank() != bci.Bank {
		t.Errorf("parsers = %v, want only BCI", parsers)
	}
}

func TestCodeSenderFor(t *testing.T) {
	cases := []struct {
		name      string
		config    Config
		wantKind  string
		wantError bool
	}{
		{"production uses Resend", Config{IsProduction: true, Mailer: MailerConfig{ResendAPIKey: "re_key", MailFrom: "Fokko <codigos@example.com>"}}, "resend", false},
		{"development without key prints the code", Config{}, "console", false},
		{"development with a key uses Resend", Config{Mailer: MailerConfig{ResendAPIKey: "re_key", MailFrom: "f@example.com"}}, "resend", false},
		{"production never prints the code", Config{IsProduction: true}, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sender, err := codeSenderFor(tc.config, io.Discard)

			if (err != nil) != tc.wantError {
				t.Fatalf("codeSenderFor error = %v, want error %v", err, tc.wantError)
			}
			kind := ""
			switch sender.(type) {
			case mailer.ResendSender:
				kind = "resend"
			case mailer.ConsoleSender:
				kind = "console"
			}
			if kind != tc.wantKind {
				t.Errorf("sender = %T, want %s", sender, tc.wantKind)
			}
		})
	}
}

// countingCleanup counts runs and fails the first one, to check the loop goes on.
type countingCleanup struct {
	mu   sync.Mutex
	runs int
}

func (c *countingCleanup) clean(context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.runs++
	if c.runs == 1 {
		return errors.New("database is locked")
	}
	return nil
}

func (c *countingCleanup) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.runs
}

func TestCleanupLoop_RunsEveryIntervalUntilCanceled(t *testing.T) {
	cleanup := &countingCleanup{}
	var logs bytes.Buffer
	loop := cleanupLoop{clean: cleanup.clean, every: 5 * time.Millisecond, logger: slog.New(slog.NewTextHandler(&logs, nil))}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() { loop.run(ctx); close(done) }()
	deadline := time.Now().Add(2 * time.Second)
	for cleanup.count() < 3 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("cleanup loop did not stop after cancel")
	}

	if cleanup.count() < 3 {
		t.Errorf("cleanup ran %d times, want it to keep running after a failure", cleanup.count())
	}
	if !strings.Contains(logs.String(), "clean expired access failed") || !strings.Contains(logs.String(), "database is locked") {
		t.Errorf("logs = %q, want the failure logged", logs.String())
	}
}
