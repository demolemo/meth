// Package sqlite stores metrics and their values in a single sqlite file.
package sqlite

import (
	"database/sql"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	// this package is the one that hardcodes the "sqlite3" driver name, so it
	// registers the driver too. Leaving that to the caller means every new
	// consumer gets `sql: unknown driver "sqlite3"` at runtime until someone
	// remembers the blank import.
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

// Open connects to the database at path, creating the file and its parent
// directory if they are missing, and applies schema.sql. Safe to call against
// an existing database: every statement in the schema is IF NOT EXISTS.
//
// Foreign keys are off by default in sqlite and have to be asked for per
// connection, which is what metric_values.metric_id REFERENCES relies on.
func Open(path string) (*sql.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("no database path")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("database directory: %w", err)
	}

	// the file: URI form, not the bare "path?params" string: with a ? in the
	// path that one splits in the wrong place, silently opening a different
	// file with foreign keys off
	dsn := url.URL{Scheme: "file", Path: path, RawQuery: "_foreign_keys=on"}
	db, err := sql.Open("sqlite3", dsn.String())
	if err != nil {
		return nil, err
	}
	// sql.Open is lazy, nothing has touched the file yet
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: applying schema: %w", path, err)
	}
	return db, nil
}
