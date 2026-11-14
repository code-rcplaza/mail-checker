package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestUserID_TravelsInTheRequestContext(t *testing.T) {
	got, err := UserIDFrom(WithUserID(context.Background(), 42))

	if err != nil || got != 42 {
		t.Errorf("UserIDFrom = %d, %v; want 42", got, err)
	}
}

func TestUserID_WithoutTheSessionMiddlewareIsAnError(t *testing.T) {
	got, err := UserIDFrom(context.Background())

	if !errors.Is(err, domain.ErrUnauthorized) || got != 0 {
		t.Errorf("UserIDFrom = %d, %v; want ErrUnauthorized, never a zero user", got, err)
	}
}

func serveWithOrigin(t *testing.T, method, origin string) *httptest.ResponseRecorder {
	t.Helper()
	router := NewRouter(discardLogger(), AlwaysHealthy)
	router.MethodFunc(method, "/api/things", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	request := httptest.NewRequest(method, "http://fokko.example/api/things", strings.NewReader(`{}`))
	if origin != "" {
		request.Header.Set("Origin", origin)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestCheckOrigin(t *testing.T) {
	cases := []struct {
		name   string
		method string
		origin string
		want   int
	}{
		{"a POST from another site is forbidden", http.MethodPost, "https://evil.example", http.StatusForbidden},
		{"a PUT from another site is forbidden", http.MethodPut, "https://evil.example", http.StatusForbidden},
		{"a DELETE from another site is forbidden", http.MethodDelete, "https://evil.example", http.StatusForbidden},
		{"a POST from the app itself passes", http.MethodPost, "http://fokko.example", http.StatusNoContent},
		{"a GET from another site is not checked", http.MethodGet, "https://evil.example", http.StatusNoContent},
		{"a POST without Origin passes: SameSite=Lax already keeps the cookie home", http.MethodPost, "", http.StatusNoContent},
		{"a malformed Origin is forbidden", http.MethodPost, "::not a url", http.StatusForbidden},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := serveWithOrigin(t, tc.method, tc.origin)

			if got.Code != tc.want {
				t.Errorf("status = %d, want %d", got.Code, tc.want)
			}
			if tc.want == http.StatusForbidden && !strings.Contains(got.Body.String(), `"common.forbidden"`) {
				t.Errorf("body = %s, want common.forbidden", got.Body.String())
			}
		})
	}
}

func TestRouter_SendsNoCORSHeaders(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, "/api/health", nil)
			request.Header.Set("Origin", "https://evil.example")
			request.Header.Set("Access-Control-Request-Method", http.MethodPost)
			recorder := httptest.NewRecorder()

			NewRouter(discardLogger(), AlwaysHealthy).ServeHTTP(recorder, request)

			if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
				t.Errorf("Access-Control-Allow-Origin = %q, want none", got)
			}
			if method == http.MethodOptions && recorder.Code < 400 {
				t.Errorf("preflight status = %d, want it not authorized", recorder.Code)
			}
		})
	}
}
