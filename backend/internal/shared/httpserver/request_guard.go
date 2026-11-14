package httpserver

import (
	"context"
	"net/http"
	"net/url"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type userIDKey struct{}

// WithUserID is how the session middleware hands the user to the handlers. The
// user travels in the context only inside the HTTP layer; handlers pass it on
// to use cases as a parameter.
func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFrom fails instead of answering a zero user when no session put one.
func UserIDFrom(ctx context.Context) (int64, error) {
	userID, found := ctx.Value(userIDKey{}).(int64)
	if !found || userID <= 0 {
		return 0, domain.Unauthorized(domain.KeyUnauthorized)
	}
	return userID, nil
}

// CheckOrigin rejects a request that changes data from another site. A request
// without Origin passes: browsers send it on every POST, and the SameSite=Lax
// cookie does not travel from another site anyway.
func (responder Responder) CheckOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if changesData(r.Method) && !isSameOrigin(r) {
			responder.WriteError(Exchange{Writer: w, Request: r}, domain.Forbidden(domain.KeyForbidden))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func changesData(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// isSameOrigin compares the Origin with the Host the request came to.
func isSameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Host == r.Host
}
