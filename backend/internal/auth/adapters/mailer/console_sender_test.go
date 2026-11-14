package mailer

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

var _ ports.CodeSender = ConsoleSender{}

func TestConsoleSender_WritesTheEmailAndTheCode(t *testing.T) {
	var console bytes.Buffer

	err := ConsoleSender{Writer: &console}.Send(context.Background(), "ana@example.com", "042917")

	if err != nil || !strings.Contains(console.String(), "ana@example.com") || !strings.Contains(console.String(), "042917") {
		t.Errorf("console = %q, err = %v; want the email and the code", console.String(), err)
	}
}

func TestConsoleSender_DoesNotLog(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(previous) })

	err := ConsoleSender{Writer: &bytes.Buffer{}}.Send(context.Background(), "ana@example.com", "042917")

	if err != nil || logs.Len() != 0 {
		t.Errorf("logs = %q, err = %v; want nothing logged", logs.String(), err)
	}
}
