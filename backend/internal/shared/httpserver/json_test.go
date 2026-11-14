package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type accountRequest struct {
	Email string `json:"email"`
}

func decodeBody(body string) (accountRequest, error) {
	var dst accountRequest
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	err := DecodeJSON(Exchange{Writer: httptest.NewRecorder(), Request: req}, &dst)
	return dst, err
}

func TestDecodeJSON_DecodesValidBody(t *testing.T) {
	got, err := decodeBody(`{"email":"user@example.com"}`)

	if err != nil {
		t.Fatalf("DecodeJSON() error = %v, want nil", err)
	}
	if got.Email != "user@example.com" {
		t.Errorf("Email = %q, want user@example.com", got.Email)
	}
}

func TestDecodeJSON_RejectsInvalidBody(t *testing.T) {
	oversized := `{"email":"` + strings.Repeat("a", 1<<20) + `"}`
	cases := map[string]string{
		"malformed":       `{"email":`,
		"not json":        `hola`,
		"empty":           ``,
		"unknown field":   `{"email":"a@b.c","nope":1}`,
		"larger than 1MB": oversized,
		"trailing data":   `{"email":"a@b.c"}{"email":"x"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := decodeBody(body)

			if !errors.Is(err, domain.ErrInvalid) {
				t.Fatalf("DecodeJSON() error = %v, want ErrInvalid", err)
			}
			var keyed domain.KeyedError
			if !errors.As(err, &keyed) || keyed.Key != domain.KeyInvalidJSON {
				t.Errorf("error key = %q, want %q", keyed.Key, domain.KeyInvalidJSON)
			}
		})
	}
}

func TestWriteJSON_SetsStatusAndContentType(t *testing.T) {
	recorder := httptest.NewRecorder()

	err := WriteJSON(recorder, Reply{Status: http.StatusCreated, Body: accountRequest{Email: "user@example.com"}})

	if err != nil {
		t.Fatalf("WriteJSON() error = %v, want nil", err)
	}
	if recorder.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body accountRequest
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not JSON: %v", recorder.Body.String(), err)
	}
	if body.Email != "user@example.com" {
		t.Errorf("body.Email = %q, want user@example.com", body.Email)
	}
}
