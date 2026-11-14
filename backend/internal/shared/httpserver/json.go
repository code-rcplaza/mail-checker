package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	maxBodyBytes = 1 << 20 // 1 MB
	contentType  = "application/json"
)

// Exchange is one HTTP request and the writer for its response; they always
// travel together.
type Exchange struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

// Reply is a response to write: its status and the value sent as JSON.
type Reply struct {
	Status int
	Body   any
}

// DecodeJSON reads a single JSON object into dst. Malformed bodies, unknown
// fields, trailing data or bodies over 1 MB fail with common.invalid_json.
func DecodeJSON(exchange Exchange, dst any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(exchange.Writer, exchange.Request.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return invalidJSON(err)
	}
	if err := ensureNoTrailingData(decoder); err != nil {
		return invalidJSON(err)
	}
	return nil
}

// WriteJSON marshals value first, so an unencodable value never leaves a
// half-written response, then writes status and body.
func WriteJSON(w http.ResponseWriter, reply Reply) error {
	body, err := json.Marshal(reply.Body)
	if err != nil {
		return fmt.Errorf("encode json response: %w", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(reply.Status)
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("write json response: %w", err)
	}
	return nil
}

func ensureNoTrailingData(decoder *json.Decoder) error {
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON object")
	}
	return nil
}

// invalidJSON keeps the cause in the message for logs; only the key reaches the user.
func invalidJSON(cause error) error {
	return fmt.Errorf("decode json: %v: %w", cause, domain.Invalid(domain.KeyInvalidJSON, nil))
}
