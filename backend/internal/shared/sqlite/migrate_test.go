package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"testing/fstest"
)

func tableNames(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		names = append(names, name)
	}
	return names
}

func currentVersion(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	var version int64
	if err := db.QueryRow(`SELECT MAX(version_id) FROM goose_db_version`).Scan(&version); err != nil {
		t.Fatalf("read goose version: %v", err)
	}
	return version
}

func TestMigrate_AppliesEmbeddedMigrationsToNewDB(t *testing.T) {
	db := openTempDB(t)

	if err := Migrate(context.Background(), db, EmbeddedMigrations()); err != nil {
		t.Fatalf("Migrate() error = %v, want nil", err)
	}

	got := tableNames(t, db)
	if len(got) != 1 || got[0] != "goose_db_version" {
		t.Errorf("tables = %v, want only [goose_db_version]", got)
	}
}

func TestMigrate_IsIdempotent(t *testing.T) {
	db := openTempDB(t)
	ctx := context.Background()
	if err := Migrate(ctx, db, EmbeddedMigrations()); err != nil {
		t.Fatalf("first Migrate() error = %v, want nil", err)
	}
	versionAfterFirst := currentVersion(t, db)

	if err := Migrate(ctx, db, EmbeddedMigrations()); err != nil {
		t.Fatalf("second Migrate() error = %v, want nil", err)
	}

	if got := currentVersion(t, db); got != versionAfterFirst {
		t.Errorf("version after second run = %d, want %d", got, versionAfterFirst)
	}
}

func TestMigrate_FailsNamingTheBrokenMigration(t *testing.T) {
	db := openTempDB(t)
	migrations := fstest.MapFS{
		"00001_broken.sql": {Data: []byte("-- +goose Up\nTHIS IS NOT SQL;\n")},
	}

	err := Migrate(context.Background(), db, migrations)

	if err == nil {
		t.Fatalf("Migrate() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "00001_broken.sql") {
		t.Errorf("error %q does not name 00001_broken.sql", err)
	}
}

func TestMigrate_RollsBackFailedMigration(t *testing.T) {
	db := openTempDB(t)
	migrations := fstest.MapFS{
		"00001_applied.sql": {Data: []byte("-- +goose Up\nCREATE TABLE applied (id INTEGER);\n")},
		"00002_partial.sql": {Data: []byte("-- +goose Up\nCREATE TABLE partial (id INTEGER);\nTHIS IS NOT SQL;\n")},
	}

	if err := Migrate(context.Background(), db, migrations); err == nil {
		t.Fatalf("Migrate() error = nil, want error")
	}

	tables := strings.Join(tableNames(t, db), ",")
	if !strings.Contains(tables, "applied") {
		t.Errorf("tables = [%s], want 'applied' (the valid migration before the broken one)", tables)
	}
	if strings.Contains(tables, "partial") {
		t.Errorf("tables = [%s], want 'partial' rolled back", tables)
	}
}
