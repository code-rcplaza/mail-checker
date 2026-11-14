package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

type LoginCodes struct {
	queries *queries.Queries
}

func NewLoginCodes(db *sql.DB) LoginCodes {
	return LoginCodes{queries: queries.New(db)}
}

// Save is one upsert, so replacing the previous code is atomic.
func (c LoginCodes) Save(ctx context.Context, code domain.LoginCode) error {
	err := c.queries.SaveLoginCode(ctx, queries.SaveLoginCodeParams{
		UserID:    code.UserID,
		CodeHash:  code.CodeHash,
		ExpiresAt: code.ExpiresAt.Unix(),
		Attempts:  int64(code.Attempts),
	})
	if err != nil {
		return fmt.Errorf("save login code: %w", err)
	}
	return nil
}

func (c LoginCodes) Find(ctx context.Context, userID int64) (domain.LoginCode, error) {
	row, err := c.queries.FindLoginCode(ctx, userID)
	if err != nil {
		return domain.LoginCode{}, lookupError("find login code", err)
	}
	return domain.LoginCode{UserID: row.UserID, CodeHash: row.CodeHash, ExpiresAt: sharedsqlite.TimeFromUnix(row.ExpiresAt), Attempts: int(row.Attempts)}, nil
}

// AddFailedAttempt adds in the database, not read-then-write, so two requests at
// once cannot both see 4 and skip the limit.
func (c LoginCodes) AddFailedAttempt(ctx context.Context, userID int64) (int, error) {
	attempts, err := c.queries.AddFailedAttempt(ctx, userID)
	if err != nil {
		return 0, lookupError("add failed attempt", err)
	}
	return int(attempts), nil
}

func (c LoginCodes) Delete(ctx context.Context, userID int64) error {
	if err := c.queries.DeleteLoginCode(ctx, userID); err != nil {
		return fmt.Errorf("delete login code: %w", err)
	}
	return nil
}
