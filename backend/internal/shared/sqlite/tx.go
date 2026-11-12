package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WithTx runs fn inside a transaction: commit if fn returns nil, rollback if it
// returns an error or panics (the panic is re-raised).
//
// Inside fn always use tx, never db: the database has a single connection, so a
// query on db would wait forever for the connection the transaction holds.
func WithTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackAndRepanic(tx)
	if err := fn(tx); err != nil {
		return rollbackWith(tx, err)
	}
	return commit(tx)
}

// rollbackAndRepanic re-raises the original panic; if the rollback also fails,
// the panic carries both so no error is lost.
func rollbackAndRepanic(tx *sql.Tx) {
	recovered := recover()
	if recovered == nil {
		return
	}
	if err := tx.Rollback(); err != nil {
		panic(fmt.Errorf("%v (rollback also failed: %w)", recovered, err))
	}
	panic(recovered)
}

func rollbackWith(tx *sql.Tx, cause error) error {
	if err := tx.Rollback(); err != nil {
		return errors.Join(cause, fmt.Errorf("rollback tx: %w", err))
	}
	return cause
}

func commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
