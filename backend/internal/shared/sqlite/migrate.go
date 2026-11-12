package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

// EmbeddedMigrations returns the app migrations with the "migrations/" prefix
// stripped, as goose expects them at the root.
func EmbeddedMigrations() fs.FS {
	migrations, err := fs.Sub(embeddedMigrations, "migrations")
	if err != nil {
		panic(fmt.Sprintf("embedded migrations: %v", err)) // only if the embed path is wrong: a build bug
	}
	return migrations
}

// Migrate applies every pending migration. Each SQL migration runs in its own
// transaction; a failure stops the app and names the broken file.
func Migrate(ctx context.Context, db *sql.DB, migrations fs.FS) error {
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrations)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return wrapMigrationError(err)
	}
	return nil
}

func wrapMigrationError(err error) error {
	var partial *goose.PartialError
	if errors.As(err, &partial) && partial.Failed != nil {
		return fmt.Errorf("migrate %q: %w", partial.Failed.Source.Path, err)
	}
	return fmt.Errorf("migrate: %w", err)
}
