package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

type Sessions struct {
	queries *queries.Queries
}

func NewSessions(db *sql.DB) Sessions {
	return Sessions{queries: queries.New(db)}
}

func (s Sessions) Create(ctx context.Context, session domain.Session) error {
	err := s.queries.CreateSession(ctx, queries.CreateSessionParams{UserID: session.UserID, TokenHash: session.TokenHash, ExpiresAt: session.ExpiresAt.Unix()})
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (s Sessions) FindByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error) {
	row, err := s.queries.FindSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return domain.Session{}, lookupError("find session", err)
	}
	return domain.Session{UserID: row.UserID, TokenHash: row.TokenHash, ExpiresAt: sharedsqlite.TimeFromUnix(row.ExpiresAt)}, nil
}

func (s Sessions) Delete(ctx context.Context, tokenHash string) error {
	if err := s.queries.DeleteSession(ctx, tokenHash); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
