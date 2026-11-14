package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
)

func setCookie(t *testing.T, cookies SessionCookies, write func(SessionCookies, http.ResponseWriter)) *http.Cookie {
	t.Helper()
	recorder := httptest.NewRecorder()
	write(cookies, recorder)
	set := recorder.Result().Cookies()
	if len(set) != 1 {
		t.Fatalf("cookies set = %d, want 1", len(set))
	}
	return set[0]
}

func TestSessionCookies_Set(t *testing.T) {
	cookie := setCookie(t, SessionCookies{}, func(c SessionCookies, w http.ResponseWriter) { c.Set(w, "token-123") })

	if cookie.Name != SessionCookieName || cookie.Value != "token-123" || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" || cookie.MaxAge != int(domain.SessionLifetime.Seconds()) {
		t.Errorf("cookie = %+v, want HttpOnly, SameSite=Lax, Path=/ and 30 days", cookie)
	}
}

func TestSessionCookies_SecureOnlyInProduction(t *testing.T) {
	for _, isProduction := range []bool{true, false} {
		cookie := setCookie(t, SessionCookies{IsProduction: isProduction}, func(c SessionCookies, w http.ResponseWriter) { c.Set(w, "token") })

		if cookie.Secure != isProduction {
			t.Errorf("production %v: Secure = %v, want %v", isProduction, cookie.Secure, isProduction)
		}
	}
}

func TestSessionCookies_ClearExpiresItWithTheSameAttributes(t *testing.T) {
	cookie := setCookie(t, SessionCookies{IsProduction: true}, func(c SessionCookies, w http.ResponseWriter) { c.Clear(w) })

	if cookie.Name != SessionCookieName || cookie.Value != "" || cookie.MaxAge >= 0 || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" || !cookie.Secure {
		t.Errorf("cleared cookie = %+v, want it expired with the same attributes", cookie)
	}
}
