package sqlite_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/demolemo/meth/sqlite"
)

func TestOpen(t *testing.T) {
	cases := []struct {
		name string
		file string
		abs  bool
	}{
		{name: "relative", file: "w.db"},
		{name: "absolute", file: "w.db", abs: true},
		{name: "question mark", file: "w?y.db"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			chdir(t, dir)

			path := c.file
			if c.abs {
				path = filepath.Join(dir, c.file)
			}
			db, err := sqlite.Open(path)
			if err != nil {
				t.Fatalf("Open(%q): %v", path, err)
			}
			defer db.Close()

			want, err := os.Stat(filepath.Join(dir, c.file))
			if err != nil {
				t.Fatalf("the file meant to be created: %v", err)
			}
			got, err := os.Stat(openedFile(t, db))
			if err != nil {
				t.Fatalf("the file sqlite opened: %v", err)
			}
			if !os.SameFile(want, got) {
				t.Errorf("sqlite opened %s, want %s", got.Name(), want.Name())
			}

			var fk int
			if err := db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk); err != nil {
				t.Fatal(err)
			}
			if fk != 1 {
				t.Errorf("PRAGMA foreign_keys = %d, want 1", fk)
			}
		})
	}
}

// openedFile asks the connection which file it is actually on.
func openedFile(t *testing.T, db *sql.DB) string {
	t.Helper()
	var seq int
	var name, file string
	if err := db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &file); err != nil {
		t.Fatal(err)
	}
	return file
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Fatal(err)
		}
	})
}
