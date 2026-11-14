package mailer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

var _ ports.CodeSender = ResendSender{}

const (
	testAPIKey = "re_test_key"
	testFrom   = "Fokko <codigos@example.com>"
)

type receivedEmail struct {
	method, path, authorization, contentType string
	body                                     resendEmail
}

// resendServer answers every request with status and keeps what it received.
func resendServer(t *testing.T, status int) (*httptest.Server, *[]receivedEmail) {
	t.Helper()
	received := &[]receivedEmail{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body resendEmail
		_ = json.NewDecoder(r.Body).Decode(&body)
		*received = append(*received, receivedEmail{r.Method, r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("Content-Type"), body})
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"message":"secret detail ana@example.com 042917"}`))
	}))
	t.Cleanup(server.Close)
	return server, received
}

func senderFor(baseURL string, timeout time.Duration) ResendSender {
	return NewResendSender(ResendSettings{BaseURL: baseURL, APIKey: testAPIKey, From: testFrom, Timeout: timeout})
}

func TestResendSender_PostsTheEmailWithTheCode(t *testing.T) {
	server, received := resendServer(t, http.StatusOK)

	err := senderFor(server.URL, time.Second).Send(context.Background(), "ana@example.com", "042917")

	if err != nil || len(*received) != 1 {
		t.Fatalf("err = %v, requests = %d; want one request without error", err, len(*received))
	}
	got := (*received)[0]
	if got.method != http.MethodPost || got.path != "/emails" || got.authorization != "Bearer "+testAPIKey || got.contentType != "application/json" {
		t.Errorf("request = %s %s auth=%q type=%q; want POST /emails with the bearer key and JSON", got.method, got.path, got.authorization, got.contentType)
	}
	if got.body.From != testFrom || len(got.body.To) != 1 || got.body.To[0] != "ana@example.com" || got.body.Subject != codeSubject || !strings.Contains(got.body.Text, "042917") {
		t.Errorf("body = %+v; want from, to, subject and a text with the code", got.body)
	}
}

func TestResendSender_AnyTwoHundredIsSuccess(t *testing.T) {
	server, _ := resendServer(t, http.StatusAccepted)

	if err := senderFor(server.URL, time.Second).Send(context.Background(), "ana@example.com", "042917"); err != nil {
		t.Errorf("Send error = %v, want nil", err)
	}
}

func TestResendSender_RejectedSendReportsOnlyTheStatus(t *testing.T) {
	server, _ := resendServer(t, http.StatusUnprocessableEntity)

	err := senderFor(server.URL, time.Second).Send(context.Background(), "ana@example.com", "042917")

	if err == nil || !strings.Contains(err.Error(), "422") {
		t.Fatalf("err = %v, want an error with the status 422", err)
	}
	for _, private := range []string{"secret detail", "ana@example.com", "042917", testAPIKey} {
		if strings.Contains(err.Error(), private) {
			t.Errorf("error %q contains %q", err, private)
		}
	}
}

func TestResendSender_StopsAtTheTimeout(t *testing.T) {
	// The handler never answers; it is released when the test ends so Close does not wait for it.
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	started := time.Now()

	err := senderFor(server.URL, 50*time.Millisecond).Send(context.Background(), "ana@example.com", "042917")

	if err == nil || time.Since(started) > 2*time.Second {
		t.Errorf("err = %v after %v; want a timeout error quickly", err, time.Since(started))
	}
}

func TestResendSender_CanceledContextSendsNothing(t *testing.T) {
	server, received := resendServer(t, http.StatusOK)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := senderFor(server.URL, time.Second).Send(ctx, "ana@example.com", "042917")

	if !errors.Is(err, context.Canceled) || len(*received) != 0 {
		t.Errorf("err = %v, requests = %d; want context.Canceled and nothing sent", err, len(*received))
	}
}
