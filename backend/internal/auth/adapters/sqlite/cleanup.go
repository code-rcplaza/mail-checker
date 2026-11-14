package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

type ExpiredAccess struct {
	db      *sql.DB
	queries *queries.Queries
}

func NewExpiredAccess(db *sql.DB) ExpiredAccess {
	return ExpiredAccess{db: db, queries: queries.New(db)}
}

// DeleteExpired clears the three tables in one transaction.
func (e ExpiredAccess) DeleteExpired(ctx context.Context, cutoff domain.CleanupCutoff) error {
	return sharedsqlite.WithTx(ctx, e.db, func(tx *sql.Tx) error {
		q := e.queries.WithTx(tx)
		if err := q.DeleteExpiredLoginCodes(ctx, cutoff.ExpiredBy.Unix()); err != nil {
			return fmt.Errorf("delete expired login codes: %w", err)
		}
		if err := q.DeleteExpiredSessions(ctx, cutoff.ExpiredBy.Unix()); err != nil {
			return fmt.Errorf("delete expired sessions: %w", err)
		}
		if err := q.DeleteOldCodeSends(ctx, cutoff.SentBy.Unix()); err != nil {
			return fmt.Errorf("delete old code sends: %w", err)
		}
		return nil
	})
}
