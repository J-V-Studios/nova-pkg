// Package config loads ~/.config/nova-pkg/config.toml.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the parsed config file.
type Config struct {
	// Priority orders backends for search output and install/remove
	// preference. Unlisted backends append after in default order.
	Priority []string `toml:"priority"`
}

// defaultFile is written when no config exists yet.
const defaultFile = `# nova-pkg configuration

# Backend priority for search output order and install/remove
# preference. Unlisted backends append after in default order.
# Known backends: pacman, aur, flatpak
priority = ["pacman", "aur", "flatpak"]
`

// Load reads the user's config file, auto-generating defaults when
// missing. See loadFrom.
func Load() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	return loadFrom(filepath.Join(dir, "nova-pkg", "config.toml"))
}

// loadFrom reads the config at path. A missing file is auto-generated
// with defaults, then parsed. A write failure is not fatal: defaults
// apply. A parse failure of an existing file is an error.
func loadFrom(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config: %w", err)
		}
		if werr := writeDefault(path); werr != nil {
			return Config{}, fmt.Errorf("config: cannot write %s: %w", path, werr)
		}
		return Parse(strings.NewReader(defaultFile))
	}
	defer f.Close()
	return Parse(f)
}

// writeDefault creates path (and parent dirs) with defaultFile content.
func writeDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(defaultFile), 0o644)
}

// Parse decodes config from r. Pure function for tests.
func Parse(r io.Reader) (Config, error) {
	var c Config
	if _, err := toml.NewDecoder(r).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("config parse: %w", err)
	}
	return c, nil
}
