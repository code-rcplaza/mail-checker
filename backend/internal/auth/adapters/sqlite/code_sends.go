package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

type CodeSends struct {
	queries *queries.Queries
}

func NewCodeSends(db *sql.DB) CodeSends {
	return CodeSends{queries: queries.New(db)}
}

func (s CodeSends) ListSince(ctx context.Context, email string, since time.Time) ([]time.Time, error) {
	rows, err := s.queries.ListCodeSendsSince(ctx, queries.ListCodeSendsSinceParams{Email: email, SentAt: since.Unix()})
	if err != nil {
		return nil, fmt.Errorf("list code sends: %w", err)
	}
	sends := make([]time.Time, 0, len(rows))
	for _, sentAt := range rows {
		sends = append(sends, sharedsqlite.TimeFromUnix(sentAt))
	}
	return sends, nil
}

func (s CodeSends) Record(ctx context.Context, send domain.CodeSend) error {
	if err := s.queries.RecordCodeSend(ctx, queries.RecordCodeSendParams{Email: send.Email, SentAt: send.SentAt.Unix()}); err != nil {
		return fmt.Errorf("record code send: %w", err)
	}
	return nil
}
