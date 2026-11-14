// Package http is the auth feature's HTTP edge: the session cookie, the
// middleware that requires a session and the auth routes.
package http

import (
	"net/http"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
)

const SessionCookieName = "session"

// SessionCookies writes the session token where JavaScript cannot read it.
// Secure only in production: development runs on plain http.
type SessionCookies struct {
	IsProduction bool
}

func (c SessionCookies) Set(w http.ResponseWriter, token string) {
	http.SetCookie(w, c.cookie(token, int(domain.SessionLifetime.Seconds())))
}

// Clear expires the cookie with the same attributes, or the browser keeps it.
func (c SessionCookies) Clear(w http.ResponseWriter) {
	http.SetCookie(w, c.cookie("", -1))
}

func (c SessionCookies) cookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   c.IsProduction,
		SameSite: http.SameSiteLaxMode,
	}
}
