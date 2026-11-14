// Package httpserver holds the shared HTTP pieces: router, JSON helpers and the
// error shape. Each feature only registers its routes on the router.
package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// NewRouter returns the base router: request log, panic recovery and the
// origin check on every route, /api/health (answered by check) and JSON 404s.
// It sends no CORS headers: the front is served from the same origin.
func NewRouter(logger *slog.Logger, check HealthCheck) *chi.Mux {
	responder := Responder{Logger: logger}
	router := chi.NewRouter()
	router.Use(responder.LogRequests, responder.RecoverPanics, responder.CheckOrigin)
	router.Get("/api/health", responder.healthHandler(check))
	router.NotFound(responder.routeNotFoundHandler())
	return router
}

// Responder writes responses and logs what cannot reach the client.
type Responder struct {
	Logger *slog.Logger
}

func (responder Responder) routeNotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responder.WriteError(Exchange{Writer: w, Request: r}, domain.NotFound(domain.KeyRouteNotFound))
	}
}

// Respond writes JSON; if that fails (e.g. client gone), the error is logged:
// there is no one left to answer.
func (responder Responder) Respond(exchange Exchange, reply Reply) {
	if err := WriteJSON(exchange.Writer, reply); err != nil {
		responder.Logger.Warn("write response failed", "path", exchange.Request.URL.Path, "err", err)
	}
}
