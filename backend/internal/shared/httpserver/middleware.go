package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// RecoverPanics turns a panic into a 500 with common.internal; the detail and
// the stack go only to the log, and the server keeps serving.
func (responder Responder) RecoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}
			if recovered == http.ErrAbortHandler {
				panic(recovered) // net/http's own way to abort a response
			}
			responder.Logger.Error("handler panicked", "path", r.URL.Path, "panic", fmt.Sprint(recovered), "stack", string(debug.Stack()))
			responder.WriteError(Exchange{Writer: w, Request: r}, fmt.Errorf("panic: %v", recovered))
		}()
		next.ServeHTTP(w, r)
	})
}

// LogRequests logs method, path, status and duration of each request; never the
// query or the body, which can carry line keys and emails. A 5xx logs as error.
func (responder Responder) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		responder.Logger.Log(r.Context(), levelOf(recorder.status), "request",
			"method", r.Method, "path", r.URL.Path, "status", recorder.status, "duration", time.Since(started))
	})
}

func levelOf(status int) slog.Level {
	if status >= http.StatusInternalServerError {
		return slog.LevelError
	}
	return slog.LevelInfo
}

// statusRecorder remembers the status the handler wrote; 200 when it wrote none.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

// Unwrap lets http.ResponseController reach the connection (deadlines, flush)
// through the recorder.
func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }
