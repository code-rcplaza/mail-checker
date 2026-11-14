package http

import (
	"net/http"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
)

// RequireSession lets a request through only with a valid session cookie and
// hands its user to the handler. It goes on the protected route group only, so
// /api/health, /api/auth/* and unknown routes (404) never ask for a session.
type RequireSession struct {
	Validate  usecases.ValidateSession
	Responder httpserver.Responder
}

func (s RequireSession) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.userOf(r)
		if err != nil {
			s.Responder.WriteError(httpserver.Exchange{Writer: w, Request: r}, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(httpserver.WithUserID(r.Context(), userID)))
	})
}

// userOf answers the same 401 with no cookie, an expired one or an unknown one.
func (s RequireSession) userOf(r *http.Request) (int64, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return 0, domain.Unauthorized(domain.KeyUnauthorized)
	}
	return s.Validate.Execute(r.Context(), cookie.Value)
}
