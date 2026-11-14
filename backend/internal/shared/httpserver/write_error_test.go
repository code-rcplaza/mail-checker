package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type errorBody struct {
	Error struct {
		Key    string         `json:"key"`
		Params map[string]any `json:"params"`
	} `json:"error"`
}

func writeErrorWithLog(t *testing.T, err error) (*httptest.ResponseRecorder, errorBody, string) {
	t.Helper()
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	recorder := httptest.NewRecorder()

	Responder{Logger: logger}.WriteError(Exchange{Writer: recorder, Request: httptest.NewRequest(http.MethodGet, "/api/x", nil)}, err)

	var body errorBody
	if jsonErr := json.Unmarshal(recorder.Body.Bytes(), &body); jsonErr != nil {
		t.Fatalf("body %q is not JSON: %v", recorder.Body.String(), jsonErr)
	}
	return recorder, body, logs.String()
}

func assertStatusAndKey(t *testing.T, recorder *httptest.ResponseRecorder, body errorBody, status int, key domain.MessageKey) {
	t.Helper()
	if recorder.Code != status {
		t.Errorf("status = %d, want %d", recorder.Code, status)
	}
	if body.Error.Key != string(key) {
		t.Errorf("error.key = %q, want %q", body.Error.Key, key)
	}
}

func TestWriteError_InvalidReturns400WithKeyAndParams(t *testing.T) {
	recorder, body, _ := writeErrorWithLog(t, domain.Invalid("statements.bank_not_detected", map[string]any{"count": 2}))

	assertStatusAndKey(t, recorder, body, http.StatusBadRequest, "statements.bank_not_detected")
	if body.Error.Params["count"] != float64(2) {
		t.Errorf("error.params = %v, want {count: 2}", body.Error.Params)
	}
}

func TestWriteError_KeyedNotFoundAndConflictKeepTheirKey(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		status int
		key    domain.MessageKey
	}{
		{"NotFound", domain.NotFound("accounts.not_found"), http.StatusNotFound, "accounts.not_found"},
		{"Conflict", domain.Conflict("accounts.email_taken"), http.StatusConflict, "accounts.email_taken"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder, body, _ := writeErrorWithLog(t, tc.err)
			assertStatusAndKey(t, recorder, body, tc.status, tc.key)
		})
	}
}

func TestWriteError_BareSentinelsUseCommonKeys(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		status int
		key    domain.MessageKey
	}{
		{"ErrNotFound", domain.ErrNotFound, http.StatusNotFound, domain.KeyNotFound},
		{"ErrConflict", domain.ErrConflict, http.StatusConflict, domain.KeyConflict},
		{"ErrInvalid", domain.ErrInvalid, http.StatusBadRequest, domain.KeyValidationFailed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder, body, _ := writeErrorWithLog(t, tc.err)
			assertStatusAndKey(t, recorder, body, tc.status, tc.key)
		})
	}
}

func TestWriteError_WrappedDomainErrorKeepsStatusAndKey(t *testing.T) {
	wrapped := fmt.Errorf("remove account: %w", domain.NotFound("accounts.not_found"))

	recorder, body, _ := writeErrorWithLog(t, wrapped)

	assertStatusAndKey(t, recorder, body, http.StatusNotFound, "accounts.not_found")
}

func TestWriteError_UnexpectedErrorReturns500WithoutDetailsAndLogsIt(t *testing.T) {
	recorder, body, logs := writeErrorWithLog(t, errors.New("disk exploded at /secret/path"))

	assertStatusAndKey(t, recorder, body, http.StatusInternalServerError, domain.KeyInternal)
	if strings.Contains(recorder.Body.String(), "disk exploded") {
		t.Errorf("body %q leaks the internal error", recorder.Body.String())
	}
	if !strings.Contains(logs, "disk exploded at /secret/path") {
		t.Errorf("logs %q do not contain the internal error", logs)
	}
}

func TestWriteError_OmitsParamsWhenEmpty(t *testing.T) {
	recorder, body, _ := writeErrorWithLog(t, domain.NotFound("accounts.not_found"))

	assertStatusAndKey(t, recorder, body, http.StatusNotFound, "accounts.not_found")
	if strings.Contains(recorder.Body.String(), "params") {
		t.Errorf("body %q includes params, want them omitted", recorder.Body.String())
	}
}
