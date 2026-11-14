package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
)

var testNow = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

const (
	validToken   = "valid-token"
	expiredToken = "expired-token"
)

type sessionsByHash map[string]domain.Session

func (s sessionsByHash) Create(_ context.Context, session domain.Session) error {
	s[session.TokenHash] = session
	return nil
}

func (s sessionsByHash) FindByTokenHash(_ context.Context, hash string) (domain.Session, error) {
	session, found := s[hash]
	if !found {
		return domain.Session{}, fmt.Errorf("find session: %w", shared.ErrNotFound)
	}
	return session, nil
}

func (s sessionsByHash) Delete(_ context.Context, hash string) error {
	delete(s, hash)
	return nil
}

// protectedRouter has a public auth route and a protected one that echoes the user.
func protectedRouter() http.Handler {
	sessions := sessionsByHash{
		domain.HashSecret(validToken):   {UserID: 7, TokenHash: domain.HashSecret(validToken), ExpiresAt: testNow.Add(time.Hour)},
		domain.HashSecret(expiredToken): {UserID: 7, TokenHash: domain.HashSecret(expiredToken), ExpiresAt: testNow},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := httpserver.NewRouter(logger, httpserver.AlwaysHealthy)
	requireSession := RequireSession{Validate: usecases.NewValidateSession(sessions, func() time.Time { return testNow }), Responder: httpserver.Responder{Logger: logger}}
	router.Post("/api/auth/login", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	router.Group(func(protected chi.Router) {
		protected.Use(requireSession.Handler)
		protected.Get("/api/me", func(w http.ResponseWriter, r *http.Request) {
			userID, err := httpserver.UserIDFrom(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusTeapot)
				return
			}
			fmt.Fprint(w, userID)
		})
	})
	return router
}

func get(method, path, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	if token != "" {
		request.AddCookie(&http.Cookie{Name: SessionCookieName, Value: token})
	}
	recorder := httptest.NewRecorder()
	protectedRouter().ServeHTTP(recorder, request)
	return recorder
}

func TestRequireSession(t *testing.T) {
	cases := []struct {
		name     string
		method   string
		path     string
		token    string
		want     int
		wantBody string
	}{
		{"a valid session reaches the handler with its user", http.MethodGet, "/api/me", validToken, http.StatusOK, "7"},
		{"no cookie is unauthorized", http.MethodGet, "/api/me", "", http.StatusUnauthorized, `"common.unauthorized"`},
		{"an expired session is unauthorized", http.MethodGet, "/api/me", expiredToken, http.StatusUnauthorized, `"common.unauthorized"`},
		{"an unknown session is unauthorized", http.MethodGet, "/api/me", "guessed", http.StatusUnauthorized, `"common.unauthorized"`},
		{"health needs no session", http.MethodGet, "/api/health", "", http.StatusOK, `"ok"`},
		{"auth routes need no session", http.MethodPost, "/api/auth/login", "", http.StatusNoContent, ""},
		{"an unknown route is a JSON 404, not a 401", http.MethodGet, "/api/nothing-here", "", http.StatusNotFound, `"common.route_not_found"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := get(tc.method, tc.path, tc.token)

			if got.Code != tc.want || !strings.Contains(got.Body.String(), tc.wantBody) {
				t.Errorf("%s %s = %d %s, want %d with %s", tc.method, tc.path, got.Code, got.Body.String(), tc.want, tc.wantBody)
			}
		})
	}
}
