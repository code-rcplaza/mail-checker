package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func openTempDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := Open(context.Background(), filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatalf("Open() error = %v, want nil", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestOpen_AppliesPragmas(t *testing.T) {
	db := openTempDB(t)
	cases := []struct {
		pragma string
		want   string
	}{
		{"journal_mode", "wal"},
		{"foreign_keys", "1"},
		{"busy_timeout", "5000"},
		{"synchronous", "2"},
	}
	for _, tc := range cases {
		t.Run(tc.pragma, func(t *testing.T) {
			var got string
			if err := db.QueryRow(fmt.Sprintf("PRAGMA %s", tc.pragma)).Scan(&got); err != nil {
				t.Fatalf("PRAGMA %s error = %v", tc.pragma, err)
			}
			if got != tc.want {
				t.Errorf("PRAGMA %s = %q, want %q", tc.pragma, got, tc.want)
			}
		})
	}
}

func TestOpen_LimitsToSingleConnection(t *testing.T) {
	db := openTempDB(t)

	if got := db.Stats().MaxOpenConnections; got != 1 {
		t.Errorf("MaxOpenConnections = %d, want 1", got)
	}
}

func TestOpen_CreatesMissingParentDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "app.db")

	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open() error = %v, want nil", err)
	}
	defer db.Close()

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("parent dir not created: %v", err)
	}
}

func TestOpen_FailsWithPathWhenNotWritable(t *testing.T) {
	readOnlyDir := filepath.Join(t.TempDir(), "readonly")
	if err := os.Mkdir(readOnlyDir, 0o500); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(readOnlyDir, "app.db")

	db, err := Open(context.Background(), path)
	if err == nil {
		db.Close()
		t.Fatalf("Open() error = nil, want error for read-only dir")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %q does not include path %q", err, path)
	}
}
