package sqlite

import (
	"database/sql"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// wrapped adds the step to a store error and leaves nil alone, so a write reads as one line.
func wrapped(step string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", step, err)
	}
	return nil
}

// notFoundWhenNone turns "no row matched" into ErrNotFound: the row is unknown or another user's.
func notFoundWhenNone(step string, rows int64) error {
	if rows == 0 {
		return fmt.Errorf("%s: %w", step, shared.ErrNotFound)
	}
	return nil
}

// nullWhenZero stores "not given" as NULL, the way the domain uses the zero value.
func nullWhenZero(value int64) sql.NullInt64 { return sql.NullInt64{Int64: value, Valid: value != 0} }

func nullWhenEmpty(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}
