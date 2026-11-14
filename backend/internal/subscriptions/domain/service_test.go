package domain

import (
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestNewService(t *testing.T) {
	cases := []struct {
		name      string
		cancelURL string
		wantErr   error
	}{
		{"https link is kept", "https://example.com/account/cancel", nil},
		{"missing link is not an error", "", nil},
		{"http link is rejected", "http://example.com/account/cancel", shared.ErrInvalid},
		{"odd scheme is rejected", "javascript:alert(1)", shared.ErrInvalid},
		{"https without a host is rejected", "https://", shared.ErrInvalid},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewService(Service{Name: "Streaming Co", CancelURL: tc.cancelURL})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewService error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				assertKey(t, err, KeyCancelURLInvalid)
				return
			}
			want := Service{Name: "Streaming Co", CancelURL: tc.cancelURL}
			if got != want {
				t.Errorf("NewService = %+v, want %+v", got, want)
			}
		})
	}
}

func assertKey(t *testing.T, err error, want shared.MessageKey) {
	t.Helper()
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) {
		t.Fatalf("error %v is not a KeyedError", err)
	}
	if keyed.Key != want {
		t.Errorf("key = %q, want %q", keyed.Key, want)
	}
}
