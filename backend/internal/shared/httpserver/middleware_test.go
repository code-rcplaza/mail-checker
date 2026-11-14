package httpserver

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// routerWithLogs is the base router plus a route that panics and one that fails.
func routerWithLogs(t *testing.T) (http.Handler, *bytes.Buffer) {
	t.Helper()
	var logs bytes.Buffer
	router := NewRouter(slog.New(slog.NewTextHandler(&logs, nil)), AlwaysHealthy)
	router.Get("/api/boom", func(http.ResponseWriter, *http.Request) { panic("secret detail") })
	router.Post("/api/echo", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	return router, &logs
}

func TestRecoverPanics_AnswersInternalAndKeepsServing(t *testing.T) {
	router, logs := routerWithLogs(t)

	panicked := httptest.NewRecorder()
	router.ServeHTTP(panicked, httptest.NewRequest(http.MethodGet, "/api/boom", nil))
	next := httptest.NewRecorder()
	router.ServeHTTP(next, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	if panicked.Code != http.StatusInternalServerError || !strings.Contains(panicked.Body.String(), `"common.internal"`) {
		t.Errorf("panic answered %d %s, want 500 with common.internal", panicked.Code, panicked.Body.String())
	}
	if strings.Contains(panicked.Body.String(), "secret detail") {
		t.Errorf("the panic detail reached the response: %s", panicked.Body.String())
	}
	if !strings.Contains(logs.String(), "secret detail") {
		t.Errorf("logs = %q, want the panic detail logged", logs.String())
	}
	if next.Code != http.StatusOK {
		t.Errorf("request after the panic = %d, want 200", next.Code)
	}
}

func TestLogRequests_LogsMethodPathStatusAndDurationOnly(t *testing.T) {
	router, logs := routerWithLogs(t)

	request := httptest.NewRequest(http.MethodPost, "/api/echo?lineKey=service:Netflix", strings.NewReader(`{"email":"ana@example.com"}`))
	router.ServeHTTP(httptest.NewRecorder(), request)

	line := logs.String()
	for _, want := range []string{"method=POST", "path=/api/echo", "status=204", "duration="} {
		if !strings.Contains(line, want) {
			t.Errorf("log %q does not contain %q", line, want)
		}
	}
	for _, private := range []string{"lineKey", "Netflix", "ana@example.com"} {
		if strings.Contains(line, private) {
			t.Errorf("log %q contains %q", line, private)
		}
	}
	if !strings.Contains(line, "level=INFO") {
		t.Errorf("log %q, want INFO for a 2xx", line)
	}
}

func TestLogRequests_ServerErrorsLogAsErrors(t *testing.T) {
	router, logs := routerWithLogs(t)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/boom", nil))

	if !strings.Contains(logs.String(), "level=ERROR") || !strings.Contains(logs.String(), "status=500") {
		t.Errorf("logs = %q, want the 500 request logged at ERROR", logs.String())
	}
}

func TestLogRequests_KeepsTheResponseControllerWorking(t *testing.T) {
	var deadlineErr error
	router := NewRouter(discardLogger(), AlwaysHealthy)
	router.Get("/api/deadline", func(w http.ResponseWriter, _ *http.Request) {
		deadlineErr = http.NewResponseController(w).SetReadDeadline(time.Now().Add(time.Second))
		w.WriteHeader(http.StatusNoContent)
	})
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL + "/api/deadline")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	response.Body.Close()

	if deadlineErr != nil {
		t.Errorf("SetReadDeadline behind the request log = %v, want it to reach the connection", deadlineErr)
	}
}
