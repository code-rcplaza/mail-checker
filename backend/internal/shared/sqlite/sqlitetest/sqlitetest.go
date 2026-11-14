// Package sqlitetest gives repository tests a real, migrated database in a
// temporary folder, so every repo is tested against the true schema.
package sqlitetest

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/migrations"
)

func OpenMigrated(t testing.TB) *sql.DB {
	t.Helper()
	ctx := context.Background()
	db, err := sqlite.Open(ctx, filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatalf("open temporary database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := sqlite.Migrate(ctx, db, migrations.FS); err != nil {
		t.Fatalf("migrate temporary database: %v", err)
	}
	return db
}

// Users inserts the owners a repository test needs, so features other than
// auth do not import auth to get one.
type Users struct{ DB *sql.DB }

// Insert adds a user with that email and returns its id.
func (u Users) Insert(t testing.TB, email string) int64 {
	t.Helper()
	result, err := u.DB.ExecContext(context.Background(),
		`INSERT INTO users (public_id, name, email, accepted_terms_at, created_at) VALUES (?, 'Test', ?, 1, 1)`, "public-"+email, email)
	if err != nil {
		t.Fatalf("insert user %s: %v", email, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("user id of %s: %v", email, err)
	}
	return id
}
