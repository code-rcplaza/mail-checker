package mailer

import (
	"context"
	"fmt"
	"io"
)

// ConsoleSender is for development only: it writes the code to the console
// instead of mailing it. It never goes through slog, because logs must not
// carry codes or emails; main.go refuses it in production.
type ConsoleSender struct {
	Writer io.Writer
}

func (s ConsoleSender) Send(_ context.Context, email, code string) error {
	if _, err := fmt.Fprintf(s.Writer, "login code for %s: %s\n", email, code); err != nil {
		return fmt.Errorf("write login code to console: %w", err)
	}
	return nil
}
