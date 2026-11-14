package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ResendBaseURL is the Resend API; tests point the sender to an httptest server.
const ResendBaseURL = "https://api.resend.com"

const (
	codeSubject   = "Tu código de acceso a Fokko"
	codeTextShape = "Tu código es %s. Vence en 10 minutos. Si no lo pediste, ignora este correo."
)

type resendEmail struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

type ResendSettings struct {
	BaseURL string
	APIKey  string
	From    string
	Timeout time.Duration
}

// ResendSender mails the code through the Resend HTTP API, without SDK.
type ResendSender struct {
	settings ResendSettings
	client   *http.Client
}

func NewResendSender(settings ResendSettings) ResendSender {
	return ResendSender{settings: settings, client: &http.Client{Timeout: settings.Timeout}}
}

func (s ResendSender) Send(ctx context.Context, email, code string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("send code with resend: %w", err)
	}
	request, err := s.newRequest(ctx, s.codeEmail(email, code))
	if err != nil {
		return err
	}
	return s.post(request)
}

func (s ResendSender) codeEmail(to, code string) resendEmail {
	return resendEmail{From: s.settings.From, To: []string{to}, Subject: codeSubject, Text: fmt.Sprintf(codeTextShape, code)}
}

func (s ResendSender) newRequest(ctx context.Context, email resendEmail) (*http.Request, error) {
	body, err := json.Marshal(email)
	if err != nil {
		return nil, fmt.Errorf("encode resend email: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.settings.BaseURL+"/emails", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build resend request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+s.settings.APIKey)
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

// post reports a rejection only by its status: the body may echo the email or the code.
func (s ResendSender) post(request *http.Request) error {
	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("send code with resend: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("send code with resend: status %d", response.StatusCode)
	}
	return nil
}
