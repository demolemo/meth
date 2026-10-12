// Package config resolves where meth keeps its data. Nothing here touches the
// database, so `meth config` can answer "where is my data" without creating it.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	appDir     = "meth"
	configFile = "config.toml"
	dbFile     = "meth.db"
)

// Config is the whole config file. Add a field here and a line to the template
// in Write, or `meth` will reject the key as unknown.
type Config struct {
	// DB is the sqlite file holding every metric and value.
	DB string `toml:"db"`

	// File is where this config was read from, empty when no file exists.
	File string `toml:"-"`
}

// Path is the config file location, $XDG_CONFIG_HOME/meth/config.toml falling
// back to ~/.config/meth/config.toml.
//
// Not os.UserConfigDir: on darwin that returns ~/Library/Application Support
// and ignores XDG_CONFIG_HOME, which is not where anyone goes looking for the
// dotfiles of a terminal tool.
func Path() (string, error) {
	return xdgDir("XDG_CONFIG_HOME", ".config", configFile)
}

// DefaultDB is the database used when no config file sets one,
// $XDG_DATA_HOME/meth/meth.db falling back to ~/.local/share/meth/meth.db.
func DefaultDB() (string, error) {
	return xdgDir("XDG_DATA_HOME", filepath.Join(".local", "share"), dbFile)
}

// xdgDir joins $envVar, or ~/fallback when it is unset, with meth/name.
func xdgDir(envVar, fallback, name string) (string, error) {
	dir := os.Getenv(envVar)
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, fallback)
	}
	return filepath.Join(dir, appDir, name), nil
}

// Load reads the config file. A missing file is not an error, it means
// defaults. An unknown key is an error: a typo that silently kept the old
// value would be worse than a failed command.
func Load() (Config, error) {
	cfg := Config{}

	path, err := Path()
	if err != nil {
		return cfg, err
	}

	md, err := toml.DecodeFile(path, &cfg)
	switch {
	case os.IsNotExist(err):
		// no file, every field stays at its zero value and picks up a default below
	case err != nil:
		return cfg, fmt.Errorf("%s: %w", path, err)
	default:
		cfg.File = path
		if undecoded := md.Undecoded(); len(undecoded) > 0 {
			keys := make([]string, len(undecoded))
			for i, k := range undecoded {
				keys[i] = k.String()
			}
			return cfg, fmt.Errorf("%s: unknown keys: %s", path, strings.Join(keys, ", "))
		}
	}

	if cfg.DB == "" {
		if cfg.DB, err = DefaultDB(); err != nil {
			return cfg, err
		}
	}
	if cfg.DB, err = expandHome(cfg.DB); err != nil {
		return cfg, fmt.Errorf("%s: db: %w", path, err)
	}
	return cfg, nil
}

// Write creates the config file, with comments, and returns where it landed.
// It refuses to clobber an existing file.
func Write(cfg Config) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	body := fmt.Sprintf(`# meth config. Delete this file to go back to the defaults.

# The sqlite file holding every metric and value. A leading ~ is expanded.
# Override for a single command with: meth --db ./scratch.db ...
db = %q
`, cfg.DB)

	// O_EXCL so `config init` never eats a config you had already edited
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return path, err
	}
	defer f.Close()
	if _, err := f.WriteString(body); err != nil {
		return path, err
	}
	return path, nil
}

// expandHome turns a leading ~ into the home directory. Anything else is
// handed to the filesystem as written, so a relative path stays relative.
func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~")), nil
}
