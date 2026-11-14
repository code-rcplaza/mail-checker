package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/migrations"
)

func openMigratedDB(t *testing.T) *sql.DB {
	t.Helper()
	db := openTempDB(t)
	if err := Migrate(context.Background(), db, migrations.FS); err != nil {
		t.Fatalf("Migrate() error = %v, want nil", err)
	}
	return db
}

func TestUsersTable_RejectsTheSameEmailInOtherCasing(t *testing.T) {
	db := openMigratedDB(t)
	ctx := context.Background()
	insert := `INSERT INTO users (public_id, name, email, accepted_terms_at, created_at) VALUES (?, 'Ana', ?, 1, 1)`
	if _, err := db.ExecContext(ctx, insert, "user-1", "ana@example.com"); err != nil {
		t.Fatalf("first insert error = %v, want nil", err)
	}

	if _, err := db.ExecContext(ctx, insert, "user-2", "ANA@Example.com"); err == nil {
		t.Errorf("second insert with another casing error = nil, want a UNIQUE violation")
	}
}

func TestUsersTable_EveryColumnIsRequired(t *testing.T) {
	db := openMigratedDB(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `INSERT INTO users (public_id, name, email, accepted_terms_at, created_at) VALUES ('ok', 'Ana', 'ok@x.cl', 1, 1)`); err != nil {
		t.Fatalf("complete insert error = %v, want nil", err)
	}
	cases := map[string]string{
		"public_id":         `INSERT INTO users (name, email, accepted_terms_at, created_at) VALUES ('Ana', 'a@x.cl', 1, 1)`,
		"name":              `INSERT INTO users (public_id, email, accepted_terms_at, created_at) VALUES ('u', 'a@x.cl', 1, 1)`,
		"email":             `INSERT INTO users (public_id, name, accepted_terms_at, created_at) VALUES ('u', 'Ana', 1, 1)`,
		"accepted_terms_at": `INSERT INTO users (public_id, name, email, created_at) VALUES ('u', 'Ana', 'a@x.cl', 1)`,
		"created_at":        `INSERT INTO users (public_id, name, email, accepted_terms_at) VALUES ('u', 'Ana', 'a@x.cl', 1)`,
	}
	for column, insert := range cases {
		t.Run(column, func(t *testing.T) {
			if _, err := db.ExecContext(ctx, insert); err == nil {
				t.Errorf("insert without %s error = nil, want NOT NULL violation", column)
			}
		})
	}
}
