// Package httpserver holds the shared HTTP pieces: router, JSON helpers and the
// error shape. Each feature only registers its routes on the router.
package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type healthResponse struct {
	Status string `json:"status"`
}

// NewRouter returns the base router with /api/health and JSON 404s.
func NewRouter(logger *slog.Logger) *chi.Mux {
	responder := Responder{Logger: logger}
	router := chi.NewRouter()
	router.Get("/api/health", responder.healthHandler())
	router.NotFound(responder.routeNotFoundHandler())
	return router
}

// Responder writes responses and logs what cannot reach the client.
type Responder struct {
	Logger *slog.Logger
}

func (responder Responder) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responder.Respond(Exchange{Writer: w, Request: r}, Reply{Status: http.StatusOK, Body: healthResponse{Status: "ok"}})
	}
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
