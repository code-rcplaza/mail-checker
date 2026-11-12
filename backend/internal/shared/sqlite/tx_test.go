package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
)

func openDBWithItemsTable(t *testing.T) *sql.DB {
	t.Helper()
	db := openTempDB(t)
	if _, err := db.Exec(`CREATE TABLE items (name TEXT NOT NULL)`); err != nil {
		t.Fatalf("create items table: %v", err)
	}
	return db
}

func insertItem(tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO items (name) VALUES ('netflix')`)
	return err
}

func countItems(t *testing.T, db *sql.DB) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM items`).Scan(&count); err != nil {
		t.Fatalf("count items: %v", err)
	}
	return count
}

func TestWithTx_CommitsWhenFnSucceeds(t *testing.T) {
	db := openDBWithItemsTable(t)

	err := WithTx(context.Background(), db, insertItem)

	if err != nil {
		t.Fatalf("WithTx() error = %v, want nil", err)
	}
	if got := countItems(t, db); got != 1 {
		t.Errorf("items = %d, want 1", got)
	}
}

func TestWithTx_RollsBackAndReturnsFnError(t *testing.T) {
	db := openDBWithItemsTable(t)
	errFn := errors.New("fn failed")

	err := WithTx(context.Background(), db, func(tx *sql.Tx) error {
		if err := insertItem(tx); err != nil {
			return err
		}
		return errFn
	})

	if !errors.Is(err, errFn) {
		t.Fatalf("WithTx() error = %v, want errors.Is(err, errFn)", err)
	}
	if got := countItems(t, db); got != 0 {
		t.Errorf("items = %d, want 0 (rolled back)", got)
	}
}

func TestWithTx_RollsBackAndRepanics(t *testing.T) {
	db := openDBWithItemsTable(t)
	const panicValue = "boom"

	recovered := func() (r any) {
		defer func() { r = recover() }()
		_ = WithTx(context.Background(), db, func(tx *sql.Tx) error {
			if err := insertItem(tx); err != nil {
				return err
			}
			panic(panicValue)
		})
		return nil
	}()

	if recovered != panicValue {
		t.Fatalf("recovered = %v, want %q (panic must propagate)", recovered, panicValue)
	}
	if got := countItems(t, db); got != 0 {
		t.Errorf("items = %d, want 0 (rolled back)", got)
	}
}

func TestWithTx_RepanicsWithRollbackErrorWhenRollbackFails(t *testing.T) {
	db := openDBWithItemsTable(t)
	const panicValue = "boom"

	recovered := func() (r any) {
		defer func() { r = recover() }()
		_ = WithTx(context.Background(), db, func(tx *sql.Tx) error {
			_ = tx.Rollback() // forces the deferred rollback to fail with sql.ErrTxDone
			panic(panicValue)
		})
		return nil
	}()

	err, isError := recovered.(error)
	if !isError {
		t.Fatalf("recovered = %#v, want an error carrying the rollback failure", recovered)
	}
	if !errors.Is(err, sql.ErrTxDone) {
		t.Errorf("recovered error %q does not wrap sql.ErrTxDone", err)
	}
	if !strings.Contains(err.Error(), panicValue) {
		t.Errorf("recovered error %q does not keep the original panic %q", err, panicValue)
	}
}

func TestWithTx_DoesNotRunFnWhenContextCanceled(t *testing.T) {
	db := openDBWithItemsTable(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fnRan := false

	err := WithTx(ctx, db, func(tx *sql.Tx) error {
		fnRan = true
		return nil
	})

	if err == nil {
		t.Errorf("WithTx() error = nil, want error for canceled context")
	}
	if fnRan {
		t.Errorf("fn ran, want not run for canceled context")
	}
}
