// Package sqlite opens the single app database, already configured, so no
// repository has to care about durability, foreign keys or concurrency.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite" // registers the "sqlite" driver
)

const (
	driverName = "sqlite"
	dirPerm    = 0o750
)

// pragmas are applied by the driver to every new connection, not just one.
var pragmas = []string{
	"journal_mode(WAL)",
	"foreign_keys(1)",
	"busy_timeout(5000)",
	"synchronous(FULL)",
}

// Open creates the parent directory if needed, opens the database with the
// pragmas applied and a single writer, and checks it is actually usable.
func Open(ctx context.Context, path string) (*sql.DB, error) {
	if err := ensureParentDir(path); err != nil {
		return nil, wrapOpenError(path, err)
	}
	db, err := sql.Open(driverName, dsnWithPragmas(path))
	if err != nil {
		return nil, wrapOpenError(path, err)
	}
	limitToSingleWriter(db)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, wrapOpenError(path, err)
	}
	return db, nil
}

func ensureParentDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), dirPerm)
}

func dsnWithPragmas(path string) string {
	return "file:" + uriPathEscaper.Replace(path) + "?_pragma=" + strings.Join(pragmas, "&_pragma=")
}

// uriPathEscaper keeps the path a plain file name inside the "file:" URI: an
// unescaped "?" or "#" would cut it short and open another file.
var uriPathEscaper = strings.NewReplacer("%", "%25", "?", "%3f", "#", "%23")

func limitToSingleWriter(db *sql.DB) {
	db.SetMaxOpenConns(1)
}

func wrapOpenError(path string, err error) error {
	return fmt.Errorf("open sqlite %q: %w", path, err)
}
