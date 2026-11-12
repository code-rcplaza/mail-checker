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
	router := chi.NewRouter()
	router.Get("/api/health", healthHandler(logger))
	router.NotFound(routeNotFoundHandler(logger))
	return router
}

func healthHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, logger, http.StatusOK, healthResponse{Status: "ok"})
	}
}

func routeNotFoundHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := errorResponse{Error: errorPayload{Key: domain.KeyRouteNotFound}}
		respond(w, r, logger, http.StatusNotFound, body)
	}
}

// respond writes JSON; if that fails (e.g. client gone), the error is logged:
// there is no one left to answer.
func respond(w http.ResponseWriter, r *http.Request, logger *slog.Logger, status int, value any) {
	if err := WriteJSON(w, status, value); err != nil {
		logger.Warn("write response failed", "path", r.URL.Path, "err", err)
	}
}
