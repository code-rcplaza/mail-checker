package httpserver

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// errorKind maps a domain sentinel to its HTTP status and the key used when
// the error carries no key of its own.
type errorKind struct {
	sentinel    error
	status      int
	fallbackKey domain.MessageKey
}

var errorKinds = []errorKind{
	{domain.ErrInvalid, http.StatusBadRequest, domain.KeyValidationFailed},
	{domain.ErrNotFound, http.StatusNotFound, domain.KeyNotFound},
	{domain.ErrConflict, http.StatusConflict, domain.KeyConflict},
}

var internalErrorKind = errorKind{status: http.StatusInternalServerError, fallbackKey: domain.KeyInternal}

// WriteError answers with the status and i18n key of a domain error. Any other
// error becomes a 500 with common.internal: the real error goes only to the log.
func WriteError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	kind := matchErrorKind(err)
	if kind.status == http.StatusInternalServerError {
		logger.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
	}
	respond(w, r, logger, kind.status, errorResponse{Error: payloadFor(err, kind)})
}

func matchErrorKind(err error) errorKind {
	for _, kind := range errorKinds {
		if errors.Is(err, kind.sentinel) {
			return kind
		}
	}
	return internalErrorKind
}

func payloadFor(err error, kind errorKind) errorPayload {
	var keyed domain.KeyedError
	if kind != internalErrorKind && errors.As(err, &keyed) {
		return errorPayload{Key: keyed.Key, Params: keyed.Params}
	}
	return errorPayload{Key: kind.fallbackKey}
}
