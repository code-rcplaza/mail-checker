package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// healthTimeout is how long the database gets to answer before the app counts as down.
const healthTimeout = time.Second

// HealthCheck tells whether the app can serve; in the app it pings the database.
type HealthCheck func(ctx context.Context) error

// AlwaysHealthy is the check of a router whose health is not what is being tested.
func AlwaysHealthy(context.Context) error { return nil }

type healthResponse struct {
	Status string `json:"status"`
}

// healthHandler answers 503 when the check fails or runs out of time; the
// cause goes only to the log, never the database path or its error.
func (responder Responder) healthHandler(check HealthCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchange := Exchange{Writer: w, Request: r}
		ctx, cancel := context.WithTimeout(r.Context(), healthTimeout)
		defer cancel()
		if err := check(ctx); err != nil {
			responder.Logger.Warn("health check failed", "err", err)
			responder.WriteError(exchange, domain.Unavailable(domain.KeyUnavailable))
			return
		}
		responder.Respond(exchange, Reply{Status: http.StatusOK, Body: healthResponse{Status: "ok"}})
	}
}
